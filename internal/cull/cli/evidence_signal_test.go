package cli

import (
	"os"
	"os/signal"
	"path/filepath"
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
