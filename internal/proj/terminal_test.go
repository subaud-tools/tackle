package proj

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// terminalServer starts a tmux server on a socket given by full path (as
// $TMUX gives it) with one detached session, and returns the socket, the
// session name and its pane id.
func terminalServer(t *testing.T) (socket, name, pane string) {
	t.Helper()
	requireTmux(t)
	// Kept short: macOS caps unix socket paths at 104 bytes.
	tmp, err := os.MkdirTemp("/tmp", "pjt")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmp) })
	socket = filepath.Join(tmp, "proj-terminal-test")
	name = "terminal-test/w"
	if _, err := Run(socket, "new-session", "-d", "-s", name, "-c", tmp); err != nil {
		t.Fatalf("new-session: %v", err)
	}
	t.Cleanup(func() { _, _ = Run(socket, "kill-server") })
	pane, err = Run(socket, "list-panes", "-t", "="+name, "-F", "#{pane_id}")
	if err != nil || pane == "" {
		t.Fatalf("pane id: %q %v", pane, err)
	}
	return socket, name, pane
}

// fakeClients answers tmux's list-clients with the ttys *ttys holds (none
// while it is empty) and records every other tmux call, so a test can attach
// a client without a terminal.
func fakeClients(t *testing.T, ttys *[]string) *[][]string {
	t.Helper()
	var calls [][]string
	orig := execCommand
	t.Cleanup(func() { execCommand = orig })
	execCommand = func(name string, arg ...string) *exec.Cmd {
		for _, a := range arg {
			if a == "list-clients" {
				return orig("printf", "%s", strings.Join(*ttys, "\n"))
			}
		}
		calls = append(calls, arg)
		return orig(name, arg...)
	}
	return &calls
}

func selected(calls [][]string, pane string) bool {
	var win, p bool
	for _, c := range calls {
		s := strings.Join(c, " ")
		win = win || strings.Contains(s, "select-window -t "+pane)
		p = p || strings.Contains(s, "select-pane -t "+pane)
	}
	return win && p
}

func TestTerminalStateJSON(t *testing.T) {
	b, err := json.Marshal(TerminalState{State: "attached", Session: "a/1"})
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"state":"attached","session":"a/1"}` {
		t.Fatalf("json = %s", b)
	}
}

func TestTerminalEnded(t *testing.T) {
	requireTmux(t)
	got := Terminal(filepath.Join(t.TempDir(), "gone"), "%1")
	if got != (TerminalState{State: "ended"}) {
		t.Fatalf("Terminal = %+v, want ended", got)
	}
}

func TestTerminalDetached(t *testing.T) {
	socket, name, pane := terminalServer(t)
	if got := Terminal(socket, pane); got != (TerminalState{State: "detached", Session: name}) {
		t.Fatalf("Terminal = %+v, want detached %s", got, name)
	}
}

func TestTerminalAttached(t *testing.T) {
	socket, name, pane := terminalServer(t)
	fakeClients(t, &[]string{"/dev/ttys042"})
	if got := Terminal(socket, pane); got != (TerminalState{State: "attached", Session: name}) {
		t.Fatalf("Terminal = %+v, want attached %s", got, name)
	}
}

func TestTerminalGoEnded(t *testing.T) {
	requireTmux(t)
	calls := fakeGhostty(t, true)
	if _, err := TerminalGo(filepath.Join(t.TempDir(), "gone"), "%1"); err == nil {
		t.Fatal("TerminalGo on an ended session: want an error")
	}
	if len(*calls) != 0 {
		t.Fatalf("drove Ghostty for an ended session: %+v", *calls)
	}
}

func TestTerminalGoAttachedFocuses(t *testing.T) {
	socket, name, pane := terminalServer(t)
	scripts := fakeGhostty(t, true)
	tmuxCalls := fakeClients(t, &[]string{"/dev/ttys042", "/dev/ttys043"})
	osascript = func(script string, args ...string) (string, error) {
		*scripts = append(*scripts, scriptCall{script, args})
		return "focused\n", nil
	}
	got, err := TerminalGo(socket, pane)
	if err != nil {
		t.Fatal(err)
	}
	if got != (TerminalState{State: "attached", Session: name}) {
		t.Fatalf("state = %+v", got)
	}
	if len(*scripts) != 1 || !strings.Contains((*scripts)[0].script, "focus") ||
		strings.Join((*scripts)[0].args, " ") != "/dev/ttys042" {
		t.Fatalf("want one focus on the first client tty, got %+v", *scripts)
	}
	if !selected(*tmuxCalls, pane) {
		t.Fatalf("pane not selected: %v", *tmuxCalls)
	}
}

func TestTerminalGoAttachedNoGhosttyTab(t *testing.T) {
	socket, _, pane := terminalServer(t)
	fakeGhostty(t, true)
	fakeClients(t, &[]string{"/dev/ttys042"})
	osascript = func(script string, args ...string) (string, error) { return "\n", nil }
	if _, err := TerminalGo(socket, pane); err == nil {
		t.Fatal("no Ghostty terminal has the client's tty: want an error")
	}
}

func TestTerminalGoDetachedOpensWindow(t *testing.T) {
	socket, name, pane := terminalServer(t)
	scripts := fakeGhostty(t, true)
	var ttys []string
	tmuxCalls := fakeClients(t, &ttys)
	hasClient = func(socket, name string) bool { return len(ttys) > 0 }
	osascript = func(script string, args ...string) (string, error) {
		*scripts = append(*scripts, scriptCall{script, args})
		ttys = []string{"/dev/ttys077"} // the new window attaches
		return "win-1\n", nil
	}
	got, err := TerminalGo(socket, pane)
	if err != nil {
		t.Fatal(err)
	}
	if got != (TerminalState{State: "attached", Session: name}) {
		t.Fatalf("state = %+v", got)
	}
	want := "'/opt/homebrew/bin/tmux' -S '" + socket + "' attach -t '=" + name + "'"
	if len(*scripts) != 1 || !strings.Contains((*scripts)[0].script, "new window") || (*scripts)[0].args[0] != want {
		t.Fatalf("want one new window running %q, got %+v", want, *scripts)
	}
	if !selected(*tmuxCalls, pane) {
		t.Fatalf("pane not selected: %v", *tmuxCalls)
	}
}
