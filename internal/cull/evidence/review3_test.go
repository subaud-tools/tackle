package evidence

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
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
