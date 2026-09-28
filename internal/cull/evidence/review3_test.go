package evidence

import (
	"bytes"
	"context"
	"fmt"
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
	p := planModule(abs, sbs[0].Dir, mod, 30, func(string, ...any) {})
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
