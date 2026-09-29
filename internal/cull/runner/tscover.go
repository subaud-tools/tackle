package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/schuettc/tackle/internal/cull/verify"
)

// v8Coverage is the part of a NODE_V8_COVERAGE file this package reads.
type v8Coverage struct {
	Result []struct {
		URL       string `json:"url"`
		Functions []struct {
			Ranges []v8Range `json:"ranges"`
		} `json:"functions"`
	} `json:"result"`
}

type v8Range struct {
	Start int `json:"startOffset"`
	End   int `json:"endOffset"`
	Count int `json:"count"`
}

// Cover measures per-test-file coverage of module (node:test only) with
// NODE_V8_COVERAGE, running each test file alone. node:test can't
// attribute coverage to single tests, so a unit is a test file: every
// test in it covers what the file ran. A file that never loaded the
// module covers nothing. When V8's offsets provably match the module's
// source (plain JavaScript, or TypeScript run by node's own type
// stripping, which keeps every offset), a mutant is covered if the
// innermost V8 range of any of its bytes ran; otherwise any file that
// loaded the module covers all of it. A test file that can start a
// process (see tsFileSpawns), or did start one while its tests ran (see
// tsSpawnPreload), covers all of it too. With node's test isolation off
// (files share one process), Cover fails. Nothing is global:
// module-level code runs in every file that loads the module.
func (r *tsRunner) Cover(ctx context.Context, dir string, files map[string]map[string]string, module string, timeout time.Duration) (*Coverage, error) {
	if r.kind != "node" {
		return nil, fmt.Errorf("per-file coverage is only measured for node:test (this project runs %s)", r.kind)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if nodeIsolationOff(r.nodeCmd) || nodeIsolationOff(os.Getenv("NODE_OPTIONS")) {
		// Every test file then runs in one process: a file's coverage
		// alone isn't what it runs after the others.
		return nil, fmt.Errorf("node test isolation is off (--test-isolation=none): test files share one process")
	}
	target := filepath.Join(dir, filepath.FromSlash(module))
	src, err := os.ReadFile(target)
	if err != nil {
		return nil, err
	}
	real, err := filepath.EvalSymlinks(target)
	if err != nil {
		return nil, err
	}
	u16 := utf16Offsets(src)
	transformed := r.transforms()
	tmp, err := os.MkdirTemp("", "cull-v8cov-*")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.RemoveAll(tmp) }()
	preload := filepath.Join(tmp, "spawn-mark.cjs")
	if err := os.WriteFile(preload, []byte(tsSpawnPreload), 0o600); err != nil {
		return nil, err
	}
	dctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cov := &Coverage{global: func(int, int) bool { return false }, Passed: map[string]bool{}}
	loaded := false
	for i, f := range sortedFiles(files) {
		covDir := filepath.Join(tmp, fmt.Sprint(i))
		if err := os.MkdirAll(covDir, 0o700); err != nil {
			return nil, err
		}
		mark := filepath.Join(tmp, fmt.Sprintf("%d.spawned", i))
		env := tsCoverEnv(covDir, preload, mark)
		junit := filepath.Join(tmp, fmt.Sprintf("%d.xml", i))
		argv := append(strings.Fields(r.nodeCmd), "--test-reporter=junit", "--test-reporter-destination="+junit, f)
		res := runProcessEnv(dctx, dir, argv, env, timeout)
		if err := procErr(ctx, res, ""); err != nil {
			return nil, err
		}
		o := newOutcome()
		readNodeJUnit(junit, []string{f}, files[f], &o)
		if !o.Invalid {
			maps.Copy(cov.Passed, o.Passed)
		}
		groups, found, err := readV8Dir(covDir, real)
		if err != nil {
			return nil, err
		}
		u := Unit{File: f, hit: func(int, int) bool { return false }}
		for _, id := range files[f] {
			u.IDs = append(u.IDs, id)
		}
		spawns, err := tsFileSpawns(filepath.Join(dir, filepath.FromSlash(f)))
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(mark); err == nil {
			spawns = true
		}
		if spawns {
			// A child that isn't node, or whose env sets NODE_V8_COVERAGE,
			// escapes this coverage.
			u.hit = func(int, int) bool { return true }
		} else if found {
			loaded = true
			if !transformed && allAligned(groups, u16[len(src)], "file://"+real) {
				u.hit = func(start, end int) bool { return v8GroupsSiteRan(groups, u16, start, end) }
			} else {
				u.hit = func(int, int) bool { return true }
			}
		}
		cov.Units = append(cov.Units, u)
	}
	if !loaded {
		return nil, fmt.Errorf("no test file loaded %s in the coverage runs", module)
	}
	return cov, nil
}

// nodeIsolationOff reports whether args (a command or NODE_OPTIONS)
// turn node --test's per-file process isolation off.
func nodeIsolationOff(args string) bool {
	f := strings.Fields(args)
	for i, a := range f {
		switch a {
		case "--test-isolation=none", "--experimental-test-isolation=none":
			return true
		case "--test-isolation", "--experimental-test-isolation":
			if i+1 < len(f) && f[i+1] == "none" {
				return true
			}
		}
	}
	return false
}

// tsSpawnPreload is loaded (NODE_OPTIONS --require) into every node
// process of a coverage run. In a test file's process (NODE_TEST_CONTEXT
// set: node --test's own spawns of test files don't count; or a run with
// test isolation off, where node --test spawns nothing), a call that
// starts a process or a worker creates the file $CULL_TS_SPAWNED.
const tsSpawnPreload = `'use strict';
const mark = process.env.CULL_TS_SPAWNED;
const inTest = !!process.env.NODE_TEST_CONTEXT ||
  process.execArgv.some((a) => /^--(experimental-)?test-isolation=none$/.test(a));
if (mark && inTest) {
  const fs = require('node:fs');
  const hit = () => { try { fs.writeFileSync(mark, ''); } catch {} };
  const cp = require('node:child_process');
  for (const n of ['spawn', 'spawnSync', 'exec', 'execSync', 'execFile', 'execFileSync', 'fork']) {
    const orig = cp[n];
    if (typeof orig !== 'function') continue;
    cp[n] = function (...args) { hit(); return orig.apply(this, args); };
  }
  const wt = require('node:worker_threads');
  const W = wt.Worker;
  wt.Worker = class extends W { constructor(...args) { hit(); super(...args); } };
  require('node:module').syncBuiltinESMExports();
}
`

// tsCoverEnv is the environment of a coverage run: V8 coverage into
// covDir, preload loaded, its mark file at mark.
func tsCoverEnv(covDir, preload, mark string) []string {
	opts := ""
	var env []string
	for _, kv := range verify.FilteredEnv() {
		switch {
		case strings.HasPrefix(kv, "NODE_OPTIONS="):
			opts = strings.TrimPrefix(kv, "NODE_OPTIONS=") + " "
		case strings.HasPrefix(kv, "NODE_V8_COVERAGE="), strings.HasPrefix(kv, "CULL_TS_SPAWNED="):
		default:
			env = append(env, kv)
		}
	}
	q := strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(preload)
	return append(env, "NODE_V8_COVERAGE="+covDir, "CULL_TS_SPAWNED="+mark, "NODE_OPTIONS="+opts+`--require "`+q+`"`)
}

// tsSpawnMarks are source fragments that can start a process or a
// worker.
var tsSpawnMarks = []string{"child_process", "worker_threads", "spawn(", "spawnSync(", "exec(", "execSync(", "execFile(", "execFileSync(", "fork("}

// tsRelImport matches a relative module specifier in an import, export
// ... from, dynamic import or require.
var tsRelImport = regexp.MustCompile(`(?:\bfrom|\bimport|\brequire)\s*\(?\s*['"](\.{1,2}/[^'"]*)['"]`)

// tsFileSpawns reports whether the file at path, or a file it imports
// through relative specifiers (transitively), mentions any of
// tsSpawnMarks. Bare specifiers (npm packages) aren't read: the runtime
// preload sees those.
func tsFileSpawns(path string) (bool, error) {
	seen := map[string]bool{}
	queue := []string{path}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		if seen[p] {
			continue
		}
		seen[p] = true
		src, err := os.ReadFile(p)
		if err != nil {
			if p == path {
				return false, err
			}
			continue
		}
		for _, m := range tsSpawnMarks {
			if strings.Contains(string(src), m) {
				return true, nil
			}
		}
		for _, sm := range tsRelImport.FindAllStringSubmatch(string(src), -1) {
			if r := tsResolve(filepath.Join(filepath.Dir(p), filepath.FromSlash(sm[1]))); r != "" {
				queue = append(queue, r)
			}
		}
	}
	return false, nil
}

// tsResolve finds the file a relative specifier (joined to its
// importer's directory) names: itself, with a script extension added or
// its .js-style extension swapped for a TypeScript one, or a directory's
// index; "" if none.
func tsResolve(base string) string {
	exts := []string{".ts", ".tsx", ".mts", ".cts", ".js", ".jsx", ".mjs", ".cjs"}
	cands := []string{base}
	for _, e := range exts {
		cands = append(cands, base+e)
	}
	swap := map[string][]string{".js": {".ts", ".tsx"}, ".jsx": {".tsx"}, ".mjs": {".mts"}, ".cjs": {".cts"}}
	if ts, ok := swap[filepath.Ext(base)]; ok {
		for _, e := range ts {
			cands = append(cands, strings.TrimSuffix(base, filepath.Ext(base))+e)
		}
	}
	for _, e := range exts {
		cands = append(cands, filepath.Join(base, "index"+e))
	}
	for _, c := range cands {
		if st, err := os.Stat(c); err == nil && st.Mode().IsRegular() {
			return c
		}
	}
	return ""
}

// transforms reports whether the test command (or NODE_OPTIONS) loads a
// transpiler or asks node to transform TypeScript: offsets then differ.
func (r *tsRunner) transforms() bool {
	fields := strings.Fields(r.nodeCmd)
	if len(fields) == 0 || filepath.Base(fields[0]) != "node" {
		return true // tsx, ts-node, ...
	}
	for _, f := range append(fields[1:], strings.Fields(os.Getenv("NODE_OPTIONS"))...) {
		for _, flag := range []string{"--import", "--loader", "--experimental-loader", "--require", "-r", "--experimental-transform-types"} {
			if f == flag || strings.HasPrefix(f, flag+"=") {
				return true
			}
		}
	}
	return false
}

// readV8Dir collects the function ranges V8 reported for the script at
// real (a resolved path) in the coverage files in dir, one group per
// (coverage file, script entry): each process (the test's, a node child
// inheriting NODE_V8_COVERAGE) and each copy of the script reports its
// own counts, which can't be merged range by range.
func readV8Dir(dir, real string) ([][]v8Range, bool, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, false, err
	}
	if len(ents) == 0 {
		return nil, false, fmt.Errorf("node wrote no V8 coverage")
	}
	var out [][]v8Range
	found := false
	for _, e := range ents {
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, false, err
		}
		var c v8Coverage
		if err := json.Unmarshal(b, &c); err != nil {
			return nil, false, fmt.Errorf("unreadable V8 coverage %s: %v", e.Name(), err)
		}
		for _, s := range c.Result {
			u, err := url.Parse(s.URL)
			if err != nil || u.Scheme != "file" {
				continue
			}
			p, err := filepath.EvalSymlinks(u.Path)
			if err != nil || p != real {
				continue
			}
			found = true
			var g []v8Range
			for _, fn := range s.Functions {
				g = append(g, fn.Ranges...)
			}
			out = append(out, g)
		}
	}
	return out, found, nil
}

// aligned reports whether V8's script is the module's source (length n
// UTF-16 units), possibly followed by node's "//# sourceURL=<url>" line.
func aligned(ranges []v8Range, n int, url string) bool {
	end := -1
	for _, r := range ranges {
		if r.Start == 0 && r.End > end {
			end = r.End
		}
	}
	extra := end - n
	tail := len("//# sourceURL=" + url)
	return extra == 0 || (extra >= tail && extra <= tail+2)
}

// allAligned reports whether every group is aligned (see aligned).
func allAligned(groups [][]v8Range, n int, url string) bool {
	for _, g := range groups {
		if !aligned(g, n, url) {
			return false
		}
	}
	return true
}

// v8GroupsSiteRan reports whether the site ran (see v8SiteRan) in any
// group.
func v8GroupsSiteRan(groups [][]v8Range, u16 []int, start, end int) bool {
	for _, g := range groups {
		if v8SiteRan(g, u16, start, end) {
			return true
		}
	}
	return false
}

// v8SiteRan reports whether any byte of [start, end) of the module (its
// start when empty) ran (see v8Ran), given u16 (see utf16Offsets).
func v8SiteRan(ranges []v8Range, u16 []int, start, end int) bool {
	for off := start; off < max(end, start+1); off++ {
		if v8Ran(ranges, u16[off]) {
			return true
		}
	}
	return false
}

// v8Ran reports whether the innermost range holding UTF-16 offset off ran
// (true when no range holds it: not provably unexecuted).
func v8Ran(ranges []v8Range, off int) bool {
	best := -1
	for i, r := range ranges {
		if off >= r.Start && off < r.End && (best < 0 || r.End-r.Start <= ranges[best].End-ranges[best].Start) {
			best = i
		}
	}
	return best < 0 || ranges[best].Count > 0
}

// utf16Offsets maps each byte offset of src (and len(src)) to its UTF-16
// code unit offset; bytes inside a rune map to the rune's start.
func utf16Offsets(src []byte) []int {
	out := make([]int, len(src)+1)
	u := 0
	for i := 0; i < len(src); {
		r, size := utf8.DecodeRune(src[i:])
		for k := 0; k < size; k++ {
			out[i+k] = u
		}
		if n := utf16.RuneLen(r); n > 0 {
			u += n
		} else {
			u++
		}
		i += size
	}
	out[len(src)] = u
	return out
}
