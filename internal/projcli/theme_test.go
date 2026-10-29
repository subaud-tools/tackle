package projcli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/subaud-tools/tackle/internal/proj"
)

func TestThemeDispatch(t *testing.T) {
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
			t.Setenv("PROJ_HOME", t.TempDir())
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
