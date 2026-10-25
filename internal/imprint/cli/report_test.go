package cli

import "testing"

// imprint report, which the regenerated Claude Code plugin runs, exits 0 in a
// session: imprint sets no report token, so there is nothing to deliver.
func TestReport(t *testing.T) {
	t.Setenv("CLAUDE_CODE_SESSION_ID", "s1")
	t.Setenv("AGENT_SESSION_ID", "")
	t.Setenv("AGENT_SESSION_CHILD", "")
	for _, args := range [][]string{{"landed", "tok"}, {"turn-ended"}} {
		if code, out, errw := run(t, append([]string{"report"}, args...)...); code != 0 {
			t.Errorf("report %v: code %d: %s %s", args, code, out, errw)
		}
	}
}
