package proj

import (
	"errors"
	"net/url"
	"strings"
	"testing"
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
	origScript, origPages := osascript, ensurePages
	t.Cleanup(func() { osascript, ensurePages = origScript, origPages })
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
		case chromeSplitScript:
			return "88\n", nil
		case chromeShowTabScript:
			return "shown\n", nil
		case chromeNewTabScript:
			return "99\n", nil
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

func TestWindowOpensAndSplitsWithNoStoredID(t *testing.T) {
	socket, name, pane := terminalServer(t)
	calls := fakeChrome(t, "", nil)
	got, err := Window(socket, pane, "http://right.example/")
	if err != nil {
		t.Fatal(err)
	}
	if got != (WindowResult{Window: "opened", ID: "77"}) {
		t.Fatalf("got %+v", got)
	}
	sameScripts(t, *calls, chromeOpenScript, chromeSplitScript, chromeSetRightScript, chromeRaiseScript)
	open := (*calls)[0].args[0]
	for _, want := range []string{"127.0.0.1:9/term?", "t=tok", "session=terminal-test%2Fw"} {
		if !strings.Contains(open, want) {
			t.Fatalf("opened %q, want it to contain %q", open, want)
		}
	}
	if a := (*calls)[2].args; a[0] != "77" || a[1] != "http://right.example/" {
		t.Fatalf("set right with %q", a)
	}
	if got := Query(socket, "="+name+":", "#{@chrome_window}"); got != "77" {
		t.Fatalf("@chrome_window = %q", got)
	}
	if got := Query(socket, "="+name+":", "#{@chrome_right}"); got != "88" {
		t.Fatalf("@chrome_right = %q", got)
	}
}

func TestWindowLeavesTheRightHalfEmptyWithoutRight(t *testing.T) {
	socket, name, pane := terminalServer(t)
	calls := fakeChrome(t, "", nil)
	if _, err := Window(socket, pane, ""); err != nil {
		t.Fatal(err)
	}
	sameScripts(t, *calls, chromeOpenScript, chromeSplitScript, chromeRaiseScript)
	if got := Query(socket, "="+name+":", "#{@chrome_right}"); got != "88" {
		t.Fatalf("@chrome_right = %q", got)
	}
}

func TestWindowRaisesAStoredWindowChromeStillHas(t *testing.T) {
	socket, name, pane := terminalServer(t)
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

// storedWindow gives the session the stored window 55 and returns its
// terminal page's URL as Chrome shows it.
func storedWindow(t *testing.T, socket, name string) string {
	t.Helper()
	if _, err := Run(socket, "set-option", "-t", "="+name+":", "@chrome_window", "55"); err != nil {
		t.Fatal(err)
	}
	return "http://127.0.0.1:9/term?" + url.Values{"session": {name}, "socket": {socket}}.Encode() + "\n"
}

// A Playwright tab beside the terminal is not the right half: only the tab
// stored as the right half is navigated, and it is selected.
func TestWindowRaisePointsTheStoredRightTabAtRight(t *testing.T) {
	socket, name, pane := terminalServer(t)
	term := storedWindow(t, socket, name)
	if _, err := Run(socket, "set-option", "-t", "="+name+":", "@chrome_right", "88"); err != nil {
		t.Fatal(err)
	}
	calls := fakeChrome(t, term+"https://playwright.example/\nchrome://newtab/\n", nil)
	got, err := Window(socket, pane, "http://right.example/?t=tok&f=a")
	if err != nil {
		t.Fatal(err)
	}
	if got != (WindowResult{Window: "raised", ID: "55"}) {
		t.Fatalf("got %+v", got)
	}
	sameScripts(t, *calls, chromeHasWindowScript, chromeShowTabScript, chromeRaiseScript)
	if a := (*calls)[1].args; strings.Join(a, " ") != "55 88 http://right.example/?t=tok&f=a" {
		t.Fatalf("show tab with %q", a)
	}
}

// With the right half closed, or never recorded, a new tab takes the URL and
// becomes the right half; no tab that is there changes.
func TestWindowRaiseOpensANewTabWhenTheRightTabIsGone(t *testing.T) {
	for name, stored := range map[string]string{"closed": "88", "never stored": ""} {
		t.Run(name, func(t *testing.T) {
			socket, session, pane := terminalServer(t)
			term := storedWindow(t, socket, session)
			if stored != "" {
				if _, err := Run(socket, "set-option", "-t", "="+session+":", "@chrome_right", stored); err != nil {
					t.Fatal(err)
				}
			}
			calls := fakeChrome(t, "https://first.example/\n"+term, nil)
			fake := osascript
			osascript = func(script string, args ...string) (string, error) {
				out, err := fake(script, args...)
				if script == chromeShowTabScript {
					return "", err // the tab is gone
				}
				return out, err
			}
			if _, err := Window(socket, pane, "http://right.example/"); err != nil {
				t.Fatal(err)
			}
			want := []string{chromeHasWindowScript}
			if stored != "" {
				want = append(want, chromeShowTabScript)
			}
			sameScripts(t, *calls, append(want, chromeNewTabScript, chromeRaiseScript)...)
			if got := Query(socket, "="+session+":", "#{@chrome_right}"); got != "99" {
				t.Fatalf("@chrome_right = %q", got)
			}
		})
	}
}

// A page that took its token shows its address without it, and may add a
// section; that is still the page asked for.
func TestWindowRaiseSelectsATabAlreadyShowingRight(t *testing.T) {
	socket, name, pane := terminalServer(t)
	term := storedWindow(t, socket, name)
	calls := fakeChrome(t, term+"chrome://newtab/\nhttp://right.example/?f=a#decisions\n", nil)
	if _, err := Window(socket, pane, "http://right.example/?t=tok&f=a"); err != nil {
		t.Fatal(err)
	}
	sameScripts(t, *calls, chromeHasWindowScript, chromeSelectTabScript, chromeRaiseScript)
	if a := (*calls)[1].args; strings.Join(a, " ") != "55 3" {
		t.Fatalf("selected %q", a)
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
	sameScripts(t, *calls, chromeHasWindowScript, chromeOpenScript, chromeSplitScript, chromeRaiseScript)
	if got := Query(socket, "="+name+":", "#{@chrome_window}"); got != "77" {
		t.Fatalf("@chrome_window = %q", got)
	}
}

func TestWindowSplitFailureNamesTheStepAndAccessibility(t *testing.T) {
	socket, name, pane := terminalServer(t)
	macos := errors.New("System Events got an error: osascript is not allowed assistive access. (-1719)")
	calls := fakeChrome(t, "", map[string]error{chromeSplitScript: macos})
	_, err := Window(socket, pane, "http://right.example/")
	if err == nil {
		t.Fatal("no error")
	}
	for _, want := range []string{"split", "assistive access", "Accessibility"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("err = %q, want it to contain %q", err, want)
		}
	}
	sameScripts(t, *calls, chromeOpenScript, chromeSplitScript)
	if got := Query(socket, "="+name+":", "#{@chrome_window}"); got != "" {
		t.Fatalf("@chrome_window = %q, want nothing stored", got)
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

// The AppleScript itself is checked live; this holds the two guards in place.
func TestSplitScriptActsOnlyOnTheTargetWindow(t *testing.T) {
	if strings.Contains(chromeSplitScript, "of tabs of front window") {
		t.Error("the split waits on the front window, which may not be the target")
	}
	if !strings.Contains(chromeSplitScript, "(id of front window as text) is not (item 1 of argv)") {
		t.Error("the split does not confirm the target is in front before the click")
	}
}
