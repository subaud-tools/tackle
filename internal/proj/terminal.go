package proj

import (
	"errors"
	"fmt"
	"strings"
)

// TerminalState is what `proj terminal --json` prints, and what diazo decodes:
// State is "attached" (a client shows the session), "detached" (the session
// runs with no client) or "ended" (the socket or pane is gone; Session is
// then empty). Chrome is true when the session's stored Chrome window still
// holds its terminal page.
type TerminalState struct {
	State   string `json:"state"`
	Session string `json:"session"`
	Chrome  bool   `json:"chrome,omitempty"`
}

// focusTTYScript focuses the Ghostty terminal whose tty is argv[1], which
// brings its window to the front and selects its tab, and prints "focused";
// it prints nothing when no terminal has that tty. Ghostty lists every
// terminal of every window and tab at the application level.
const focusTTYScript = `on run argv
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
	st.Chrome = st.State != "ended" && windowShows(socket, pane, st.Session)
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

// TerminalGo takes Court to pane's session: it focuses the Ghostty terminal
// attached to it (the first client's), or opens a new Ghostty window attached
// to it when none is, then selects pane's window and pane, and returns the
// state after. An ended session is an error.
func TerminalGo(socket, pane string) (TerminalState, error) {
	st, ttys := terminalState(socket, pane)
	switch st.State {
	case "ended":
		return st, errors.New("the session has ended")
	case "attached":
		out, err := osascript(focusTTYScript, ttys[0])
		if err != nil {
			return st, fmt.Errorf("focus Ghostty terminal %s: %w", ttys[0], err)
		}
		if strings.TrimSpace(out) != "focused" {
			return st, fmt.Errorf("no Ghostty terminal shows %s (tty %s)", st.Session, ttys[0])
		}
	case "detached":
		s := SavedSession{Name: st.Session, Entry: Entry{Socket: socket}}
		if _, err := osascript(newWindowScript, attachCommand(s)); err != nil {
			return st, fmt.Errorf("open window for %s: %w", s.Name, err)
		}
		_ = waitAttached(s) // the state printed after says whether it did
	}
	_, _ = Run(socket, "select-window", "-t", pane)
	_, _ = Run(socket, "select-pane", "-t", pane)
	return Terminal(socket, pane), nil
}
