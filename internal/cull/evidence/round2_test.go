package evidence

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/schuettc/tackle/internal/cull/cases"
	"github.com/schuettc/tackle/internal/cull/mutate"
)

// A target left evidence_unmeasured is retried by the next run: once it
// is measurable, it settles.
func TestRetryUnmeasured(t *testing.T) {
	needTools(t)
	root := newRepo(t, map[string]string{
		"go.mod":       "module fixture\n\ngo 1.22\n",
		"calc.go":      strings.Replace(calcSrc, "a + b", "a - b", 1), // TestAddSole fails: unmeasurable
		"calc_test.go": settleTest,
	})
	writeLast(t, root, []string{"calc_test.go"}, map[string]verdict{"TestAddSole": {"review", "review_band"}})
	if _, err := Run(context.Background(), root, settleOpts()); err != nil {
		t.Fatal(err)
	}
	_, by := loadLast(t, root)
	if tr := by["TestAddSole"]; tr.Verdict != "review" || tr.Rule != RuleUnmeasured {
		t.Fatalf("first run: %s/%s", tr.Verdict, tr.Rule)
	}
	write(t, root, "calc.go", calcSrc) // fixed
	res, err := Run(context.Background(), root, settleOpts())
	if err != nil {
		t.Fatal(err)
	}
	if res.Settled["go:calc_test.go:TestAddSole"] != "keep" {
		t.Errorf("second run Settled = %v Unmeasured = %v %v", res.Settled, res.Unmeasured, res.Reasons)
	}
	_, by = loadLast(t, root)
	tr := by["TestAddSole"]
	if tr.Verdict != "keep" || tr.Rule != RuleSoleGuard {
		t.Errorf("second run: %s/%s", tr.Verdict, tr.Rule)
	}
	n := 0
	for _, r := range tr.Reasons {
		if strings.HasPrefix(r, "evidence: ") {
			n++
		}
	}
	if n != 1 {
		t.Errorf("reasons = %q, want one evidence line", tr.Reasons)
	}
}

func TestTestSupportNeverMutated(t *testing.T) {
	for rel, want := range map[string]bool{
		"conftest.py": true, "pkg/conftest.py": true, "tests/helpers.py": true,
		"a/test/util.go": true, "src/__tests__/h.ts": true, "testdata/x.go": true,
		"spec/h.ts": true, "calc.py": false, "pkg/tests_util.py": false, "contest.py": false,
	} {
		if got := testSupport(rel); got != want {
			t.Errorf("testSupport(%q) = %v, want %v", rel, got, want)
		}
	}

	needPython(t)
	root := newRepo(t, map[string]string{
		"calc.py":           "def add(a, b):\n    return a + b\n",
		"conftest.py":       "def triple(x):\n    return x * 3\n",
		"tests/__init__.py": "",
		"tests/helpers.py":  "def double(x):\n    return x * 2\n",
		"tests/test_calc.py": "from calc import add\nfrom conftest import triple\nfrom tests.helpers import double\n\n\n" +
			"def test_all():\n    assert add(double(1), triple(1)) == 5\n",
	})
	writeLast(t, root, []string{"tests/test_calc.py"}, map[string]verdict{"test_all": {"review", "review_band"}})
	if _, err := Run(context.Background(), root, settleOpts()); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, ".cull", "evidence.json"))
	if err != nil {
		t.Fatal(err)
	}
	var ev struct {
		Modules []struct {
			Rel string `json:"rel"`
		} `json:"modules"`
	}
	if err := json.Unmarshal(data, &ev); err != nil {
		t.Fatal(err)
	}
	if len(ev.Modules) != 1 || ev.Modules[0].Rel != "calc.py" {
		t.Errorf("modules = %+v, want calc.py only", ev.Modules)
	}
}

func TestMatrixMeasuredByModule(t *testing.T) {
	needTools(t)
	m, err := RunMatrix(context.Background(), calcRepo(t), []Module{calcModule()}, opts(2))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(m.MeasuredByModule, map[string]map[string]bool{"calc.go": m.Baseline}) {
		t.Errorf("MeasuredByModule = %v, Baseline = %v", m.MeasuredByModule, m.Baseline)
	}
}

func mu(mod, desc string, killed ...string) Mutant {
	if killed == nil {
		killed = []string{}
	}
	return Mutant{Module: mod, Site: mutate.Site{Desc: desc}, KilledBy: killed}
}

// A target measured in module a but not b: b's mutants don't count.
func TestMeasurePerModule(t *testing.T) {
	byModule := map[string][]Mutant{
		"a.go": {mu("a.go", "a1", "T", "X")},
		"b.go": {mu("b.go", "b1", "X"), mu("b.go", "b2")},
	}
	measuredBy := map[string]map[string]bool{"a.go": {"T": true, "X": true}, "b.go": {"X": true}}
	baseline := map[string]bool{"T": true, "X": true}
	ev := measure("T", []string{"a.go", "b.go"}, byModule, measuredBy, baseline)
	if !ev.Measured || ev.Planted != 1 || ev.Caught != 1 || ev.Sole != 0 ||
		!reflect.DeepEqual(ev.OtherTests, []string{"X"}) {
		t.Errorf("ev = %+v", ev)
	}
}

func TestMeasureTable(t *testing.T) {
	root := t.TempDir()
	const id = "go:a_test.go:TestT"
	byID := map[string]cases.TestCase{id: {ID: id, Lang: "go", File: "a_test.go"}}
	mods := []string{"a.go"}
	all := map[string]map[string]bool{"a.go": {id: true, "X": true}}
	none := map[string]map[string]bool{"a.go": {"X": true}}
	yes := map[string]bool{id: true, "X": true}
	noT := map[string]bool{"X": true}
	excluded := Mutant{Module: "a.go", Site: mutate.Site{Desc: "e"}}
	for _, c := range []struct {
		name       string
		muts       []Mutant
		measuredBy map[string]map[string]bool
		baseline   map[string]bool
		flaky      bool
		verdict    string
		rule       string
		reason     string
	}{
		{"all excluded", []Mutant{excluded, excluded}, all, yes, false, "review", RuleUnmeasured, reasonExcluded},
		{"flaky", []Mutant{mu("a.go", "m", "X")}, none, noT, true, "review", RuleUnmeasured, reasonFlaky},
		{"not in baseline", []Mutant{mu("a.go", "m", "X")}, none, noT, false, "review", RuleUnmeasured, reasonNoBaseline},
		{"catches nothing", []Mutant{mu("a.go", "m", "X"), mu("a.go", "n")}, all, yes, false, "cut", RuleCatchesNothing, ""},
		{"redundant", []Mutant{mu("a.go", "m", id, "X")}, all, yes, false, "cut", RuleRedundant, ""},
		{"sole", []Mutant{mu("a.go", "m", id), mu("a.go", "n", id, "X")}, all, yes, false, "keep", RuleSoleGuard, ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			byModule := map[string][]Mutant{"a.go": c.muts}
			ev := measure(id, mods, byModule, c.measuredBy, c.baseline)
			if !ev.Measured {
				ev.Reason = unmeasuredReason(root, id, byID, mods, byModule, c.baseline,
					map[string]bool{id: c.flaky}, nil, func(string) bool { return true })
			}
			v, r := settle(ev)
			if v != c.verdict || r != c.rule || ev.Reason != c.reason {
				t.Errorf("got %s/%s %q, want %s/%s %q (ev %+v)", v, r, ev.Reason, c.verdict, c.rule, c.reason, ev)
			}
		})
	}
}

// last.json rewritten while evidence ran (e.g. a concurrent cull check):
// refuse, write nothing.
func TestLastJSONChangedRefused(t *testing.T) {
	needTools(t)
	root := goSettleRepo(t)
	last := filepath.Join(root, ".cull", "last.json")
	var edited []byte
	afterMatrixHook = func() {
		data, _ := os.ReadFile(last)
		edited = append(data, '\n')
		if err := os.WriteFile(last, edited, 0o600); err != nil {
			t.Error(err)
		}
	}
	defer func() { afterMatrixHook = nil }()
	_, err := Run(context.Background(), root, settleOpts())
	if err == nil || !strings.Contains(err.Error(), "last.json") {
		t.Fatalf("err = %v", err)
	}
	if now, _ := os.ReadFile(last); !bytes.Equal(now, edited) {
		t.Error("last.json rewritten")
	}
	if _, err := os.Stat(filepath.Join(root, ".cull", "evidence.json")); !os.IsNotExist(err) {
		t.Error("evidence.json written")
	}
}
