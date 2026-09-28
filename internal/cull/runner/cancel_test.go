package runner

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// A cancelled ctx starts no process: Err is set and nothing ran.
func TestRunCancelledStartsNothing(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "ran")
	write(t, dir, "go.mod", "module m\n\ngo 1.22\n")
	write(t, dir, "a_test.go", "package m\n\nimport (\n\t\"os\"\n\t\"testing\"\n)\n\nfunc TestA(t *testing.T) { os.WriteFile("+"`"+marker+"`"+", nil, 0o644) }\n")
	write(t, dir, "test_a.py", "open("+`"`+marker+`"`+", 'w').close()\n\ndef test_a():\n    pass\n")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for name, c := range map[string]struct {
		r    Runner
		file string
	}{
		"go":     {golangRunner{}, "a_test.go"},
		"python": {pythonRunner{}, "test_a.py"},
		"node":   {&tsRunner{kind: "node", nodeCmd: "touch " + marker}, "a.test.js"},
	} {
		o := c.r.Run(ctx, dir, []string{c.file}, map[string]string{}, 10*time.Second)
		if o.Err == "" {
			t.Errorf("%s: Err empty", name)
		}
		time.Sleep(100 * time.Millisecond)
		if _, err := os.Stat(marker); err == nil {
			t.Errorf("%s: a process ran", name)
			os.Remove(marker)
		}
	}
}
