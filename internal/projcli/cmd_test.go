package projcli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// captureStdout runs Dispatch(args, ...) and returns everything written to
// its out writer, plus the resulting exit code.
func captureStdout(t *testing.T, args []string) (string, int) {
	t.Helper()
	var out, errw bytes.Buffer
	code := Dispatch(args, &out, &errw)
	return out.String(), code
}

func TestRunListRequiresJSON(t *testing.T) {
	if _, code := captureStdout(t, []string{"list"}); code != 2 {
		t.Fatalf("proj list: got exit %d, want 2", code)
	}
}

func TestRunListJSON(t *testing.T) {
	out, code := captureStdout(t, []string{"list", "--json"})
	if code != 0 {
		t.Fatalf("proj list --json: got exit %d, want 0", code)
	}
	trimmed := strings.TrimSpace(out)
	if !strings.HasPrefix(trimmed, "{") || !strings.Contains(trimmed, `"projects"`) {
		t.Fatalf("proj list --json: output does not look like JSON: %q", out)
	}
}

func TestRunNewRequiresTarget(t *testing.T) {
	if _, code := captureStdout(t, []string{"new"}); code != 2 {
		t.Fatalf("proj new: got exit %d, want 2", code)
	}
}

func TestRunNewMalformedTarget(t *testing.T) {
	// Fails at parseNewTarget (missing '/') before any LoadRoots/tmux call, so
	// this needs no live server. A malformed target is a usage error (exit 2),
	// matching pre-migration behavior.
	if _, code := captureStdout(t, []string{"new", "noslash"}); code != 2 {
		t.Fatalf("proj new noslash: got exit %d, want 2", code)
	}
}

func TestRunListHelp(t *testing.T) {
	out, code := captureStdout(t, []string{"list", "-h"})
	if code != 0 {
		t.Fatalf("proj list -h: got exit %d, want 0", code)
	}
	if !strings.Contains(out, "Usage: proj list") {
		t.Fatalf("proj list -h: output missing usage line: %q", out)
	}
}

func TestRunNewHelp(t *testing.T) {
	out, code := captureStdout(t, []string{"new", "-h"})
	if code != 0 {
		t.Fatalf("proj new -h: got exit %d, want 0", code)
	}
	if !strings.Contains(out, "Usage: proj new") {
		t.Fatalf("proj new -h: output missing usage line: %q", out)
	}
}

func TestSnapshotRoute(t *testing.T) {
	t.Setenv("PROJ_HOME", t.TempDir())
	out, code := captureStdout(t, []string{"__snapshot", "--socket", "proj-restoretest-none"})
	if code != 0 || out != "" {
		t.Fatalf("__snapshot: exit %d output %q; want 0 and silence", code, out)
	}
}

func TestRunTerminalRequiresSocketAndPane(t *testing.T) {
	if _, code := captureStdout(t, []string{"terminal", "--json"}); code != 2 {
		t.Fatalf("proj terminal --json: got exit %d, want 2", code)
	}
}

func TestRunTerminalEndedJSON(t *testing.T) {
	sock := t.TempDir() + "/gone"
	out, code := captureStdout(t, []string{"terminal", "--socket", sock, "--pane", "%1", "--json"})
	if code != 0 {
		t.Fatalf("proj terminal: got exit %d, want 0", code)
	}
	var got map[string]string
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	if len(got) != 2 || got["state"] != "ended" || got["session"] != "" {
		t.Fatalf("got %v, want state ended and an empty session", got)
	}
}

func TestRunTerminalGoEndedExits1(t *testing.T) {
	sock := t.TempDir() + "/gone"
	var out, errw bytes.Buffer
	code := Dispatch([]string{"terminal", "--socket", sock, "--pane", "%1", "--go", "--json"}, &out, &errw)
	if code != 1 {
		t.Fatalf("proj terminal --go on an ended session: got exit %d, want 1", code)
	}
	if out.Len() != 0 {
		t.Fatalf("stdout = %q, want nothing", out.String())
	}
}
