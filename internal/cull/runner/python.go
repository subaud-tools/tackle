package runner

import (
	"context"
	"encoding/xml"
	"os"
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

	tmp, err := os.CreateTemp("", "cull-pytest-*.xml")
	if err != nil {
		out.Err = err.Error()
		return out
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	defer os.Remove(tmpPath)

	files := append([]string(nil), testFiles...)
	sort.Strings(files)

	interp := verify.PythonInterpreter(dir)
	argv := append([]string{interp, "-m", "pytest", "-q", "-p", "no:cacheprovider", "--junitxml=" + tmpPath}, files...)

	res := runProcess(ctx, dir, argv, timeout)
	if res.StartErr != nil {
		out.Err = res.StartErr.Error()
		return out
	}
	if res.TimedOut {
		out.TimedOut = true
		return out
	}

	data, err := os.ReadFile(tmpPath)
	if err != nil || len(data) == 0 {
		return out
	}
	var suites junitTestSuites
	if err := xml.Unmarshal(data, &suites); err != nil {
		return out
	}
	for _, s := range suites.Suites {
		for _, c := range s.Cases {
			if c.Skipped != nil {
				continue
			}
			native := pythonNativeKey(c.Classname, c.Name)
			passed := c.Failure == nil && c.Error == nil
			attributePython(index, native, passed, &out)
		}
	}
	return out
}

// pythonNativeKey builds the native lookup key pytest's junitxml classname
// (a dotted module.Class path, or just the module for a bare function) and
// name (the function name, plus "[params]" for a parametrized case) imply:
// "Class::func" for a class method, else "func".
func pythonNativeKey(classname, name string) string {
	parts := strings.Split(classname, ".")
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
