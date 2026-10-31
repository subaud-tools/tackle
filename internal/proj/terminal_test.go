package proj

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/creack/pty"
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

// pageShows makes the page the session's one client: list-clients answers 7,
// whatever the format, and 7 is the stored page client.
func pageShows(t *testing.T, socket, name string) {
	t.Helper()
	fakeClients(t, &[]string{"7"})
	if _, err := Run(socket, "set-option", "-t", "="+name+":", "@chrome_client", "7"); err != nil {
		t.Fatal(err)
	}
}

// realClient attaches a real tmux client to name in its own pty.
func realClient(t *testing.T, socket, name string) {
	t.Helper()
	cmd := exec.Command("tmux", socketFlag(socket), socket, "attach", "-t", "="+name) //nolint:gosec // G204: a test's own socket
	cmd.Env = append(slices.DeleteFunc(os.Environ(), func(kv string) bool { return strings.HasPrefix(kv, "TMUX") }), "TERM=xterm-256color")
	tty, err := pty.Start(cmd)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		_ = tty.Close()
	})
}

func clientTTYs(socket string, args ...string) []string {
	out, _ := Run(socket, append([]string{"list-clients", "-F", "#{client_tty}"}, args...)...)
	return splitLines(out)
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

func TestTerminalChromeOnlyWhenTheStoredWindowHoldsThisSession(t *testing.T) {
	socket, name, pane := terminalServer(t)
	term := storedWindow(t, socket, name)
	pageShows(t, socket, name)
	for tabs, want := range map[string]bool{
		term + "https://example.com/\n": true,
		"":                              false,
		strings.ReplaceAll(term, "terminal-test", "someone-else"): false,
	} {
		fakeChrome(t, tabs, nil)
		if got := Terminal(socket, pane); got.Chrome != want {
			t.Errorf("tabs %q: chrome = %v, want %v", tabs, got.Chrome, want)
		}
	}
}

// The window still holds the page, but the page was detached: Ghostty shows it.
func TestTerminalChromeIsFalseWhenAnotherClientShowsTheSession(t *testing.T) {
	socket, name, pane := terminalServer(t)
	term := storedWindow(t, socket, name)
	pageShows(t, socket, name)
	if _, err := Run(socket, "set-option", "-t", "="+name+":", "@chrome_client", "8"); err != nil {
		t.Fatal(err)
	}
	calls := fakeChrome(t, term, nil)
	if got := Terminal(socket, pane); got.Chrome {
		t.Fatalf("Terminal = %+v", got)
	}
	if len(*calls) != 0 {
		t.Fatalf("ran %q", scripts(*calls))
	}
}

func TestTerminalChromeAsksChromeNothingWithoutAStoredWindow(t *testing.T) {
	socket, _, pane := terminalServer(t)
	calls := fakeChrome(t, "", nil)
	if got := Terminal(socket, pane); got.Chrome {
		t.Fatalf("Terminal = %+v", got)
	}
	if len(*calls) != 0 {
		t.Fatalf("ran %q", scripts(*calls))
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

// A client Ghostty doesn't hold, the Chrome page's: --go opens Ghostty, whose
// attach -d detaches that client, and leaves Chrome alone.
func TestTerminalGoOpensGhosttyWhenOnlyChromeShowsIt(t *testing.T) {
	socket, name, pane := terminalServer(t)
	storedWindow(t, socket, name)
	calls := fakeGhostty(t, true)
	ttys := []string{"/dev/ttys042"}
	fakeClients(t, &ttys)
	osascript = func(script string, args ...string) (string, error) {
		*calls = append(*calls, scriptCall{script, args})
		if script == newWindowScript {
			ttys = []string{"/dev/ttys077"} // the new window took the session
		}
		return "", nil
	}
	got, err := TerminalGo(socket, pane)
	if err != nil {
		t.Fatal(err)
	}
	if got != (TerminalState{State: "attached", Session: name}) {
		t.Fatalf("state = %+v", got)
	}
	if len(*calls) != 2 || (*calls)[0].script != focusTTYScript || (*calls)[1].script != newWindowScript {
		t.Fatalf("scripts = %+v", *calls)
	}
	if want := "attach -d -t '=" + name + "'"; !strings.Contains((*calls)[1].args[0], want) {
		t.Fatalf("command = %q, want %q", (*calls)[1].args[0], want)
	}
}

// Focusing a Ghostty tab detaches this session's other clients, and no one
// else's: tmux's detach-client -a would take every other session's on the
// server too.
func TestTerminalGoFocusingGhosttyDetachesOnlyThisSessionsOtherClients(t *testing.T) {
	socket, name, pane := terminalServer(t)
	other := "terminal-test/other"
	if _, err := Run(socket, "new-session", "-d", "-s", other); err != nil {
		t.Fatal(err)
	}
	realClient(t, socket, name)
	realClient(t, socket, name)
	realClient(t, socket, other)
	var mine, theirs []string
	for i := 0; i < 200 && (len(mine) < 2 || len(theirs) < 1); i++ {
		time.Sleep(25 * time.Millisecond)
		mine, theirs = clientTTYs(socket, "-t", "="+name), clientTTYs(socket, "-t", "="+other)
	}
	if len(mine) != 2 || len(theirs) != 1 {
		t.Fatalf("clients: mine %q, theirs %q", mine, theirs)
	}
	fakeGhostty(t, true)
	osascript = func(script string, args ...string) (string, error) {
		if script == focusTTYScript && args[0] == mine[1] {
			return "focused\n", nil
		}
		return "", nil
	}
	if _, err := TerminalGo(socket, pane); err != nil {
		t.Fatal(err)
	}
	var left []string
	for i := 0; i < 200; i++ {
		if left = clientTTYs(socket); len(left) == 2 {
			break
		}
		time.Sleep(25 * time.Millisecond)
	}
	slices.Sort(left)
	want := []string{mine[1], theirs[0]}
	slices.Sort(want)
	if !slices.Equal(left, want) {
		t.Fatalf("clients left = %q, want %q", left, want)
	}
}

// The Chrome window's client may come first: a Ghostty terminal on a later
// client is still the one shown.
func TestTerminalGoFocusesGhosttyOnALaterClient(t *testing.T) {
	socket, _, pane := terminalServer(t)
	calls := fakeChrome(t, "", nil)
	fakeClients(t, &[]string{"/dev/ttys042", "/dev/ttys043"})
	fake := osascript
	osascript = func(script string, args ...string) (string, error) {
		out, err := fake(script, args...)
		if script == focusTTYScript && args[0] == "/dev/ttys043" {
			return "focused\n", nil
		}
		return out, err
	}
	if _, err := TerminalGo(socket, pane); err != nil {
		t.Fatal(err)
	}
	sameScripts(t, *calls, focusTTYScript, focusTTYScript)
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
	want := "'/opt/homebrew/bin/tmux' -S '" + socket + "' attach -d -t '=" + name + "'"
	if len(*scripts) != 1 || !strings.Contains((*scripts)[0].script, "new window") || (*scripts)[0].args[0] != want {
		t.Fatalf("want one new window running %q, got %+v", want, *scripts)
	}
	if !selected(*tmuxCalls, pane) {
		t.Fatalf("pane not selected: %v", *tmuxCalls)
	}
}
