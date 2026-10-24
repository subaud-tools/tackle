package channel

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/internal/cull/serve"
	"github.com/subaud-tools/tools-common/channelmcp"
)

// `cull watch` is the channel's presence and wait loop with a JSON line on
// stdout for each Send, which is what the Claude Code plugin's mod reads.
// Court's Send for a review this session owns arrives as one line,
// {"content": ..., "meta": {...}}, with the same meta the channel sends.
func TestWatchPrintsASendAsOneLine(t *testing.T) {
	e := newEnv(t, true)
	root := project(t)
	e.jev.set(map[string]float64{`"test_name":"TestMaybe"`: 0.4})
	id := Identity{Session: "s-1", Harness: "claude", Label: "claude · one", CWD: root}

	// The tools, as `cull channel` serves them under Claude Code.
	cl := NewClient(nil)
	cl.Find = serve.Running
	tools := New(id, cl, "test")
	tools.ServeOnly = true
	a := e.connectChannel(tools)
	var ck checkResult
	a.toolJSON("cull_check", map[string]any{}, &ck)
	a.toolJSON("cull_review", map[string]any{}, &struct{}{})

	// The events, as `cull watch` takes them.
	r, w := io.Pipe()
	wcl := NewClient(nil)
	wcl.Find = serve.Running
	watch := New(id, wcl, "test")
	watch.Log = e.logs
	watch.Retry, watch.Presence, watch.Poll = 50*time.Millisecond, 200*time.Millisecond, time.Second
	watch.Push = channelmcp.NewLinesFor(w, "cull")
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { _ = watch.Watch(ctx); close(done) }()
	t.Cleanup(func() { cancel(); <-done })

	answerAndSend(t, e, root)

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
	if ev.Meta["source"] != "cull" || ev.Meta["project"] != root || ev.Meta["send_id"] == "" {
		t.Fatalf("meta = %v", ev.Meta)
	}
	if !strings.Contains(ev.Content, "Court sent his answers") {
		t.Fatalf("content = %q", ev.Content)
	}
	if got, ok := a.expectEvent(300 * time.Millisecond); ok {
		t.Fatalf("the serve-only channel took the send too: %s", got.Content)
	}
}

// Under Claude Code the plugin's mod runs `cull watch`, so `cull channel` must
// not take Court's answers too: it serves the tools and never waits. Under pi
// it waits and pushes as it always has.
func TestChannelWaitsOnlyOutsideClaudeCode(t *testing.T) {
	cases := []struct {
		name       string
		claude, pi string
		waits      bool
	}{
		{"pi", "", "s-1", true},
		{"claude code", "s-1", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := newEnv(t, true)
			root := project(t)
			e.jev.set(map[string]float64{`"test_name":"TestMaybe"`: 0.4})
			t.Setenv("CLAUDE_CODE_SESSION_ID", tc.claude)
			t.Setenv("AGENT_SESSION_ID", tc.pi)
			t.Setenv("AGENT_SESSION_CHILD", "")
			t.Chdir(root)

			cl := NewClient(nil)
			cl.Find = serve.Running
			c := e.connectChannel(ForSession(cl, "test"))
			var ck checkResult
			c.toolJSON("cull_check", map[string]any{}, &ck)
			c.toolJSON("cull_review", map[string]any{}, &struct{}{})
			answerAndSend(t, e, root)

			if tc.waits {
				if _, ok := c.expectEvent(10 * time.Second); !ok {
					t.Fatal("the channel under pi never got the send")
				}
				return
			}
			if ev, ok := c.expectEvent(1500 * time.Millisecond); ok {
				t.Fatalf("the channel took a send under Claude Code; only watch may: %s", ev.Content)
			}
			if s, _ := c.tool("cull_status", map[string]any{}); !strings.Contains(s, "cull watch") {
				t.Fatalf("status does not say answers go through cull watch:\n%s", s)
			}
		})
	}
}
