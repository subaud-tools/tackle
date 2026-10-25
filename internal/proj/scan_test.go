package proj

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// scanLine builds one list-panes -a line in scanFormat()'s column order.
func scanLine(f ...string) string { return strings.Join(f, "\x1f") }

func TestParseScanGroupsPanesBySession(t *testing.T) {
	out := strings.Join([]string{
		scanLine("a/w", "1", "0", "10", "%1", "zsh", "/a/short", "0", "100", "conv-a", "pi", "lab-a"),
		scanLine("a/w", "1", "0", "40", "%2", "node", "/a/tall", "1", "200", "conv-a", "pi", "lab-a"),
		scanLine("a/w", "1", "81", "50", "%3", "yazi", "/a/right", "0", "300", "conv-a", "pi", "lab-a"),
		scanLine("b/w", "0", "0", "50", "%4", "2.1.248", "/b", "0", "400", "conv-b", "", "lab-b"),
	}, "\n")
	got := parseScan(out)
	if len(got) != 2 || got[0].Name != "a/w" || got[1].Name != "b/w" {
		t.Fatalf("sessions = %+v, want a/w then b/w", got)
	}
	want := ScannedSession{Name: "a/w", Attached: true, Pane: "%2", Command: "node", Dir: "/a/tall",
		Bell: true, Activity: 200, Conversation: "conv-a", Agent: "pi", Label: "lab-a"}
	if got[0] != want {
		t.Fatalf("a/w = %+v\nwant  %+v", got[0], want)
	}
	if got[1].Pane != "%4" || got[1].Command != "2.1.248" || got[1].Activity != 400 {
		t.Fatalf("b/w = %+v", got[1])
	}
}

func TestParseScanEmptyFields(t *testing.T) {
	got := parseScan(scanLine("a/w", "0", "0", "50", "%1", "", "/a", "1", "123", "", "", ""))
	want := ScannedSession{Name: "a/w", Pane: "%1", Dir: "/a", Bell: true, Activity: 123}
	if len(got) != 1 || got[0] != want {
		t.Fatalf("parseScan = %+v\nwant         %+v", got, want)
	}
}

func TestParseScanAttached(t *testing.T) {
	got := parseScan(strings.Join([]string{
		scanLine("a/w", "2", "0", "50", "%1", "zsh", "/a", "0", "0", "", "", ""),
		scanLine("b/w", "0", "0", "50", "%2", "zsh", "/b", "0", "0", "", "", ""),
	}, "\n"))
	if len(got) != 2 || !got[0].Attached || got[1].Attached {
		t.Fatalf("attached = %+v, want a/w true, b/w false", got)
	}
}

// resolved returns dir with symlinks resolved, as tmux reports pane paths.
func resolved(t *testing.T, dir string) string {
	t.Helper()
	if r, err := filepath.EvalSymlinks(dir); err == nil {
		return r
	}
	return dir
}

func TestScanServerFolderIsMainPane(t *testing.T) {
	sock, name := "proj-scantest-folder", "scantest-folder/w"
	dirA := testSession(t, sock, name)
	dirB := resolved(t, t.TempDir())
	// Split right into B; the new pane is left active.
	if _, err := Run(sock, "split-window", "-h", "-t", name, "-c", dirB); err != nil {
		t.Fatal(err)
	}
	got, err := ScanServer(sock)
	if err != nil {
		t.Fatalf("ScanServer: %v", err)
	}
	if len(got) != 1 || got[0].Dir != dirA {
		t.Fatalf("ScanServer = %+v, want one session in %s", got, dirA)
	}
}

func TestLiveSessionsOneCallPerServer(t *testing.T) {
	requireTmux(t)
	// A private socket dir, so Servers() sees only this test's server. Kept
	// short: macOS caps unix socket paths at 104 bytes.
	tmp, err := os.MkdirTemp("/tmp", "pjs")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmp) })
	t.Setenv("TMUX_TMPDIR", tmp)
	sock := "proj-scantest-calls"
	t.Cleanup(func() { _, _ = Run(sock, "kill-server") })
	for _, n := range []string{"scantest-calls/a", "scantest-calls/b", "scantest-calls/c"} {
		if _, err := EnsureSession(sock, n, t.TempDir(), "none", "", nil); err != nil {
			t.Fatalf("EnsureSession: %v", err)
		}
	}

	calls := 0
	orig := execCommand
	t.Cleanup(func() { execCommand = orig })
	execCommand = func(name string, arg ...string) *exec.Cmd {
		if name == "tmux" {
			calls++
		}
		return orig(name, arg...)
	}
	got := LiveSessions()
	if len(got) != 3 {
		t.Fatalf("LiveSessions = %+v, want 3", got)
	}
	if calls != 2 {
		t.Fatalf("tmux calls = %d, want 2 (the Servers() probe and the scan)", calls)
	}
}
