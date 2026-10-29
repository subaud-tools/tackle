package proj

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/subaud-tools/tools-common/sessiontheme"
)

func TestTerminalSequencesOnlyColorControls(t *testing.T) {
	p := sessiontheme.Palette{Background: "#112233", Foreground: "#ddeeff", Cursor: "#aabbcc", CursorText: "#001122"}
	for i := range p.ANSI {
		p.ANSI[i] = "#445566"
	}
	got := string(TerminalSequences(p))
	if !strings.HasPrefix(got, "\x1b]10;#ddeeff\x1b\\\x1b]11;#112233\x1b\\\x1b]12;#aabbcc\x1b\\") {
		t.Fatalf("defaults/cursor: %q", got)
	}
	if strings.Count(got, "\x1b]4;") != 16 || !strings.Contains(got, "\x1b]4;15;#445566\x1b\\") {
		t.Fatalf("ANSI palette: %q", got)
	}
	for _, part := range strings.Split(strings.TrimSuffix(got, "\x1b\\"), "\x1b\\") {
		if !strings.HasPrefix(part, "\x1b]") || strings.Contains(part, "\n") {
			t.Fatalf("not a color control: %q", part)
		}
	}
}

func TestTerminalWriterRejectsNonTTYWithoutChangingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ordinary")
	if err := os.WriteFile(path, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := writeTerminalDevice(path, []byte("overwrite")); err == nil {
		t.Fatal("accepted ordinary file")
	}
	b, _ := os.ReadFile(path)
	if string(b) != "keep" {
		t.Fatal("changed ordinary file")
	}
}

func TestTmuxThemeWritesExplicitRolesAndWindowStyles(t *testing.T) {
	old := runner
	t.Cleanup(func() { runner = old })
	got := map[string]string{}
	runner = func(socket string, args ...string) (string, error) {
		if socket != "proj-test" {
			t.Fatalf("socket=%s", socket)
		}
		for i := 0; i < len(args); i += 6 {
			got[args[i+2]+"/"+args[i+3]] = args[i+4]
		}
		return "", nil
	}
	p := sessiontheme.Palette{Background: "#101010", Foreground: "#eeeeee", Panel: "#202020", SelectedBackground: "#303030", Muted: "#aaaaaa", Border: "#bbbbbb", Accent: "#cccccc", AccentText: "#010101", Success: "#00ff00", Warning: "#ffff00", Error: "#ff0000"}
	s := themeSession{socket: "proj-test", id: "$1", options: make([]string, len(themeSessionOptions))}
	windows := []themeWindow{{socket: "proj-test", session: "$1", id: "@1", options: make([]string, len(themeWindowOptions))}, {socket: "proj-test", session: "$2", id: "@2"}}
	if err := applyTmuxTheme(s, windows, "stable", p); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"$1/@proj_session_id": "stable", "$1/@proj_theme_bg": "#101010", "$1/@proj_theme_fg": "#eeeeee", "$1/@proj_theme_panel": "#202020", "$1/@proj_theme_selected_bg": "#303030", "$1/@proj_theme_muted": "#aaaaaa", "$1/@proj_theme_border": "#bbbbbb", "$1/@proj_theme_accent": "#cccccc", "$1/@proj_theme_accent_text": "#010101", "$1/@proj_theme_success": "#00ff00", "$1/@proj_theme_warning": "#ffff00", "$1/@proj_theme_error": "#ff0000", "$1/status-style": "bg=#101010,fg=#eeeeee",
		"@1/window-status-style": "bg=#202020,fg=#eeeeee", "@1/window-status-current-style": "bg=#303030,fg=#eeeeee", "@1/pane-border-style": "fg=#bbbbbb", "@1/pane-active-border-style": "fg=#cccccc",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tmux options=%v\nwant=%v", got, want)
	}
}
