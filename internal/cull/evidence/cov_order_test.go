package evidence

import (
	"context"
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
