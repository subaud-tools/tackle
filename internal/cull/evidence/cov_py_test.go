package evidence

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// needCoverage skips unless a Python with pytest and coverage.py is
// available: python3, or the venv $CULL_PYCOV names (VIRTUAL_ENV is then
// pointed at it for the test).
func needCoverage(t *testing.T) {
	t.Helper()
	needTools(t)
	py := "python3"
	if v := os.Getenv("CULL_PYCOV"); v != "" {
		t.Setenv("VIRTUAL_ENV", v)
		py = filepath.Join(v, "bin", "python")
	}
	if err := exec.Command(py, "-c", "import coverage, pytest").Run(); err != nil {
		t.Skip("no Python with pytest and coverage.py (set CULL_PYCOV to a venv that has them)")
	}
}

// runLog records, per planted mutant (by its first changed line), the
// tests that ran on it. Tests append their name to the log file; the
// plant hook appends "--- <line>" before each mutant run.
type runLog struct {
	path string
	src  string
}

func newRunLog(t *testing.T, src string) *runLog {
	l := &runLog{path: filepath.Join(t.TempDir(), "runs"), src: src}
	plantHook = func(file string) {
		b, _ := os.ReadFile(file)
		line := firstDiffLine(l.src, string(b))
		f, _ := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		_, _ = fmt.Fprintf(f, "--- %d\n", line)
		_ = f.Close()
	}
	t.Cleanup(func() { plantHook = nil })
	return l
}

func firstDiffLine(a, b string) int {
	al, bl := strings.Split(a, "\n"), strings.Split(b, "\n")
	for i := range al {
		if i >= len(bl) || al[i] != bl[i] {
			return i + 1
		}
	}
	return len(al) + 1
}

// byLine maps each mutated line to the sorted set of tests run on its
// mutants (the union over mutants on that line).
func (l *runLog) byLine(t *testing.T) map[int][]string {
	t.Helper()
	b, _ := os.ReadFile(l.path)
	sets := map[int]map[string]bool{}
	cur := -1
	for _, s := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		if n, ok := strings.CutPrefix(s, "--- "); ok {
			_, _ = fmt.Sscan(n, &cur)
			if sets[cur] == nil {
				sets[cur] = map[string]bool{}
			}
			continue
		}
		if cur >= 0 && s != "" {
			sets[cur][s] = true
		}
	}
	out := map[int][]string{}
	for line, set := range sets {
		out[line] = sortedKeys(set)
	}
	return out
}

// pyLogTest is a test function that logs its name to path when it runs.
func pyLogTest(path, name, body string) string {
	return fmt.Sprintf("def %s():\n    with open(%s, 'a') as fh:\n        fh.write('%s\\n')\n    %s\n\n\n", name, pyQuote(path), name, body)
}

const covPySrc = "LIMIT = 3\n\n\ndef f(x):\n    return x + 1\n\n\ndef g(x):\n    return x > LIMIT\n\n\ndef h(x):\n    return x * 2\n"

func covPyRepo(t *testing.T, log string, extra map[string]string) string {
	files := map[string]string{
		"calc.py": covPySrc,
		"test_calc.py": "from calc import f, g\n\n\n" +
			pyLogTest(log, "test_f", "assert f(1) == 2") +
			pyLogTest(log, "test_g1", "assert g(5)") +
			pyLogTest(log, "test_g2", "assert not g(1)"),
	}
	for k, v := range extra {
		files[k] = v
	}
	return newRepo(t, files)
}

func covPyModule() Module {
	return Module{Rel: "calc.py", Lang: "python", TestFiles: []string{"test_calc.py"},
		Tests: []string{"py:test_calc.py:test_f", "py:test_calc.py:test_g1", "py:test_calc.py:test_g2"}}
}

// Each mutant runs the whole file of the tests whose coverage includes
// its line (see TestCoverageSelectsWholeFiles for a file it skips); a
// module-level mutant runs them all; a mutant on a line no test runs is
// not run (uncovered, killed by nobody). The kill matrix equals the one
// running every test on every mutant.
func TestCoverageSelectsTests(t *testing.T) {
	needCoverage(t)
	l := newRunLog(t, covPySrc)
	root := covPyRepo(t, l.path, nil)
	var log bytes.Buffer
	o := opts(1)
	o.Stderr = &log
	m, err := RunMatrix(context.Background(), root, []Module{covPyModule()}, o)
	if err != nil {
		t.Fatal(err)
	}
	all := []string{"test_f", "test_g1", "test_g2"}
	want := map[int][]string{1: all, 5: all, 9: all}
	if got := l.byLine(t); !reflect.DeepEqual(got, want) {
		t.Errorf("tests run per mutated line:\n got %v\nwant %v\n%s", got, want, log.String())
	}
	uncovered := 0
	for _, mu := range m.Mutants {
		if strings.HasPrefix(mu.Site.Desc, "line 13:") {
			if !mu.Uncovered || mu.KilledBy == nil || len(mu.KilledBy) != 0 {
				t.Errorf("%s: uncovered %v, killed by %#v", mu.Site.Desc, mu.Uncovered, mu.KilledBy)
			}
			uncovered++
		} else if mu.Uncovered {
			t.Errorf("%s: uncovered", mu.Site.Desc)
		}
	}
	if uncovered == 0 {
		t.Error("no mutant on line 13")
	}

	// Same kill matrix as running every test on every mutant.
	plantHook = nil
	o.NoCoverage = true
	full, err := RunMatrix(context.Background(), root, []Module{covPyModule()}, o)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(rows(m), rows(full)) {
		t.Errorf("selected:\n %v\nfull:\n %v", rows(m), rows(full))
	}
}

// Coverage that can't be collected (here the module is imported before
// coverage starts, by a pytest plugin named in addopts) falls back to running every test
// on every mutant, and says why once on stderr.
func TestCoverageFallback(t *testing.T) {
	needPython(t)
	l := newRunLog(t, covPySrc)
	root := covPyRepo(t, l.path, map[string]string{
		"pytest.ini": "[pytest]\naddopts = -p early\n",
		"early.py":   "try:\n    import calc\nexcept Exception:\n    pass\n",
	})
	var log bytes.Buffer
	o := opts(1)
	o.Stderr = &log
	if _, err := RunMatrix(context.Background(), root, []Module{covPyModule()}, o); err != nil {
		t.Fatal(err)
	}
	all := []string{"test_f", "test_g1", "test_g2"}
	got := l.byLine(t)
	lines := []int{}
	for line, tests := range got {
		lines = append(lines, line)
		if !reflect.DeepEqual(tests, all) {
			t.Errorf("line %d ran %v", line, tests)
		}
	}
	sort.Ints(lines)
	if !reflect.DeepEqual(lines, []int{1, 5, 9, 13}) {
		t.Errorf("mutated lines %v", lines)
	}
	if n := strings.Count(log.String(), "running every test file per mutant"); n != 1 {
		t.Errorf("fallback logged %d times:\n%s", n, log.String())
	}
}

// Uncovered mutants count as planted for nobody.
func TestUncoveredNotPlanted(t *testing.T) {
	unc := Mutant{Module: "a.go", KilledBy: []string{}, Uncovered: true}
	byModule := map[string][]Mutant{"a.go": {mu("a.go", "m", "T"), unc, unc}}
	measuredBy := map[string]map[string]bool{"a.go": {"T": true, "X": true}}
	baseline := map[string]bool{"T": true, "X": true}
	if ev := measure("T", []string{"a.go"}, byModule, measuredBy, baseline); ev.Planted != 1 || ev.Sole != 1 {
		t.Errorf("T: %+v", ev)
	}
	if ev := measure("X", []string{"a.go"}, byModule, measuredBy, baseline); ev.Planted != 1 || ev.Caught != 0 {
		t.Errorf("X: %+v", ev)
	}
	only := map[string][]Mutant{"a.go": {unc}}
	ev := measure("X", []string{"a.go"}, only, measuredBy, baseline)
	if ev.Measured || ev.Planted != 0 {
		t.Errorf("only uncovered: %+v", ev)
	}
}
