package runner

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/schuettc/tackle/internal/cull/verify"
)

// pythonRunner runs pytest with a JUnit XML report: `<interp> -m pytest -q
// -p no:cacheprovider --junitxml=<tmp> <files>`.
type pythonRunner struct{}

// junitTestSuites is the root <testsuites> element pytest's --junitxml
// writes.
type junitTestSuites struct {
	Suites []junitTestSuite `xml:"testsuite"`
}

type junitTestSuite struct {
	Cases []junitTestCase `xml:"testcase"`
}

type junitTestCase struct {
	Classname string    `xml:"classname,attr"`
	Name      string    `xml:"name,attr"`
	Failure   *struct{} `xml:"failure"`
	Error     *struct{} `xml:"error"`
	Skipped   *struct{} `xml:"skipped"`
}

func (pythonRunner) Run(ctx context.Context, dir string, testFiles []string, index map[string]string, timeout time.Duration) Outcome {
	out := newOutcome()
	if len(testFiles) == 0 {
		return out
	}
	if err := ctx.Err(); err != nil {
		out.Err = err.Error()
		return out
	}

	tmp, err := os.CreateTemp("", "cull-pytest-*.xml")
	if err != nil {
		out.Err = err.Error()
		return out
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmpPath) }()

	files := append([]string(nil), testFiles...)
	sort.Strings(files)

	interp := verify.PythonInterpreter(dir)
	argv := append([]string{interp, "-m", "pytest", "-q", "-p", "no:cacheprovider", "--junitxml=" + tmpPath}, files...)

	res := runProcessEnv(ctx, dir, argv, pythonEnv(dir), timeout)
	if res.StartErr != nil {
		out.Err = res.StartErr.Error()
		return out
	}
	if res.TimedOut {
		out.TimedOut = true
		return out
	}
	// pytest exits 2 (usage/collection error), 3 (internal error), or 4
	// (command-line usage error) when it didn't run the requested tests
	// as a whole 	most commonly a collection/import error, which also
	// shows up as <error> testcases below, but the exit code alone is
	// enough even if junitxml still wrote a (partial or empty) report.
	if res.ExitCode == 2 || res.ExitCode == 3 || res.ExitCode == 4 {
		out.Invalid = true
		out.InvalidReason = fmt.Sprintf("pytest exited %d (collection or import error)", res.ExitCode)
	}

	readPyJUnit(dir, tmpPath, index, &out)
	return out
}

// RunEach runs one pytest per test file (see EachRunner).
func (r pythonRunner) RunEach(ctx context.Context, dir string, testFiles []string, index map[string]string, timeout time.Duration) Outcome {
	files := append([]string(nil), testFiles...)
	sort.Strings(files)
	return runEach(ctx, files, timeout, func(ctx context.Context, f string) Outcome {
		return r.Run(ctx, dir, []string{f}, index, timeout)
	})
}

// readPyJUnit adds the results of pytest's junitxml report at path to
// out (nothing when it's missing or unreadable).
func readPyJUnit(dir, path string, index map[string]string, out *Outcome) {
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		return
	}
	var suites junitTestSuites
	if err := xml.Unmarshal(data, &suites); err != nil {
		return
	}
	for _, s := range suites.Suites {
		for _, c := range s.Cases {
			if c.Error != nil {
				// A collection/import error is reported as an <error>
				// testcase, not a real test result: it doesn't belong in
				// Passed/Failed/Unmapped, and it means the run as a whole
				// didn't execute the requested tests.
				out.Invalid = true
				if out.InvalidReason == "" {
					out.InvalidReason = "pytest collection error: " + c.Classname + c.Name
				}
				continue
			}
			if c.Skipped != nil {
				continue
			}
			native := pythonNativeKey(dir, c.Classname, c.Name)
			passed := c.Failure == nil
			attributePython(index, native, passed, out)
		}
	}
}

// pythonNativeKey builds the native lookup key pytest's junitxml classname
// and name imply. classname is the dotted path from the run dir to the
// test: some number of leading segments are the module's file path (e.g.
// "sub.dir.test_b" for sub/dir/test_b.py, whether or not there's an
// __init__.py 	 pytest's junitxml classname is always dir-dotted, not
// package-dotted), and any segments after that are enclosing class names.
// Naively treating the last dotted segment as a class (as if classname
// were always "module.Class") misreads a plain function in a nested
// directory as Class::func. Instead, find the longest prefix of segments
// that resolves to an actual .py file under dir; whatever's left is the
// class chain.
func pythonNativeKey(dir, classname, name string) string {
	if classname == "" {
		return name
	}
	parts := strings.Split(classname, ".")
	for k := len(parts); k >= 1; k-- {
		candidate := filepath.Join(append([]string{dir}, parts[:k]...)...) + ".py"
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			remaining := parts[k:]
			if len(remaining) == 0 {
				return name
			}
			return strings.Join(remaining, ".") + "::" + name
		}
	}
	// No segment prefix resolved to a file under dir (e.g. dir isn't the
	// actual pytest rootdir for some reason): fall back to the old
	// last-segment-is-a-class heuristic rather than mapping nothing.
	if len(parts) >= 2 {
		return parts[len(parts)-1] + "::" + name
	}
	return name
}

// attributePython maps a pytest junitxml native key to a cull id via
// index, falling back to the parametrize-stripped function unit ("func" or
// "Class::func" with any trailing "[...]" removed) when the exact
// (possibly parametrized) key isn't mapped: a parametrized pytest function
// is one unit in extraction, not one per parameter set.
func attributePython(index map[string]string, native string, passed bool, out *Outcome) {
	if id, ok := index[native]; ok {
		markResult(out, id, passed)
		return
	}
	if base := stripParam(native); base != native {
		if id, ok := index[base]; ok {
			markResult(out, id, passed)
			return
		}
	}
	out.Unmapped++
}

// stripParam removes a pytest parametrize suffix like "[1]" or
// "[a-b]" from a native key.
func stripParam(s string) string {
	i := strings.IndexByte(s, '[')
	if i < 0 {
		return s
	}
	return s[:i]
}

// pythonEnv is the filtered environment with bytecode caching off:
// CPython validates a cached .pyc by whole-second source mtime and size,
// so a same-size mutant planted within a second of an earlier version
// could otherwise run that version's bytecode. PYTHONPYCACHEPREFIX is
// dropped so no shared out-of-tree cache is consulted either.
// PYTHONPATH starts with <dir>/src (when it exists) and <dir>, ahead of
// any inherited entries, so the run imports dir's copy of the project
// rather than an installed one (e.g. an editable install of another
// checkout).
func pythonEnv(dir string) []string {
	var env, path []string
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	}
	if fi, err := os.Stat(filepath.Join(dir, "src")); err == nil && fi.IsDir() {
		path = append(path, filepath.Join(dir, "src"))
	}
	path = append(path, dir)
	for _, kv := range verify.FilteredEnv() {
		if strings.HasPrefix(kv, "PYTHONDONTWRITEBYTECODE=") || strings.HasPrefix(kv, "PYTHONPYCACHEPREFIX=") {
			continue
		}
		if old, ok := strings.CutPrefix(kv, "PYTHONPATH="); ok {
			if old != "" {
				path = append(path, old)
			}
			continue
		}
		env = append(env, kv)
	}
	return append(env, "PYTHONDONTWRITEBYTECODE=1", "PYTHONPATH="+strings.Join(path, string(os.PathListSeparator)))
}
