package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/schuettc/tackle/internal/cull/check"
	"github.com/schuettc/tackle/internal/cull/extract"
)

const evCalc = "package calc\n\nfunc Add(a, b int) int { return a + b }\n"

const evCalcTest = `package calc

import "testing"

func TestAdd(t *testing.T) {
	if Add(1, 2) != 3 {
		t.Fatal("bad")
	}
}

func TestNothing(t *testing.T) {
	_ = Add(1, 2)
}

func TestConst(t *testing.T) {
	if 1+1 != 2 {
		t.Fatal("bad")
	}
}
`

// evidenceProject makes a git Go project with a hand-written last.json
// whose three tests are all review-band, and chdirs into it.
func evidenceProject(t *testing.T) string {
	t.Helper()
	for _, bin := range []string{"go", "git"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("%s not on PATH", bin)
		}
	}
	t.Setenv("CULL_HOME", t.TempDir())
	root := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	writeFixture(t, root, "go.mod", "module example.com/fixture\n\ngo 1.22\n")
	writeFixture(t, root, "calc.go", evCalc)
	writeFixture(t, root, "calc_test.go", evCalcTest)
	rel := "calc_test.go"
	res, err := extract.ForFile(rel).Extract(root, []string{rel}, 0)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(evCalcTest))
	fi := check.FileInfo{SHA256: hex.EncodeToString(sum[:])}
	r := check.Report{Root: root, Mode: "suite", Files: map[string]check.FileInfo{}}
	for _, c := range res.Cases {
		r.Tests = append(r.Tests, check.TestResult{TestCase: c, Verdict: "review", Rule: "review_band"})
		fi.Tests = append(fi.Tests, check.FileTest{ID: c.ID, Hash: c.Hash})
	}
	r.Files[rel] = fi
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, ".cull/last.json", string(data))
	t.Chdir(root)
	return root
}

func TestEvidenceExitCodes(t *testing.T) {
	t.Run("0 settled", func(t *testing.T) {
		evidenceProject(t)
		code, out, errw := run(t, "", "evidence", "--workers", "2")
		if code != 0 {
			t.Fatalf("code %d, errw %q, out %q", code, errw, out)
		}
		for _, want := range []string{
			"kept as sole guards: 1", "cut as redundant: 0", "cut as catching nothing: 1",
			"unmeasured: 1", "go:calc_test.go:TestConst: no project code to plant bugs in",
			"modules: 1", "mutants: 1",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("output missing %q:\n%s", want, out)
			}
		}
	})
	t.Run("0 json", func(t *testing.T) {
		evidenceProject(t)
		code, out, errw := run(t, "", "evidence", "--json", "--max-mutants", "5", "--timeout", "30s")
		if code != 0 {
			t.Fatalf("code %d, errw %q", code, errw)
		}
		var got struct {
			Settled    map[string]string `json:"settled"`
			Unmeasured []string          `json:"unmeasured"`
			Modules    int               `json:"modules"`
			Mutants    int               `json:"mutants"`
		}
		if err := json.Unmarshal([]byte(out), &got); err != nil {
			t.Fatalf("%v: %q", err, out)
		}
		if got.Settled["go:calc_test.go:TestAdd"] != "keep" || got.Settled["go:calc_test.go:TestNothing"] != "cut" ||
			len(got.Unmeasured) != 1 || got.Modules != 1 || got.Mutants != 1 {
			t.Errorf("json = %s", out)
		}
		// --json changes stdout only: progress still goes to stderr.
		if !strings.Contains(errw, "cull evidence: baseline calc.go") {
			t.Errorf("stderr lacks progress in --json mode: %q", errw)
		}
	})
	t.Run("2 changed file refused", func(t *testing.T) {
		root := evidenceProject(t)
		writeFixture(t, root, "calc_test.go", evCalcTest+"\n// edited\n")
		code, _, errw := run(t, "", "evidence")
		if code != 2 || !strings.Contains(errw, "changed since cull check") {
			t.Fatalf("code %d, errw %q", code, errw)
		}
		if _, err := os.Stat(filepath.Join(root, ".cull", "evidence.json")); !os.IsNotExist(err) {
			t.Error("evidence.json written")
		}
	})
	t.Run("2 no last.json", func(t *testing.T) {
		root := evidenceProject(t)
		os.Remove(filepath.Join(root, ".cull", "last.json"))
		if code, _, errw := run(t, "", "evidence"); code != 2 || !strings.Contains(errw, "cull check") {
			t.Fatalf("code %d, errw %q", code, errw)
		}
	})
	t.Run("2 bad flag", func(t *testing.T) {
		evidenceProject(t)
		if code, _, _ := run(t, "", "evidence", "--timeout", "0s"); code != 2 {
			t.Fatalf("code %d", code)
		}
	})
}
