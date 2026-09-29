package runner

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

// goBlock is one coverprofile block: [start, end) byte offsets in the
// module and its execution count.
type goBlock struct {
	start, end int
	count      int
}

// Cover measures per-test coverage of module with Go's coverage
// instrumentation: per package directory, the test binary is built once
// (go test -c -coverpkg=<module's package>) and run once per top-level
// test with its own coverprofile (isolated), and once per test file
// with that file's top-level tests together, in order (a coverprofile
// can't tell those tests apart: each covers the file's union). A unit
// is a top-level test (its subtests share its coverage), covering what
// either run ran. Code in no coverage block (package-level
// declarations and var initializers, signatures) is global; init bodies
// run in every test binary run, so every test covers them. A test
// passed only if both its runs exited 0. Every test of a package whose test
// binary can start a process (see goTestsSpawn) covers every line: the
// child's work isn't in the test's profile. When that can't be told
// (go list fails), Cover fails: the module runs every test.
func (golangRunner) Cover(ctx context.Context, dir string, files map[string]map[string]string, module string, timeout time.Duration) (*Coverage, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	src, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(module)))
	if err != nil {
		return nil, err
	}
	starts := lineStarts(src)
	tmp, err := os.MkdirTemp("", "cull-gocov-*")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.RemoveAll(tmp) }()
	dctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	type test struct{ name, file string }
	bySub := map[string][]test{}
	for _, f := range sortedFiles(files) {
		names, err := goTopLevelTestNames(filepath.Join(dir, filepath.FromSlash(f)))
		if err != nil {
			return nil, err
		}
		sub := path.Dir(filepath.ToSlash(f))
		for _, n := range names {
			bySub[sub] = append(bySub[sub], test{n, f})
		}
	}
	subs := make([]string, 0, len(bySub))
	for s := range bySub {
		subs = append(subs, s)
	}
	sort.Strings(subs)
	coverPkg := "./" + path.Dir(module)
	if path.Dir(module) == "." {
		coverPkg = "./"
	}
	var all []goBlock // every block of the module (any count)
	cov := &Coverage{Passed: map[string]bool{}}
	for i, sub := range subs {
		bin := filepath.Join(tmp, fmt.Sprintf("p%d.test", i))
		pkg := "./"
		if sub != "." {
			pkg = "./" + sub
		}
		spawns, err := goTestsSpawn(dctx, dir, pkg, timeout)
		if err != nil {
			return nil, err
		}
		res := runProcess(dctx, dir, []string{"go", "test", "-c", "-o", bin, "-covermode=atomic", "-coverpkg=" + coverPkg, pkg}, timeout)
		if err := procErr(ctx, res, "go test -c"); err != nil {
			return nil, err
		}
		// profile runs the named top-level tests together, in the
		// binary's order, with their own coverprofile.
		k := 0
		profile := func(names []string) ([]goBlock, bool, error) {
			k++
			prof := filepath.Join(tmp, fmt.Sprintf("p%d-%d.out", i, k))
			res := runProcess(dctx, filepath.Join(dir, filepath.FromSlash(sub)),
				[]string{bin, "-test.run", "^(" + strings.Join(names, "|") + ")$", "-test.count=1", "-test.coverprofile=" + prof}, timeout)
			if err := procErr(ctx, res, ""); err != nil {
				return nil, false, err
			}
			blocks, err := readGoProfile(prof, module, starts)
			return blocks, res.ExitCode == 0, err
		}
		// The in-order pass: each file's tests together, in order.
		byFile := map[string][]string{}
		for _, tc := range bySub[sub] {
			byFile[tc.file] = append(byFile[tc.file], tc.name)
		}
		inOrder := map[string][]goBlock{}
		inOrderOK := map[string]bool{}
		for _, f := range sortedKeys(byFile) {
			blocks, ok, err := profile(byFile[f])
			if err != nil {
				return nil, err
			}
			inOrder[f], inOrderOK[f] = blocks, ok
		}
		if err := goCrossFile(dir, byFile, inOrder, profile); err != nil {
			return nil, err
		}
		// The isolated pass: each test alone.
		for _, tc := range bySub[sub] {
			blocks, ok, err := profile([]string{tc.name})
			if err != nil {
				return nil, err
			}
			if all == nil {
				all = blocks
			}
			hit := append(hitBlocks(blocks), hitBlocks(inOrder[tc.file])...)
			u := Unit{Sel: tc.name, File: tc.file, hit: func(start, end int) bool {
				for _, b := range hit {
					if b.start < max(end, start+1) && start < b.end {
						return true
					}
				}
				return false
			}}
			if spawns {
				u.hit = func(int, int) bool { return true }
			}
			for native, id := range files[tc.file] {
				if native == tc.name || strings.HasPrefix(native, tc.name+"/") {
					u.IDs = append(u.IDs, id)
				}
			}
			sort.Strings(u.IDs)
			if ok && inOrderOK[tc.file] {
				for _, id := range u.IDs {
					cov.Passed[id] = true
				}
			}
			cov.Units = append(cov.Units, u)
		}
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("no coverage block of %s in any test's profile", module)
	}
	cov.global = func(start, end int) bool {
		for _, b := range all {
			if start >= b.start && end <= b.end {
				return false
			}
		}
		return true
	}
	return cov, nil
}

// goCrossFile checks that each test file of one package directory runs
// the same blocks alone (alone, from profile(its tests)) as after the
// files the test binary runs before it: the blocks whose count grows
// from profile(the tests of the files up to it) to profile(those before
// it). A file that doesn't may depend on (or be changed by) another
// file's state; the error names it and its tests.
func goCrossFile(dir string, byFile map[string][]string, alone map[string][]goBlock, profile func([]string) ([]goBlock, bool, error)) error {
	if len(byFile) < 2 {
		return nil
	}
	// The binary runs the package's own test files, then its external
	// (package x_test) ones, each set in file name order.
	var own, ext []string
	for _, f := range sortedKeys(byFile) {
		pf, err := parser.ParseFile(token.NewFileSet(), filepath.Join(dir, filepath.FromSlash(f)), nil, parser.PackageClauseOnly)
		if err != nil {
			return err
		}
		if strings.HasSuffix(pf.Name.Name, "_test") {
			ext = append(ext, f)
		} else {
			own = append(own, f)
		}
	}
	counts := func(blocks []goBlock) map[[2]int]int {
		m := map[[2]int]int{}
		for _, b := range blocks {
			m[[2]int{b.start, b.end}] += b.count
		}
		return m
	}
	ran := func(m map[[2]int]int) map[[2]int]bool {
		out := map[[2]int]bool{}
		for k, n := range m {
			if n > 0 {
				out[k] = true
			}
		}
		return out
	}
	prev := map[[2]int]int{}
	var names []string
	for k, f := range append(own, ext...) {
		names = append(names, byFile[f]...)
		cur := counts(alone[f])
		if k > 0 {
			blocks, _, err := profile(names)
			if err != nil {
				return err
			}
			cur = counts(blocks)
		}
		after := map[[2]int]int{}
		for key, n := range cur {
			after[key] = n - prev[key]
		}
		if !maps.Equal(ran(after), ran(counts(alone[f]))) {
			return fmt.Errorf("%s (%s) runs other lines after the test files before it than alone", f, strings.Join(byFile[f], ", "))
		}
		prev = cur
	}
	return nil
}

// hitBlocks are the blocks that ran.
func hitBlocks(blocks []goBlock) []goBlock {
	var out []goBlock
	for _, b := range blocks {
		if b.count > 0 {
			out = append(out, b)
		}
	}
	return out
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// goTestsSpawn reports whether the test binary of package directory pkg
// (./-relative, in the module at dir) can start a process: a
// non-standard package in its test dependency graph (go list -deps
// -test) imports os/exec, or a file of a package of this module in that
// graph calls os.StartProcess, syscall's Exec, ForkExec or StartProcess,
// or golang.org/x/sys/unix's Exec or ForkExec (qualified or through a
// dot import). An error (go list failed) means it can't be told.
func goTestsSpawn(ctx context.Context, dir, pkg string, timeout time.Duration) (bool, error) {
	res := runProcess(ctx, dir, []string{"go", "list", "-deps", "-test", "-f",
		"{{if not .Standard}}{{.Dir}}\t{{with .Module}}{{.Main}}{{end}}\t{{join .Imports \" \"}}\t{{join .GoFiles \" \"}} {{join .CgoFiles \" \"}}{{end}}", pkg}, timeout)
	if err := procErr(ctx, res, "go list -deps -test "+pkg); err != nil {
		return false, err
	}
	scanned := map[string]bool{}
	for ln := range strings.SplitSeq(string(res.Stdout), "\n") {
		f := strings.Split(ln, "\t")
		if len(f) != 4 {
			continue
		}
		if slices.Contains(strings.Fields(f[2]), "os/exec") {
			return true, nil
		}
		if f[1] != "true" {
			continue
		}
		for _, name := range strings.Fields(f[3]) {
			if filepath.IsAbs(name) {
				continue // generated (the test main)
			}
			p := filepath.Join(f[0], name)
			if scanned[p] {
				continue
			}
			scanned[p] = true
			spawns, err := goFileSpawns(p)
			if err != nil || spawns {
				return spawns, err
			}
		}
	}
	return false, nil
}

// goSpawnCalls are the functions (per import path) that start a process.
var goSpawnCalls = map[string][]string{
	"os":                    {"StartProcess"},
	"syscall":               {"Exec", "ForkExec", "StartProcess"},
	"golang.org/x/sys/unix": {"Exec", "ForkExec"},
}

// goFileSpawns reports whether the Go file at name imports os/exec or
// calls one of goSpawnCalls.
func goFileSpawns(name string) (bool, error) {
	f, err := parser.ParseFile(token.NewFileSet(), name, nil, parser.SkipObjectResolution)
	if err != nil {
		return false, err
	}
	local := map[string]string{} // local package name -> import path
	var dot []string             // functions reachable unqualified
	for _, im := range f.Imports {
		p, _ := strconv.Unquote(im.Path.Value)
		if p == "os/exec" {
			return true, nil
		}
		if _, ok := goSpawnCalls[p]; !ok {
			continue
		}
		n := path.Base(p)
		if im.Name != nil {
			n = im.Name.Name
		}
		if n == "." {
			dot = append(dot, goSpawnCalls[p]...)
		} else {
			local[n] = p
		}
	}
	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.SelectorExpr:
			if x, ok := n.X.(*ast.Ident); ok && slices.Contains(goSpawnCalls[local[x.Name]], n.Sel.Name) {
				found = true
			}
		case *ast.CallExpr:
			if id, ok := n.Fun.(*ast.Ident); ok && slices.Contains(dot, id.Name) {
				found = true
			}
		}
		return !found
	})
	return found, nil
}

// procErr explains why a coverage subprocess can't be used, or is nil.
// A test binary run may fail (a failing test still writes its profile);
// what (non-empty) names a step that must exit 0.
func procErr(ctx context.Context, res procResult, what string) error {
	switch {
	case ctx.Err() != nil:
		return ctx.Err()
	case res.StartErr != nil:
		return res.StartErr
	case res.TimedOut:
		return fmt.Errorf("the coverage run timed out")
	case what != "" && res.ExitCode != 0:
		return fmt.Errorf("%s failed: %s", what, lastLineOf(res.Stderr, res.Stdout))
	}
	return nil
}

// readGoProfile reads the blocks of module (root-relative) from a
// coverprofile: those whose file is <import path>/<module>, i.e. ends
// with "/"+module.
func readGoProfile(prof, module string, starts []int) ([]goBlock, error) {
	data, err := os.ReadFile(prof)
	if err != nil {
		return nil, fmt.Errorf("the test binary wrote no coverprofile")
	}
	off := func(line, col int) int {
		if line < 1 || line > len(starts) {
			return -1
		}
		return starts[line-1] + col - 1
	}
	var out []goBlock
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		ln := sc.Text()
		i := strings.LastIndex(ln, ".go:")
		if i < 0 || !strings.HasSuffix("/"+ln[:i+3], "/"+module) {
			continue
		}
		var l1, c1, l2, c2, n, count int
		if _, err := fmt.Sscanf(ln[i+4:], "%d.%d,%d.%d %d %d", &l1, &c1, &l2, &c2, &n, &count); err != nil {
			return nil, fmt.Errorf("unreadable coverprofile line %q", ln)
		}
		s, e := off(l1, c1), off(l2, c2)
		if s < 0 || e < s {
			return nil, fmt.Errorf("coverprofile block %q is outside %s", ln, module)
		}
		out = append(out, goBlock{s, e, count})
	}
	return out, nil
}

// RunUnits runs the units' top-level tests, per package directory.
func (golangRunner) RunUnits(ctx context.Context, dir string, units []Unit, index map[string]string, timeout time.Duration) Outcome {
	groups := map[string]map[string]bool{}
	for _, u := range units {
		sub := path.Dir(filepath.ToSlash(u.File))
		if groups[sub] == nil {
			groups[sub] = map[string]bool{}
		}
		groups[sub][u.Sel] = true
	}
	return runGoGroups(ctx, dir, groups, index, timeout)
}
