package proj

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/subaud-tools/tackle/internal/proj/pages"
	"github.com/subaud-tools/tools-common/localweb"
	"github.com/subaud-tools/tools-common/pageserve"
	"github.com/subaud-tools/tools-common/splitpage"
)

// chromeCall is one AppleScript the fake Chrome ran.
type chromeCall struct {
	script string
	args   []string
}

// fakeChrome swaps the osascript seam and the page server's start: scripts in
// fail return their error, and the has-window script answers with tabs, the
// stored window's URLs. A new window gets the id 77.
func fakeChrome(t *testing.T, tabs string, fail map[string]error) *[]chromeCall {
	t.Helper()
	origScript, origPages, origPoll := osascript, ensurePages, attachPoll
	t.Cleanup(func() { osascript, ensurePages, attachPoll = origScript, origPages, origPoll })
	attachPoll = 0
	ensurePages = func() (string, string, error) { return "127.0.0.1:9", "tok", nil }
	var calls []chromeCall
	osascript = func(script string, args ...string) (string, error) {
		calls = append(calls, chromeCall{script, args})
		if err := fail[script]; err != nil {
			return "", err
		}
		switch script {
		case chromeHasWindowScript:
			return tabs, nil
		case chromeOpenScript:
			return "77\n", nil
		}
		return "", nil
	}
	return &calls
}

func scripts(calls []chromeCall) []string {
	var out []string
	for _, c := range calls {
		out = append(out, c.script)
	}
	return out
}

func sameScripts(t *testing.T, calls []chromeCall, want ...string) {
	t.Helper()
	got := scripts(calls)
	if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("ran %d scripts %q, want %d %q", len(got), got, len(want), want)
	}
}

// fakePageServer stands in for proj's page server, after fakeChrome: it
// records each tab it is asked to open, with the token it was sent.
func fakePageServer(t *testing.T) *[]splitpage.OpenRequest {
	t.Helper()
	var opened []splitpage.OpenRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req splitpage.OpenRequest
		if r.URL.Path != "/api/split/open" || r.Header.Get(localweb.TokenHeader) != "tok" || json.NewDecoder(r.Body).Decode(&req) != nil {
			http.Error(w, "unexpected", http.StatusBadRequest)
			return
		}
		opened = append(opened, req)
		_ = json.NewEncoder(w).Encode(splitpage.TabList{Session: req.Session, Tabs: []splitpage.Tab{}})
	}))
	t.Cleanup(srv.Close)
	orig := ensurePages
	t.Cleanup(func() { ensurePages = orig })
	ensurePages = func() (string, string, error) { return strings.TrimPrefix(srv.URL, "http://"), "tok", nil }
	return &opened
}

func TestWindowOpensTheTerminalPageAloneWithNoStoredID(t *testing.T) {
	socket, name, pane := terminalServer(t)
	calls := fakeChrome(t, "", nil)
	opened := fakePageServer(t)
	got, err := Window(socket, pane, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != (WindowResult{Window: "opened", ID: "77"}) {
		t.Fatalf("got %+v", got)
	}
	sameScripts(t, *calls, chromeOpenScript, chromeRaiseScript)
	open := (*calls)[0].args[0]
	for _, want := range []string{"/term?", "t=tok", "session=terminal-test%2Fw"} {
		if !strings.Contains(open, want) {
			t.Fatalf("opened %q, want it to contain %q", open, want)
		}
	}
	if len(*opened) != 0 {
		t.Fatalf("opened tabs %+v, want none", *opened)
	}
	if got := Query(socket, "="+name+":", "#{@chrome_window}"); got != "77" {
		t.Fatalf("@chrome_window = %q", got)
	}
}

func TestWindowOpensRightAsATabOfTheNewWindowsSession(t *testing.T) {
	socket, name, pane := terminalServer(t)
	calls := fakeChrome(t, "", nil)
	opened := fakePageServer(t)
	if _, err := Window(socket, pane, "http://127.0.0.1:4321/plan?t=x"); err != nil {
		t.Fatal(err)
	}
	sameScripts(t, *calls, chromeOpenScript, chromeRaiseScript)
	if len(*opened) != 1 || (*opened)[0].Session != name || (*opened)[0].Address != "http://127.0.0.1:4321/plan?t=x" {
		t.Fatalf("opened tabs %+v", *opened)
	}
}

func TestWindowRaisesAStoredWindowChromeStillHas(t *testing.T) {
	socket, name, pane := terminalServer(t)
	pageShows(t, socket, name)
	if _, err := Run(socket, "set-option", "-t", "="+name+":", "@chrome_window", "55"); err != nil {
		t.Fatal(err)
	}
	calls := fakeChrome(t, "http://127.0.0.1:9/term?"+url.Values{"session": {name}, "socket": {socket}}.Encode()+"\nhttps://example.com/\n", nil)
	got, err := Window(socket, pane, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != (WindowResult{Window: "raised", ID: "55"}) {
		t.Fatalf("got %+v", got)
	}
	sameScripts(t, *calls, chromeHasWindowScript, chromeRaiseScript)
	if a := (*calls)[1].args; a[0] != "55" {
		t.Fatalf("raised %q", a)
	}
}

// Open in Chrome on a window whose page Ghostty detached: the page loads again,
// which attaches it and so detaches Ghostty.
func TestWindowRaiseReloadsThePageWhenItIsNotAttached(t *testing.T) {
	socket, name, pane := terminalServer(t)
	term := storedWindow(t, socket, name)
	calls := fakeChrome(t, "https://example.com/\n"+term, nil)
	got, err := Window(socket, pane, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != (WindowResult{Window: "raised", ID: "55"}) {
		t.Fatalf("got %+v", got)
	}
	sameScripts(t, *calls, chromeHasWindowScript, chromeSetTabScript, chromeRaiseScript)
	a := (*calls)[1].args
	if a[0] != "55" || a[1] != "2" || !pages.TermMatches(a[2], socket, name) || !strings.Contains(a[2], "t=tok") {
		t.Fatalf("set tab %q", a)
	}
}

// diazo reads the state as soon as the press returns, so Window waits for the
// page to attach before it does.
func TestWindowWaitsForThePageToAttach(t *testing.T) {
	socket, name, pane := terminalServer(t)
	term := storedWindow(t, socket, name)
	if _, err := Run(socket, "set-option", "-t", "="+name+":", "@chrome_client", "7"); err != nil {
		t.Fatal(err)
	}
	fakeChrome(t, term, nil)
	lists := 0
	orig := execCommand
	t.Cleanup(func() { execCommand = orig })
	execCommand = func(name string, arg ...string) *exec.Cmd {
		if slices.Contains(arg, "list-clients") {
			if lists++; lists >= 3 {
				return orig("printf", "7")
			}
			return orig("printf", "")
		}
		return orig(name, arg...)
	}
	if _, err := Window(socket, pane, ""); err != nil {
		t.Fatal(err)
	}
	if lists != 3 {
		t.Fatalf("looked at the clients %d times, want 3", lists)
	}
}

// storedWindow gives the session the stored window 55 and returns its
// terminal page's URL as Chrome shows it.
func storedWindow(t *testing.T, socket, name string) string {
	t.Helper()
	if _, err := Run(socket, "set-option", "-t", "="+name+":", "@chrome_window", "55"); err != nil {
		t.Fatal(err)
	}
	return "http://127.0.0.1:9/term?" + url.Values{"session": {name}, "socket": {socket}}.Encode() + "\n"
}

// On a window that exists, right opens as a tab of the session on proj's page
// server, then the window comes to the front; no Chrome tab is touched.
func TestWindowRaiseOpensRightAsATabOfTheSession(t *testing.T) {
	socket, name, pane := terminalServer(t)
	pageShows(t, socket, name)
	term := storedWindow(t, socket, name)
	calls := fakeChrome(t, term+"https://playwright.example/\n", nil)
	opened := fakePageServer(t)
	got, err := Window(socket, pane, "http://127.0.0.1:4321/plan?t=tok&f=a")
	if err != nil {
		t.Fatal(err)
	}
	if got != (WindowResult{Window: "raised", ID: "55"}) {
		t.Fatalf("got %+v", got)
	}
	sameScripts(t, *calls, chromeHasWindowScript, chromeRaiseScript)
	if len(*opened) != 1 || (*opened)[0].Session != name || (*opened)[0].Address != "http://127.0.0.1:4321/plan?t=tok&f=a" {
		t.Fatalf("opened tabs %+v", *opened)
	}
}

func TestWindowRefusesARightThatIsNotOurPageAndTouchesNothing(t *testing.T) {
	socket, name, pane := terminalServer(t)
	storedWindow(t, socket, name)
	calls := fakeChrome(t, "", nil)
	opened := fakePageServer(t)
	_, err := Window(socket, pane, "https://example.com/")
	if err == nil || err.Error() != "only pages on 127.0.0.1 open in the split; give it as a link" {
		t.Fatalf("err = %v", err)
	}
	if len(*calls) != 0 || len(*opened) != 0 {
		t.Fatalf("ran %q and opened %+v", scripts(*calls), *opened)
	}
}

func TestWindowReplacesAStoredWindowChromeNoLongerHas(t *testing.T) {
	socket, name, pane := terminalServer(t)
	if _, err := Run(socket, "set-option", "-t", "="+name+":", "@chrome_window", "55"); err != nil {
		t.Fatal(err)
	}
	calls := fakeChrome(t, "", nil)
	got, err := Window(socket, pane, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != (WindowResult{Window: "opened", ID: "77"}) {
		t.Fatalf("got %+v", got)
	}
	sameScripts(t, *calls, chromeHasWindowScript, chromeOpenScript, chromeRaiseScript)
	if got := Query(socket, "="+name+":", "#{@chrome_window}"); got != "77" {
		t.Fatalf("@chrome_window = %q", got)
	}
}

func TestWindowOpensANewOneWhenTheStoredWindowShowsTheSameNameOnAnotherSocket(t *testing.T) {
	socket, name, pane := terminalServer(t)
	if _, err := Run(socket, "set-option", "-t", "="+name+":", "@chrome_window", "55"); err != nil {
		t.Fatal(err)
	}
	other := "http://127.0.0.1:9/term?" + url.Values{"session": {name}, "socket": {socket + "-other"}}.Encode() + "\n"
	calls := fakeChrome(t, other, nil)
	got, err := Window(socket, pane, "")
	if err != nil {
		t.Fatal(err)
	}
	if got.Window != "opened" {
		t.Fatalf("got %+v; ran %q", got, scripts(*calls))
	}
}

// osascriptOnPath puts a stand-in for osascript on PATH that logs each call as
// ping or real, and answers -600 to the first failPing pings and failReal real
// scripts. It returns the log's path.
func osascriptOnPath(t *testing.T, failPing, failReal int) string {
	t.Helper()
	dir := t.TempDir()
	log := filepath.Join(dir, "log")
	body := "#!/bin/sh\n" +
		"case \"$*\" in *\"count windows\"*) kind=ping; limit=" + strconv.Itoa(failPing) + ";; *) kind=real; limit=" + strconv.Itoa(failReal) + ";; esac\n" +
		"echo $kind >> '" + log + "'\n" +
		"[ \"$(grep -c \"^$kind$\" '" + log + "')\" -le $limit ] && { echo 'execution error: Application isn\u2019t running. (-600)' >&2; exit 1; }\n" +
		"echo ok\n"
	if err := os.WriteFile(filepath.Join(dir, "osascript"), []byte(body), 0o755); err != nil { //nolint:gosec // G306: the stand-in must run
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return log
}

func osascriptLog(t *testing.T, log string) string {
	t.Helper()
	b, err := os.ReadFile(log) //nolint:gosec // G304: the test's own temp file
	if err != nil {
		t.Fatal(err)
	}
	return strings.Join(strings.Fields(string(b)), " ")
}

func TestPingFailingOnceStillRunsTheScriptOnce(t *testing.T) {
	log := osascriptOnPath(t, 1, 0)
	if out, err := osascript(chromeOpenScript, "http://x/"); err != nil || strings.TrimSpace(out) != "ok" {
		t.Fatalf("out %q, err %v", out, err)
	}
	if got := osascriptLog(t, log); got != "ping ping real" {
		t.Fatalf("calls = %q", got)
	}
}

func TestPingFailingTwiceRunsNoScript(t *testing.T) {
	log := osascriptOnPath(t, 2, 0)
	if _, err := osascript(chromeOpenScript, "http://x/"); err == nil || !strings.Contains(failureText(err), "-600") {
		t.Fatalf("err = %v", err)
	}
	if got := osascriptLog(t, log); got != "ping ping" {
		t.Fatalf("calls = %q", got)
	}
}

// A script that failed may have done its work: a second window, a second tab.
func TestAScriptFailingWith600IsNotRunAgain(t *testing.T) {
	for name, script := range map[string]string{"Chrome": chromeOpenScript, "Ghostty": newWindowScript} {
		t.Run(name, func(t *testing.T) {
			log := osascriptOnPath(t, 0, 1)
			if _, err := osascript(script, "x"); err == nil {
				t.Fatal("want the error")
			}
			if got := osascriptLog(t, log); got != "ping real" {
				t.Fatalf("calls = %q", got)
			}
		})
	}
}

// Reading a session's state must not launch an app that isn't running.
func TestPingsLaunchNothing(t *testing.T) {
	for _, script := range []string{chromeHasWindowScript, focusTTYScript} {
		if ping := pingFor(script); !strings.HasPrefix(ping, "if application ") || !strings.Contains(ping, " is running then") {
			t.Errorf("ping = %q", ping)
		}
	}
}

// diazo_open places the plan's page without taking the session from Ghostty.
func TestWindowKeepingViewLeavesADetachedPageDetached(t *testing.T) {
	socket, name, pane := terminalServer(t)
	storedWindow(t, socket, name)
	calls := fakeChrome(t, "", nil)
	opened := fakePageServer(t)
	addr, _, _ := ensurePages()
	chromeShows(t, "http://"+addr+"/term?"+url.Values{"session": {name}, "socket": {socket}, "i": {pages.Instance("tok")}}.Encode()+"\n")
	got, err := WindowKeepingView(socket, pane, "http://127.0.0.1:4321/plan")
	if err != nil {
		t.Fatal(err)
	}
	if got != (WindowResult{Window: "raised", ID: "55"}) {
		t.Fatalf("got %+v", got)
	}
	sameScripts(t, *calls, chromeHasWindowScript, chromeRaiseScript)
	if len(*opened) != 1 {
		t.Fatalf("opened tabs %+v", *opened)
	}
}

func TestWindowOpensThePageToTakeTheSession(t *testing.T) {
	socket, name, pane := terminalServer(t)
	calls := fakeChrome(t, "", nil)
	if _, err := Window(socket, pane, ""); err != nil {
		t.Fatal(err)
	}
	if a := (*calls)[0].args[0]; !strings.Contains(a, "take=1") || !pages.TermMatches(a, socket, name) {
		t.Fatalf("opened %q", a)
	}
}

// chromeShows makes the fake Chrome's stored window show tabs, after
// fakePageServer, so a tab can be on the fake server's port.
func chromeShows(t *testing.T, tabs string) {
	t.Helper()
	inner := osascript
	osascript = func(script string, args ...string) (string, error) {
		out, err := inner(script, args...)
		if script == chromeHasWindowScript && err == nil {
			return tabs, nil
		}
		return out, err
	}
}

// After the page server restarts, the window's terminal page is still on the
// old port, saying its server stopped. Keeping the view, it loads from the new
// server, offering to show the session rather than taking it from Ghostty.
func TestWindowKeepingViewReloadsATerminalPageOnAStoppedServer(t *testing.T) {
	socket, name, pane := terminalServer(t)
	term := storedWindow(t, socket, name)
	calls := fakeChrome(t, term, nil)
	opened := fakePageServer(t)
	addr, _, _ := ensurePages()
	if _, err := WindowKeepingView(socket, pane, "http://127.0.0.1:4321/plan"); err != nil {
		t.Fatal(err)
	}
	sameScripts(t, *calls, chromeHasWindowScript, chromeSetTabScript, chromeRaiseScript)
	a := (*calls)[1].args
	if a[0] != "55" || a[1] != "1" || !strings.HasPrefix(a[2], "http://"+addr+"/term?") || !pages.TermMatches(a[2], socket, name) || !strings.Contains(a[2], "t=tok") || strings.Contains(a[2], "take=") {
		t.Fatalf("set tab %q, want the page on %s without take", a, addr)
	}
	if len(*opened) != 1 {
		t.Fatalf("opened tabs %+v", *opened)
	}
}

// A new window whose tab fails to open still opened, with the terminal in it:
// it reports opened, so the caller does not open the page somewhere else too,
// and says on stderr why the tab is missing.
func TestWindowOpensTheNewWindowWhenItsTabFails(t *testing.T) {
	socket, _, pane := terminalServer(t)
	calls := fakeChrome(t, "", nil)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "no split here", http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)
	ensurePages = func() (string, string, error) { return strings.TrimPrefix(srv.URL, "http://"), "tok", nil }
	var warned strings.Builder
	orig := warnings
	t.Cleanup(func() { warnings = orig })
	warnings = &warned
	got, err := Window(socket, pane, "http://127.0.0.1:4321/plan")
	if err != nil || got != (WindowResult{Window: "opened", ID: "77"}) {
		t.Fatalf("Window = %+v, %v; want opened", got, err)
	}
	sameScripts(t, *calls, chromeOpenScript, chromeRaiseScript)
	if !strings.Contains(warned.String(), "open the tab") {
		t.Fatalf("warned %q", warned.String())
	}
}

// A page server restarted on the same port leaves the window's terminal page
// dead on the same address: the page names the server it came from, so it
// loads from the new one, keeping the view.
func TestWindowKeepingViewReloadsAPageFromAnEarlierServerOnTheSamePort(t *testing.T) {
	socket, name, pane := terminalServer(t)
	storedWindow(t, socket, name)
	calls := fakeChrome(t, "", nil)
	opened := fakePageServer(t)
	addr, _, _ := ensurePages()
	chromeShows(t, "http://"+addr+"/term?"+url.Values{"session": {name}, "socket": {socket}, "i": {pages.Instance("an-earlier-token")}}.Encode()+"\n")
	if _, err := WindowKeepingView(socket, pane, "http://127.0.0.1:4321/plan"); err != nil {
		t.Fatal(err)
	}
	sameScripts(t, *calls, chromeHasWindowScript, chromeSetTabScript, chromeRaiseScript)
	if a := (*calls)[1].args[2]; strings.Contains(a, "take=") || !strings.Contains(a, "i="+pages.Instance("tok")) {
		t.Fatalf("set tab %q, want the page from the current server without take", a)
	}
	if len(*opened) != 1 {
		t.Fatalf("opened tabs %+v", *opened)
	}
}

// usePagesClient has ensurePages find the page server through c.
func usePagesClient(t *testing.T, c *pageserve.Client) {
	t.Helper()
	orig := pagesClient
	t.Cleanup(func() { pagesClient = orig })
	pagesClient = func() *pageserve.Client { return c }
}

// A server an older proj left running marks itself stale once proj is
// replaced: proj window stops it and starts this binary's before it opens.
func TestEnsurePagesReplacesAStaleServer(t *testing.T) {
	live := pageserve.Entry{Base: "http://127.0.0.1:1", Token: "old", Stale: true}
	running, stopped := true, false
	usePagesClient(t, &pageserve.Client{
		Find: func() (pageserve.Entry, error) {
			if !running {
				return pageserve.Entry{}, pageserve.ErrNotRunning
			}
			return live, nil
		},
		Stop: func(context.Context) error { running, stopped = false, true; return nil },
		Start: func() (pageserve.Entry, error) {
			running, live = true, pageserve.Entry{Base: "http://127.0.0.1:2", Token: "new"}
			return live, nil
		},
	})
	addr, token, err := ensurePages()
	if err != nil || addr != "127.0.0.1:2" || token != "new" || !stopped {
		t.Fatalf("ensurePages = %q, %q, %v (stopped %v); want the new server", addr, token, err, stopped)
	}
}

func TestEnsurePagesReusesAnUpToDateServer(t *testing.T) {
	usePagesClient(t, &pageserve.Client{
		Find: func() (pageserve.Entry, error) { return pageserve.Entry{Base: "http://127.0.0.1:1", Token: "tok"}, nil },
		Stop: func(context.Context) error { t.Error("stopped an up-to-date server"); return nil },
		Start: func() (pageserve.Entry, error) {
			t.Error("started a server with one running")
			return pageserve.Entry{}, nil
		},
	})
	if addr, token, err := ensurePages(); err != nil || addr != "127.0.0.1:1" || token != "tok" {
		t.Fatalf("ensurePages = %q, %q, %v; want the running server", addr, token, err)
	}
}
