package projcli

import (
	"os/exec"
	"testing"

	"github.com/subaud-tools/tackle/internal/proj"
)

// withReapRecord isolates the record under a temp PROJ_HOME holding names, and
// clears $TMUX so the ambient server is the test's choice.
func withReapRecord(t *testing.T, names ...string) {
	t.Helper()
	t.Setenv("PROJ_HOME", t.TempDir())
	t.Setenv("TMUX", "")
	if err := proj.UpdateRecord(func(r *proj.Record) bool {
		for _, n := range names {
			r.Sessions[n] = proj.Entry{Socket: "proj-reaptest-none"}
		}
		return true
	}); err != nil {
		t.Fatal(err)
	}
}

// startReapServer starts a throwaway tmux server hosting session name.
func startReapServer(t *testing.T, sock, name string) {
	t.Helper()
	if _, err := exec.LookPath("tmux"); err != nil {
		t.Skip("tmux not installed")
	}
	t.Cleanup(func() { _ = exec.Command("tmux", "-L", sock, "kill-server").Run() })
	if out, err := exec.Command("tmux", "-L", sock, "-f", "/dev/null", "new-session", "-d", "-s", name, "sleep 60").CombinedOutput(); err != nil {
		t.Fatalf("tmux new-session: %v\n%s", err, out)
	}
	// A second session keeps the server up after name is killed.
	if out, err := exec.Command("tmux", "-L", sock, "new-session", "-d", "-s", "_keep", "sleep 60").CombinedOutput(); err != nil {
		t.Fatalf("tmux new-session: %v\n%s", err, out)
	}
}

func recorded(t *testing.T, name string) bool {
	t.Helper()
	rec, err := proj.LoadRecord()
	if err != nil {
		t.Fatal(err)
	}
	_, ok := rec.Sessions[name]
	return ok
}

func TestReapRequiresSession(t *testing.T) {
	if _, code := captureStdout(t, []string{"reap"}); code != 2 {
		t.Fatalf("proj reap: got exit %d, want 2", code)
	}
}

func TestReapUnknownSession(t *testing.T) {
	withReapRecord(t)
	if _, code := captureStdout(t, []string{"reap", "nope/w", "--socket", "proj-reaptest-none"}); code != 1 {
		t.Fatalf("proj reap nope/w: got exit %d, want 1", code)
	}
}

func TestReapForgetsSavedSession(t *testing.T) {
	withReapRecord(t, "reaptest/saved")
	out, code := captureStdout(t, []string{"reap", "reaptest/saved", "--socket", "proj-reaptest-none"})
	if code != 0 || out != "forgot reaptest/saved\n" {
		t.Fatalf("proj reap: exit %d output %q; want 0 and \"forgot reaptest/saved\"", code, out)
	}
	if recorded(t, "reaptest/saved") {
		t.Fatal("reaptest/saved still in the record")
	}
}

func TestReapKillsLiveSession(t *testing.T) {
	sock, name := "proj-reaptest-live", "reaptest/live"
	startReapServer(t, sock, name)
	withReapRecord(t, name)
	out, code := captureStdout(t, []string{"reap", name, "--socket", sock})
	if code != 0 || out != "reaped "+name+"\n" {
		t.Fatalf("proj reap: exit %d output %q; want 0 and \"reaped %s\"", code, out, name)
	}
	if hasSession(sock, name) {
		t.Fatal("session still running")
	}
	if recorded(t, name) {
		t.Fatal("session still in the record")
	}
}

// The binding runs reap from inside the session it reaps: the socket comes
// from $TMUX and the current session is not refused.
func TestReapSessionItRunsIn(t *testing.T) {
	sock, name := "proj-reaptest-self", "reaptest/self"
	startReapServer(t, sock, name)
	withReapRecord(t, name)
	t.Setenv("TMUX", "/tmp/tmux-test/"+sock+",1,0")
	out, code := captureStdout(t, []string{"reap", name})
	if code != 0 || out != "reaped "+name+"\n" {
		t.Fatalf("proj reap: exit %d output %q; want 0 and \"reaped %s\"", code, out, name)
	}
	if hasSession(sock, name) {
		t.Fatal("session still running")
	}
}
