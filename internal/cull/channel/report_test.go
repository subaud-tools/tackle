package channel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/internal/cull/serve"
	"github.com/subaud-tools/tools-common/channelmcp"
	"github.com/subaud-tools/tools-common/msgstate"
)

// notify sends the channel a notification, as channels.tools does.
func (c *client) notify(method string, params any) {
	c.t.Helper()
	b, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "method": method, "params": params})
	if _, err := c.w.Write(append(b, '\n')); err != nil {
		c.t.Fatal(err)
	}
}

// pi's reports reach serve in the order sent, and a slow serve never holds
// up the channel's answers to the agent's tool calls.
func TestReportsReachServeInOrderWithoutBlockingToolCalls(t *testing.T) {
	release := make(chan struct{})
	posted := make(chan map[string]string, 4)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/agent/report":
			var b map[string]string
			_ = json.NewDecoder(r.Body).Decode(&b)
			<-release
			posted <- b
		case "/api/agent/wait":
			select {
			case <-r.Context().Done():
			case <-time.After(100 * time.Millisecond):
			}
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)
	t.Cleanup(func() {
		select {
		case <-release:
		default:
			close(release)
		}
	})
	cl := NewClient(nil)
	cl.Find = func() (serve.Advert, error) { return serve.Advert{Base: srv.URL, Token: "tok", PID: os.Getpid()}, nil }
	e := &env{t: t, logs: &syncBuf{}}
	c := e.connectChannel(New(Identity{Session: "s1", Harness: "pi", Label: "pi · s1", CWD: project(t)}, cl, "test"))
	c.notify("notifications/channels/report", map[string]string{"event": "landed", "report": "0123456789abcdef"})
	c.notify("notifications/channels/report", map[string]string{"event": "turn_ended"})
	done := make(chan struct{})
	go func() {
		c.call("tools/list", map[string]any{})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("a tools call waited on serve's report answer")
	}
	close(release)
	want := []map[string]string{
		{"session": "s1", "event": "landed", "token": "0123456789abcdef"},
		{"session": "s1", "event": "turn_ended", "token": ""},
	}
	for i, w := range want {
		select {
		case got := <-posted:
			if got["session"] != w["session"] || got["event"] != w["event"] || got["token"] != w["token"] {
				t.Fatalf("report %d posted %v, want %v", i+1, got, w)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("report %d never reached serve", i+1)
		}
	}
}

// The event for a Send carries its claim's report token, which the harness
// names when it reports the send held or landed; the agent settles the send
// with cull_reply, and a session that did not get the send is refused.
func TestASendsEventCarriesItsTokenAndCullReplySettlesIt(t *testing.T) {
	e := newEnv(t, true)
	root := project(t)
	e.jev.set(map[string]float64{`"test_name":"TestMaybe"`: 0.4})
	a := e.connect(ident{"s-1", "claude · a", root})
	var ck checkResult
	a.toolJSON("cull_check", map[string]any{}, &ck)
	a.toolJSON("cull_review", map[string]any{}, &struct{}{})
	answerAndSend(t, e, root)
	ev, ok := a.expectEvent(10 * time.Second)
	if !ok {
		t.Fatal("no event")
	}
	if !regexp.MustCompile(`^[0-9a-f]{16}$`).MatchString(ev.Meta[channelmcp.MetaReport]) {
		t.Fatalf("meta %v: no report token", ev.Meta)
	}
	id, _ := strconv.ParseInt(ev.Meta["send_id"], 10, 64)
	b := e.connect(ident{"s-2", "claude · b", root})
	if text, isErr := b.tool("cull_reply", map[string]any{"ids": []int64{id}, "state": "answered", "text": "x"}); !isErr || !strings.Contains(text, "another session") {
		t.Fatalf("s-2's reply: %v %s", isErr, text)
	}
	if text, isErr := a.tool("cull_reply", map[string]any{"ids": []int64{id}, "state": "answered", "text": "1 test cut"}); isErr {
		t.Fatalf("reply: %s", text)
	}
	p, _ := e.st.Project(context.Background(), root)
	sends, err := e.st.RecentSends(context.Background(), p.ID, 1)
	if err != nil || len(sends) != 1 || sends[0].State != msgstate.Answered || sends[0].Reply != "1 test cut" {
		t.Fatalf("%+v %v", sends, err)
	}
}
