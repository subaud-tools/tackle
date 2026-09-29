package evidence

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/schuettc/tackle/internal/cull/runner"
)

// The reviewer's registry: check() runs only once something is
// registered.
const orderPySrc = "_reg = []\n\n\ndef register(x):\n    _reg.append(x)\n\n\ndef check(x):\n    return x > 0\n\n\ndef all_ok():\n    return all(check(x) for x in _reg)\n"

// test_c reaches check() only through test_b's registration.
const orderPyTests = "from calc import all_ok, register\n\n\ndef test_a():\n    assert all_ok()\n\n\ndef test_b():\n    register(5)\n\n\ndef test_c():\n    assert all_ok()\n"

const orderGoSrc = "package calc\n\nvar reg []int\n\nfunc Register(x int) { reg = append(reg, x) }\n\nfunc check(x int) bool { return x > 0 }\n\nfunc AllOK() bool {\n\tfor _, x := range reg {\n\t\tif !check(x) {\n\t\t\treturn false\n\t\t}\n\t}\n\treturn true\n}\n"

const orderGoTests = "package calc\n\nimport \"testing\"\n\nfunc TestA(t *testing.T) {\n\tif !AllOK() {\n\t\tt.Fatal(\"a\")\n\t}\n}\n\nfunc TestB(t *testing.T) { Register(5) }\n\nfunc TestC(t *testing.T) {\n\tif !AllOK() {\n\t\tt.Fatal(\"c\")\n\t}\n}\n"

// wantSole fails t unless the last settle left name a sole guard.
func wantSole(t *testing.T, root, name string) {
	t.Helper()
	_, by := loadLast(t, root)
	if by[name].Rule != RuleSoleGuard {
		t.Errorf("%s: %s/%s, want a sole guard", name, by[name].Verdict, by[name].Rule)
	}
}

// The in-order coverage pass: a test covers what it ran after the tests
// before it, not only what it runs alone.
func TestCoverageInOrderPassPython(t *testing.T) {
	needCoverage(t)
	root := newRepo(t, map[string]string{"calc.py": orderPySrc, "test_calc.py": orderPyTests})
	ids := []string{"py:test_calc.py:test_a", "py:test_calc.py:test_b", "py:test_calc.py:test_c"}
	checkCovered(t, root, "python", "test_calc.py", "calc.py", ids, orderPySrc, "x > 0", "py:test_calc.py:test_c")
}

func TestCoverageInOrderPassGo(t *testing.T) {
	needTools(t)
	root := newRepo(t, map[string]string{"go.mod": "module fixture\n\ngo 1.22\n", "calc.go": orderGoSrc, "calc_test.go": orderGoTests})
	ids := []string{"go:calc_test.go:TestA", "go:calc_test.go:TestB", "go:calc_test.go:TestC"}
	checkCovered(t, root, "go", "calc_test.go", "calc.go", ids, orderGoSrc, "x > 0", "go:calc_test.go:TestC")
}

// checkCovered fails t unless Cover says the unit of id covers frag.
func checkCovered(t *testing.T, root, lang, file, module string, ids []string, src, frag, id string) {
	t.Helper()
	idx := nativeIndex(lang, file, ids)
	r, err := runner.For(root, lang, file)
	if err != nil {
		t.Fatal(err)
	}
	cov, err := r.(runner.Coverer).Cover(context.Background(), root, map[string]map[string]string{file: idx}, module, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	at := strings.Index(src, frag)
	for _, u := range cov.Units {
		for _, uid := range u.IDs {
			if uid == id && !u.Covers(at, at+len(frag)) {
				t.Errorf("%s doesn't cover %q", id, frag)
			}
		}
	}
}

// The reviewer's repro: verdicts equal the full run's; test_c, alone
// catching check()'s mutants, stays a sole guard.
func TestCoverageOrderStatePython(t *testing.T) {
	needCoverage(t)
	root := newRepo(t, map[string]string{"calc.py": orderPySrc, "test_calc.py": orderPyTests})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"test_calc.py"}, map[string]verdict{"test_a": band, "test_b": band, "test_c": band})
	settleExpect(t, root, "")
	wantSole(t, root, "test_c")
}

// Go: a package-level registry shared by the tests of one file.
func TestCoverageOrderStateGo(t *testing.T) {
	needTools(t)
	root := newRepo(t, map[string]string{"go.mod": "module fixture\n\ngo 1.22\n", "calc.go": orderGoSrc, "calc_test.go": orderGoTests})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"calc_test.go"}, map[string]verdict{"TestA": band, "TestB": band, "TestC": band})
	settleExpect(t, root, "")
	wantSole(t, root, "TestC")
}

// File 1 registers; file 2's test reaches check() only after it: the
// module falls back, naming the file and the test.
func TestCoverageOrderCrossFilePython(t *testing.T) {
	needCoverage(t)
	root := newRepo(t, map[string]string{
		"calc.py":   orderPySrc,
		"test_1.py": "from calc import register\n\n\ndef test_reg():\n    register(5)\n",
		"test_2.py": "from calc import all_ok\n\n\ndef test_c():\n    assert all_ok()\n",
	})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"test_1.py", "test_2.py"}, map[string]verdict{"test_reg": band, "test_c": band})
	settleExpect(t, root, "test_2.py::test_c")
}

func TestCoverageOrderCrossFileGo(t *testing.T) {
	needTools(t)
	root := newRepo(t, map[string]string{
		"go.mod":    "module fixture\n\ngo 1.22\n",
		"calc.go":   orderGoSrc,
		"a_test.go": "package calc\n\nimport \"testing\"\n\nfunc TestReg(t *testing.T) { Register(5) }\n",
		"b_test.go": "package calc\n\nimport \"testing\"\n\nfunc TestC(t *testing.T) {\n\tif !AllOK() {\n\t\tt.Fatal(\"c\")\n\t}\n}\n",
	})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"a_test.go", "b_test.go"}, map[string]verdict{"TestReg": band, "TestC": band})
	settleExpect(t, root, "b_test.go (TestC)")
}

// A mutant runs every test of each file holding a test that covers its
// line (test_h0 covers nothing, yet runs with test_f), and no file
// holding none.
func TestCoverageSelectsWholeFiles(t *testing.T) {
	needCoverage(t)
	l := newRunLog(t, covPySrc)
	root := newRepo(t, map[string]string{
		"calc.py": covPySrc,
		"test_f.py": "from calc import f\n\n\n" +
			pyLogTest(l.path, "test_f", "assert f(1) == 2") +
			pyLogTest(l.path, "test_h0", "assert True"),
		"test_g.py": "from calc import g\n\n\n" +
			pyLogTest(l.path, "test_g1", "assert g(5)") +
			pyLogTest(l.path, "test_g2", "assert not g(1)"),
	})
	mod := Module{Rel: "calc.py", Lang: "python", TestFiles: []string{"test_f.py", "test_g.py"},
		Tests: []string{"py:test_f.py:test_f", "py:test_f.py:test_h0", "py:test_g.py:test_g1", "py:test_g.py:test_g2"}}
	o := opts(1)
	if _, err := RunMatrix(context.Background(), root, []Module{mod}, o); err != nil {
		t.Fatal(err)
	}
	want := map[int][]string{1: {"test_f", "test_g1", "test_g2", "test_h0"}, 5: {"test_f", "test_h0"}, 9: {"test_g1", "test_g2"}}
	if got := l.byLine(t); !reflect.DeepEqual(got, want) {
		t.Errorf("tests run per mutated line:\n got %v\nwant %v", got, want)
	}
}
