package channel

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/internal/sift/serve"
	"github.com/subaud-tools/tools-common/channelmcp"
	"github.com/subaud-tools/tools-common/pageserve"
)

// The event for a Send carries its claim's report token, which the harness
// names when it reports the send held or landed; the agent settles the send
// with sift_reply, and a session that did not get the send is refused.
func TestASendsEventCarriesItsTokenAndSiftReplySettlesIt(t *testing.T) {
	e := newEnv(t)
	c := e.connectChannel(func() *Channel {
		cl := pageserve.NewClient(serve.Live(), nil)
		return pageserve.New(e.spec(), pageserve.Identity{Session: "s1", Harness: "pi", Label: "pi · s1", CWD: e.repo}, cl)
	}())
	sendARound(t, e, c)
	var ev event
	select {
	case ev = <-c.events:
	case <-time.After(10 * time.Second):
		t.Fatal("no event")
	}
	if !regexp.MustCompile(`^[0-9a-f]{16}$`).MatchString(ev.Meta[channelmcp.MetaReport]) {
		t.Fatalf("meta %v: no report token", ev.Meta)
	}
	id, _ := strconv.ParseInt(ev.Meta["send_id"], 10, 64)
	other := e.connect("s2")
	if text, isErr := other.tool("sift_reply", map[string]any{"ids": []int64{id}, "state": "answered", "text": "x"}); !isErr || !strings.Contains(text, "another session") {
		t.Fatalf("s2's reply: %v %s", isErr, text)
	}
	if text, isErr := c.tool("sift_reply", map[string]any{"ids": []int64{id}, "state": "answered", "text": "applied: one PR"}); isErr {
		t.Fatalf("reply: %s", text)
	}
	_, body := e.api("GET", "/api/review", "")
	if !strings.Contains(body, `"state":"answered"`) || !strings.Contains(body, `"reply":"applied: one PR"`) {
		t.Fatalf("review %s", body)
	}
}
