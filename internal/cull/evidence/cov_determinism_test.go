package evidence

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// settleBoth runs Run on root with coverage selection and without
// (restoring last.json in between) and fails t unless the settled
// verdicts and rules are identical.
func settleBoth(t *testing.T, root string) {
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
		if !noCov && strings.Contains(log.String(), "no per-test coverage") {
			t.Fatalf("coverage not used:\n%s", log.String())
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

// Settled verdicts with coverage selection equal those of the full run.
func TestCoverageDeterminismGo(t *testing.T) {
	needTools(t)
	settleBoth(t, goSettleRepo(t))
}

func TestCoverageDeterminismPython(t *testing.T) {
	needCoverage(t)
	root := covPyRepo(t, filepath.Join(t.TempDir(), "runs"), map[string]string{})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"test_calc.py"}, map[string]verdict{"test_f": band, "test_g1": band, "test_g2": band})
	settleBoth(t, root)
}
