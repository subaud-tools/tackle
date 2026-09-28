package evidence

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/schuettc/tackle/internal/cull/cases"
	"github.com/schuettc/tackle/internal/cull/mutate"
)

// bigPy is a function with many mutation sites: a module-wide sample of
// a module holding it mostly lands here.
func bigPy() string {
	var b strings.Builder
	b.WriteString("def big(x):\n    t = 0\n")
	for i := 0; i < 40; i++ {
		fmt.Fprintf(&b, "    if x > %d:\n        t = t + %d\n", i, i)
	}
	b.WriteString("    return t\n")
	return b.String()
}

// Mutants are planted inside the functions the targets call: a target's
// one sole-guard bug in a small function next to a large one is found.
func TestTargetedSampling(t *testing.T) {
	needPython(t)
	root := newRepo(t, map[string]string{
		"calc.py": bigPy() + "\n\ndef ident(x):\n    return x\n",
		"test_calc.py": "from calc import big, ident\n\n\ndef test_big():\n    assert big(5) == 10\n\n\n" +
			"def test_ident():\n    assert ident(3) == 3\n",
	})
	writeLast(t, root, []string{"test_calc.py"}, map[string]verdict{"test_ident": {"review", "review_band"}})
	res, err := Run(context.Background(), root, settleOpts())
	if err != nil {
		t.Fatal(err)
	}
	_, by := loadLast(t, root)
	if tr := by["test_ident"]; tr.Verdict != "keep" || tr.Rule != RuleSoleGuard {
		t.Fatalf("test_ident: %s/%s (%d mutants) %v", tr.Verdict, tr.Rule, res.Mutants, tr.Reasons)
	}
}

// A callee whose Source isn't in the module falls back to the whole
// module; sites shared by callees are planted once.
func TestTargetSitesFallbackDedupe(t *testing.T) {
	src := []byte("aaaa bbbb")
	sites := []mutate.Site{{Start: 0, End: 1, Desc: "a"}, {Start: 5, End: 6, Desc: "b"}}
	got := targetSites(src, sites, []cases.Callee{{Symbol: "b", Source: "bbbb"}, {Symbol: "b2", Source: "bbbb"}})
	if len(got) != 1 || got[0].Desc != "b" {
		t.Errorf("span: %v", got)
	}
	got = targetSites(src, sites, []cases.Callee{{Symbol: "x", Source: "zzz"}, {Symbol: "b", Source: "bbbb"}})
	if len(got) != 2 || got[0].Desc != "a" || got[1].Desc != "b" {
		t.Errorf("fallback: %v", got)
	}
}

// Mutant runs get 3× the slower baseline run: slow tests (here ~2 s,
// with a 1 s floor) are measured, not cut off.
func TestTimeoutScalesWithBaseline(t *testing.T) {
	needPython(t)
	root := newRepo(t, map[string]string{
		"calc.py":      pyCalc,
		"test_calc.py": "import time\nfrom calc import add\n\n\ndef test_add():\n    time.sleep(2)\n    assert add(2, 3) == 5\n",
	})
	id := "py:test_calc.py:test_add"
	mod := Module{Rel: "calc.py", Lang: "python", TestFiles: []string{"test_calc.py"}, Tests: []string{id}}
	o := opts(2)
	o.RunTimeout = time.Second
	m, err := RunMatrix(context.Background(), root, []Module{mod}, o)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Baseline[id] || len(m.Mutants) == 0 {
		t.Fatalf("baseline %v, %d mutants", m.Baseline, len(m.Mutants))
	}
	for _, mu := range m.Mutants {
		if mu.KilledBy == nil {
			t.Errorf("mutant %s excluded", mu.Site.Desc)
		}
	}
}

// A broken test file (an import error) sinks only its own tests: the
// module's other test files are still measured.
func TestBrokenTestFileIsolated(t *testing.T) {
	needPython(t)
	root := newRepo(t, map[string]string{
		"calc.py":      pyCalc,
		"test_good.py": "from calc import add\n\n\ndef test_add():\n    assert add(2, 3) == 5\n",
		"test_bad.py":  "import no_such_module_cull\nfrom calc import add\n\n\ndef test_add_bad():\n    assert add(1, 1) == 2\n",
	})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"test_bad.py", "test_good.py"}, map[string]verdict{"test_add": band, "test_add_bad": band})
	res, err := Run(context.Background(), root, settleOpts())
	if err != nil {
		t.Fatal(err)
	}
	if res.Settled["py:test_good.py:test_add"] != "keep" {
		t.Errorf("Settled = %v, reasons %v", res.Settled, res.Reasons)
	}
	if r := res.Reasons["py:test_bad.py:test_add_bad"]; r != "their test file did not run cleanly" {
		t.Errorf("test_add_bad reason = %q", r)
	}
}
