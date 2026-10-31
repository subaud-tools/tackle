package projcli

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/internal/proj"
	"github.com/subaud-tools/tackle/internal/projtui"
)

// pickerCalls replaces the picker's actions with recorders for one test.
type pickerCalls struct {
	calls  []string
	socket string
	pane   string
	right  string
}

func recordPicker(t *testing.T, ensure func(socket, name string)) *pickerCalls {
	t.Helper()
	c := &pickerCalls{}
	origEnsure, origSpawn, origGoto, origWindow := ensureSession, spawnSidebar, gotoSession, showWindow
	t.Cleanup(func() {
		ensureSession, spawnSidebar, gotoSession, showWindow = origEnsure, origSpawn, origGoto, origWindow
	})
	ensureSession = func(socket, name, _, _, _ string, _ []string) (bool, error) {
		c.calls = append(c.calls, "ensure")
		if ensure != nil {
			ensure(socket, name)
		}
		return true, nil
	}
	spawnSidebar = func(_, _, _ string) { c.calls = append(c.calls, "sidebar") }
	gotoSession = func(_, _ string) error {
		c.calls = append(c.calls, "goto")
		return nil
	}
	showWindow = func(socket, pane, right string) (proj.WindowResult, error) {
		c.calls = append(c.calls, "window")
		c.socket, c.pane, c.right = socket, pane, right
		return proj.WindowResult{Window: "opened", ID: "1"}, nil
	}
	return c
}

// tempSession starts a session on a throwaway tmux server and returns its pane
// id and the socket path a process in that pane sees in $TMUX, the form
// `proj window` is given by every other caller.
func tempSession(t *testing.T, sock, name string) (pane, tmuxSocket string) {
	t.Helper()
	if _, err := exec.LookPath("tmux"); err != nil {
		t.Skip("tmux not installed")
	}
	seen := filepath.Join(t.TempDir(), "tmux-env")
	out, err := exec.Command("tmux", "-L", sock, "-f", "/dev/null", "new-session", "-d", "-s", name, "-P", "-F", "#{pane_id}",
		"sh", "-c", `echo "$TMUX" > "$0"; exec sleep 60`, seen).Output()
	if err != nil {
		t.Fatalf("tmux new-session: %v", err)
	}
	t.Cleanup(func() { _ = exec.Command("tmux", "-L", sock, "kill-server").Run() })
	for range 100 {
		if b, err := os.ReadFile(seen); err == nil && len(b) > 0 {
			t.Setenv("TMUX", strings.TrimSpace(string(b)))
			tmuxSocket = proj.SocketFromEnv()
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if tmuxSocket == "" {
		t.Fatal("the pane never reported $TMUX")
	}
	return strings.TrimSpace(string(out)), tmuxSocket
}

func tempSockets(t *testing.T) {
	t.Helper()
	dir, err := os.MkdirTemp("/tmp", "proj-picker-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	t.Setenv("TMUX_TMPDIR", dir)
}

func TestChromeJumpShowsTheWindowForTheSessionsSocketAndPane(t *testing.T) {
	tempSockets(t)
	sock, name := "proj-pickertest-jump", "pickertest-jump/w"
	pane, tmuxSocket := tempSession(t, sock, name)
	c := recordPicker(t, nil)

	if code := actOnResult(projtui.Result{Kind: "jump", Socket: sock, Name: name, Chrome: true}); code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !slices.Equal(c.calls, []string{"window"}) {
		t.Fatalf("calls = %v, want the window and no terminal", c.calls)
	}
	if c.pane != pane {
		t.Fatalf("pane = %q, want %q", c.pane, pane)
	}
	// A window is found again by its page's socket value, so the name the picker
	// holds would open a second window beside the one diazo's path opened.
	if c.socket != tmuxSocket {
		t.Fatalf("socket = %q, want the $TMUX path %q", c.socket, tmuxSocket)
	}
}

func TestChromeNewWorkEnsuresThenShowsTheWindowWithoutTheSidebar(t *testing.T) {
	tempSockets(t)
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	if err := os.MkdirAll(filepath.Join(home, ".config", "proj"), 0o755); err != nil {
		t.Fatal(err)
	}
	roots := "project:" + home + "\n"
	if err := os.WriteFile(filepath.Join(home, ".config", "proj", "roots"), []byte(roots), 0o644); err != nil {
		t.Fatal(err)
	}
	project := filepath.Base(home)

	sock, name := "proj-pickertest-new", "pickertest-new/w"
	var pane, tmuxSocket string
	c := recordPicker(t, func(socket, name string) {
		if pane == "" {
			pane, tmuxSocket = tempSession(t, socket, name)
		}
	})

	res := projtui.Result{Kind: "new", Project: project, Socket: sock, Name: name, Chrome: true}
	if code := actOnResult(res); code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !slices.Equal(c.calls, []string{"ensure", "window"}) {
		t.Fatalf("calls = %v, want ensure then window", c.calls)
	}
	if c.pane != pane || c.socket != tmuxSocket {
		t.Fatalf("window(%q, %q), want (%q, %q)", c.socket, c.pane, tmuxSocket, pane)
	}

	// ^s can still turn the sidebar on for a Chrome session.
	c.calls = nil
	res.Sidebar = true
	if code := actOnResult(res); code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !slices.Equal(c.calls, []string{"ensure", "sidebar", "window"}) {
		t.Fatalf("calls = %v, want the sidebar before the window", c.calls)
	}
}

func TestTerminalResultAttachesAndLeavesTheWindowAlone(t *testing.T) {
	c := recordPicker(t, nil)
	if code := actOnResult(projtui.Result{Kind: "jump", Socket: "proj-x", Name: "x/w"}); code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !slices.Equal(c.calls, []string{"goto"}) {
		t.Fatalf("calls = %v, want only the terminal attach", c.calls)
	}
}

func TestChromeOnAnEndedSessionFailsAndOpensNoTerminal(t *testing.T) {
	tempSockets(t)
	c := recordPicker(t, nil)
	showWindow = proj.Window
	if code := actOnResult(projtui.Result{Kind: "jump", Socket: "proj-pickertest-gone", Name: "gone/w", Chrome: true}); code != 1 {
		t.Fatalf("exit %d, want 1", code)
	}
	if slices.Contains(c.calls, "goto") {
		t.Fatalf("calls = %v, want no terminal", c.calls)
	}
}
