package proj

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/subaud-tools/tackle/internal/proj/pages"
)

// WindowResult is what `proj window --json` prints: Window is "opened" or
// "raised", ID is Chrome's id of the window.
type WindowResult struct {
	Window string `json:"window"`
	ID     string `json:"id"`
}

// windowOption is the tmux session option holding the session's Chrome window.
const windowOption = "@chrome_window"

// pagesSocket and pagesSession are where proj window runs the page server.
const (
	pagesSocket  = "proj-pages"
	pagesSession = "pages"
)

// chromeHasWindowScript prints "yes" when Chrome is running and its window
// argv[1] holds a tab whose URL contains argv[2]. The URL check matters
// because Chrome reuses window ids after a restart; a window that only has the
// old id is not the session's.
const chromeHasWindowScript = `on run argv
if application "Google Chrome" is not running then return ""
tell application "Google Chrome"
repeat with w in windows
if (id of w as text) is (item 1 of argv) then
repeat with t in tabs of w
if (URL of t) contains (item 2 of argv) then return "yes"
end repeat
end if
end repeat
end tell
return ""
end run`

// chromeOpenScript opens a new window on the URL argv[1] and prints its id.
const chromeOpenScript = `on run argv
tell application "Google Chrome"
set w to make new window
set URL of active tab of w to (item 1 of argv)
return id of w as text
end tell
end run`

// chromeSplitScript brings window argv[1] to the front, clicks Tab > Add Tab to
// New Split View, and waits for the new tab, since a disabled or absent item
// raises no error of its own.
const chromeSplitScript = `on run argv
tell application "Google Chrome"
activate
repeat with w in windows
if (id of w as text) is (item 1 of argv) then
set index of w to 1
exit repeat
end if
end repeat
end tell
delay 0.3
tell application "System Events" to tell process "Google Chrome"
click menu item "Add Tab to New Split View" of menu "Tab" of menu bar 1
end tell
tell application "Google Chrome"
repeat 20 times
if (count of tabs of front window) > 1 then return
delay 0.1
end repeat
end tell
error "Tab > Add Tab to New Split View made no new tab"
end run`

// chromeSetRightScript points window argv[1]'s active tab, the split's new
// half, at the URL argv[2].
const chromeSetRightScript = `on run argv
tell application "Google Chrome"
repeat with w in windows
if (id of w as text) is (item 1 of argv) then
set URL of active tab of w to (item 2 of argv)
return
end if
end repeat
end tell
error "Chrome has no window " & (item 1 of argv)
end run`

// chromeRaiseScript brings window argv[1] to the front.
const chromeRaiseScript = `on run argv
tell application "Google Chrome"
repeat with w in windows
if (id of w as text) is (item 1 of argv) then
set minimized of w to false
set index of w to 1
activate
return
end if
end repeat
end tell
error "Chrome has no window " & (item 1 of argv)
end run`

// ensurePages returns the running page server's address and token, starting
// it first in its own tmux session when nothing answers. A seam for tests.
var ensurePages = func() (addr, token string, err error) {
	if addr, token, ok := pages.Address(); ok {
		return addr, token, nil
	}
	exe, err := os.Executable()
	if err != nil {
		return "", "", err
	}
	// A failure here is most often another proj window starting it at the same
	// moment, so it only says what went wrong if the server never answers.
	_, startErr := Run(pagesSocket, "new-session", "-d", "-s", pagesSession, exe, "pages")
	for range 50 {
		if addr, token, ok := pages.Address(); ok {
			return addr, token, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	if startErr != nil {
		return "", "", fmt.Errorf("proj pages did not start: %s", failureText(startErr))
	}
	return "", "", errors.New("proj pages did not answer; see tmux session pages on socket proj-pages")
}

// Window brings pane's session's Chrome window to the front, or opens and
// splits a new one with its right half on right (empty when right is ""). It
// never closes or moves a tab it did not open.
func Window(socket, pane, right string) (WindowResult, error) {
	name := Query(socket, pane, "#{session_name}")
	if name == "" {
		return WindowResult{}, errors.New("the session has ended")
	}
	if id := Query(socket, pane, "#{"+windowOption+"}"); id != "" {
		out, err := osascript(chromeHasWindowScript, id, pages.TermMatch(name))
		if err != nil {
			return WindowResult{}, stepError("find the window", err)
		}
		if strings.TrimSpace(out) == "yes" {
			if _, err := osascript(chromeRaiseScript, id); err != nil {
				return WindowResult{}, stepError("raise the window", err)
			}
			return WindowResult{Window: "raised", ID: id}, nil
		}
	}

	addr, token, err := ensurePages()
	if err != nil {
		return WindowResult{}, stepError("start the page server", err)
	}
	out, err := osascript(chromeOpenScript, pages.TermURL(addr, token, socket, name))
	if err != nil {
		return WindowResult{}, stepError("open the window", err)
	}
	id := strings.TrimSpace(out)
	if id == "" {
		return WindowResult{}, stepError("open the window", errors.New("no window id from Chrome"))
	}
	if _, err := osascript(chromeSplitScript, id); err != nil {
		return WindowResult{}, stepError("split the window", err)
	}
	if right != "" {
		if _, err := osascript(chromeSetRightScript, id, right); err != nil {
			return WindowResult{}, stepError("set the right half", err)
		}
	}
	if _, err := Run(socket, "set-option", "-t", "="+name+":", windowOption, id); err != nil {
		return WindowResult{}, stepError("store the window id", err)
	}
	if _, err := osascript(chromeRaiseScript, id); err != nil {
		return WindowResult{}, stepError("raise the window", err)
	}
	return WindowResult{Window: "opened", ID: id}, nil
}

// stepError names the step that failed and carries Chrome's or macOS's message,
// and says what to do when that message is the missing Accessibility permission.
func stepError(step string, err error) error {
	msg := failureText(err)
	if strings.Contains(msg, "assistive access") || strings.Contains(msg, "-25211") || strings.Contains(msg, "-1719") {
		msg += " (macOS Accessibility permission is missing: allow the app running proj in System Settings > Privacy & Security > Accessibility)"
	}
	return fmt.Errorf("%s: %s", step, msg)
}

// failureText is a command's stderr when it has one, else the error itself.
func failureText(err error) string {
	var ee *exec.ExitError
	if errors.As(err, &ee) && len(ee.Stderr) > 0 {
		return strings.TrimSpace(string(ee.Stderr))
	}
	return err.Error()
}
