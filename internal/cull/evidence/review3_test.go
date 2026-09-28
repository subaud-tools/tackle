package evidence

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/schuettc/tackle/internal/cull/cases"
	"github.com/schuettc/tackle/internal/cull/mutate"
)

// After isolation the kept files run together as the real baseline: a
// test that fails only in that combination (test_b, polluted by test_a's
// import-time env var) is flaky, not a killer of every mutant, so
// test_add stays the sole guard of the flipped +.
func TestIsolatedKeptFilesCombinedBaseline(t *testing.T) {
	needPython(t)
	root := newRepo(t, map[string]string{
		"calc.py":     pyCalc,
		"test_a.py":   "import os\n\nos.environ[\"CULL_POLLUTE\"] = \"1\"\nfrom calc import add\n\n\ndef test_add():\n    assert add(2, 3) == 5\n",
		"test_b.py":   "import os\nfrom calc import add\n\n\ndef test_b():\n    assert \"CULL_POLLUTE\" not in os.environ\n    assert add(1, 0) == 1\n",
		"test_bad.py": "import no_such_module_cull\nfrom calc import add\n\n\ndef test_bad():\n    assert add(1, 0) == 1\n",
	})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"test_a.py", "test_b.py", "test_bad.py"}, map[string]verdict{"test_add": band})
	res, err := Run(context.Background(), root, settleOpts())
	if err != nil {
		t.Fatal(err)
	}
	_, by := loadLast(t, root)
	tr := by["test_add"]
	if tr.Verdict == "cut" {
		t.Fatalf("sole guard cut: %s/%s %v (settled %v)", tr.Verdict, tr.Rule, tr.Reasons, res.Settled)
	}
	if tr.Verdict != "keep" && tr.Rule != RuleUnmeasured {
		t.Fatalf("test_add: %s/%s %v", tr.Verdict, tr.Rule, tr.Reasons)
	}
}

// planBaseline prepares one sandbox of root and runs mod's baseline.
func planBaseline(t *testing.T, root string, mod Module) *modPlan {
	t.Helper()
	sbs, err := Prepare(root, 1)
	if err != nil {
		t.Fatal(err)
	}
	abs, _ := filepath.Abs(root)
	p, _ := planModule(abs, sbs[0].Dir, mod, 30, func(string, ...any) {})
	if p == nil {
		t.Fatal("no plan")
	}
	p.baseline(context.Background(), sbs[0], time.Minute)
	return p
}

// A dropped file's slow runs don't scale the mutant timeout: only the
// kept files' runs do.
func TestSlowestFromKeptFilesOnly(t *testing.T) {
	needPython(t)
	root := newRepo(t, map[string]string{
		"calc.py":      pyCalc,
		"test_good.py": "from calc import add\n\n\ndef test_add():\n    assert add(2, 3) == 5\n",
		"test_slow.py": "import time\n\ntime.sleep(3)\nimport no_such_module_cull\n\n\ndef test_slow():\n    pass\n",
	})
	mod := Module{Rel: "calc.py", Lang: "python", TestFiles: []string{"test_good.py", "test_slow.py"},
		Tests: []string{"py:test_good.py:test_add", "py:test_slow.py:test_slow"}}
	p := planBaseline(t, root, mod)
	if !p.ok || !reflect.DeepEqual(p.dropped, []string{"test_slow.py"}) {
		t.Fatalf("ok %v dropped %v", p.ok, p.dropped)
	}
	if p.slowest >= 3*time.Second {
		t.Errorf("slowest = %v, includes the dropped file's runs", p.slowest)
	}
}

// The combined baseline stops at its first unusable run.
func TestCombinedBaselineStopsAtFirstBadRun(t *testing.T) {
	needPython(t)
	counter := filepath.Join(t.TempDir(), "runs")
	root := newRepo(t, map[string]string{
		"calc.py": pyCalc,
		"test_bad.py": "with open(" + pyQuote(counter) + ", \"a\") as f:\n    f.write(\"x\")\nimport no_such_module_cull\n" +
			"from calc import add\n\n\ndef test_bad():\n    assert add(1, 1) == 2\n",
	})
	mod := Module{Rel: "calc.py", Lang: "python", TestFiles: []string{"test_bad.py"}, Tests: []string{"py:test_bad.py:test_bad"}}
	p := planBaseline(t, root, mod)
	if p.ok {
		t.Fatal("baseline ok")
	}
	// One combined run, then one isolated run (which stops at its bad run).
	if b, _ := os.ReadFile(counter); len(b) != 2 {
		t.Errorf("test_bad.py ran %d times, want 2", len(b))
	}
}

func nolog(string, ...any) {}

// A module with CRLF line endings: the extractor's LF Source still finds
// its span (offsets in the CRLF bytes), so the target's sole catch in a
// small function next to a big one is planted.
func TestTargetedSamplingCRLF(t *testing.T) {
	needPython(t)
	crlf := func(s string) string { return strings.ReplaceAll(s, "\n", "\r\n") }
	root := newRepo(t, map[string]string{
		"calc.py": crlf(bigPy() + "\n\ndef ident(x):\n    return x\n"),
		"test_calc.py": "from calc import big, ident\n\n\ndef test_big():\n    assert big(5) == 10\n\n\n" +
			"def test_ident():\n    assert ident(3) == 3\n",
	})
	writeLast(t, root, []string{"test_calc.py"}, map[string]verdict{"test_ident": {"review", "review_band"}})
	o := settleOpts()
	var log bytes.Buffer
	o.Stderr = &log
	if _, err := Run(context.Background(), root, o); err != nil {
		t.Fatal(err)
	}
	_, by := loadLast(t, root)
	if tr := by["test_ident"]; tr.Verdict != "keep" || tr.Rule != RuleSoleGuard {
		t.Fatalf("test_ident: %s/%s %v\n%s", tr.Verdict, tr.Rule, tr.Reasons, log.String())
	}
	if strings.Contains(log.String(), "whole module") {
		t.Errorf("fell back:\n%s", log.String())
	}
}

// Spans: CRLF offsets map back to the original bytes; every occurrence
// of a duplicated definition is included; a fallback is logged.
func TestTargetSitesSpans(t *testing.T) {
	src := []byte("def f():\r\n    return 1\r\nx = 2\r\n")
	sites := []mutate.Site{{Start: 21, End: 22, Desc: "in f"}, {Start: 28, End: 29, Desc: "x"}}
	got := targetSites("m.py", src, sites, []cases.Callee{{Symbol: "f", Source: "def f():\n    return 1\n"}}, nolog)
	if len(got) != 1 || got[0].Desc != "in f" {
		t.Errorf("crlf: %v", got)
	}

	dup := []byte("def g():\n    return 1\n\n\ndef g():\n    return 1\n")
	sites = []mutate.Site{{Start: 20, End: 21, Desc: "first"}, {Start: 44, End: 45, Desc: "second"}}
	got = targetSites("m.py", dup, sites, []cases.Callee{{Symbol: "g", Source: "def g():\n    return 1\n"}}, nolog)
	if len(got) != 2 {
		t.Errorf("duplicate: %v", got)
	}

	var logged []string
	logf := func(f string, a ...any) { logged = append(logged, fmt.Sprintf(f, a...)) }
	targetSites("m.py", dup, sites, []cases.Callee{{Symbol: "h", Source: "def h(): pass\n"}}, logf)
	if len(logged) != 1 || !strings.Contains(logged[0], "m.py") || !strings.Contains(logged[0], "h") ||
		!strings.Contains(logged[0], "whole module") {
		t.Errorf("log = %q", logged)
	}
}

// A method name on several receivers (A.Get, B.Get) resolves to one of
// them at extraction: evidence spans every same-named method in the
// package, so TestA's unique catch in A.Get is planted.
func TestGoAmbiguousMethodSpansAll(t *testing.T) {
	needTools(t)
	root := newRepo(t, map[string]string{
		"go.mod": "module fixture\n\ngo 1.22\n",
		"a.go":   "package calc\n\ntype A struct{}\n\nfunc (A) Get(x int) int { return x + 1 }\n",
		"b.go":   "package calc\n\ntype B struct{}\n\nfunc (B) Get(x int) int { return x * 2 }\n",
		"calc_test.go": "package calc\n\nimport \"testing\"\n\n" +
			"func TestA(t *testing.T) {\n\tif (A{}).Get(1) != 2 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n\n" +
			"func TestB(t *testing.T) {\n\tif (B{}).Get(3) != 6 {\n\t\tt.Fatal(\"b\")\n\t}\n}\n",
	})
	writeLast(t, root, []string{"calc_test.go"}, map[string]verdict{"TestA": {"review", "review_band"}})
	o := settleOpts()
	var log bytes.Buffer
	o.Stderr = &log
	if _, err := Run(context.Background(), root, o); err != nil {
		t.Fatal(err)
	}
	_, by := loadLast(t, root)
	if tr := by["TestA"]; tr.Verdict != "keep" || tr.Rule != RuleSoleGuard {
		t.Fatalf("TestA: %s/%s %v\n%s", tr.Verdict, tr.Rule, tr.Reasons, log.String())
	}
	if !strings.Contains(log.String(), "Get") || !strings.Contains(log.String(), "receivers") {
		t.Errorf("ambiguity not logged:\n%s", log.String())
	}
}

// A measured test that passes a valid canary run provably doesn't load
// this copy of the module: it leaves the module's measured set.
func TestCanaryPasserUnmeasured(t *testing.T) {
	needPython(t)
	root := newRepo(t, map[string]string{
		"calc.py":   pyCalc,
		"test_a.py": "def test_add():\n    from calc import add\n    assert add(2, 3) == 5\n",
		"test_b.py": "def test_b():\n    assert True\n",
	})
	a, b := "py:test_a.py:test_add", "py:test_b.py:test_b"
	mod := Module{Rel: "calc.py", Lang: "python", TestFiles: []string{"test_a.py", "test_b.py"}, Tests: []string{a, b}}
	m, err := RunMatrix(context.Background(), root, []Module{mod}, opts(1))
	if err != nil {
		t.Fatal(err)
	}
	if got := m.MeasuredByModule["calc.py"]; !reflect.DeepEqual(got, map[string]bool{a: true}) {
		t.Errorf("measured = %v, want only %s", got, a)
	}
	if !reflect.DeepEqual(m.Unloaded["calc.py"], []string{b}) {
		t.Errorf("Unloaded = %v", m.Unloaded)
	}
}

// A target measured in one module but unmeasured in another it exercises
// (here util.py is shadowed) is unmeasured, never settled on the first
// module alone: its sole catch could be in the unmeasured one.
func TestPartialMeasurementUnmeasured(t *testing.T) {
	needPython(t)
	outside := t.TempDir()
	write(t, outside, "util.py", "def dbl(x):\n    return x * 2\n")
	root := newRepo(t, map[string]string{
		"calc.py":     pyCalc,
		"util.py":     "def dbl(x):\n    return x * 2\n",
		"conftest.py": "import sys\nsys.path.insert(0, " + pyQuote(outside) + ")\n",
		"tests/test_calc.py": "from calc import add\nfrom util import dbl\n\n\n" +
			"def test_both():\n    assert add(2, 3) == 5\n    assert dbl(2) == 4\n\n\n" +
			"def test_add_only():\n    assert add(2, 3) == 5\n",
	})
	writeLast(t, root, []string{"tests/test_calc.py"}, map[string]verdict{"test_both": {"review", "review_band"}})
	res, err := Run(context.Background(), root, settleOpts())
	if err != nil {
		t.Fatal(err)
	}
	id := "py:tests/test_calc.py:test_both"
	if len(res.Settled) != 0 {
		t.Errorf("Settled = %v", res.Settled)
	}
	if res.Reasons[id] != reasonShadowed {
		t.Errorf("reason = %q", res.Reasons[id])
	}
	_, by := loadLast(t, root)
	if tr := by["test_both"]; tr.Verdict != "review" || tr.Rule != RuleUnmeasured {
		t.Errorf("test_both: %s/%s %v", tr.Verdict, tr.Rule, tr.Reasons)
	}
}

// moduleGap: any module the target exercises that is unmeasured for it
// gives the reason; a module with no sites to plant says nothing.
func TestModuleGapTable(t *testing.T) {
	const id, file = "py:t.py:test_x", "t.py"
	ok := []Mutant{mu("a.py", "m", id)}
	excl := []Mutant{{Module: "a.py"}}
	yes := map[string]map[string]bool{"a.py": {id: true}}
	none := map[string]map[string]bool{}
	for _, c := range []struct {
		name string
		v    matrixView
		want string
	}{
		{"measured", matrixView{byModule: map[string][]Mutant{"a.py": ok}, measuredBy: yes}, ""},
		{"shadowed", matrixView{shadowed: map[string]string{"a.py": reasonShadowed}, measuredBy: none}, reasonShadowed},
		{"unconfirmed", matrixView{shadowed: map[string]string{"a.py": reasonUnconfirmed}, measuredBy: none}, reasonUnconfirmed},
		{"dropped", matrixView{dropped: map[string][]string{"a.py": {file}}, measuredBy: none}, reasonDropped},
		{"unloaded", matrixView{unloaded: map[string][]string{"a.py": {id}}, measuredBy: none}, reasonShadowed},
		{"all excluded", matrixView{byModule: map[string][]Mutant{"a.py": excl}, measuredBy: yes}, reasonExcluded},
		{"no baseline", matrixView{measuredBy: none}, reasonNoBaseline},
		{"flaky", matrixView{measuredBy: none, flaky: map[string]bool{id: true}}, reasonFlaky},
		{"no sites", matrixView{measuredBy: none, noSites: map[string]bool{"a.py": true}}, ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := moduleGap(id, file, "a.py", c.v); got != c.want {
				t.Errorf("moduleGap = %q, want %q", got, c.want)
			}
			// A second, fully measured module doesn't hide the gap.
			v := c.v
			v.byModule = map[string][]Mutant{"a.py": v.byModule["a.py"], "b.py": {mu("b.py", "n", id)}}
			mb := map[string]map[string]bool{"b.py": {id: true}}
			maps.Copy(mb, v.measuredBy)
			v.measuredBy = mb
			ev := measure(id, []string{"a.py", "b.py"}, v.byModule, v.measuredBy, map[string]bool{id: true})
			if !ev.Measured {
				t.Fatalf("b.py alone not measured: %+v", ev)
			}
			ev = gateModules(ev, id, file, []string{"a.py", "b.py"}, v)
			if want := c.want == ""; ev.Measured != want || ev.Reason != c.want {
				t.Errorf("gated: measured %v reason %q, want %q", ev.Measured, ev.Reason, c.want)
			}
		})
	}
}

// A same-package func Get shares its bare callee Symbol with the methods
// (A) Get and (B) Get. Whichever callee is seen first must not decide for
// the others: the function keeps its own span (TestFn, its sole guard, is
// kept) and A.Get's caller is spanned across receivers (TestA is kept).
func TestGoFuncAndSameNamedMethods(t *testing.T) {
	needTools(t)
	testA := "func TestA(t *testing.T) {\n\tif (A{}).Get(1) != 2 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n\n"
	testFn := "func TestFn(t *testing.T) {\n\tif Get(3) != 2 {\n\t\tt.Fatal(\"fn\")\n\t}\n}\n\n"
	for _, tc := range []struct{ name, tests string }{
		{"function first", testFn + testA},
		{"method first", testA + testFn},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t, map[string]string{
				"go.mod":       "module fixture\n\ngo 1.22\n",
				"a.go":         "package calc\n\ntype A struct{}\n\nfunc (A) Get(x int) int { return x + 1 }\n",
				"b.go":         "package calc\n\ntype B struct{}\n\nfunc (B) Get(x int) int { return x * 2 }\n",
				"c.go":         "package calc\n\nfunc Get(x int) int { return x - 1 }\n",
				"calc_test.go": "package calc\n\nimport \"testing\"\n\n" + tc.tests,
			})
			writeLast(t, root, []string{"calc_test.go"}, map[string]verdict{
				"TestA":  {"review", "review_band"},
				"TestFn": {"review", "review_band"},
			})
			o := settleOpts()
			var log bytes.Buffer
			o.Stderr = &log
			if _, err := Run(context.Background(), root, o); err != nil {
				t.Fatal(err)
			}
			_, by := loadLast(t, root)
			for _, id := range []string{"TestFn", "TestA"} {
				if tr := by[id]; tr.Verdict != "keep" || tr.Rule != RuleSoleGuard {
					t.Errorf("%s: %s/%s %v", id, tr.Verdict, tr.Rule, tr.Reasons)
				}
			}
			if !strings.Contains(log.String(), "A.Get in a.go") {
				t.Errorf("A.Get not spanned:\n%s", log.String())
			}
			if t.Failed() {
				t.Log(log.String())
			}
		})
	}
}
