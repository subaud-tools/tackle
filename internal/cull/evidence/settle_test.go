package evidence

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/schuettc/tackle/internal/cull/check"
	"github.com/schuettc/tackle/internal/cull/extract"
)

const settleTest = `package calc

import "testing"

func TestAddSole(t *testing.T) {
	if Add(2, 3) != 5 {
		t.Fatal("add")
	}
}

func TestPosA(t *testing.T) {
	if !Pos(1) || Pos(0) {
		t.Fatal("pos")
	}
}

func TestPosB(t *testing.T) {
	if !Pos(1) || Pos(0) {
		t.Fatal("pos b")
	}
}

func TestPosC(t *testing.T) {
	if !Pos(1) || Pos(0) {
		t.Fatal("pos c")
	}
}

func TestPosNoAssert(t *testing.T) {
	_ = Pos(5)
}

func TestConst(t *testing.T) {
	if 1+1 != 2 {
		t.Fatal("math")
	}
}

func TestTrivial(t *testing.T) {
	t.Log("hi")
}
`

// verdict is a hand-written last.json entry: verdict and rule by test name.
type verdict struct{ v, rule string }

// writeLast writes a hand-written .cull/last.json for root: every test
// extracted from files, with the verdicts named in vs (by test name; the
// rest keep/"keep_confident") and the files' current hashes.
func writeLast(t *testing.T, root string, files []string, vs map[string]verdict) {
	t.Helper()
	r := check.Report{Root: root, Mode: "suite", Files: map[string]check.FileInfo{}}
	for _, rel := range files {
		res, err := extract.ForFile(rel).Extract(root, []string{rel}, 0)
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(data)
		fi := check.FileInfo{SHA256: hex.EncodeToString(sum[:])}
		for _, c := range res.Cases {
			v, ok := vs[c.Name]
			if !ok {
				v = verdict{"keep", "keep_confident"}
			}
			r.Tests = append(r.Tests, check.TestResult{TestCase: c, Verdict: v.v, Rule: v.rule, Reasons: []string{"jev"}})
			fi.Tests = append(fi.Tests, check.FileTest{ID: c.ID, Hash: c.Hash})
		}
		r.Files[rel] = fi
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	write(t, root, ".cull/last.json", string(data))
}

func goSettleRepo(t *testing.T) string {
	root := newRepo(t, map[string]string{
		"go.mod":       "module fixture\n\ngo 1.22\n",
		"calc.go":      calcSrc,
		"calc_test.go": settleTest,
	})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"calc_test.go"}, map[string]verdict{
		"TestAddSole":     band,
		"TestPosB":        band,
		"TestPosNoAssert": band,
		"TestConst":       band,
		"TestPosA":        {"keep", "keep_confident"},
		"TestPosC":        {"review", "truncated"},
		"TestTrivial":     {"cut", "cut_confident"},
	})
	return root
}

func loadLast(t *testing.T, root string) (check.Report, map[string]check.TestResult) {
	t.Helper()
	r, err := check.LoadReport(root)
	if err != nil {
		t.Fatal(err)
	}
	by := map[string]check.TestResult{}
	for _, tr := range r.Tests {
		by[tr.Name] = tr
	}
	return r, by
}

func settleOpts() Options {
	return Options{Workers: 2, MaxMutants: 30, RunTimeout: 60 * time.Second, Stderr: &bytes.Buffer{}, NoCache: true}
}

func perm(t *testing.T, p string) os.FileMode {
	t.Helper()
	fi, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	return fi.Mode().Perm()
}

func TestSettleGo(t *testing.T) {
	needTools(t)
	root := goSettleRepo(t)
	before := hashTree(t, root)
	res, err := Run(context.Background(), root, settleOpts())
	if err != nil {
		t.Fatal(err)
	}
	after := hashTree(t, root)
	for f, h := range before {
		if !strings.HasPrefix(f, filepath.Join(root, ".cull")) && after[f] != h {
			t.Errorf("working tree file %s changed", f)
		}
	}

	wantSettled := map[string]string{
		"go:calc_test.go:TestAddSole":     "keep",
		"go:calc_test.go:TestPosB":        "cut",
		"go:calc_test.go:TestPosNoAssert": "cut",
	}
	if !reflect.DeepEqual(res.Settled, wantSettled) {
		t.Errorf("Settled = %v", res.Settled)
	}
	if !reflect.DeepEqual(res.Unmeasured, []string{"go:calc_test.go:TestConst"}) {
		t.Errorf("Unmeasured = %v", res.Unmeasured)
	}
	if res.Modules != 1 || res.Mutants != 3 {
		t.Errorf("Modules %d Mutants %d, want 1 3", res.Modules, res.Mutants)
	}

	r, by := loadLast(t, root)
	want := map[string]struct{ v, rule, reason string }{
		"TestAddSole":     {"keep", "evidence_sole_guard", "evidence: caught 1/3, sole 1"},
		"TestPosB":        {"cut", "evidence_redundant", "evidence: caught 2/3, sole 0"},
		"TestPosNoAssert": {"cut", "evidence_catches_nothing", "evidence: caught 0/3, sole 0"},
		"TestConst":       {"review", "evidence_unmeasured", "evidence: not measured (no project code to plant bugs in)"},
		"TestPosA":        {"keep", "keep_confident", "jev"},
		"TestPosC":        {"review", "truncated", "jev"},
		"TestTrivial":     {"cut", "cut_confident", "jev"},
	}
	for name, w := range want {
		tr := by[name]
		last := ""
		if len(tr.Reasons) > 0 {
			last = tr.Reasons[len(tr.Reasons)-1]
		}
		if tr.Verdict != w.v || tr.Rule != w.rule || last != w.reason {
			t.Errorf("%s: %s/%s %q, want %s/%s %q", name, tr.Verdict, tr.Rule, last, w.v, w.rule, w.reason)
		}
	}
	for _, name := range []string{"TestPosA", "TestPosC", "TestTrivial"} {
		if by[name].Evidence != nil || len(by[name].Reasons) != 1 {
			t.Errorf("%s touched: %+v", name, by[name])
		}
	}

	sole := by["TestAddSole"].Evidence
	if sole == nil || !sole.Measured || sole.Planted != 3 || sole.Caught != 1 || sole.Sole != 1 ||
		!reflect.DeepEqual(sole.CaughtBy, []string{"calc.go line 3: flip +"}) ||
		!reflect.DeepEqual(sole.SoleFor, []string{"calc.go line 3: flip +"}) ||
		!reflect.DeepEqual(sole.OtherTests, []string{"go:calc_test.go:TestPosA", "go:calc_test.go:TestPosB", "go:calc_test.go:TestPosC"}) {
		t.Errorf("TestAddSole evidence = %+v", sole)
	}
	red := by["TestPosB"].Evidence
	if red == nil || !red.Measured || red.Caught != 2 || red.Sole != 0 ||
		!reflect.DeepEqual(red.OtherTests, []string{"go:calc_test.go:TestAddSole", "go:calc_test.go:TestPosA", "go:calc_test.go:TestPosC"}) {
		t.Errorf("TestPosB evidence = %+v", red)
	}
	un := by["TestConst"].Evidence
	if un == nil || un.Measured || un.Reason != "no project code to plant bugs in" {
		t.Errorf("TestConst evidence = %+v", un)
	}

	wantSum := map[string]int{"keep": 2, "cut": 3, "review": 2}
	for k, n := range wantSum {
		if r.Summary[k] != n {
			t.Errorf("summary[%s] = %d, want %d (%v)", k, r.Summary[k], n, r.Summary)
		}
	}

	last := filepath.Join(root, ".cull", "last.json")
	evj := filepath.Join(root, ".cull", "evidence.json")
	if p := perm(t, last); p != 0o600 {
		t.Errorf("last.json perm %o", p)
	}
	if p := perm(t, evj); p != 0o600 {
		t.Errorf("evidence.json perm %o", p)
	}
	var ev struct {
		Modules []struct {
			Rel   string   `json:"rel"`
			Tests []string `json:"tests"`
		} `json:"modules"`
		Mutants []Mutant `json:"mutants"`
	}
	data, err := os.ReadFile(evj)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &ev); err != nil {
		t.Fatal(err)
	}
	if len(ev.Mutants) != 3 || len(ev.Modules) != 1 || ev.Modules[0].Rel != "calc.go" || len(ev.Modules[0].Tests) != 7 {
		t.Errorf("evidence.json = %s", data)
	}
	ents, _ := os.ReadDir(filepath.Join(root, ".cull"))
	if len(ents) != 2 {
		t.Errorf(".cull has %d entries, want last.json and evidence.json only", len(ents))
	}
}

func TestSettleChangedFileRefused(t *testing.T) {
	needTools(t)
	root := goSettleRepo(t)
	write(t, root, "calc_test.go", settleTest+"\n// edited\n")
	last := filepath.Join(root, ".cull", "last.json")
	orig, _ := os.ReadFile(last)
	_, err := Run(context.Background(), root, settleOpts())
	if err == nil || !strings.Contains(err.Error(), "changed since cull check") {
		t.Fatalf("err = %v", err)
	}
	if now, _ := os.ReadFile(last); !bytes.Equal(now, orig) {
		t.Error("last.json rewritten")
	}
	if _, err := os.Stat(filepath.Join(root, ".cull", "evidence.json")); !os.IsNotExist(err) {
		t.Error("evidence.json written")
	}
}

// A module edited while the matrix runs: nothing is written [RF 5].
func TestSettleChangedDuringRunRefused(t *testing.T) {
	needTools(t)
	root := goSettleRepo(t)
	afterMatrixHook = func() { write(t, root, "calc.go", calcSrc+"\n// edited\n") }
	defer func() { afterMatrixHook = nil }()
	last := filepath.Join(root, ".cull", "last.json")
	orig, _ := os.ReadFile(last)
	_, err := Run(context.Background(), root, settleOpts())
	if err == nil || !strings.Contains(err.Error(), "calc.go") {
		t.Fatalf("err = %v", err)
	}
	if now, _ := os.ReadFile(last); !bytes.Equal(now, orig) {
		t.Error("last.json rewritten")
	}
	if _, err := os.Stat(filepath.Join(root, ".cull", "evidence.json")); !os.IsNotExist(err) {
		t.Error("evidence.json written")
	}
}

func TestSettleNoTargets(t *testing.T) {
	needTools(t)
	root := newRepo(t, map[string]string{
		"go.mod":       "module fixture\n\ngo 1.22\n",
		"calc.go":      calcSrc,
		"calc_test.go": settleTest,
	})
	writeLast(t, root, []string{"calc_test.go"}, nil)
	last := filepath.Join(root, ".cull", "last.json")
	orig, _ := os.ReadFile(last)
	res, err := Run(context.Background(), root, settleOpts())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Settled) != 0 || len(res.Unmeasured) != 0 || res.Modules != 0 {
		t.Errorf("res = %+v", res)
	}
	if now, _ := os.ReadFile(last); !bytes.Equal(now, orig) {
		t.Error("last.json rewritten")
	}
}

func TestSettleNoLastJSON(t *testing.T) {
	root := t.TempDir()
	if _, err := Run(context.Background(), root, settleOpts()); err == nil || !strings.Contains(err.Error(), "cull check") {
		t.Fatalf("err = %v", err)
	}
}

func TestSettlePython(t *testing.T) {
	needPython(t)
	root := newRepo(t, map[string]string{
		"calc.py": "def add(a, b):\n    return a + b\n\n\ndef pos(n):\n    return n > 0\n",
		"test_calc.py": "from calc import add, pos\n\n\ndef test_add():\n    assert add(2, 3) == 5\n\n\n" +
			"def test_pos_noassert():\n    pos(5)\n\n\ndef test_const():\n    assert 1 + 1 == 2\n",
	})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"test_calc.py"}, map[string]verdict{
		"test_add": band, "test_pos_noassert": band, "test_const": band,
	})
	res, err := Run(context.Background(), root, settleOpts())
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"py:test_calc.py:test_add": "keep", "py:test_calc.py:test_pos_noassert": "cut"}
	if !reflect.DeepEqual(res.Settled, want) {
		t.Errorf("Settled = %v", res.Settled)
	}
	_, by := loadLast(t, root)
	if by["test_add"].Rule != "evidence_sole_guard" || by["test_pos_noassert"].Rule != "evidence_catches_nothing" ||
		by["test_const"].Rule != "evidence_unmeasured" {
		t.Errorf("rules: %s %s %s", by["test_add"].Rule, by["test_pos_noassert"].Rule, by["test_const"].Rule)
	}
}
