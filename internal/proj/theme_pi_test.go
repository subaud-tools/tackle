package proj

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Launch, restore, and enabled refresh fixtures must never inherit live paths.
func isolateSessionThemes(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("SHELL", "/bin/sh")
	t.Setenv("PROJ_HOME", t.TempDir())
	t.Setenv("PI_CODING_AGENT_DIR", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
}

func TestPiLaunchSelectsThemeOnlyWhenEnabled(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		for _, resume := range []bool{false, true} {
			t.Run(fmtLaunchCase(enabled, resume), func(t *testing.T) {
				calls := withFakeRun(t)
				stubBins(t, "pi")
				if enabled {
					writeConfig(t, "session_themes = true\n")
				}
				if err := UpdateRecord(func(r *Record) bool {
					r.Groups = map[string]ThemeGroup{"saved": {Hue: "green"}}
					r.Sessions["x/w"] = Entry{ID: "stable-pi", GroupID: "saved", Shade: "muted", Agent: "pi", Socket: "proj-x"}
					return true
				}); err != nil {
					t.Fatal(err)
				}
				oldRunner := runner
				runner = func(socket string, args ...string) (string, error) {
					if enabled && args[0] == "send-keys" {
						if _, err := os.Stat(filepath.Join(os.Getenv("PI_CODING_AGENT_DIR"), "themes", "proj-session-stable-pi.json")); err != nil {
							t.Errorf("theme not ready before launch: %v", err)
						}
					}
					return oldRunner(socket, args...)
				}
				var err error
				if resume {
					_, err = RestoreSession("proj-x", "x/w", "/tmp/dir", "pi", "id1")
				} else {
					_, err = EnsureSession("proj-x", "x/w", "/tmp/dir", "pi", "provider/model", nil)
				}
				if err != nil {
					t.Fatal(err)
				}
				want := "pi --model 'provider/model' --name 'x/w'"
				if resume {
					want = "pi --session 'id1'"
				}
				if enabled {
					want += " --use-theme 'proj-session-stable-pi'"
				}
				if got := sendKeysText(*calls); got != want {
					t.Fatalf("typed %q, want %q", got, want)
				}
				paths, err := filepath.Glob(filepath.Join(os.Getenv("PI_CODING_AGENT_DIR"), "themes", "*"))
				if err != nil {
					t.Fatal(err)
				}
				if !enabled && len(paths) != 0 {
					t.Fatalf("disabled launch created %v", paths)
				}
				if enabled {
					r, err := LoadRecord()
					if err != nil {
						t.Fatal(err)
					}
					if opt := findStableOption(*calls); opt != "stable-pi" || r.Sessions["x/w"].ID != opt {
						t.Fatalf("launch identity stamp=%q record=%+v", opt, r.Sessions["x/w"])
					}
					assertPiColor(t, "proj-session-stable-pi", "export", "pageBg", "#2d3b2f")
				}
			})
		}
	}
}

func fmtLaunchCase(enabled, resume bool) string {
	s := "disabled"
	if enabled {
		s = "enabled"
	}
	if resume {
		return s + "/resume"
	}
	return s + "/new"
}

func findStableOption(calls []recordedCall) string {
	for _, c := range calls {
		if len(c.args) == 5 && c.args[0] == "set-option" && c.args[3] == "@proj_session_id" {
			return c.args[4]
		}
	}
	return ""
}

func assertPiColor(t *testing.T, name, section, role, want string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(os.Getenv("PI_CODING_AGENT_DIR"), "themes", name+".json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	var values map[string]string
	if err := json.Unmarshal(doc[section], &values); err != nil {
		t.Fatal(err)
	}
	if got := values[role]; got != want {
		t.Fatalf("%s.%s=%s, want %s", section, role, got, want)
	}
}

func TestDetachedPiThemeUsesAllocatedIdentityThenFirstAttachment(t *testing.T) {
	_, _, _, _ = themeFixture(t)
	external, oldNew := runner, newSessionHome
	t.Cleanup(func() { newSessionHome = oldNew })
	newSessionHome = func(string, string) error { return nil }
	var calls []recordedCall
	runner = func(socket string, args ...string) (string, error) {
		calls = append(calls, recordedCall{socket, args})
		if args[0] == "has-session" {
			return "", errors.New("missing")
		}
		return "", nil
	}
	stubBins(t, "pi")
	if _, err := EnsureSession("proj-test", "a", "/tmp", "pi", "", nil); err != nil {
		t.Fatal(err)
	}
	r, err := LoadRecord()
	if err != nil {
		t.Fatal(err)
	}
	e := r.Sessions["a"]
	if e.ID == "" || e.ID != findStableOption(calls) || e.GroupID != "" || e.Shade != "" || len(r.Groups) != 0 {
		t.Fatalf("detached identity invented an assignment: %+v / %+v", e, r.Groups)
	}
	name := "proj-session-" + e.ID
	assertPiColor(t, name, "export", "pageBg", "#14283f")
	runner = func(socket string, args ...string) (string, error) {
		out, err := external(socket, args...)
		if args[0] == "list-sessions" {
			out = strings.ReplaceAll(out, "stable-a", e.ID)
		}
		return out, err
	}
	if err := ReconcileThemes(true); err != nil {
		t.Fatal(err)
	}
	r, err = LoadRecord()
	if err != nil {
		t.Fatal(err)
	}
	attached := r.Sessions["a"]
	if attached.ID != e.ID || attached.GroupID == "" || attached.Shade != "midnight" {
		t.Fatalf("first attachment=%+v", attached)
	}
	assertPiColor(t, name, "export", "pageBg", "#102032")
}

func TestPreparePiThemeDirectoryAndDisabledNoWrite(t *testing.T) {
	for _, override := range []string{"", "~", "~/custom", "absolute"} {
		t.Run(override, func(t *testing.T) {
			isolateSessionThemes(t)
			home := t.TempDir()
			t.Setenv("HOME", home)
			dir := filepath.Join(home, ".pi", "agent")
			switch override {
			case "~":
				dir = home
			case "~/custom":
				dir = filepath.Join(home, "custom")
			case "absolute":
				dir = t.TempDir()
				override = dir
			}
			t.Setenv("PI_CODING_AGENT_DIR", override)
			e := Entry{ID: "stable-dir"}
			if name, err := PreparePiTheme(e); name != "" || err != nil {
				t.Fatalf("disabled prep=(%q,%v)", name, err)
			}
			if _, err := os.Stat(filepath.Join(dir, "themes")); !os.IsNotExist(err) {
				t.Fatalf("disabled prep created theme dir: %v", err)
			}
			cfg := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "proj")
			if err := os.MkdirAll(cfg, 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(cfg, "config.toml"), []byte("session_themes = true\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if name, err := PreparePiTheme(e); name != "proj-session-stable-dir" || err != nil {
				t.Fatalf("enabled prep=(%q,%v)", name, err)
			}
			b, err := os.ReadFile(filepath.Join(dir, "themes", "proj-session-stable-dir.json"))
			if err != nil {
				t.Fatal(err)
			}
			var doc struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(b, &doc); err != nil || doc.Name != "proj-session-stable-dir" {
				t.Fatalf("invalid theme %s, %v", b, err)
			}
		})
	}
}
