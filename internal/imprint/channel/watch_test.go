package channel

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Under Claude Code `imprint channel` spools each message for `imprint watch`
// (the plugin's mod) instead of pushing it, and removes the spool when the
// session ends; under pi it pushes as it always has and spools nothing.
func TestChannelSpoolsOnlyUnderClaudeCode(t *testing.T) {
	cases := []struct {
		name       string
		claude, pi string
		spools     bool
	}{
		{"pi", "", "s-1", false},
		{"claude code", "s-1", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("CLAUDE_CODE_SESSION_ID", tc.claude)
			t.Setenv("AGENT_SESSION_ID", tc.pi)
			t.Setenv("AGENT_SESSION_CHILD", "")
			base, _, scripts := newChannel(t)
			if err := os.WriteFile(filepath.Join(scripts, "whoami"), []byte("print demo-user\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			dir := t.TempDir()
			spool, _ := SpoolFile(dir, "s-1")
			c := drive(t, ForSession(base.Ops, "test", dir))
			if text, isErr := c.tool("imprint_login", nil); isErr {
				t.Fatalf("login: %q", text)
			}

			if !tc.spools {
				if msg, meta := c.note(); msg != "logged in as demo-user" || meta["op"] != "imprint_login" {
					t.Fatalf("%q %v", msg, meta)
				}
				if _, err := os.Stat(spool); !os.IsNotExist(err) {
					t.Fatalf("the channel spooled under pi: %v", err)
				}
				return
			}
			deadline := time.Now().Add(5 * time.Second)
			var got string
			for time.Now().Before(deadline) {
				b, _ := os.ReadFile(spool)
				if got = string(b); strings.HasSuffix(got, "\n") {
					break
				}
				time.Sleep(20 * time.Millisecond)
			}
			if got != `{"content":"logged in as demo-user","meta":{"op":"imprint_login"}}`+"\n" {
				t.Fatalf("spool = %q", got)
			}
			select {
			case p := <-c.notes:
				t.Fatalf("the channel pushed under Claude Code; only watch may deliver: %v", p)
			case <-time.After(300 * time.Millisecond):
			}
			if s, _ := c.tool("imprint_status", nil); !strings.Contains(s, "imprint watch") {
				t.Fatalf("status does not say messages go through imprint watch:\n%s", s)
			}
			_ = c.closeIn()
			select {
			case <-c.done:
			case <-time.After(5 * time.Second):
				t.Fatal("the channel did not end with its session")
			}
			if _, err := os.Stat(spool); !os.IsNotExist(err) {
				t.Fatalf("the spool outlived the channel: %v", err)
			}
		})
	}
}

// Claude Code ends a channel with SIGTERM while stdin may still be open; the
// channel must return then, so its spool is removed rather than left behind.
func TestChannelSignalledRemovesItsSpool(t *testing.T) {
	t.Setenv("CLAUDE_CODE_SESSION_ID", "s-1")
	t.Setenv("AGENT_SESSION_ID", "")
	t.Setenv("AGENT_SESSION_CHILD", "")
	base, _, _ := newChannel(t)
	dir := t.TempDir()
	spool, _ := SpoolFile(dir, "s-1")
	c := drive(t, ForSession(base.Ops, "test", dir))
	c.tool("imprint_status", nil)
	if _, err := os.Stat(spool); err != nil {
		t.Fatalf("no spool: %v", err)
	}
	c.cancel()
	select {
	case <-c.done:
	case <-time.After(5 * time.Second):
		t.Fatal("the signalled channel kept running")
	}
	if _, err := os.Stat(spool); !os.IsNotExist(err) {
		t.Fatalf("the spool outlived the signalled channel: %v", err)
	}
}
