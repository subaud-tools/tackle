package proj

import (
	"fmt"
	"os/exec"
	"strings"
)

// osascript runs an AppleScript, one -e per script line, with args passed to
// its `on run argv` handler. A seam so tests never drive the real Ghostty.
var osascript = func(script string, args ...string) (string, error) {
	var argv []string
	for _, line := range strings.Split(script, "\n") {
		argv = append(argv, "-e", line)
	}
	out, err := exec.Command("osascript", append(argv, args...)...).Output()
	return string(out), err
}

// GhosttyTerminal identifies a native surface without depending on its title.
type GhosttyTerminal struct{ WindowID, TabID, TTY string }

// Ghostty 1.4 exposes windows -> tabs -> terminals, with tty on each terminal.
const readTerminalsScript = `set out to ""
if application "Ghostty" is running then
tell application "Ghostty"
repeat with w in windows
repeat with t in tabs of w
repeat with term in terminals of t
set out to out & (id of w) & (character id 9) & (id of t) & (character id 9) & (tty of term) & linefeed
end repeat
end repeat
end repeat
end tell
end if
return out`

// ObserveGhostty never launches Ghostty. A malformed observation is an error,
// not a successful empty desktop.
func ObserveGhostty() ([]GhosttyTerminal, error) {
	out, err := osascript(readTerminalsScript)
	if err != nil {
		return nil, err
	}
	var terminals []GhosttyTerminal
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 3 || fields[0] == "" || fields[1] == "" || fields[2] == "" {
			return nil, fmt.Errorf("proj: malformed Ghostty terminal observation")
		}
		terminals = append(terminals, GhosttyTerminal{WindowID: fields[0], TabID: fields[1], TTY: fields[2]})
	}
	return terminals, nil
}

// mergeLayout folds the current Ghostty arrangement (cur) into the saved one
// (old). Every session with a tab now sits where its tab is. A saved session
// with no tab now (not restored yet, detached, or in another terminal) keeps
// its saved window: it is appended to the current window that holds the most
// of that saved window's sessions, or stays a window of its own when none of
// them have a tab. Like the record, the layout never shrinks on what is
// merely absent, so attaching one session before restoring the rest cannot
// overwrite the saved layout with a one-tab window.
func mergeLayout(old, cur [][]string) [][]string {
	tabbed := map[string]int{} // session -> index of its current window
	for i, w := range cur {
		for _, s := range w {
			tabbed[s] = i
		}
	}
	// Walk the saved windows in order. Each one either maps onto the current
	// window holding most of its sessions (lowest index on a tie), handing it
	// its untabbed sessions, or survives on its own. Current windows no saved
	// window maps onto follow, in Ghostty's order.
	extra := make([][]string, len(cur))
	placed := make([]bool, len(cur))
	var order [][]string // nil entry = a current window, resolved below
	var curAt []int      // for each nil entry in order, its cur index
	for _, w := range old {
		votes := map[int]int{}
		var rest []string
		for _, s := range w {
			if i, ok := tabbed[s]; ok {
				votes[i]++
			} else {
				rest = append(rest, s)
			}
		}
		best := -1
		for i, v := range votes {
			if best < 0 || v > votes[best] || (v == votes[best] && i < best) {
				best = i
			}
		}
		if best < 0 {
			if len(rest) > 0 { // an empty saved window would read as a nil marker
				order = append(order, rest)
			}
			continue
		}
		extra[best] = append(extra[best], rest...)
		if !placed[best] {
			placed[best] = true
			order = append(order, nil)
			curAt = append(curAt, best)
		}
	}
	for i := range cur {
		if !placed[i] {
			order = append(order, nil)
			curAt = append(curAt, i)
		}
	}
	var out [][]string
	seen := map[string]bool{}
	k := 0
	for _, w := range order {
		if w == nil {
			i := curAt[k]
			k++
			w = append(append([]string{}, cur[i]...), extra[i]...)
		}
		var keep []string
		for _, s := range w {
			if !seen[s] {
				seen[s] = true
				keep = append(keep, s)
			}
		}
		if len(keep) > 0 {
			out = append(out, keep)
		}
	}
	return out
}

// SaveLayout captures the current Ghostty arrangement of live sessions and
// merges it into the record's layout (see mergeLayout). It reports the windows
// and tabs saved. A Ghostty failure leaves the record untouched.
func SaveLayout(_ map[string]bool) (windows, tabs int, err error) {
	err = withThemeSyncLock(func() error {
		o, e := observeThemes()
		if e != nil {
			return e
		}
		r, e := o.save()
		if e != nil {
			return e
		}
		windows = len(r.Layout.Windows)
		for _, w := range r.Layout.Windows {
			tabs += len(w)
		}
		return nil
	})
	return windows, tabs, err
}
