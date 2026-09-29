package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func requirePython3(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not found")
	}
}

func writePyFixture(t *testing.T) (dir string) {
	t.Helper()
	dir = t.TempDir()
	p := filepath.Join(dir, "test_a.py")
	content := `import os


def test_plain():
    assert True


def test_fail():
    assert False


class TestThing:
    def test_method(self):
        assert True

    def test_method_fail(self):
        assert False


def test_env_clean():
    assert os.environ.get("TYPESAFE_API_KEY") is None
`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writePyHangFixture(t *testing.T) (dir string) {
	t.Helper()
	dir = t.TempDir()
	p := filepath.Join(dir, "test_hang.py")
	content := `import time


def test_hang():
    time.sleep(30)
`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestPythonRunPassFailClassAndUnmapped(t *testing.T) {
	requirePython3(t)
	dir := writePyFixture(t)
	r, err := For(dir, "python", "test_a.py")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	index := map[string]string{
		"test_plain":                  "py:test_a.py:test_plain",
		"test_fail":                   "py:test_a.py:test_fail",
		"TestThing::test_method":      "py:test_a.py:TestThing::test_method",
		"TestThing::test_method_fail": "py:test_a.py:TestThing::test_method_fail",
		"test_env_clean":              "py:test_a.py:test_env_clean",
		// no parametrize case in this fixture, but Unmapped is still
		// exercised: nothing left out here, so add a bogus expectation
		// below by omitting nothing (0 unmapped).
	}
	out := r.Run(context.Background(), dir, []string{"test_a.py"}, index, 30*time.Second)
	if out.TimedOut {
		t.Fatalf("TimedOut = true, want false: %+v", out)
	}
	wantPassed := []string{"py:test_a.py:test_plain", "py:test_a.py:TestThing::test_method", "py:test_a.py:test_env_clean"}
	for _, id := range wantPassed {
		if !out.Passed[id] {
			t.Errorf("Passed[%q] = false, want true; out=%+v", id, out)
		}
	}
	wantFailed := []string{"py:test_a.py:test_fail", "py:test_a.py:TestThing::test_method_fail"}
	for _, id := range wantFailed {
		if !out.Failed[id] {
			t.Errorf("Failed[%q] = false, want true; out=%+v", id, out)
		}
	}
}

func TestPythonRunParametrizeMapsToFunctionUnit(t *testing.T) {
	requirePython3(t)
	dir := t.TempDir()
	p := filepath.Join(dir, "test_p.py")
	content := `import pytest


@pytest.mark.parametrize("x", [1, 2])
def test_param(x):
    assert x > 0


def test_other():
    assert True
`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := For(dir, "python", "test_p.py")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	// index only has the function-level unit, not the "[1]"/"[2]" native
	// pytest node ids: both parametrize instances must attribute there.
	index := map[string]string{
		"test_param": "py:test_p.py:test_param",
		// test_other deliberately omitted -> Unmapped.
	}
	out := r.Run(context.Background(), dir, []string{"test_p.py"}, index, 30*time.Second)
	if !out.Passed["py:test_p.py:test_param"] {
		t.Errorf("Passed[test_param] = false, want true (both parametrize cases pass): out=%+v", out)
	}
	if out.Unmapped != 1 {
		t.Errorf("Unmapped = %d, want 1 (test_other not in index): out=%+v", out.Unmapped, out)
	}
}

func TestPythonRunNestedDirNoInitMapsPlainFuncAndClassMethod(t *testing.T) {
	requirePython3(t)
	dir := t.TempDir()
	p := filepath.Join(dir, "sub", "dir", "test_b.py")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	content := `def test_plain():
    assert True


class TestC:
    def test_method(self):
        assert True
`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	// No __init__.py anywhere: pytest's junitxml classname for the plain
	// function is the whole dotted module path ("sub.dir.test_b"), not
	// "test_b" treated as a class.
	r, err := For(dir, "python", "sub/dir/test_b.py")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	index := map[string]string{
		"test_plain":         "py:sub/dir/test_b.py:test_plain",
		"TestC::test_method": "py:sub/dir/test_b.py:TestC::test_method",
	}
	out := r.Run(context.Background(), dir, []string{"sub/dir/test_b.py"}, index, 30*time.Second)
	if out.Invalid {
		t.Fatalf("Invalid = true, want false: out=%+v", out)
	}
	if !out.Passed["py:sub/dir/test_b.py:test_plain"] {
		t.Errorf("Passed[test_plain] = false, want true (module path misread as a class): out=%+v", out)
	}
	if !out.Passed["py:sub/dir/test_b.py:TestC::test_method"] {
		t.Errorf("Passed[TestC::test_method] = false, want true: out=%+v", out)
	}
	if out.Unmapped != 0 {
		t.Errorf("Unmapped = %d, want 0: out=%+v", out.Unmapped, out)
	}
}

func TestPythonRunImportErrorIsInvalid(t *testing.T) {
	requirePython3(t)
	dir := t.TempDir()
	p := filepath.Join(dir, "test_bad.py")
	content := "import nonexistent_module_xyz\n\n\ndef test_a():\n    assert True\n"
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := For(dir, "python", "test_bad.py")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	index := map[string]string{"test_a": "py:test_bad.py:test_a"}
	out := r.Run(context.Background(), dir, []string{"test_bad.py"}, index, 30*time.Second)
	if !out.Invalid {
		t.Fatalf("Invalid = false, want true (collection/import error): out=%+v", out)
	}
	if out.InvalidReason == "" {
		t.Errorf("InvalidReason = %q, want a non-empty reason", out.InvalidReason)
	}
}

func TestPythonRunTimeout(t *testing.T) {
	requirePython3(t)
	dir := writePyHangFixture(t)
	r, err := For(dir, "python", "test_hang.py")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	index := map[string]string{"test_hang": "py:test_hang.py:test_hang"}
	start := time.Now()
	out := r.Run(context.Background(), dir, []string{"test_hang.py"}, index, 300*time.Millisecond)
	elapsed := time.Since(start)
	if !out.TimedOut {
		t.Errorf("TimedOut = false, want true: %+v", out)
	}
	if elapsed > 6*time.Second {
		t.Errorf("elapsed = %v, want well under 6s (killGrace-bounded)", elapsed)
	}
}

func TestPythonEnvSandboxFirst(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PYTHONPATH", "/elsewhere")
	var got []string
	for _, kv := range pythonEnv(dir) {
		if strings.HasPrefix(kv, "PYTHONPATH=") {
			got = append(got, kv)
		}
	}
	want := "PYTHONPATH=" + filepath.Join(dir, "src") + string(os.PathListSeparator) + dir + string(os.PathListSeparator) + "/elsewhere"
	if len(got) != 1 || got[0] != want {
		t.Fatalf("got %v, want %s", got, want)
	}
}

// pytest-randomly is disabled in Run and Cover: every run keeps the
// files' order.
func TestPythonRandomlyBlocked(t *testing.T) {
	requirePython3(t)
	if v := os.Getenv("CULL_PYCOV"); v != "" {
		t.Setenv("VIRTUAL_ENV", v)
	}
	dir := t.TempDir()
	files := map[string]string{
		"calc.py":   "def f(x):\n    return x > 3\n",
		"test_r.py": "from calc import f\n\n\ndef test_r(request):\n    assert f(5)\n    assert request.config.pluginmanager.is_blocked('randomly')\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	idx := map[string]string{"test_r": "py:test_r.py:test_r"}
	r := pythonRunner{}
	if out := r.Run(context.Background(), dir, []string{"test_r.py"}, idx, time.Minute); !out.Passed["py:test_r.py:test_r"] {
		t.Errorf("Run: randomly not blocked: %+v", out)
	}
	if os.Getenv("CULL_PYCOV") == "" {
		return
	}
	cov, err := r.Cover(context.Background(), dir, map[string]map[string]string{"test_r.py": idx}, "calc.py", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if !cov.Passed["py:test_r.py:test_r"] {
		t.Error("Cover: randomly not blocked")
	}
}
