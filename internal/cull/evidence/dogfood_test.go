package evidence

import (
	"context"
	"fmt"
	"strings"
	"testing"

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
