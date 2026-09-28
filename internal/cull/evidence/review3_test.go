package evidence

import (
	"context"
	"testing"
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
