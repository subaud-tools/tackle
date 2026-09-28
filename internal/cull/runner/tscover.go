package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"os"
	"path/filepath"
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
// stripping, which keeps every offset), a mutant is covered if its
// innermost V8 range ran; otherwise any file that loaded the module
// covers all of it. A test file that can start a process (see
// tsFileSpawns) covers all of it too. Nothing is global: module-level code runs in every
// file that loads the module.
func (r *tsRunner) Cover(ctx context.Context, dir string, files map[string]map[string]string, module string, timeout time.Duration) (*Coverage, error) {
	if r.kind != "node" {
		return nil, fmt.Errorf("per-file coverage is only measured for node:test (this project runs %s)", r.kind)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
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
	defer os.RemoveAll(tmp)
	dctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cov := &Coverage{global: func(int, int) bool { return false }, Passed: map[string]bool{}}
	loaded := false
	for i, f := range sortedFiles(files) {
		covDir := filepath.Join(tmp, fmt.Sprint(i))
		if err := os.MkdirAll(covDir, 0o700); err != nil {
			return nil, err
		}
		env := append(verify.FilteredEnv(), "NODE_V8_COVERAGE="+covDir)
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
		ranges, found, err := readV8Dir(covDir, real)
		if err != nil {
			return nil, err
		}
		u := Unit{Sel: f, File: f, hit: func(int, int) bool { return false }}
		for _, id := range files[f] {
			u.IDs = append(u.IDs, id)
		}
		spawns, err := tsFileSpawns(filepath.Join(dir, filepath.FromSlash(f)))
		if err != nil {
			return nil, err
		}
		if spawns {
			// A child started with its own env escapes NODE_V8_COVERAGE.
			u.hit = func(int, int) bool { return true }
		} else if found {
			loaded = true
			if !transformed && aligned(ranges, u16[len(src)], "file://"+real) {
				u.hit = func(start, end int) bool {
					return v8Ran(ranges, u16[start]) || (end > start+1 && v8Ran(ranges, u16[end-1]))
				}
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

// tsSpawnMarks are source fragments of a test file that can start a
// process or a worker.
var tsSpawnMarks = []string{"child_process", "worker_threads", "spawn(", "spawnSync(", "exec(", "execSync(", "execFile(", "execFileSync(", "fork("}

// tsFileSpawns reports whether the test file at path mentions any of
// tsSpawnMarks.
func tsFileSpawns(path string) (bool, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	for _, m := range tsSpawnMarks {
		if strings.Contains(string(src), m) {
			return true, nil
		}
	}
	return false, nil
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

// readV8Dir collects every function range V8 reported for the script at
// real (a resolved path) in the coverage files in dir.
func readV8Dir(dir, real string) ([]v8Range, bool, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, false, err
	}
	if len(ents) == 0 {
		return nil, false, fmt.Errorf("node wrote no V8 coverage")
	}
	var out []v8Range
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
			for _, fn := range s.Functions {
				out = append(out, fn.Ranges...)
			}
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

// RunUnits runs the units' test files like Run.
func (r *tsRunner) RunUnits(ctx context.Context, dir string, units []Unit, index map[string]string, timeout time.Duration) Outcome {
	seen := map[string]bool{}
	var files []string
	for _, u := range units {
		if !seen[u.File] {
			seen[u.File] = true
			files = append(files, u.File)
		}
	}
	return r.Run(ctx, dir, files, index, timeout)
}
