package channel

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/subaud-tools/tools-common/channelmcp"
	"github.com/subaud-tools/tools-common/msgstate"
)

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
