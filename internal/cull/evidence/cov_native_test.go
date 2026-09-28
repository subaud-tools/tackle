package evidence

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// nativeTool builds a do-nothing native binary named culltool (Go,
// CGO_ENABLED=0) in dir and returns its path.
func nativeTool(t *testing.T, dir string) string {
	t.Helper()
	src := t.TempDir()
	write(t, src, "go.mod", "module culltool\n\ngo 1.22\n")
	write(t, src, "main.go", "package main\n\nfunc main() {}\n")
	out := filepath.Join(dir, "culltool")
	cmd := exec.Command("go", "build", "-o", out, ".")
	cmd.Dir = src
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, b)
	}
	return out
}

// configureSpawn is a conftest whose pytest_configure runs argv (Python
// expressions) with subprocess.run(..., check=True), plus extra keyword
// arguments.
func configureSpawn(argv, kw string) string {
	return "import os, subprocess\n\n\ndef pytest_configure(config):\n    subprocess.run([" + argv + "], check=True" + kw + ")\n"
}

// A native binary outside the project run by pytest_configure (a bundled
// solver probe) doesn't disable coverage selection; it's logged once.
func TestCoverageNativeSetupSpawnIgnored(t *testing.T) {
	needCoverage(t)
	tool := nativeTool(t, t.TempDir())
	root := pySettleRepo(t, soundCalcFG, soundFGTests)
	write(t, root, "conftest.py", configureSpawn(pyQuote(tool), ""))
	log := settleExpect(t, root, "")
	if n := strings.Count(log, "culltool"); n != 1 || !strings.Contains(log, "ignor") {
		t.Errorf("ignored spawn logged %d times:\n%s", n, log)
	}
}

// The same found through the PATH passed to the call.
func TestCoverageNativeSetupSpawnViaPathIgnored(t *testing.T) {
	needCoverage(t)
	dir := t.TempDir()
	nativeTool(t, dir)
	root := pySettleRepo(t, soundCalcFG, soundFGTests)
	write(t, root, "conftest.py", configureSpawn("'culltool'", ", env={**os.environ, 'PATH': "+pyQuote(dir)+"}"))
	log := settleExpect(t, root, "")
	if !strings.Contains(log, "culltool") {
		t.Errorf("ignored spawn not logged:\n%s", log)
	}
}

// A shebang script outside the project isn't native: it falls back.
func TestCoverageScriptSetupSpawnFallsBack(t *testing.T) {
	needCoverage(t)
	script := filepath.Join(t.TempDir(), "culltool")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	root := pySettleRepo(t, soundCalcFG, soundFGTests)
	write(t, root, "conftest.py", configureSpawn(pyQuote(script), ""))
	settleExpect(t, root, "spawned")
}

// A native binary inside the project root falls back.
func TestCoverageProjectNativeSetupSpawnFallsBack(t *testing.T) {
	needCoverage(t)
	root := pySettleRepo(t, soundCalcFG, soundFGTests)
	nativeTool(t, filepath.Join(root, "bin"))
	write(t, root, "conftest.py", configureSpawn("os.path.join(os.path.dirname(os.path.abspath(__file__)), 'bin', 'culltool')", ""))
	settleExpect(t, root, "spawned")
}

// A shell (native, outside the project, but it runs anything) falls back.
func TestCoverageShellSetupSpawnFallsBack(t *testing.T) {
	needCoverage(t)
	root := pySettleRepo(t, soundCalcFG, soundFGTests)
	write(t, root, "conftest.py", configureSpawn("'/bin/sh', '-c', 'true'", ""))
	settleExpect(t, root, "spawned")
}
