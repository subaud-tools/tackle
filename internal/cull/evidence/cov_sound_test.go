package evidence

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/schuettc/tackle/internal/cull/runner"
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

// C1: test_b alone checks g, in a subprocess coverage can't see: a unit
// that spawns covers the whole module, so g's mutants run test_b.
func TestCoverageSubprocessCoversAll(t *testing.T) {
	needCoverage(t)
	root := pySettleRepo(t, soundCalcFG, "import subprocess, sys\nfrom calc import f, g\n\n\ndef test_a():\n    assert f(10)\n    assert not f(3)\n    g(0)\n\n\ndef test_b():\n    assert f(10)\n    subprocess.run([sys.executable, '-c', 'import calc; assert calc.g(3) == 6'], check=True)\n")
	settleExpect(t, root, "")
}

// C1: a spawn outside any test's call phase (a fixture) can't be
// attributed: the module falls back.
func TestCoverageFixtureSpawnFallsBack(t *testing.T) {
	needCoverage(t)
	root := pySettleRepo(t, soundCalcFG, "import subprocess, sys\nimport pytest\nfrom calc import f, g\n\n\n@pytest.fixture\ndef checked():\n    subprocess.run([sys.executable, '-c', 'import calc; assert calc.g(3) == 6'], check=True)\n\n\ndef test_a():\n    assert f(10)\n    assert not f(3)\n    g(0)\n\n\ndef test_b(checked):\n    assert f(10)\n")
	settleExpect(t, root, "spawned")
}

// C1 (Go): TestB alone checks G, in a re-run of its own test binary: a
// test package that uses os/exec covers the whole module.
func TestGoCoverageSubprocessCoversAll(t *testing.T) {
	needTools(t)
	src := "package calc\n\nfunc F(x int) bool { return x > 3 }\n\nfunc G(x int) int { return x * 2 }\n"
	test := `package calc

import (
	"os"
	"os/exec"
	"testing"
)

func TestA(t *testing.T) {
	if !F(10) || F(3) {
		t.Fatal("f")
	}
	G(0)
}

func TestB(t *testing.T) {
	if !F(10) {
		t.Fatal("f")
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestHelper$")
	cmd.Env = append(os.Environ(), "CALC_HELPER=1")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestHelper(t *testing.T) {
	if os.Getenv("CALC_HELPER") == "" {
		return
	}
	if G(3) != 6 {
		os.Exit(1)
	}
}
`
	root := newRepo(t, map[string]string{"go.mod": "module fixture\n\ngo 1.22\n", "calc.go": src, "calc_test.go": test})
	mod := Module{Rel: "calc.go", Lang: "go", TestFiles: []string{"calc_test.go"}, Tests: []string{"go:calc_test.go:TestA", "go:calc_test.go:TestB", "go:calc_test.go:TestHelper"}}
	matrixMatchesFull(t, root, mod)
}

// C1 (Go): as above, but the re-run lives in a helper package the test
// imports: its os/exec import (in the test's dependency graph) makes
// every test of the package cover the whole module.
func TestGoCoverageHelperSpawnCoversAll(t *testing.T) {
	needTools(t)
	src := "package calc\n\nfunc F(x int) bool { return x > 3 }\n\nfunc G(x int) int { return x * 2 }\n"
	helper := "package testutil\n\nimport (\n\t\"os\"\n\t\"os/exec\"\n)\n\nfunc ReRun(name string) error {\n\tcmd := exec.Command(os.Args[0], \"-test.run=^\"+name+\"$\")\n\tcmd.Env = append(os.Environ(), \"CALC_HELPER=1\")\n\treturn cmd.Run()\n}\n"
	test := "package calc\n\nimport (\n\t\"os\"\n\t\"testing\"\n\n\t\"fixture/testutil\"\n)\n\nfunc TestA(t *testing.T) {\n\tif !F(10) || F(3) {\n\t\tt.Fatal(\"f\")\n\t}\n\tG(0)\n}\n\nfunc TestB(t *testing.T) {\n\tif !F(10) {\n\t\tt.Fatal(\"f\")\n\t}\n\tif err := testutil.ReRun(\"TestHelper\"); err != nil {\n\t\tt.Fatal(err)\n\t}\n}\n\nfunc TestHelper(t *testing.T) {\n\tif os.Getenv(\"CALC_HELPER\") == \"\" {\n\t\treturn\n\t}\n\tif G(3) != 6 {\n\t\tos.Exit(1)\n\t}\n}\n"
	root := newRepo(t, map[string]string{"go.mod": "module fixture\n\ngo 1.22\n", "calc.go": src, "calc_test.go": test, "testutil/testutil.go": helper})
	mod := Module{Rel: "calc.go", Lang: "go", TestFiles: []string{"calc_test.go"}, Tests: []string{"go:calc_test.go:TestA", "go:calc_test.go:TestB", "go:calc_test.go:TestHelper"}}
	matrixMatchesFull(t, root, mod)
}

// C1 (TS): b.test.ts alone checks g, in a node child with an empty env
// (NODE_V8_COVERAGE not inherited): a test file mentioning
// child_process covers the whole module.
func TestTSCoverageSubprocessCoversAll(t *testing.T) {
	needTools(t)
	if os.Getenv("CULL_TS") == "" {
		t.Skip("CULL_TS not set")
	}
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not on PATH")
	}
	src := "export function f(x: number): boolean {\n  return x > 3;\n}\n\nexport function g(x: number): number {\n  return x * 2;\n}\n"
	a := "import { test } from 'node:test';\nimport assert from 'node:assert';\nimport { f, g } from './calc.ts';\n\ntest('a', () => {\n  assert.ok(f(10));\n  assert.ok(!f(3));\n  g(0);\n});\n"
	b := "import { test } from 'node:test';\nimport assert from 'node:assert';\nimport { execFileSync } from 'node:child_process';\nimport { f } from './calc.ts';\n\ntest('b', () => {\n  assert.ok(f(10));\n  execFileSync(process.execPath, ['--input-type=module', '-e', \"const m = await import('./calc.ts'); if (m.g(3) !== 6) process.exit(1);\"], { env: {} });\n});\n"
	root := newRepo(t, map[string]string{
		"package.json": `{"type": "module", "scripts": {"test": "node --test"}}`,
		"calc.ts":      src,
		"a.test.ts":    a,
		"b.test.ts":    b,
	})
	mod := Module{Rel: "calc.ts", Lang: "typescript", TestFiles: []string{"a.test.ts", "b.test.ts"}, Tests: []string{"ts:a.test.ts:a", "ts:b.test.ts:b"}}
	matrixMatchesFull(t, root, mod)
}

// TS: b.test.ts checks f in its own process and g in a node child
// started by an imported helper (the child inherits NODE_V8_COVERAGE, so
// its coverage file holds g=1, f=0 while b's holds f=1, g=0).
func TestTSCoverageHelperSpawnMatchesFull(t *testing.T) {
	needTools(t)
	if os.Getenv("CULL_TS") == "" {
		t.Skip("CULL_TS not set")
	}
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not on PATH")
	}
	src := "export function f(x: number): boolean {\n  return x > 3;\n}\n\nexport function g(x: number): number {\n  return x * 2;\n}\n"
	a := "import { test } from 'node:test';\nimport assert from 'node:assert';\nimport { f, g } from './calc.ts';\n\ntest('a', () => {\n  assert.ok(f(10));\n  assert.ok(!f(3));\n  g(0);\n});\n"
	h := "import * as cp from 'node:child_process';\n\nexport function runScript(code: string): void {\n  cp.execFileSync(process.execPath, ['--input-type=module', '-e', code], { env: { PATH: process.env.PATH } });\n}\n"
	b := "import { test } from 'node:test';\nimport assert from 'node:assert';\nimport { runScript } from './helper.ts';\nimport { f } from './calc.ts';\n\ntest('b', () => {\n  assert.ok(f(10));\n  runScript(\"const m = await import('./calc.ts'); if (m.g(3) !== 6) process.exit(1);\");\n});\n"
	root := newRepo(t, map[string]string{
		"package.json": `{"type": "module", "scripts": {"test": "node --test"}}`,
		"calc.ts":      src, "helper.ts": h, "a.test.ts": a, "b.test.ts": b,
	})
	mod := Module{Rel: "calc.ts", Lang: "typescript", TestFiles: []string{"a.test.ts", "b.test.ts"}, Tests: []string{"ts:a.test.ts:a", "ts:b.test.ts:b"}}
	matrixMatchesFull(t, root, mod)
}

// matrixMatchesFull runs mod's matrix with coverage selection (which
// must be used) and without, and fails t unless they're identical.
func matrixMatchesFull(t *testing.T, root string, mod Module) {
	t.Helper()
	var log bytes.Buffer
	o := opts(1)
	o.Stderr = &log
	m, err := RunMatrix(context.Background(), root, []Module{mod}, o)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(log.String(), "no per-test coverage") {
		t.Fatalf("coverage not used:\n%s", log.String())
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

// C2: table() is memoized, so only the first test to call it records its
// body; test_b alone checks 'b'. Each test's coverage is measured in
// its own forked child, so both record it.
func TestCoverageMemoIsolated(t *testing.T) {
	needCoverage(t)
	root := pySettleRepo(t, "import functools\n\n\n@functools.lru_cache(maxsize=None)\ndef table():\n    return {'a': 1, 'b': 5}\n\n\ndef get(k):\n    return table()[k]\n",
		"from calc import get\n\n\ndef test_a():\n    assert get('a') == 1\n\n\ndef test_b():\n    assert get('b') == 5\n")
	settleExpect(t, root, "")
}

// C3: both tests import the module inside their body, so its
// import-time lines (the def and its default) run in the first test
// only; test_b alone checks limit=3 (f(4)).
func TestCoverageLazyImportIsolated(t *testing.T) {
	needCoverage(t)
	root := pySettleRepo(t, "def f(x, limit=3):\n    return x > limit\n",
		"def test_a():\n    import calc\n    assert calc.f(10)\n\n\ndef test_b():\n    import calc\n    assert calc.f(4)\n")
	settleExpect(t, root, "")
}

// Isolating each test in a forked child doesn't change which tests pass
// (fixtures of every scope, classes, parametrization, skips, xfail, a
// failure, output and a teardown finalizer).
func TestCoverageIsolationKeepsResults(t *testing.T) {
	needCoverage(t)
	calc := "def f(x):\n    return x > 3\n"
	tests := `import pytest
from calc import f

STATE = []


@pytest.fixture(scope="session")
def sess():
    return {"n": 0}


@pytest.fixture(scope="module")
def mod():
    yield [1]
    STATE.append("down")


@pytest.fixture
def fn(sess):
    sess["n"] += 1
    return sess["n"]


def test_ok(fn, mod):
    print("out")
    assert f(10) and fn == 1 and mod == [1]


def test_fail():
    assert f(1)


@pytest.mark.parametrize("x,want", [(10, True), (1, False)])
def test_param(x, want):
    assert f(x) == want


@pytest.mark.parametrize("x", [10, 1])
def test_param_bad(x):
    assert f(x)


@pytest.mark.skip
def test_skipped():
    assert False


@pytest.mark.xfail
def test_xfail():
    assert f(1)


class TestK:
    def test_m(self, fn):
        assert f(5) and fn >= 1

    def test_bad(self):
        assert not f(5)
`
	root := newRepo(t, map[string]string{"calc.py": calc, "test_calc.py": tests})
	ids := []string{"py:test_calc.py:test_ok", "py:test_calc.py:test_fail", "py:test_calc.py:test_param", "py:test_calc.py:test_param_bad", "py:test_calc.py:test_skipped", "py:test_calc.py:test_xfail", "py:test_calc.py:TestK::test_m", "py:test_calc.py:TestK::test_bad"}
	idx := nativeIndex("python", "test_calc.py", ids)
	r, err := runner.For(root, "python", "test_calc.py")
	if err != nil {
		t.Fatal(err)
	}
	full := r.Run(context.Background(), root, []string{"test_calc.py"}, idx, time.Minute)
	cov, err := r.(runner.Coverer).Cover(context.Background(), root, map[string]map[string]string{"test_calc.py": idx}, "calc.py", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"py:test_calc.py:test_ok": true, "py:test_calc.py:test_param": true, "py:test_calc.py:TestK::test_m": true}
	if !reflect.DeepEqual(full.Passed, want) {
		t.Fatalf("full run passed %v", full.Passed)
	}
	if !reflect.DeepEqual(cov.Passed, full.Passed) {
		t.Errorf("coverage run passed %v, full run %v", cov.Passed, full.Passed)
	}
}
