package proj

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"slices"
	"strconv"
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

// windowOption and rightOption are the tmux session options holding the
// session's Chrome window and the id of the tab that is its right half.
const (
	windowOption = "@chrome_window"
	rightOption  = "@chrome_right"
)

// pagesSocket and pagesSession are where proj window runs the page server.
const (
	pagesSocket  = "proj-pages"
	pagesSession = "pages"
)

// chromeHasWindowScript prints the URL of each tab of window argv[1], one a
// line, and nothing when Chrome is not running or has no such window. Chrome
// reuses window ids after a restart, so an id alone does not make a window the
// session's; the caller looks for the session's terminal page among the URLs.
const chromeHasWindowScript = `on run argv
if application "Google Chrome" is not running then return ""
set out to ""
tell application "Google Chrome"
repeat with w in windows
if (id of w as text) is (item 1 of argv) then
repeat with t in tabs of w
set out to out & (URL of t) & linefeed
end repeat
end if
end repeat
end tell
return out
end run`

// chromeOpenScript opens a new window on the URL argv[1] and prints its id.
const chromeOpenScript = `on run argv
tell application "Google Chrome"
set w to make new window
set URL of active tab of w to (item 1 of argv)
return id of w as text
end tell
end run`

// chromeSplitScript brings window argv[1] to the front, confirms it got there,
// clicks Tab > Add Tab to New Split View, waits for that window to gain a tab,
// and prints the id of the new half, its active tab: a disabled or absent item raises no error of its own, and the click
// acts on whichever window is in front, so another window's tabs prove nothing.
const chromeSplitScript = `on run argv
set target to item 1 of argv
tell application "Google Chrome"
activate
repeat with w in windows
if (id of w as text) is target then
set index of w to 1
exit repeat
end if
end repeat
end tell
delay 0.3
tell application "Google Chrome"
if (id of front window as text) is not (item 1 of argv) then error "window " & target & " is not the front window"
end tell
tell application "System Events" to tell process "Google Chrome"
click menu item "Add Tab to New Split View" of menu "Tab" of menu bar 1
end tell
tell application "Google Chrome"
repeat 20 times
repeat with w in windows
if (id of w as text) is target and (count of tabs of w) > 1 then return id of active tab of w as text
end repeat
delay 0.1
end repeat
end tell
error "Tab > Add Tab to New Split View made no new tab in window " & target
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

// chromeShowTabScript points the tab whose id is argv[2] in window argv[1] at
// the URL argv[3], makes it the active tab and prints "shown"; it prints nothing
// when the window has no such tab.
const chromeShowTabScript = `on run argv
tell application "Google Chrome"
repeat with w in windows
if (id of w as text) is (item 1 of argv) then
repeat with i from 1 to count of tabs of w
if (id of tab i of w as text) is (item 2 of argv) then
set URL of tab i of w to (item 3 of argv)
set active tab index of w to i
return "shown"
end if
end repeat
return ""
end if
end repeat
end tell
error "Chrome has no window " & (item 1 of argv)
end run`

// chromeNewTabScript opens a new tab on the URL argv[2] at the end of window
// argv[1], makes it the active tab and prints its id.
const chromeNewTabScript = `on run argv
tell application "Google Chrome"
repeat with w in windows
if (id of w as text) is (item 1 of argv) then
set t to make new tab at end of tabs of w with properties {URL:(item 2 of argv)}
set active tab index of w to (count of tabs of w)
return id of t as text
end if
end repeat
end tell
error "Chrome has no window " & (item 1 of argv)
end run`

// chromeSelectTabScript makes tab number argv[2] the active tab of window
// argv[1].
const chromeSelectTabScript = `on run argv
tell application "Google Chrome"
repeat with w in windows
if (id of w as text) is (item 1 of argv) then
set active tab index of w to ((item 2 of argv) as integer)
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
// splits a new one with its right half on right (empty when right is ""). On a
// window that exists, right comes into view (see pointRight). It never closes
// or moves a tab it did not open, and navigates none but the right half's.
func Window(socket, pane, right string) (WindowResult, error) {
	name := Query(socket, pane, "#{session_name}")
	if name == "" {
		return WindowResult{}, errors.New("the session has ended")
	}
	if id := Query(socket, pane, "#{"+windowOption+"}"); id != "" {
		out, err := osascript(chromeHasWindowScript, id)
		if err != nil {
			return WindowResult{}, stepError("find the window", err)
		}
		tabs := tabURLs(out)
		if terminalTab(tabs, socket, name) >= 0 {
			if right != "" {
				if err := pointRight(socket, pane, name, id, tabs, right); err != nil {
					return WindowResult{}, stepError("set the right half", err)
				}
			}
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
	half, err := osascript(chromeSplitScript, id)
	if err != nil {
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
	if _, err := Run(socket, "set-option", "-t", "="+name+":", rightOption, strings.TrimSpace(half)); err != nil {
		return WindowResult{}, stepError("store the right half's id", err)
	}
	if _, err := osascript(chromeRaiseScript, id); err != nil {
		return WindowResult{}, stepError("raise the window", err)
	}
	return WindowResult{Window: "opened", ID: id}, nil
}

// pointRight brings right into view in window id, whose tab URLs are tabs: it
// selects a tab already showing right, else points the session's right-half
// tab at it and selects that, else opens a new tab on it, which becomes the
// right half.
func pointRight(socket, pane, name, id string, tabs []string, right string) error {
	for i, u := range tabs {
		if samePage(u, right) {
			_, err := osascript(chromeSelectTabScript, id, strconv.Itoa(i+1))
			return err
		}
	}
	if half := Query(socket, pane, "#{"+rightOption+"}"); half != "" {
		out, err := osascript(chromeShowTabScript, id, half, right)
		if err != nil {
			return err
		}
		if strings.TrimSpace(out) == "shown" {
			return nil
		}
	}
	out, err := osascript(chromeNewTabScript, id, right)
	if err != nil {
		return err
	}
	_, err = Run(socket, "set-option", "-t", "="+name+":", rightOption, strings.TrimSpace(out))
	return err
}

// samePage says whether a tab's URL is the page asked for: a page takes its
// token off its address when it loads, and may add a #section.
func samePage(tab, asked string) bool {
	a, errA := pageKey(tab)
	b, errB := pageKey(asked)
	return errA == nil && errB == nil && a == b
}

func pageKey(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Del("t")
	u.RawQuery, u.Fragment = q.Encode(), ""
	return u.String(), nil
}

// tabURLs is the has-window script's output as a list.
func tabURLs(out string) []string {
	var urls []string
	for _, u := range splitLines(strings.TrimSpace(out)) {
		urls = append(urls, strings.TrimSpace(u))
	}
	return urls
}

// terminalTab is the index of the session's terminal page in urls, or -1.
func terminalTab(urls []string, socket, name string) int {
	return slices.IndexFunc(urls, func(u string) bool { return pages.TermMatches(u, socket, name) })
}

// shownWindow is the id of pane's session's stored Chrome window when it still
// holds the session's terminal page, else "". A failed look says no.
func shownWindow(socket, pane, name string) string {
	id := Query(socket, pane, "#{"+windowOption+"}")
	if id == "" {
		return ""
	}
	out, err := osascript(chromeHasWindowScript, id)
	if err != nil || terminalTab(tabURLs(out), socket, name) < 0 {
		return ""
	}
	return id
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
