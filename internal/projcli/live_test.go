package projcli

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/subaud-tools/tackle/internal/proj"
)

// proj live --json lists the running sessions only, each with the
// conversation its agent announced, else the record's.
func TestLiveJSONListsRunningSessions(t *testing.T) {
	withFakeRestore(t, map[string]proj.Entry{
		"b/1":    {Socket: "proj-b", Project: "b", Dir: "/b", Agent: "claude", Conversation: "rec-b"},
		"gone/1": {Socket: "proj-gone", Project: "gone", Dir: "/g", Agent: "pi", Conversation: "rec-g"},
	}, nil, nil)
	orig := liveSessions
	t.Cleanup(func() { liveSessions = orig })
	liveSessions = func() []proj.Session {
		return []proj.Session{
			{Name: "a/1", Socket: "proj-a", Dir: "/a", Agent: "pi", Conversation: "harness-a"},
			{Name: "b/1", Socket: "proj-b", Dir: "/b", Agent: "claude"},
			{Name: "c/1", Socket: "proj-c", Dir: "/c", Agent: "shell"},
		}
	}

	out, code := captureStdout(t, []string{"live", "--json"})
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	var got []liveJSON
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not a JSON array: %v\n%s", err, out)
	}
	want := []liveJSON{
		{Socket: "proj-a", Name: "a/1", Dir: "/a", Agent: "pi", Conversation: "harness-a"},
		{Socket: "proj-b", Name: "b/1", Dir: "/b", Agent: "claude", Conversation: "rec-b"},
		{Socket: "proj-c", Name: "c/1", Dir: "/c", Agent: "shell", Conversation: ""},
	}
	if len(got) != len(want) {
		t.Fatalf("live = %+v\nwant   %+v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("live[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
	for _, k := range []string{`"socket"`, `"name"`, `"dir"`, `"agent"`, `"conversation"`} {
		if !strings.Contains(out, k) {
			t.Fatalf("output lacks %s:\n%s", k, out)
		}
	}
}

// withRestoreFile writes sessions as the live file and fakes the platform.
func withRestoreFile(t *testing.T, sessions []liveJSON, goos string) string {
	t.Helper()
	b, err := json.Marshal(sessions)
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "live-sessions.json")
	if err := os.WriteFile(p, b, 0o644); err != nil {
		t.Fatal(err)
	}
	orig := restoreGOOS
	t.Cleanup(func() { restoreGOOS = orig })
	restoreGOOS = goos
	return p
}

// proj restore --from resumes pi with its conversation (pi --session <id>),
// skips a session already running, and opens no Ghostty window on Linux.
func TestRestoreFromFile(t *testing.T) {
	dir := t.TempDir()
	f := withFakeRestore(t, nil, map[string]error{"c/1": errors.New("boom")}, nil)
	liveState = func() proj.LiveState {
		return proj.LiveState{Running: map[string]bool{"b/1": true}, Attached: map[string]bool{}}
	}
	ghosttyCalled := false
	ghosttyOpen = func([][]proj.SavedSession, io.Writer) error { ghosttyCalled = true; return nil }
	p := withRestoreFile(t, []liveJSON{
		{Socket: "proj-a", Name: "a/1", Dir: dir, Agent: "pi", Conversation: "id1"},
		{Socket: "proj-b", Name: "b/1", Dir: dir, Agent: "pi", Conversation: "id1"},
		{Socket: "proj-c", Name: "c/1", Dir: dir, Agent: "pi"},
	}, "linux")

	out, code := captureStdout(t, []string{"restore", "--from", p, "--json"})
	if code != 0 {
		t.Fatalf("exit %d\n%s", code, out)
	}
	var got restoreJSON
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	if strings.Join(got.Restored, ",") != "a/1" || strings.Join(got.Skipped, ",") != "b/1" ||
		len(got.Failed) != 1 || got.Failed[0].Name != "c/1" || got.Failed[0].Error != "boom" {
		t.Fatalf("result = %+v", got)
	}
	for _, c := range f.restored {
		if c.name == "b/1" {
			t.Fatal("a running session was restored")
		}
	}
	if len(f.restored) == 0 || f.restored[0] != (restoreCall{"proj-a", "a/1", dir, "pi", "id1"}) {
		t.Fatalf("restore calls = %+v, want a/1 resumed with id1", f.restored)
	}
	if ghosttyCalled {
		t.Fatal("Ghostty opened on linux")
	}
}

// A file that can't be read exits 1; restore prints nothing as JSON then.
func TestRestoreFromUnreadableFile(t *testing.T) {
	withFakeRestore(t, nil, nil, nil)
	var out, errw bytes.Buffer
	code := Dispatch([]string{"restore", "--from", filepath.Join(t.TempDir(), "missing.json"), "--json"}, &out, &errw)
	if code != 1 {
		t.Fatalf("exit %d, want 1", code)
	}
}

// proj list --json carries each session's busy panes, an empty array for none.
func TestListJSONBusyPanes(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	orig := liveSessions
	t.Cleanup(func() { liveSessions = orig })
	liveSessions = func() []proj.Session {
		return []proj.Session{
			{Name: "a/1", Socket: "proj-a", BusyPanes: []proj.BusyPane{{Pane: "%4", Command: "sleep"}}},
			{Name: "b/1", Socket: "proj-b"},
		}
	}
	out, code := captureStdout(t, []string{"list", "--json"})
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	var got struct {
		Sessions []struct {
			Name      string          `json:"name"`
			BusyPanes []proj.BusyPane `json:"busy_panes"`
		} `json:"sessions"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Sessions) != 2 || len(got.Sessions[0].BusyPanes) != 1 ||
		got.Sessions[0].BusyPanes[0] != (proj.BusyPane{Pane: "%4", Command: "sleep"}) {
		t.Fatalf("sessions = %+v", got.Sessions)
	}
	if !strings.Contains(out, `"busy_panes": []`) {
		t.Fatalf("b/1 has no empty busy_panes array:\n%s", out)
	}
}
