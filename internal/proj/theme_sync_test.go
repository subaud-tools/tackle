package proj

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// All external boundaries are stubbed; no desktop, sockets or terminal devices.
func themeFixture(t *testing.T) (*string, *int, *[][]byte, *bool) {
	t.Helper()
	recordHome(t)
	writeConfig(t, "session_themes = true\n")
	oldScript, oldRunner, oldSockets, oldWrite, oldNow := osascript, runner, themeSockets, writeTerminal, themeNow
	t.Cleanup(func() {
		osascript, runner, themeSockets, writeTerminal, themeNow = oldScript, oldRunner, oldSockets, oldWrite, oldNow
	})
	observation := "w1\tt1\t/dev/ttys067\nw1\tt2\t/dev/ttys068\nw2\tt3\t/dev/ttys069\n"
	reads := 0
	osascript = func(string, ...string) (string, error) {
		reads++
		if observation == "FAIL" {
			return "", errors.New("observation failed")
		}
		return observation, nil
	}
	themeSockets = func() []string { return []string{"proj-test"} }
	values := map[string]string{}
	runner = func(socket string, args ...string) (string, error) {
		if socket != "proj-test" {
			t.Fatalf("unexpected socket %q", socket)
		}
		switch args[0] {
		case "list-panes":
			return strings.ReplaceAll("a\t1\t0\t40\t%1\tpi\t/tmp\t0\t0\tagent-a\tpi\t\t1\t0\t0\nb\t1\t0\t40\t%2\tpi\t/tmp\t0\t0\tagent-b\tpi\t\t1\t0\t0\nc\t1\t0\t40\t%3\tpi\t/tmp\t0\t0\tagent-c\tpi\t\t1\t0\t0", "\t", "\x1f"), nil
		case "list-clients":
			return "$1\t/dev/ttys067\t101\n$2\t/dev/ttys068\t102\n$3\t/dev/ttys069\t103", nil
		case "list-sessions":
			var rows []string
			for _, s := range []struct{ id, name, stable string }{{"$1", "a", "stable-a"}, {"$2", "b", "stable-b"}, {"$3", "c", "stable-c"}} {
				row := []string{s.id, s.name, s.stable}
				for _, key := range themeSessionOptions {
					row = append(row, values[s.id+key])
				}
				rows = append(rows, strings.Join(row, "\t"))
			}
			return strings.Join(rows, "\n"), nil
		case "list-windows":
			var rows []string
			for _, id := range []string{"1", "2", "3"} {
				row := []string{"$" + id, "@" + id}
				for _, key := range themeWindowOptions {
					row = append(row, values["@"+id+key])
				}
				rows = append(rows, strings.Join(row, "\t"))
			}
			return strings.Join(rows, "\n"), nil
		case "set-option", "set-window-option":
			for i := 0; i < len(args); i += 6 {
				if i+5 > len(args) {
					t.Fatalf("bad write %q", args)
				}
				values[args[i+2]+args[i+3]] = args[i+4]
			}
			return "", nil
		}
		t.Fatalf("unexpected command %q", args)
		return "", nil
	}
	var writes [][]byte
	fail := false
	writeTerminal = func(tty string, b []byte) error {
		if !strings.HasPrefix(tty, "/dev/ttys06") {
			t.Fatalf("unmatched tty %q", tty)
		}
		if fail {
			return errors.New("write failed")
		}
		writes = append(writes, bytes.Clone(b))
		return nil
	}
	themeNow = func() time.Time { return time.Unix(100, 0) }
	return &observation, &reads, &writes, &fail
}

func TestReconcileTTYMovesKeepStableIdentityAndAbsentLayout(t *testing.T) {
	observation, _, writes, _ := themeFixture(t)
	if err := UpdateRecord(func(r *Record) bool {
		r.Sessions["old-a"] = Entry{ID: "stable-a", Socket: "proj-test", Conversation: "old-agent", GroupID: "saved", Shade: "deep"}
		r.Sessions["absent"] = Entry{ID: "absent", GroupID: "saved", Shade: "bright"}
		r.Groups = map[string]ThemeGroup{"saved": {Hue: "blue"}}
		r.NativeGroups = map[string]string{"w1": "saved"}
		r.Layout.Windows = [][]string{{"old-a", "absent"}}
		return true
	}); err != nil {
		t.Fatal(err)
	}
	if err := ReconcileThemes(true); err != nil {
		t.Fatal(err)
	}
	r, _ := LoadRecord()
	a := r.Sessions["a"]
	if a.ID != "stable-a" || a.Shade != "deep" || a.GroupID != "saved" || len(*writes) != 3 {
		t.Fatalf("first reconciliation: a=%+v writes=%d", a, len(*writes))
	}
	if _, ok := r.Sessions["old-a"]; ok {
		t.Fatal("rename left old name")
	}
	*observation = "w2\tt1\t/dev/ttys067\nw1\tt2\t/dev/ttys068\nw2\tt3\t/dev/ttys069\n"
	if err := ReconcileThemes(true); err != nil {
		t.Fatal(err)
	}
	r, _ = LoadRecord()
	moved := r.Sessions["a"]
	if moved.ID != a.ID || moved.Shade != a.Shade || moved.GroupID == a.GroupID || moved.GroupID != r.Sessions["c"].GroupID {
		t.Fatalf("move: %+v -> %+v", a, moved)
	}
	if !reflect.DeepEqual(r.Layout.Windows, [][]string{{"a", "c", "absent"}, {"b"}}) {
		t.Fatalf("layout %q", r.Layout.Windows)
	}
	if len(*writes) != 4 {
		t.Fatalf("move rewrote unchanged clients: %d writes", len(*writes))
	}
}

func TestReconcileThrottleUnchangedAndRetry(t *testing.T) {
	_, reads, writes, fail := themeFixture(t)
	*fail = true
	if err := ReconcileThemes(false); err == nil {
		t.Fatal("write failure swallowed")
	}
	*fail = false
	themeNow = func() time.Time { return time.Unix(106, 0) }
	if err := ReconcileThemes(false); err != nil {
		t.Fatal(err)
	}
	if len(*writes) != 3 {
		t.Fatalf("failed outputs not retried: %d", len(*writes))
	}
	if err := ReconcileThemes(false); err != nil {
		t.Fatal(err)
	}
	if *reads != 2 {
		t.Fatalf("throttle: %d reads", *reads)
	}
	themeNow = func() time.Time { return time.Unix(112, 0) }
	if err := ReconcileThemes(false); err != nil {
		t.Fatal(err)
	}
	if len(*writes) != 3 {
		t.Fatalf("unchanged clients rewritten: %d", len(*writes))
	}
}

func TestFailedObservationLeavesRecordAndAppliedState(t *testing.T) {
	observation, _, _, _ := themeFixture(t)
	if err := ReconcileThemes(true); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(RecordPath())
	statePath := filepath.Join(filepath.Dir(RecordPath()), "theme-sync.json")
	beforeState, _ := os.ReadFile(statePath)
	*observation = "FAIL"
	if err := ReconcileThemes(true); err == nil {
		t.Fatal("observation failure swallowed")
	}
	after, _ := os.ReadFile(RecordPath())
	if !bytes.Equal(before, after) {
		t.Fatal("failed observation changed record")
	}
	afterState, _ := os.ReadFile(statePath)
	if !bytes.Equal(beforeState, afterState) {
		t.Fatal("failed observation changed applied state")
	}
}

func TestReconcileSkipsUnchangedTmuxOptions(t *testing.T) {
	themeFixture(t)
	external := runner
	writes := 0
	runner = func(socket string, args ...string) (string, error) {
		if args[0] == "set-option" || args[0] == "set-window-option" {
			writes++
		}
		return external(socket, args...)
	}
	if err := ReconcileThemes(true); err != nil {
		t.Fatal(err)
	}
	if writes != 3 {
		t.Fatalf("initial tmux batches=%d", writes)
	}
	if err := ReconcileThemes(true); err != nil {
		t.Fatal(err)
	}
	if writes != 3 {
		t.Fatalf("unchanged refresh wrote tmux options: %d", writes)
	}
}

func TestReconcileChangedClientReappliesOnlyItsTerminal(t *testing.T) {
	_, _, writes, _ := themeFixture(t)
	if err := ReconcileThemes(true); err != nil {
		t.Fatal(err)
	}
	external := runner
	runner = func(socket string, args ...string) (string, error) {
		out, err := external(socket, args...)
		if args[0] == "list-clients" {
			out = strings.Replace(out, "101", "201", 1)
		}
		return out, err
	}
	if err := ReconcileThemes(true); err != nil {
		t.Fatal(err)
	}
	if len(*writes) != 4 {
		t.Fatalf("changed client writes=%d, want one new delivery", len(*writes))
	}
}

func TestSaveLayoutUsesTTYJoinWithoutTerminalOutput(t *testing.T) {
	_, _, writes, _ := themeFixture(t)
	writeConfig(t, "session_themes = false\n")
	w, tabs, err := SaveLayout(nil)
	if err != nil || w != 2 || tabs != 3 {
		t.Fatalf("layout=(%d,%d,%v)", w, tabs, err)
	}
	if len(*writes) != 0 {
		t.Fatal("layout observation wrote terminal colors")
	}
	r, _ := LoadRecord()
	if !reflect.DeepEqual(r.Layout.Windows, [][]string{{"a", "b"}, {"c"}}) || r.Sessions["a"].ID != "stable-a" {
		t.Fatalf("TTY layout=%+v", r)
	}
}

func TestTmuxFailureDoesNotInvalidateSuccessfulTerminalDelivery(t *testing.T) {
	_, _, writes, _ := themeFixture(t)
	if err := ReconcileThemes(true); err != nil {
		t.Fatal(err)
	}
	external := runner
	runner = func(socket string, args ...string) (string, error) {
		if args[0] == "set-option" {
			return "", errors.New("tmux option write failed")
		}
		out, err := external(socket, args...)
		if args[0] == "list-sessions" {
			lines := strings.Split(out, "\n")
			fields := strings.Split(lines[0], "\t")
			fields[3] = ""
			lines[0] = strings.Join(fields, "\t")
			out = strings.Join(lines, "\n")
		}
		return out, err
	}
	if err := ReconcileThemes(true); err == nil {
		t.Fatal("tmux failure swallowed")
	}
	runner = external
	if err := ReconcileThemes(true); err != nil {
		t.Fatal(err)
	}
	if len(*writes) != 3 {
		t.Fatalf("unchanged terminal was rewritten after a tmux-only failure: %d", len(*writes))
	}
}

func TestMultipleTerminalsForSessionWriteTmuxOptionsOnce(t *testing.T) {
	observation, _, writes, _ := themeFixture(t)
	*observation += "w1\tt4\t/dev/ttys066\nplain\tt5\t/dev/ttys099\n"
	external := runner
	batches := 0
	runner = func(socket string, args ...string) (string, error) {
		if args[0] == "set-option" || args[0] == "set-window-option" {
			batches++
		}
		out, err := external(socket, args...)
		if args[0] == "list-clients" {
			out += "\n$1\t/dev/ttys066\t104"
		}
		return out, err
	}
	if err := ReconcileThemes(true); err != nil {
		t.Fatal(err)
	}
	if batches != 3 || len(*writes) != 4 {
		t.Fatalf("tmux batches=%d terminal writes=%d", batches, len(*writes))
	}
	r, _ := LoadRecord()
	if len(r.Groups) != 2 {
		t.Fatalf("unmanaged terminal consumed a group: %+v", r.Groups)
	}
}

func TestDisabledRefreshLeavesDesktopAndRecordUntouched(t *testing.T) {
	for _, config := range []string{"missing", "session_themes = false\n"} {
		t.Run(config, func(t *testing.T) {
			_, reads, writes, _ := themeFixture(t)
			t.Setenv("XDG_CONFIG_HOME", t.TempDir())
			if config != "missing" {
				writeConfig(t, config)
			}
			if err := UpdateRecord(func(r *Record) bool {
				r.Groups = map[string]ThemeGroup{"saved": {Hue: "rose"}}
				r.Sessions["a"] = Entry{ID: "stable-a", GroupID: "saved", Shade: "deep", Conversation: "agent-a"}
				return true
			}); err != nil {
				t.Fatal(err)
			}
			before, _ := os.ReadFile(RecordPath())
			oldRunner := runner
			runner = func(string, ...string) (string, error) {
				t.Error("disabled refresh accessed tmux")
				return "", errors.New("unexpected tmux access")
			}
			t.Cleanup(func() { runner = oldRunner })
			if err := ReconcileThemes(false); err != nil {
				t.Error(err)
			}
			if *reads != 0 || len(*writes) != 0 {
				t.Errorf("disabled refresh observed %d times, wrote %d palettes", *reads, len(*writes))
			}
			q, err := ThemeForAgent("agent-a")
			if err != nil || q.Assigned || q.Theme != nil {
				t.Errorf("disabled query=%+v err=%v", q, err)
			}
			if err := SetWindowHue("native", "green"); err == nil || !strings.Contains(err.Error(), "disabled") {
				t.Fatalf("disabled hue setter: %v", err)
			}
			after, _ := os.ReadFile(RecordPath())
			if !bytes.Equal(before, after) {
				t.Fatal("disabled entrypoints changed saved assignments")
			}
		})
	}
}

func TestDisabledForcedRefreshStillSavesLayoutWithoutAppearanceWrites(t *testing.T) {
	_, _, writes, _ := themeFixture(t)
	writeConfig(t, "session_themes = false\n")
	external := runner
	runner = func(socket string, args ...string) (string, error) {
		if args[0] == "set-option" || args[0] == "set-window-option" {
			t.Fatal("disabled forced refresh wrote appearance options")
		}
		return external(socket, args...)
	}
	if err := ReconcileThemes(true); err != nil {
		t.Fatal(err)
	}
	r, err := LoadRecord()
	if err != nil || !reflect.DeepEqual(r.Layout.Windows, [][]string{{"a", "b"}, {"c"}}) {
		t.Fatalf("disabled forced layout=%+v err=%v", r.Layout, err)
	}
	if len(*writes) != 0 {
		t.Fatal("disabled forced refresh wrote terminal colors")
	}
}
