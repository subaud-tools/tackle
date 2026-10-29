package proj

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/subaud-tools/tools-common/sessiontheme"
)

func TestThemeIdentityLifecycle(t *testing.T) {
	for _, scenario := range []string{"v1 upgrade", "rename changed conversation", "move keeps shade", "shade reuse", "restore binding"} {
		t.Run(scenario, func(t *testing.T) {
			recordHome(t)
			switch scenario {
			case "v1 upgrade":
				raw := []byte(`{"version":1,"sessions":{"a/one":{"socket":"proj-a","project":"a","dir":"/work","agent":"pi","conversation":"old"},"a/two":{"socket":"proj-a","project":"a","dir":"/other","agent":"claude","conversation":"two"}},"layout":{"windows":[["a/two","a/one"]]}}`)
				if err := os.MkdirAll(filepath.Dir(RecordPath()), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(RecordPath(), raw, 0600); err != nil {
					t.Fatal(err)
				}
				r, err := LoadRecord()
				if err != nil {
					t.Fatal(err)
				}
				if r.Version != 2 || len(r.Sessions) != 2 || !reflect.DeepEqual(r.Layout.Windows, [][]string{{"a/two", "a/one"}}) {
					t.Fatalf("lost record: %+v", r)
				}
				e := r.Sessions["a/one"]
				if e.ID == "" || e.GroupID == "" || e.Shade != "deep" || e.Conversation != "old" || e.Agent != "pi" || e.Dir != "/work" {
					t.Fatalf("entry: %+v", e)
				}
				if r.Sessions["a/two"].Shade != "midnight" || r.Groups[e.GroupID].Hue != "blue" || r.Layout.GroupIDs[0] != e.GroupID {
					t.Fatalf("assignment: %+v", r)
				}
				if err := UpdateRecord(func(*Record) bool { return false }); err != nil {
					t.Fatal(err)
				}
				again, _ := LoadRecord()
				if !reflect.DeepEqual(r, again) {
					t.Fatal("second load allocated identities")
				}
			case "rename changed conversation":
				r := Record{Sessions: map[string]Entry{"a/old": {ID: "stable", GroupID: "group", Shade: "muted", Socket: "proj-a", Dir: "/work", Agent: "pi", Conversation: "old"}}, Layout: GhosttyLayout{Windows: [][]string{{"a/old"}}, GroupIDs: []string{"group"}}}
				mergeLive(&r, []LiveEntry{{Name: "a/new", Entry: Entry{ID: "stable", Socket: "proj-a", Dir: "/work", Agent: "pi", Conversation: "new"}}}, time.Now())
				e := r.Sessions["a/new"]
				if len(r.Sessions) != 1 || e.ID != "stable" || e.GroupID != "group" || e.Shade != "muted" || e.Conversation != "new" || e.Agent != "pi" || e.Dir != "/work" || r.Layout.Windows[0][0] != "a/new" {
					t.Fatalf("rename lost identity: %+v", r)
				}
			case "move keeps shade", "shade reuse":
				r := Record{Sessions: map[string]Entry{}}
				for _, n := range []string{"one", "two", "three", "four", "five", "six"} {
					r.Sessions[n] = Entry{ID: n}
				}
				assignWindowGroup(&r, "native-blue", []string{"one", "two", "three", "four", "five", "six"})
				for n, want := range map[string]sessiontheme.Shade{"one": "midnight", "two": "deep", "three": "muted", "four": "bright", "five": "midnight", "six": "deep"} {
					if r.Sessions[n].Shade != want {
						t.Fatalf("%s: %+v", n, r.Sessions[n])
					}
				}
				if scenario == "move keeps shade" {
					r.Layout = GhosttyLayout{Windows: [][]string{{"one", "two", "three", "four", "five", "six"}}, GroupIDs: []string{r.Sessions["one"].GroupID}}
					assignWindowGroup(&r, "native-teal", []string{"two"})
					completeThemeRecord(&r)
					e := r.Sessions["two"]
					if e.Shade != "deep" || r.Groups[e.GroupID].Hue != "teal" {
						t.Fatalf("move: %+v", r)
					}
				}
			case "restore binding":
				fakeGhostty(t, true)
				if err := UpdateRecord(func(r *Record) bool {
					r.Groups = map[string]ThemeGroup{"saved": {Hue: "green"}}
					r.Sessions["a/one"] = Entry{ID: "stable", GroupID: "saved", Shade: "bright", Socket: "proj-a"}
					r.Layout = GhosttyLayout{Windows: [][]string{{"a/one"}}, GroupIDs: []string{"saved"}}
					return true
				}); err != nil {
					t.Fatal(err)
				}
				script := osascript
				osascript = func(text string, args ...string) (string, error) {
					if text == newWindowScript {
						f, err := os.OpenFile(filepath.Join(filepath.Dir(RecordPath()), "theme-sync.lock"), os.O_RDWR, 0600)
						if err != nil {
							t.Errorf("restore lacks shared lock: %v", err)
						} else {
							err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
							if err == nil {
								t.Error("restore did not hold lock across window creation")
								_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
							}
							_ = f.Close()
						}
					}
					return script(text, args...)
				}
				original := hasClient
				hasClient = func(socket, name string) bool {
					r, err := LoadRecord()
					if err != nil || r.NativeGroups["win-1"] != "saved" {
						t.Errorf("attachment preceded group binding: %+v %v", r, err)
					}
					return original(socket, name)
				}
				r, _ := LoadRecord()
				groups := PlanGhostty(SavedSessions(r, LiveState{}, Transcripts{}))
				for range 2 {
					if err := GhosttyOpen(groups, &bytes.Buffer{}); err != nil {
						t.Fatal(err)
					}
				}
				again, _ := LoadRecord()
				if len(again.Groups) != 1 || again.Sessions["a/one"] != r.Sessions["a/one"] {
					t.Fatalf("restore reallocated: %+v", again)
				}
			}
		})
	}
}

func TestThemeQueryReadOnlyAndHueChange(t *testing.T) {
	recordHome(t)
	if err := UpdateRecord(func(r *Record) bool {
		r.Sessions["a/one"] = Entry{ID: "stable", Conversation: "agent", GroupID: "group", Shade: "deep"}
		r.Groups = map[string]ThemeGroup{"group": {Hue: "blue"}}
		r.NativeGroups = map[string]string{"native": "group"}
		return true
	}); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(RecordPath())
	for _, id := range []string{"unknown", "agent"} {
		q, err := ThemeForAgent(id)
		if err != nil {
			t.Fatal(err)
		}
		if q.Assigned != (id == "agent") {
			t.Fatalf("query: %+v", q)
		}
	}
	after, _ := os.ReadFile(RecordPath())
	if !bytes.Equal(before, after) {
		t.Fatal("query mutated record")
	}
	if err := SetWindowHue("native", "rose"); err != nil {
		t.Fatal(err)
	}
	q, err := ThemeForAgent("agent")
	if err != nil {
		t.Fatal(err)
	}
	if q.Theme.SessionID != "stable" || q.Theme.GroupID != "group" || q.Theme.Assignment.Hue != "rose" || q.Theme.Assignment.Shade != "deep" || q.Theme.Header.Background != q.Theme.Palette.Background {
		b, _ := json.Marshal(q)
		t.Fatalf("theme: %s", b)
	}
}

// A renamed session's changed conversation must not turn it into a new entry.
// All tmux calls here are command seams, never a live server.
func TestSnapshotAdoptsAndFollowsStampedID(t *testing.T) {
	recordHome(t)
	t.Setenv("PI_CODING_AGENT_SESSION_DIR", t.TempDir())
	t.Setenv("CLAUDE_CONFIG_DIR", t.TempDir())
	original := execCommand
	t.Cleanup(func() { execCommand = original })
	name, conversation, id := "a/old", "old", ""
	writes := 0
	execCommand = func(command string, args ...string) *exec.Cmd {
		var output string
		if command != "tmux" || len(args) < 3 || args[0] != "-L" || args[1] != "proj-a" {
			t.Fatalf("unexpected command %s %q", command, args)
		}
		switch args[2] {
		case "list-panes":
			output = scanLine(name, "0", "0", "40", "%1", "pi", "/work", "0", "0", conversation, "pi", "")
		case "list-sessions":
			output = name + "\x1f" + id
		case "set-option":
			if len(args) != 7 || args[3] != "-t" || args[4] != name || args[5] != "@proj_session_id" || args[6] == "" {
				t.Fatalf("invalid identity stamp %q", args)
			}
			writes++
			id = args[6]
		default:
			t.Fatalf("unexpected tmux call %q", args)
		}
		return exec.Command("printf", "%s", output)
	}
	if err := Snapshot("proj-a"); err != nil {
		t.Fatal(err)
	}
	r, err := LoadRecord()
	if err != nil {
		t.Fatal(err)
	}
	if id == "" || writes != 1 || r.Sessions[name].ID != id {
		t.Fatalf("adoption not stamped: id=%q writes=%d record=%+v", id, writes, r)
	}
	stable := id
	name, conversation = "a/new", "changed"
	for range 2 {
		if err := Snapshot("proj-a"); err != nil {
			t.Fatal(err)
		}
	}
	r, err = LoadRecord()
	if err != nil {
		t.Fatal(err)
	}
	e := r.Sessions[name]
	if len(r.Sessions) != 1 || e.ID != stable || e.Conversation != "changed" || e.Dir != "/work" || e.Agent != "pi" || e.GroupID != "" || writes != 1 {
		t.Fatalf("snapshot lost stable detached identity: %+v writes=%d", r, writes)
	}
}
