package cli

import (
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

// A signal cancels the run: nothing is written to .cull/ and the exit
// code is 128 + the signal number (130 SIGINT, 143 SIGTERM).
func TestEvidenceSignalExit(t *testing.T) {
	for _, tc := range []struct {
		sig  os.Signal
		code int
	}{{os.Interrupt, 130}, {syscall.SIGTERM, 143}} {
		t.Run(tc.sig.String(), func(t *testing.T) {
			root := evidenceProject(t)
			last, _ := os.ReadFile(filepath.Join(root, ".cull", "last.json"))
			notifySignals = func(c chan<- os.Signal, _ ...os.Signal) { c <- tc.sig }
			defer func() { notifySignals = signal.Notify }()
			code, _, errw := run(t, "", "evidence", "--workers", "1")
			if code != tc.code {
				t.Fatalf("code %d, errw %q", code, errw)
			}
			if _, err := os.Stat(filepath.Join(root, ".cull", "evidence.json")); !os.IsNotExist(err) {
				t.Error("evidence.json written")
			}
			if b, _ := os.ReadFile(filepath.Join(root, ".cull", "last.json")); string(b) != string(last) {
				t.Error("last.json changed")
			}
		})
	}
}

// A rerun reuses the cached module results; --no-cache runs afresh.
func TestEvidenceNoCache(t *testing.T) {
	root := evidenceProject(t)
	lastPath := filepath.Join(root, ".cull", "last.json")
	last, _ := os.ReadFile(lastPath)
	for i, c := range []struct {
		args   []string
		reused bool
	}{{nil, false}, {nil, true}, {[]string{"--no-cache"}, false}} {
		if err := os.WriteFile(lastPath, last, 0o600); err != nil {
			t.Fatal(err)
		}
		code, _, errw := run(t, "", append([]string{"evidence", "--workers", "1"}, c.args...)...)
		if code != 0 {
			t.Fatalf("run %d: code %d, errw %q", i, code, errw)
		}
		if got := strings.Contains(errw, "reusing"); got != c.reused {
			t.Errorf("run %d %v: reused %v\n%s", i, c.args, got, errw)
		}
	}
}
