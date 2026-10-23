package channel

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/internal/sift/serve"
	"github.com/subaud-tools/tackle/internal/sift/store"
	"github.com/subaud-tools/tools-common/channelmcp"
)

// sendARound records a round through c's tools, recommends every file,
// opens the review for c's session, decides each file and presses Send, as
// the page would.
func sendARound(t *testing.T, e *env, c *client) {
	t.Helper()
	var chk struct {
		Round int64 `json:"round"`
	}
	c.toolJSON("sift_check", map[string]any{}, &chk)
	for {
		var n struct {
			Done       bool
			File, Base string
			Findings   []struct{ Row string }
		}
		c.toolJSON("sift_next", map[string]any{}, &n)
		if n.Done {
			break
		}
		c.toolJSON("sift_propose", map[string]any{"recommendations": []map[string]any{{"file": n.File, "base": n.Base,
			"content": "# Rewritten\n", "summary": "Guidance.",
			"findings": []map[string]string{{"row": n.Findings[0].Row, "did": "fixed", "how": "guidance"}}}}}, &struct{}{})
	}
	c.toolJSON("sift_review", map[string]any{}, &struct{}{})
	items, _ := e.st.Files(context.Background(), chk.Round)
	shown := map[string]store.Seen{}
	for _, it := range items {
		body := fmt.Sprintf(`{"round":%d,"file":%q,"action":"reject","note":"yes","prints":{%q:%q}}`, chk.Round, it.Key, it.Key, it.Fingerprint)
		if code, out := e.api("PUT", "/api/files", body); code != 200 {
			t.Fatalf("decide %d %s", code, out)
		}
	}
	items, _ = e.st.Files(context.Background(), chk.Round)
	for _, it := range items {
		shown[it.Key] = store.Seen{Fingerprint: it.Fingerprint, Decision: it.Decision.ID}
	}
	files, _ := json.Marshal(shown)
	if code, body := e.api("POST", "/api/send", fmt.Sprintf(`{"round":%d,"files":%s,"rows":{}}`, chk.Round, files)); code != 200 {
		t.Fatalf("send %d %s", code, body)
	}
}

// `sift watch` is the channel's presence and wait loop with a JSON line on
// stdout for each Send, which is what the Claude Code plugin's mod reads.
// The user's Send for a review this session owns arrives as one line,
// {"content": ..., "meta": {...}}, with the same meta the channel sends.
func TestWatchPrintsASendAsOneLine(t *testing.T) {
	e := newEnv(t)
	id := Identity{Session: "s1", Harness: "claude", Label: "claude · s1", CWD: e.repo}

	// The tools, as `sift channel` serves them under Claude Code.
	cl := NewClient(nil)
	cl.Find = serve.Running
	tools := New(id, cl, "test")
	tools.ServeOnly = true
	c := e.connectChannel(tools)

	// The events, as `sift watch` takes them.
	r, w := io.Pipe()
	wcl := NewClient(nil)
	wcl.Find = serve.Running
	watch := New(id, wcl, "test")
	watch.Log = e.logs
	watch.Retry, watch.Presence, watch.Poll = 50*time.Millisecond, 200*time.Millisecond, time.Second
	watch.Push = channelmcp.NewLines(w)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { _ = watch.Watch(ctx); close(done) }()
	t.Cleanup(func() { cancel(); <-done })

	sendARound(t, e, c)

	lines := make(chan string, 1)
	go func() {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1<<20), 1<<24)
		if sc.Scan() {
			lines <- sc.Text()
		}
	}()
	var line string
	select {
	case line = <-lines:
	case <-time.After(10 * time.Second):
		t.Fatal("no line printed for the Send")
	}
	var ev event
	if err := json.Unmarshal([]byte(line), &ev); err != nil {
		t.Fatalf("not one JSON line: %q (%v)", line, err)
	}
	if ev.Meta["source"] != "sift" || ev.Meta["round"] == "" || ev.Meta["send_id"] == "" {
		t.Fatalf("meta = %v", ev.Meta)
	}
	if !strings.Contains(ev.Content, "rejected") {
		t.Fatalf("content = %q", ev.Content)
	}
	select {
	case got := <-c.events:
		t.Fatalf("the serve-only channel took the send too: %s", got.Content)
	case <-time.After(300 * time.Millisecond):
	}
}

// Under Claude Code the plugin's mod runs `sift watch`, so `sift channel` must
// not take the user's decisions too: it serves the tools and never waits.
// Under pi it waits and pushes as it always has.
func TestChannelWaitsOnlyOutsideClaudeCode(t *testing.T) {
	cases := []struct {
		name       string
		claude, pi string
		waits      bool
	}{
		{"pi", "", "s1", true},
		{"claude code", "s1", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := newEnv(t)
			t.Setenv("CLAUDE_CODE_SESSION_ID", tc.claude)
			t.Setenv("AGENT_SESSION_ID", tc.pi)
			t.Setenv("AGENT_SESSION_CHILD", "")
			t.Chdir(e.repo)

			cl := NewClient(nil)
			cl.Find = serve.Running
			ch := ForSession(cl, "test")
			ch.LookPath = func(string) (string, error) { return "", errors.New("no gh in tests") }
			c := e.connectChannel(ch)
			sendARound(t, e, c)

			if tc.waits {
				select {
				case <-c.events:
				case <-time.After(10 * time.Second):
					t.Fatal("the channel under pi never got the send")
				}
				return
			}
			select {
			case ev := <-c.events:
				t.Fatalf("the channel took a send under Claude Code; only watch may: %s", ev.Content)
			case <-time.After(1500 * time.Millisecond):
			}
			if s, _ := c.tool("sift_status", map[string]any{}); !strings.Contains(s, "sift watch") {
				t.Fatalf("status does not say decisions go through sift watch:\n%s", s)
			}
		})
	}
}
