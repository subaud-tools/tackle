package runner

import (
	"context"
	"sort"
	"time"
)

// Coverage is which parts of one source file (the module) each test
// executed in one unmutated run of a set of test files. Offsets are UTF-8
// byte offsets into the module's bytes as Cover read them.
type Coverage struct {
	// Units are the selectable tests, at the runner's granularity: a
	// pytest function (every parametrization), a Go top-level test (its
	// subtests too), a TypeScript test file (every test in it).
	Units []Unit
	// global reports whether [start, end) is code every test runs, or
	// whose execution can't be attributed to one test (module import,
	// fixture setup/teardown, Go code outside any coverage block).
	global func(start, end int) bool
	// Passed holds the cull ids that passed in the coverage run. A test
	// that failed there (or wasn't reported) may have stopped recording
	// early: its coverage can't be trusted.
	Passed map[string]bool
	// Ignored lists the programs started outside any test (a setup hook
	// probing a bundled native tool) that coverage didn't count against
	// selection: native executables outside the project that can't run
	// its code.
	Ignored []string
}

// Global reports whether [start, end) of the module runs for every test
// (see Coverage.global). A mutant there must run every test.
func (c *Coverage) Global(start, end int) bool {
	return c.global == nil || c.global(start, end)
}

// Unit is one selectable test (see Coverage.Units).
type Unit struct {
	File string   // the test file holding it, relative to dir
	IDs  []string // the cull ids its results map to
	hit  func(start, end int) bool
}

// Covers reports whether u executed any of [start, end).
func (u Unit) Covers(start, end int) bool { return u.hit != nil && u.hit(start, end) }

// Coverer is a Runner that can measure per-test coverage of a module.
// A mutant then runs, with Run, the whole test files holding a test
// that covers it.
type Coverer interface {
	// Cover runs the test files (the keys of files, each mapped to its
	// own native-name → cull-id index) once, unmutated, in dir, and
	// reports what each test executed of module (dir-relative). An error
	// says why coverage couldn't be collected (tool missing, a bad run,
	// no data at all): the caller then runs every test on every mutant.
	Cover(ctx context.Context, dir string, files map[string]map[string]string, module string, timeout time.Duration) (*Coverage, error)
}

// lineStarts returns the byte offset of each line's start (line 1 at
// index 0).
func lineStarts(src []byte) []int {
	starts := []int{0}
	for i, b := range src {
		if b == '\n' {
			starts = append(starts, i+1)
		}
	}
	return starts
}

// lineOf returns the 1-based line holding byte offset off.
func lineOf(starts []int, off int) int {
	return sort.Search(len(starts), func(i int) bool { return starts[i] > off })
}

func sortedFiles(files map[string]map[string]string) []string {
	out := make([]string, 0, len(files))
	for f := range files {
		out = append(out, f)
	}
	sort.Strings(out)
	return out
}

// lastLine is the 1-based line of the last byte of [start, end).
func lastLine(starts []int, start, end int) int {
	if end > start {
		return lineOf(starts, end-1)
	}
	return lineOf(starts, start)
}
