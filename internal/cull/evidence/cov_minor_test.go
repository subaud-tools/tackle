package evidence

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"testing"
)

// The cache key hashes the files the sandbox ran, not the working tree:
// a helper edited after the sandboxes were copied gives the first run's
// results the old bytes' key, so a rerun on the new bytes misses.
func TestCacheKeyHashesSandbox(t *testing.T) {
	needTools(t)
	root := newRepo(t, map[string]string{
		"go.mod":      "module fixture\n\ngo 1.22\n",
		"a/a.go":      "package a\n\nfunc Add(x, y int) int { return x + y + k }\n",
		"a/k.go":      "package a\n\nconst k = 0\n",
		"a/a_test.go": "package a\n\nimport \"testing\"\n\nfunc TestAdd(t *testing.T) {\n\tif Add(2, 3) != 5 {\n\t\tt.Fatal(\"add\")\n\t}\n}\n",
	})
	mods := []Module{{Rel: "a/a.go", Lang: "go", TestFiles: []string{"a/a_test.go"}, Tests: []string{"go:a/a_test.go:TestAdd"}}}
	o := opts(1)
	o.NoCache = false
	run := func() string {
		var buf bytes.Buffer
		o.Stderr = &buf
		if _, err := RunMatrix(context.Background(), root, mods, o); err != nil {
			t.Fatal(err)
		}
		return buf.String()
	}
	preparedHook = func() { write(t, root, "a/k.go", "package a\n\nconst k = 1\n") }
	run()
	preparedHook = nil
	t.Cleanup(func() { preparedHook = nil })
	if s := run(); strings.Contains(s, "unchanged since a previous run") {
		t.Errorf("rerun on the edited helper reused the old results:\n%s", s)
	}
}

// The cache stamp holds the toolchains' own version output, so an
// in-place upgrade at the same path misses.
func TestCacheStampToolVersions(t *testing.T) {
	needTools(t)
	s := strings.Join(buildStamp(t.TempDir()), "\n")
	for _, want := range []string{"go-version=go version go1.", "node-version=", "python-version="} {
		if !strings.Contains(s, want) {
			t.Errorf("stamp lacks %q:\n%s", want, s)
		}
	}
	if _, err := exec.LookPath("node"); err == nil && !strings.Contains(s, "node-version=v") {
		t.Errorf("stamp lacks node's version:\n%s", s)
	}
	if _, err := exec.LookPath("python3"); err == nil && !strings.Contains(s, "python-version=Python 3.") {
		t.Errorf("stamp lacks Python's version:\n%s", s)
	}
}
