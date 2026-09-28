package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeGoFixture(t *testing.T) (dir string) {
	t.Helper()
	dir = t.TempDir()
	write := func(rel, content string) {
		p := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module example.com/fx\n\ngo 1.22\n")
	write("fx_test.go", `package fx

import (
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	if 1+1 != 2 {
		t.Fatal("bad")
	}
}

func TestFail(t *testing.T) {
	t.Fatal("boom")
}

func TestParse(t *testing.T) {
	t.Run("quoted_value", func(t *testing.T) {})
	t.Run("empty", func(t *testing.T) {
		t.Fatal("bad")
	})
}

func TestTable(t *testing.T) {
	cases := []struct {
		name string
		in   int
	}{
		{"one", 1},
		{"two", -1},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.in < 0 {
				t.Fatal("bad")
			}
		})
	}
}

func TestUnmappedPasses(t *testing.T) {}

func TestEnvClean(t *testing.T) {
	if os.Getenv("TYPESAFE_API_KEY") != "" {
		t.Fatal("leaked TYPESAFE_API_KEY")
	}
}
`)
	return dir
}

func writeGoHangFixture(t *testing.T) (dir string) {
	t.Helper()
	dir = t.TempDir()
	write := func(rel, content string) {
		p := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module example.com/hang\n\ngo 1.22\n")
	write("hang_test.go", `package hang

import (
	"testing"
	"time"
)

func TestHang(t *testing.T) {
	time.Sleep(30 * time.Second)
}
`)
	return dir
}

func TestGolangRunPassFailSubtestAndTableUnit(t *testing.T) {
	dir := writeGoFixture(t)
	r, err := For(dir, "go", "fx_test.go")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	index := map[string]string{
		"TestAdd":                "go:fx_test.go:TestAdd",
		"TestFail":               "go:fx_test.go:TestFail",
		"TestParse/quoted_value": "go:fx_test.go:TestParse/quoted_value",
		"TestParse/empty":        "go:fx_test.go:TestParse/empty",
		// TestTable's loop-generated subtests are one unit: only the
		// parent maps, not "TestTable/one" or "TestTable/two".
		"TestTable":    "go:fx_test.go:TestTable",
		"TestEnvClean": "go:fx_test.go:TestEnvClean",
		// TestUnmappedPasses deliberately omitted from index.
	}
	out := r.Run(context.Background(), dir, []string{"fx_test.go"}, index, 30*time.Second)

	if out.TimedOut {
		t.Fatalf("TimedOut = true, want false: %+v", out)
	}
	wantPassed := []string{"go:fx_test.go:TestAdd", "go:fx_test.go:TestParse/quoted_value", "go:fx_test.go:TestEnvClean"}
	for _, id := range wantPassed {
		if !out.Passed[id] {
			t.Errorf("Passed[%q] = false, want true; out=%+v", id, out)
		}
	}
	wantFailed := []string{"go:fx_test.go:TestFail", "go:fx_test.go:TestParse/empty", "go:fx_test.go:TestTable"}
	for _, id := range wantFailed {
		if !out.Failed[id] {
			t.Errorf("Failed[%q] = false, want true; out=%+v", id, out)
		}
	}
	if out.Passed["go:fx_test.go:TestTable"] {
		t.Errorf("Passed[TestTable] = true, want false (one subtest failed -> whole unit fails)")
	}
	if out.Unmapped != 1 {
		t.Errorf("Unmapped = %d, want 1 (TestUnmappedPasses not in index): out=%+v", out.Unmapped, out)
	}
}

func writeGoCompileErrorFixture(t *testing.T) (dir string) {
	t.Helper()
	dir = t.TempDir()
	p := filepath.Join(dir, "bad_test.go")
	content := `package bad

import "testing"

func TestBad(t *testing.T) {
	undefinedFunc()
}
`
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/bad\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestGolangRunCompileErrorIsInvalid(t *testing.T) {
	dir := writeGoCompileErrorFixture(t)
	r, err := For(dir, "go", "bad_test.go")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	index := map[string]string{"TestBad": "go:bad_test.go:TestBad"}
	out := r.Run(context.Background(), dir, []string{"bad_test.go"}, index, 30*time.Second)
	if !out.Invalid {
		t.Fatalf("Invalid = false, want true (build failure): out=%+v", out)
	}
	if out.InvalidReason == "" {
		t.Errorf("InvalidReason = %q, want a non-empty reason", out.InvalidReason)
	}
	if len(out.Passed) != 0 || len(out.Failed) != 0 {
		t.Errorf("Passed/Failed should be empty on a build failure: out=%+v", out)
	}
}

func TestGolangRunTimeoutStopsAtFirstPackageOneDeadline(t *testing.T) {
	dir := t.TempDir()
	write := func(rel, content string) {
		p := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module example.com/two\n\ngo 1.22\n")
	write("a/a_test.go", `package a

import (
	"testing"
	"time"
)

func TestHangA(t *testing.T) {
	time.Sleep(30 * time.Second)
}
`)
	marker := filepath.Join(dir, "ran_b.marker")
	write("b/b_test.go", fmt.Sprintf(`package b

import (
	"os"
	"testing"
)

func TestB(t *testing.T) {
	os.WriteFile(%q, []byte("ran"), 0o644)
}
`, marker))

	r, err := For(dir, "go", "a/a_test.go")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	index := map[string]string{
		"TestHangA": "go:a/a_test.go:TestHangA",
		"TestB":     "go:b/b_test.go:TestB",
	}
	start := time.Now()
	out := r.Run(context.Background(), dir, []string{"a/a_test.go", "b/b_test.go"}, index, 300*time.Millisecond)
	elapsed := time.Since(start)
	if !out.TimedOut {
		t.Errorf("TimedOut = false, want true: out=%+v", out)
	}
	if elapsed > 6*time.Second {
		t.Errorf("elapsed = %v, want well under 6s: one shared deadline across both packages, not one timeout per package", elapsed)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Errorf("package b ran after package a's timeout: one Run should stop, not continue to the next package")
	}
}

func TestGolangRunTimeout(t *testing.T) {
	dir := writeGoHangFixture(t)
	r, err := For(dir, "go", "hang_test.go")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	index := map[string]string{"TestHang": "go:hang_test.go:TestHang"}
	start := time.Now()
	out := r.Run(context.Background(), dir, []string{"hang_test.go"}, index, 300*time.Millisecond)
	elapsed := time.Since(start)
	if !out.TimedOut {
		t.Errorf("TimedOut = false, want true: %+v", out)
	}
	if elapsed > 6*time.Second {
		t.Errorf("elapsed = %v, want well under 6s (killGrace-bounded)", elapsed)
	}
}
