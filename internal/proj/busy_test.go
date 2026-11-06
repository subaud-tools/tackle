package proj

import (
	"strings"
	"testing"
	"time"
)

// A pane outside the main one counts as busy unless it runs a shell (a login
// shell's "-zsh" included) or a sidebar tool; a pane in another window counts.
func TestParseBusyPanes(t *testing.T) {
	out := strings.Join([]string{
		scanLine("a/w", "1", "0", "50", "%1", "node", "/a", "0", "0", "", "", ""),
		scanLine("a/w", "1", "81", "10", "%2", "yazi", "/a", "0", "0", "", "", ""),
		scanLine("a/w", "1", "81", "10", "%3", "-zsh", "/a", "0", "0", "", "", ""),
		scanLineWin("0", "a/w", "1", "0", "50", "%4", "sleep", "/a", "0", "0", "", "", ""),
		scanLineWin("0", "a/w", "1", "0", "50", "%5", "bash", "/a", "0", "0", "", "", ""),
		scanLine("b/w", "0", "0", "50", "%6", "zsh", "/b", "0", "0", "", "", ""),
	}, "\n")
	got := parseBusyPanes(out, parseScan(out))
	if len(got["a/w"]) != 1 || got["a/w"][0] != (BusyPane{Pane: "%4", Command: "sleep"}) {
		t.Fatalf("a/w busy = %+v, want only %%4 sleep", got["a/w"])
	}
	if len(got["b/w"]) != 0 {
		t.Fatalf("b/w busy = %+v, want none", got["b/w"])
	}
}

// A sleep in a second window is a busy pane in LiveSessions; a shell there is not.
func TestLiveSessionsBusyPanes(t *testing.T) {
	requireTmux(t)
	sock, name := "proj-busy-test", "proj-busy-test/w"
	defer func() { _, _ = Run(sock, "kill-server") }()

	if created, err := EnsureSession(sock, name, t.TempDir(), "none", "", nil); err != nil || !created {
		t.Fatalf("EnsureSession created=%v err=%v", created, err)
	}
	sleepPane, err := Run(sock, "new-window", "-d", "-P", "-F", "#{pane_id}", "-t", "="+name+":", "sleep", "60")
	if err != nil {
		t.Fatalf("new-window sleep: %v", err)
	}
	if _, err := Run(sock, "new-window", "-d", "-t", "="+name+":"); err != nil {
		t.Fatalf("new-window shell: %v", err)
	}

	var busy []BusyPane
	for deadline := time.Now().Add(5 * time.Second); time.Now().Before(deadline); time.Sleep(100 * time.Millisecond) {
		busy = nil
		for _, s := range LiveSessions() {
			if s.Socket == sock && s.Name == name {
				busy = s.BusyPanes
			}
		}
		if len(busy) == 1 && busy[0].Pane == sleepPane && busy[0].Command == "sleep" {
			return
		}
	}
	t.Fatalf("BusyPanes = %+v, want only %s running sleep", busy, sleepPane)
}
