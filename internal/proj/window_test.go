package proj

import (
	"errors"
	"strings"
	"testing"
)

// chromeCall is one AppleScript the fake Chrome ran.
type chromeCall struct {
	script string
	args   []string
}

// fakeChrome swaps the osascript seam and the page server's start: scripts in
// fail return their error, and the has-window script answers hasWindow. A new
// window gets the id 77.
func fakeChrome(t *testing.T, hasWindow bool, fail map[string]error) *[]chromeCall {
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
			if hasWindow {
				return "yes\n", nil
			}
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

func TestWindowOpensAndSplitsWithNoStoredID(t *testing.T) {
	socket, name, pane := terminalServer(t)
	calls := fakeChrome(t, false, nil)
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
}

func TestWindowLeavesTheRightHalfEmptyWithoutRight(t *testing.T) {
	socket, _, pane := terminalServer(t)
	calls := fakeChrome(t, false, nil)
	if _, err := Window(socket, pane, ""); err != nil {
		t.Fatal(err)
	}
	sameScripts(t, *calls, chromeOpenScript, chromeSplitScript, chromeRaiseScript)
}

func TestWindowRaisesAStoredWindowChromeStillHas(t *testing.T) {
	socket, name, pane := terminalServer(t)
	if _, err := Run(socket, "set-option", "-t", "="+name+":", "@chrome_window", "55"); err != nil {
		t.Fatal(err)
	}
	calls := fakeChrome(t, true, nil)
	got, err := Window(socket, pane, "http://right.example/")
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

func TestWindowReplacesAStoredWindowChromeNoLongerHas(t *testing.T) {
	socket, name, pane := terminalServer(t)
	if _, err := Run(socket, "set-option", "-t", "="+name+":", "@chrome_window", "55"); err != nil {
		t.Fatal(err)
	}
	calls := fakeChrome(t, false, nil)
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
	calls := fakeChrome(t, false, map[string]error{chromeSplitScript: macos})
	_, err := Window(socket, pane, "")
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
