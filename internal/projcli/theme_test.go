package projcli

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/subaud-tools/tackle/internal/proj"
)

// A regressed hidden route must not escape the seam and reach the desktop.
// Keep subprocess tripwires as well as the normal in-process sync stub.
func isolateThemeDesktop(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("PROJ_HOME", home)
	bin := filepath.Join(home, "bin")
	if err := os.MkdirAll(bin, 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"tmux", "osascript"} {
		if err := os.WriteFile(filepath.Join(bin, name), []byte("#!/bin/sh\necho 'test refused unmocked desktop access' >&2\nexit 99\n"), 0700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func TestThemeDispatch(t *testing.T) {
	old := reconcileThemes
	reconcileThemes = func(bool) error { return nil }
	t.Cleanup(func() { reconcileThemes = old })
	for _, tt := range []struct {
		name   string
		args   []string
		record string
		want   int
		output string
	}{
		{"unknown", []string{"--session", "missing", "--json"}, "", 0, `"assigned":false`},
		{"unassigned", []string{"--session", "agent", "--json"}, `{"version":2,"sessions":{"a/w":{"id":"stable","conversation":"agent"}}}`, 0, `"assigned":false`},
		{"resolved", []string{"--session", "agent", "--json"}, `{"version":2,"sessions":{"a/w":{"id":"stable","conversation":"agent","group_id":"group","shade":"deep"}},"groups":{"group":{"hue":"blue"}}}`, 0, `"session_id":"stable"`},
		{"set hue", []string{"--window", "native", "--hue", "green", "--json"}, `{"version":2,"groups":{"group":{"hue":"blue"}},"native_groups":{"native":"group"}}`, 0, `"hue":"green"`},
		{"both", []string{"--session", "a", "--window", "b", "--hue", "blue", "--json"}, "", 2, ""},
		{"invalid hue", []string{"--window", "b", "--hue", "invalid", "--json"}, "", 2, ""},
		{"missing hue", []string{"--window", "b", "--json"}, "", 2, ""},
		{"extra", []string{"--session", "a", "--json", "extra"}, "", 2, ""},
		{"missing json", []string{"--session", "a"}, "", 2, ""},
		{"unknown window", []string{"--window", "b", "--hue", "blue", "--json"}, "", 1, ""},
		{"io error", []string{"--session", "a", "--json"}, "broken json", 1, ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			isolateThemeDesktop(t)
			if tt.record != "" {
				if err := os.MkdirAll(filepath.Dir(proj.RecordPath()), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(proj.RecordPath(), []byte(tt.record), 0600); err != nil {
					t.Fatal(err)
				}
			}
			before, _ := os.ReadFile(proj.RecordPath())
			var out, errw bytes.Buffer
			code := Dispatch(append([]string{"theme"}, tt.args...), &out, &errw)
			var compact bytes.Buffer
			_ = json.Compact(&compact, out.Bytes())
			if code != tt.want || (tt.output != "" && !strings.Contains(compact.String(), tt.output)) {
				t.Fatalf("exit=%d out=%s err=%s", code, &out, &errw)
			}
			if tt.args[0] == "--session" {
				after, _ := os.ReadFile(proj.RecordPath())
				if !bytes.Equal(before, after) {
					t.Fatal("query wrote record")
				}
			}
		})
	}
}

func TestThemeSyncHiddenDispatch(t *testing.T) {
	old := reconcileThemes
	t.Cleanup(func() { reconcileThemes = old })
	for _, tt := range []struct {
		name  string
		args  []string
		fail  bool
		want  int
		force bool
	}{
		{"silent", []string{"__theme-sync"}, false, 0, false},
		{"forced", []string{"__theme-sync", "--force"}, false, 0, true},
		{"observation error", []string{"__theme-sync"}, true, 1, false},
		{"unknown flag", []string{"__theme-sync", "--bad"}, false, 2, false},
		{"save layout error", []string{"__save-layout"}, true, 1, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			isolateThemeDesktop(t)
			var calls []bool
			reconcileThemes = func(force bool) error {
				calls = append(calls, force)
				if tt.fail {
					return errors.New("cannot observe Ghostty")
				}
				return nil
			}
			var out, errw bytes.Buffer
			code := Dispatch(tt.args, &out, &errw)
			if code != tt.want || out.Len() != 0 {
				t.Fatalf("exit=%d out=%s err=%s", code, &out, &errw)
			}
			if tt.want == 2 {
				if len(calls) != 0 {
					t.Fatal("invalid flags ran sync")
				}
				return
			}
			if !reflect.DeepEqual(calls, []bool{tt.force}) {
				t.Fatalf("sync calls %v", calls)
			}
			if tt.fail && !strings.Contains(errw.String(), "cannot observe Ghostty") {
				t.Fatalf("missing diagnostic: %s", &errw)
			}
		})
	}
}

func TestHueMutationForcedReconciliationOrder(t *testing.T) {
	isolateThemeDesktop(t)
	old := reconcileThemes
	t.Cleanup(func() { reconcileThemes = old })
	calls := 0
	reconcileThemes = func(force bool) error {
		if !force {
			t.Fatal("hue reconciliation not forced")
		}
		calls++
		if calls == 1 {
			return proj.UpdateRecord(func(r *proj.Record) bool {
				r.Groups = map[string]proj.ThemeGroup{"group": {Hue: "blue"}}
				r.NativeGroups = map[string]string{"new-native": "group"}
				return true
			})
		}
		r, err := proj.LoadRecord()
		if err != nil {
			return err
		}
		if r.Groups["group"].Hue != "green" {
			t.Fatal("post reconcile ran before mutation")
		}
		return nil
	}
	var out, errw bytes.Buffer
	if code := Dispatch([]string{"theme", "--window", "new-native", "--hue", "green", "--json"}, &out, &errw); code != 0 {
		t.Fatalf("exit=%d err=%s", code, &errw)
	}
	if calls != 2 {
		t.Fatalf("forced reconciliation calls=%d", calls)
	}
}
