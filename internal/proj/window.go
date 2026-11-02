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
	"github.com/subaud-tools/tools-common/splitpage"
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

// chromeSetTabScript points tab number argv[2] of window argv[1] at the URL
// argv[3], leaving which tab is active alone.
const chromeSetTabScript = `on run argv
tell application "Google Chrome"
repeat with w in windows
if (id of w as text) is (item 1 of argv) then
set URL of tab ((item 2 of argv) as integer) of w to (item 3 of argv)
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

// Window moves pane's session to its Chrome window and brings that to the
// front; see window. It has the page take the session.
func Window(socket, pane, right string) (WindowResult, error) {
	return window(socket, pane, right, true)
}

// WindowKeepingView is Window for a page that only goes in the window: a terminal
// page another terminal detached stays detached.
func WindowKeepingView(socket, pane, right string) (WindowResult, error) {
	return window(socket, pane, right, false)
}

// errNotOurPage refuses a right that a split page's frame may not show.
var errNotOurPage = errors.New("only pages on 127.0.0.1 open in the split; give it as a link")

// window brings pane's session's Chrome window to the front, or opens a new one
// on the session's split page, and opens right (when not "") as a tab of the
// session on proj's page server. On a window that exists, with move, a terminal
// page another terminal detached loads again to take the session.
// It never closes or moves a tab it did not open, and navigates none but the
// terminal page.
func window(socket, pane, right string, move bool) (WindowResult, error) {
	if right != "" && !localPage(right) {
		return WindowResult{}, errNotOurPage
	}
	name := Query(socket, pane, "#{session_name}")
	if name == "" {
		return WindowResult{}, errors.New("the session has ended")
	}
	if id := Query(socket, pane, "#{"+windowOption+"}"); id != "" {
		out, err := osascript(chromeHasWindowScript, id)
		if err != nil {
			return WindowResult{}, stepError("find the window", err)
		}
		if tab := terminalTab(tabURLs(out), socket, name); tab >= 0 {
			if move && !pageAttached(socket, pane, name) {
				addr, token, err := ensurePages()
				if err != nil {
					return WindowResult{}, stepError("start the page server", err)
				}
				if _, err := osascript(chromeSetTabScript, id, strconv.Itoa(tab+1), pages.TermURL(addr, token, socket, name)); err != nil {
					return WindowResult{}, stepError("attach the terminal page", err)
				}
			}
			if err := openTab(name, right); err != nil {
				return WindowResult{}, err
			}
			if _, err := osascript(chromeRaiseScript, id); err != nil {
				return WindowResult{}, stepError("raise the window", err)
			}
			if move {
				waitPageAttached(socket, pane, name)
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
	if _, err := Run(socket, "set-option", "-t", "="+name+":", windowOption, id); err != nil {
		return WindowResult{}, stepError("store the window id", err)
	}
	if err := openTab(name, right); err != nil {
		return WindowResult{}, err
	}
	if _, err := osascript(chromeRaiseScript, id); err != nil {
		return WindowResult{}, stepError("raise the window", err)
	}
	waitPageAttached(socket, pane, name)
	return WindowResult{Window: "opened", ID: id}, nil
}

// openTab opens right as a tab of session's split page on proj's page server,
// which shows it beside the terminal; nothing when right is "".
func openTab(session, right string) error {
	if right == "" {
		return nil
	}
	addr, token, err := ensurePages()
	if err != nil {
		return stepError("start the page server", err)
	}
	if _, err := splitpage.Open("http://"+addr, token, session, right, ""); err != nil {
		return stepError("open the tab", err)
	}
	return nil
}

// localPage says whether a split page's frame may show raw: one of our own
// pages, http on 127.0.0.1 with a port, as splitpage.Tabs requires.
func localPage(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "http" && u.User == nil && u.Hostname() == "127.0.0.1" && u.Port() != ""
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

// pageAttached says whether the terminal page's latest client is still attached
// to the session; a client the page lost is a detached page, or no page at all.
func pageAttached(socket, pane, name string) bool {
	pid := Query(socket, pane, "#{"+pages.ClientOption+"}")
	if pid == "" {
		return false
	}
	out, _ := Run(socket, "list-clients", "-t", "="+name, "-F", "#{client_pid}")
	return slices.Contains(splitLines(out), pid)
}

// waitPageAttached gives a loading page time to attach, so the state read after
// Window says where the session is; what that state says is the verdict.
func waitPageAttached(socket, pane, name string) {
	for range 10 {
		if pageAttached(socket, pane, name) {
			return
		}
		time.Sleep(attachPoll)
	}
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

// osascriptPingDelay is the wait before the one retry of a ping an app answered
// with "Application isn't running".
const osascriptPingDelay = 300 * time.Millisecond

// pingFor is the read-only script that asks script's app whether it is
// answering, "" for a script that talks to neither. It launches nothing: reading
// a session's state must not start Chrome.
func pingFor(script string) string {
	app := ""
	switch {
	case strings.Contains(script, `"Google Chrome"`):
		app = `application id "com.google.Chrome"`
	case strings.Contains(script, `"Ghostty"`):
		app = `application "Ghostty"`
	default:
		return ""
	}
	return "if " + app + " is running then\ntell " + app + " to count windows\nend if"
}

// ping runs a ping, once more when it fails with AppleScript's -600, which Chrome
// answers now and then while running: a retry right away works. Only a ping is
// retried: a script that failed may already have opened a window.
func ping(script string) error {
	_, err := runOsascript(script)
	if err != nil && strings.Contains(failureText(err), "(-600)") {
		time.Sleep(osascriptPingDelay)
		_, err = runOsascript(script)
	}
	return err
}

// stepError names the step that failed and carries Chrome's or macOS's message.
func stepError(step string, err error) error {
	return fmt.Errorf("%s: %s", step, failureText(err))
}

// failureText is a command's stderr when it has one, else the error itself.
func failureText(err error) string {
	var ee *exec.ExitError
	if errors.As(err, &ee) && len(ee.Stderr) > 0 {
		return strings.TrimSpace(string(ee.Stderr))
	}
	return err.Error()
}
