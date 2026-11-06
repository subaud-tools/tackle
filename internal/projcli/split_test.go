package projcli

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeSplit stands in for proj.Split, recording its arguments.
func fakeSplit(t *testing.T, page string, err error) *[]string {
	t.Helper()
	var got []string
	orig := splitTab
	t.Cleanup(func() { splitTab = orig })
	splitTab = func(socket, session, right string) (string, error) {
		got = []string{socket, session, right}
		return page, err
	}
	return &got
}

func TestSplitNeedsSocketSessionAndRight(t *testing.T) {
	for _, args := range [][]string{
		{"split", "--session", "p/w", "--right", "http://127.0.0.1:1/"},
		{"split", "--socket", "/tmp/s", "--right", "http://127.0.0.1:1/"},
		{"split", "--socket", "/tmp/s", "--session", "p/w"},
	} {
		if _, code := captureStdout(t, args); code != 2 {
			t.Errorf("%q: exit %d, want 2", args, code)
		}
	}
}

func TestSplitPrintsThePage(t *testing.T) {
	const page = "https://proj.ws.example.com/term?i=1&session=p%2Fw&socket=%2Ftmp%2Fs&t=tok"
	got := fakeSplit(t, page, nil)
	out, code := captureStdout(t, []string{"split", "--socket", "/tmp/s", "--session", "p/w", "--right", "https://diazo.ws.example.com/p", "--json"})
	if code != 0 {
		t.Fatalf("exit %d: %s", code, out)
	}
	if strings.Join(*got, " ") != "/tmp/s p/w https://diazo.ws.example.com/p" {
		t.Fatalf("Split(%q)", *got)
	}
	var res map[string]string
	if err := json.Unmarshal([]byte(out), &res); err != nil || len(res) != 1 || res["split"] != page {
		t.Fatalf("--json printed %q (%v)", out, err)
	}
	out, code = captureStdout(t, []string{"split", "--socket", "/tmp/s", "--session", "p/w", "--right", "https://diazo.ws.example.com/p"})
	if code != 0 || strings.TrimSpace(out) != page {
		t.Fatalf("exit %d, printed %q; want the page", code, out)
	}
}

func TestSplitFailsWithAMessage(t *testing.T) {
	fakeSplit(t, "", errors.New("no session \"p/w\" on /tmp/s"))
	var out, errw bytes.Buffer
	if code := Dispatch([]string{"split", "--socket", "/tmp/s", "--session", "p/w", "--right", "http://127.0.0.1:1/", "--json"}, &out, &errw); code != 1 {
		t.Fatalf("exit %d, want 1", code)
	}
	if out.Len() != 0 || !strings.Contains(errw.String(), "no session") {
		t.Fatalf("out %q, err %q", out.String(), errw.String())
	}
}

// The real Split against a socket with no server: the session does not
// exist, so it exits 1 before starting proj serve.
func TestSplitOnAMissingSessionExits1(t *testing.T) {
	sock := filepath.Join(os.TempDir(), "pjsplit-none")
	var out, errw bytes.Buffer
	if code := Dispatch([]string{"split", "--socket", sock, "--session", "p/w", "--right", "http://127.0.0.1:1/"}, &out, &errw); code != 1 {
		t.Fatalf("exit %d, want 1 (err %q)", code, errw.String())
	}
	if !strings.Contains(errw.String(), "p/w") {
		t.Fatalf("err %q", errw.String())
	}
	errw.Reset()
	if code := Dispatch([]string{"split", "--socket", sock, "--session", "p/w", "--right", "https://evil.example/"}, &out, &errw); code != 1 || !strings.Contains(errw.String(), "127.0.0.1") {
		t.Fatalf("a foreign page: exit %d, err %q", code, errw.String())
	}
}

func TestCommandsJSONListsSplit(t *testing.T) {
	out, code := captureStdout(t, []string{"commands", "--json"})
	if code != 0 || !strings.Contains(out, `"name": "split"`) {
		t.Fatalf("proj commands --json: exit %d, no split: %s", code, out)
	}
}
