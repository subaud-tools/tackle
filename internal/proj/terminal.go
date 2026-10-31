package proj

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// TerminalState is what `proj terminal --json` prints, and what diazo decodes:
// State is "attached" (a client shows the session), "detached" (the session
// runs with no client) or "ended" (the socket or pane is gone; Session is
// then empty). Window is true when the session's stored Chrome window still
// holds its terminal page, attached or not; Chrome is true when that page's
// client is attached, so Chrome shows the session.
type TerminalState struct {
	State   string `json:"state"`
	Session string `json:"session"`
	Chrome  bool   `json:"chrome,omitempty"`
	Window  bool   `json:"window,omitempty"`
}

// focusTTYScript focuses the Ghostty terminal whose tty is argv[1], which
// brings its window to the front and selects its tab, and prints "focused";
// it prints nothing when no terminal has that tty. Ghostty lists every
// terminal of every window and tab at the application level.
const focusTTYScript = `on run argv
if application "Ghostty" is not running then return ""
tell application "Ghostty"
repeat with t in terminals
if (tty of t) is (item 1 of argv) then
focus t
activate
return "focused"
end if
end repeat
end tell
return ""
end run`

// Terminal reads the state of the session that owns pane on the tmux server
// at socket (a path or a server name).
func Terminal(socket, pane string) TerminalState {
	st, _ := terminalState(socket, pane)
	if st.State != "ended" {
		st.Window = shownWindow(socket, pane, st.Session) != ""
		st.Chrome = st.Window && pageAttached(socket, pane, st.Session)
	}
	return st
}

// terminalState is Terminal plus the attached clients' ttys, in tmux's order.
func terminalState(socket, pane string) (TerminalState, []string) {
	name := Query(socket, pane, "#{session_name}")
	if name == "" {
		return TerminalState{State: "ended"}, nil
	}
	out, _ := Run(socket, "list-clients", "-t", "="+name, "-F", "#{client_tty}")
	var ttys []string
	for _, l := range splitLines(out) {
		if l = strings.TrimSpace(l); l != "" {
			ttys = append(ttys, l)
		}
	}
	if len(ttys) == 0 {
		return TerminalState{State: "detached", Session: name}, nil
	}
	return TerminalState{State: "attached", Session: name}, ttys
}

// focusGhostty focuses the Ghostty terminal of the first of ttys that has one,
// and returns that tty, "" when none did; a client may be a Chrome window's. The
// error is the first focus that failed, when none succeeded.
func focusGhostty(ttys []string) (string, error) {
	var first error
	for _, tty := range ttys {
		out, err := osascript(focusTTYScript, tty)
		if err != nil {
			if first == nil {
				first = fmt.Errorf("focus Ghostty terminal %s: %w", tty, err)
			}
			continue
		}
		if strings.TrimSpace(out) == "focused" {
			return tty, nil
		}
	}
	return "", first
}

// TerminalGo moves pane's session to Ghostty: it focuses a Ghostty terminal
// attached to it and detaches the session's other clients, or opens a new
// Ghostty window attached to it, which detaches them, then selects pane's window
// and pane, and returns the state after. An ended session is an error.
func TerminalGo(socket, pane string) (TerminalState, error) {
	st, ttys := terminalState(socket, pane)
	if st.State == "ended" {
		return st, errors.New("the session has ended")
	}
	focused, err := focusGhostty(ttys)
	if err != nil {
		return st, err
	}
	if focused != "" {
		// One at a time, not detach-client -a: that takes every session's clients
		// on the server. A client already gone is no failure.
		for _, tty := range ttys {
			if tty != focused {
				_, _ = Run(socket, "detach-client", "-t", tty)
			}
		}
	} else {
		s := SavedSession{Name: st.Session, Entry: Entry{Socket: socket}}
		if _, err := osascript(newWindowScript, attachCommand(s)); err != nil {
			return st, fmt.Errorf("open window for %s: %w", s.Name, err)
		}
		waitTaken(socket, pane, ttys) // the state printed after says whether it did
	}
	_, _ = Run(socket, "select-window", "-t", pane)
	_, _ = Run(socket, "select-pane", "-t", pane)
	return Terminal(socket, pane), nil
}

// waitTaken waits for the session's clients to stop being before: the new window
// has attached, and detached them.
func waitTaken(socket, pane string, before []string) {
	for range 10 {
		if _, ttys := terminalState(socket, pane); len(ttys) > 0 && !slices.Equal(ttys, before) {
			return
		}
		time.Sleep(attachPoll)
	}
}
