package runner

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
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
// test with its own coverprofile. A unit is a top-level test (its
// subtests share its coverage). Code in no coverage block (package-level
// declarations and var initializers, signatures) is global; init bodies
// run in every test binary run, so every test covers them. A test
// passed only if its run exited 0.
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
	defer os.RemoveAll(tmp)
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
		res := runProcess(dctx, dir, []string{"go", "test", "-c", "-o", bin, "-covermode=count", "-coverpkg=" + coverPkg, pkg}, timeout)
		if err := procErr(ctx, res, "go test -c"); err != nil {
			return nil, err
		}
		for j, tc := range bySub[sub] {
			prof := filepath.Join(tmp, fmt.Sprintf("p%d-%d.out", i, j))
			res := runProcess(dctx, filepath.Join(dir, filepath.FromSlash(sub)),
				[]string{bin, "-test.run", "^" + tc.name + "$", "-test.count=1", "-test.coverprofile=" + prof}, timeout)
			if err := procErr(ctx, res, ""); err != nil {
				return nil, err
			}
			blocks, err := readGoProfile(prof, module, starts)
			if err != nil {
				return nil, err
			}
			if all == nil {
				all = blocks
			}
			var hit []goBlock
			for _, b := range blocks {
				if b.count > 0 {
					hit = append(hit, b)
				}
			}
			u := Unit{Sel: tc.name, File: tc.file, hit: func(start, end int) bool {
				for _, b := range hit {
					if b.start < max(end, start+1) && start < b.end {
						return true
					}
				}
				return false
			}}
			for native, id := range files[tc.file] {
				if native == tc.name || strings.HasPrefix(native, tc.name+"/") {
					u.IDs = append(u.IDs, id)
				}
			}
			sort.Strings(u.IDs)
			if res.ExitCode == 0 {
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
