package proj

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/subaud-tools/tackle/internal/proj/pages"
)

// fakeOsascript swaps the osascript seam for one returning out/err, and
// records every script it is handed.
func fakeOsascript(t *testing.T, out string, err error) *[]string {
	t.Helper()
	var scripts []string
	orig := osascript
	t.Cleanup(func() { osascript = orig })
	osascript = func(script string, args ...string) (string, error) {
		scripts = append(scripts, script)
		return out, err
	}
	return &scripts
}

func TestMergeLayout(t *testing.T) {
	cases := []struct {
		name          string
		old, cur, out [][]string
	}{
		{"first save", nil, [][]string{{"a"}, {"b"}}, [][]string{{"a"}, {"b"}}},
		{"ghostty closed keeps everything",
			[][]string{{"a", "b"}, {"c"}}, nil, [][]string{{"a", "b"}, {"c"}}},
		// The reboot case: one session attached by hand before the restore.
		{"one early attach cannot shrink the layout",
			[][]string{{"a", "b"}, {"c", "d"}}, [][]string{{"c"}},
			[][]string{{"a", "b"}, {"c", "d"}}},
		{"tab moved between windows",
			[][]string{{"a", "b"}, {"c"}}, [][]string{{"a"}, {"c", "b"}},
			[][]string{{"a"}, {"c", "b"}}},
		{"new window appended",
			[][]string{{"a"}}, [][]string{{"a"}, {"z"}}, [][]string{{"a"}, {"z"}}},
		{"untabbed joins the window holding most of its old window",
			[][]string{{"a", "b", "c", "x"}}, [][]string{{"a"}, {"b", "c"}},
			[][]string{{"b", "c", "x"}, {"a"}}},
		{"empty saved window dropped",
			[][]string{{}, {"a"}}, [][]string{{"a"}}, [][]string{{"a"}}},
	}
	for _, c := range cases {
		got := mergeLayout(c.old, c.cur)
		if !slices.EqualFunc(got, c.out, slices.Equal[[]string]) {
			t.Errorf("%s: mergeLayout = %q want %q", c.name, got, c.out)
		}
	}
}

func TestPagesCloseTerminalRunsTheCloseScriptForThatTTY(t *testing.T) {
	var script string
	var args []string
	orig := osascript
	t.Cleanup(func() { osascript = orig })
	failure := errors.New("Ghostty said no")
	osascript = func(s string, a ...string) (string, error) {
		script, args = s, a
		return "", failure
	}
	if err := pages.CloseTerminal("/dev/ttys042"); !errors.Is(err, failure) {
		t.Fatalf("err = %v, want the script's", err)
	}
	if !slices.Equal(args, []string{"/dev/ttys042"}) {
		t.Fatalf("args = %q", args)
	}
	if !strings.Contains(script, "(tty of t) is (item 1 of argv)") || !strings.Contains(script, "close t") {
		t.Fatalf("script does not close the matching terminal: %s", script)
	}
	if pingFor(script) == "" {
		t.Fatal("the close script gets no ping")
	}
}
