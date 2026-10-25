package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/skills"
)

// sift channel refreshes the installed sift skill before it serves: a
// stale copy is rewritten with the binary's skill and logged to stderr, and
// a copy that already matches is left alone.
func TestChannelRefreshesTheSkill(t *testing.T) {
	home := siftEnv(t)
	for _, k := range []string{"CLAUDE_CODE_SESSION_ID", "AGENT_SESSION_ID", "AGENT_SESSION_CHILD"} {
		t.Setenv(k, "")
	}
	want, err := skills.FS.ReadFile("sift/SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(home, ".claude/skills/sift/SKILL.md")
	same := filepath.Join(home, ".pi/agent/skills/sift/SKILL.md")
	for p, b := range map[string]string{stale: "---\nname: sift\n---\nold\n", same: string(want)} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(b), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	past := time.Now().Add(-time.Hour).Truncate(time.Second)
	if err := os.Chtimes(same, past, past); err != nil {
		t.Fatal(err)
	}
	code, _, errw := run(t, "", "channel")
	if code != 0 {
		t.Fatalf("code %d: %s", code, errw)
	}
	if b, _ := os.ReadFile(stale); string(b) != string(want) {
		t.Errorf("the stale skill was not rewritten:\n%s", b)
	}
	if !strings.Contains(errw, stale) {
		t.Errorf("stderr does not name %s:\n%s", stale, errw)
	}
	if st, err := os.Stat(same); err != nil || !st.ModTime().Equal(past) || strings.Contains(errw, same) {
		t.Errorf("the matching skill was rewritten (%v):\n%s", err, errw)
	}
}
