package evidence

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// settleExpect runs Run on root with coverage selection and without
// (restoring last.json in between) and fails t unless the settled
// verdicts and rules are identical. With fallback non-empty, the
// coverage run must log that the module fell back, with a reason
// containing fallback; with it empty, the module must not fall back.
func settleExpect(t *testing.T, root, fallback string) {
	t.Helper()
	lastPath := filepath.Join(root, ".cull", "last.json")
	last, err := os.ReadFile(lastPath)
	if err != nil {
		t.Fatal(err)
	}
	var got [2]Result
	for i, noCov := range []bool{false, true} {
		if err := os.WriteFile(lastPath, last, 0o600); err != nil {
			t.Fatal(err)
		}
		o := settleOpts()
		o.NoCoverage = noCov
		var log bytes.Buffer
		o.Stderr = &log
		res, err := Run(context.Background(), root, o)
		if err != nil {
			t.Fatal(err)
		}
		if !noCov {
			s := log.String()
			fell := strings.Contains(s, "no per-test coverage")
			switch {
			case fallback == "" && fell:
				t.Fatalf("coverage not used:\n%s", s)
			case fallback != "" && !fell:
				t.Errorf("coverage used, want a fallback (%s):\n%s", fallback, s)
			case fallback != "" && strings.Count(s, "no per-test coverage") != 1:
				t.Errorf("fallback logged more than once:\n%s", s)
			case fallback != "" && !strings.Contains(s, fallback):
				t.Errorf("fallback reason lacks %q:\n%s", fallback, s)
			}
		}
		got[i] = res
	}
	if !reflect.DeepEqual(got[0].Settled, got[1].Settled) || !reflect.DeepEqual(got[0].Rules, got[1].Rules) {
		t.Errorf("with coverage: %v %v\nfull run:      %v %v", got[0].Settled, got[0].Rules, got[1].Settled, got[1].Rules)
	}
	if len(got[0].Settled) == 0 {
		t.Error("nothing settled")
	}
}

// pySettleRepo is a two-test Python repo (calc.py, test_calc.py) with
// both tests in the review band.
func pySettleRepo(t *testing.T, calc, tests string) string {
	t.Helper()
	root := newRepo(t, map[string]string{"calc.py": calc, "test_calc.py": tests})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"test_calc.py"}, map[string]verdict{"test_a": band, "test_b": band})
	return root
}

const soundCalcFG = "def f(x):\n    return x > 3\n\n\ndef g(x):\n    return x * 2\n"

// C4: a test that turns tracing off loses every later test's coverage
// (test_b alone checks g): the module falls back.
func TestCoverageTracerLostFallsBack(t *testing.T) {
	needCoverage(t)
	root := pySettleRepo(t, soundCalcFG, "import sys\nfrom calc import f, g\n\n\ndef test_a():\n    assert f(10)\n    assert not f(3)\n    sys.settrace(None)\n\n\ndef test_b():\n    assert f(10)\n    assert g(3) == 6\n")
	settleExpect(t, root, "tracer")
}

// C4: a test that fails only under coverage stops recording early: the
// module falls back.
func TestCoverageFailingTestFallsBack(t *testing.T) {
	needCoverage(t)
	root := pySettleRepo(t, soundCalcFG, "import sys\nfrom calc import f, g\n\n\ndef test_a():\n    assert f(10)\n    assert not f(3)\n    g(0)\n\n\ndef test_b():\n    assert f(10)\n    assert sys.gettrace() is None\n    assert g(3) == 6\n")
	settleExpect(t, root, "did not pass in the coverage run")
}

// C4 (Go): a test that passes only after another test (so fails when
// run alone for its coverage) sends the module to the fallback, and the
// matrix equals the full run's.
func TestGoCoverageFailingAloneFallsBack(t *testing.T) {
	needTools(t)
	src := "package calc\n\nvar ready bool\n\nfunc Ready() { ready = true }\n\nfunc Pos(n int) bool { return ready && n > 0 }\n"
	test := "package calc\n\nimport \"testing\"\n\nfunc TestA(t *testing.T) { Ready() }\n\nfunc TestB(t *testing.T) {\n\tif !Pos(1) || Pos(0) {\n\t\tt.Fatal(\"pos\")\n\t}\n}\n"
	root := newRepo(t, map[string]string{"go.mod": "module fixture\n\ngo 1.22\n", "calc.go": src, "calc_test.go": test})
	mod := Module{Rel: "calc.go", Lang: "go", TestFiles: []string{"calc_test.go"}, Tests: []string{"go:calc_test.go:TestA", "go:calc_test.go:TestB"}}
	var log bytes.Buffer
	o := opts(1)
	o.Stderr = &log
	m, err := RunMatrix(context.Background(), root, []Module{mod}, o)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(log.String(), "go:calc_test.go:TestB did not pass in the coverage run") {
		t.Errorf("no fallback for TestB:\n%s", log.String())
	}
	o.NoCoverage = true
	full, err := RunMatrix(context.Background(), root, []Module{mod}, o)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(rows(m), rows(full)) {
		t.Errorf("selected:\n %v\nfull:\n %v", rows(m), rows(full))
	}
}
