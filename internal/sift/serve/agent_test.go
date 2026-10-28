package serve

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/internal/sift/store"
	"github.com/subaud-tools/tools-common/msgstate"
)

func (f *fixture) presence(session, label string) {
	f.t.Helper()
	if w := f.do("POST", "/api/agent/presence", fmt.Sprintf(`{"session":%q,"harness":"pi","label":%q}`, session, label)); w.Code != 204 {
		f.t.Fatalf("presence %d %s", w.Code, w.Body)
	}
}

func (f *fixture) wait(session string) *httptest.ResponseRecorder {
	return f.do("GET", "/api/agent/wait?session="+session+"&timeout=1", "")
}

func (f *fixture) sendNow() {
	f.t.Helper()
	f.decide("r-neg", "reject", `,"note":"keep it"`)
	if w := f.do("POST", "/api/send", fmt.Sprintf(`{"round":%d}`, f.round)); w.Code != 200 {
		f.t.Fatalf("send %d", w.Code)
	}
}

func TestReviewMakesTheCallerTheOwner(t *testing.T) {
	f := newFixture(t)
	f.presence("a", "pi · a")
	w := f.do("POST", "/api/agent/review", `{"session":"a"}`)
	var got struct {
		Round int64 `json:"round"`
		Open  int   `json:"open"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if w.Code != 200 || got.Round != f.round || got.Open != 3 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	if f.review().Round.Owner != "pi · a" {
		t.Fatalf("owner %q", f.review().Round.Owner)
	}
}

func TestSendGoesToThePresentOwner(t *testing.T) {
	f := newFixture(t)
	f.srv.waitUnit = 20 * time.Millisecond
	f.presence("owner", "pi · owner")
	f.presence("other", "pi · other")
	f.do("POST", "/api/agent/review", `{"session":"owner"}`)
	f.sendNow()
	if w := f.wait("other"); w.Code != 204 {
		t.Fatalf("another session took the owner's send: %d %s", w.Code, w.Body)
	}
	w := f.wait("owner")
	if w.Code != 200 || !strings.Contains(w.Body.String(), "keep it") {
		t.Fatalf("owner: %d %s", w.Code, w.Body)
	}
}

func TestOwnerGoneAnySessionClaims(t *testing.T) {
	f := newFixture(t)
	f.srv.waitUnit = 20 * time.Millisecond
	now := time.Now()
	f.srv.now = func() time.Time { return now }
	f.presence("owner", "")
	f.do("POST", "/api/agent/review", `{"session":"owner"}`)
	f.sendNow()
	now = now.Add(presenceTTL + time.Second)
	if w := f.wait("next"); w.Code != 200 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
}

// Two sessions wait, one send: exactly one of them gets it. The waits
// outlast the send however slow the machine (a 50 ms wait lost the race to
// a busy CI runner's send, and both timed out); once one has it, the other
// is still waiting with nothing, and is ended.
func TestTwoWaitersOneSend(t *testing.T) {
	f := newFixture(t)
	f.srv.waitUnit = 500 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	type got struct {
		i int
		w *httptest.ResponseRecorder
	}
	done := make(chan got, 2)
	for i := range 2 {
		go func() {
			req := httptest.NewRequestWithContext(ctx, "GET", fmt.Sprintf("/api/agent/wait?session=s%d&timeout=60", i), nil)
			w := httptest.NewRecorder()
			f.h.ServeHTTP(w, req)
			done <- got{i, w}
		}()
	}
	for !f.srv.present("s0") || !f.srv.present("s1") {
		time.Sleep(time.Millisecond)
	}
	f.sendNow()
	first := <-done
	if first.w.Code != 200 || !strings.Contains(first.w.Body.String(), "keep it") {
		t.Fatalf("s%d: %d %s", first.i, first.w.Code, first.w.Body)
	}
	select {
	case other := <-done:
		t.Fatalf("s%d ended too: %d %s", other.i, other.w.Code, other.w.Body)
	case <-time.After(100 * time.Millisecond):
	}
	cancel()
	if other := <-done; other.w.Body.Len() != 0 {
		t.Fatalf("s%d got a send too: %s", other.i, other.w.Body)
	}
	if left, err := f.st.Unclaimed(context.Background()); err != nil || len(left) != 0 {
		t.Fatalf("unclaimed %v %v", left, err)
	}
}

func TestStatus(t *testing.T) {
	f := newFixture(t)
	f.decide("r-neg", "accept", "")
	w := f.do("GET", "/api/agent/status", "")
	var got map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got["open"] != float64(2) || got["decided"] != float64(1) || got["state"] != "ready" || got["round"] != float64(f.round) {
		t.Fatalf("%s", w.Body)
	}
}

func TestSendTextFormat(t *testing.T) {
	got := SendText(store.Send{Round: 4, Counts: store.Counts{Accept: 2, Edit: 1, Reject: 3},
		Notes: []store.Note{{Row: "/w/a/CLAUDE.md", Note: "shorter"}}})
	want := "The user sent their decisions for sift round 4: 2 accepted, 1 edited, 3 rejected.\n" +
		"Notes:\n- /w/a/CLAUDE.md: shorter\n" +
		"Next: run sift_apply (or `sift apply`): it writes each accepted or edited file whole, on a branch per repo, and lists what it leaves to you. Then run `sift reconcile` and review each branch before it merges."
	if got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
}

// Each pending send goes to its own owner, as recorded when it was sent,
// not to the latest round's owner: two rounds, two owners present, two
// pending sends, and each owner gets only its own.
func TestEachSendGoesToItsOwnOwner(t *testing.T) {
	f := newFixture(t)
	f.srv.waitUnit = 20 * time.Millisecond
	f.presence("a", "pi · a")
	f.presence("b", "pi · b")
	f.do("POST", "/api/agent/review", `{"session":"a"}`)
	f.sendNow()
	first := f.round
	f.record()
	f.do("POST", "/api/agent/review", `{"session":"b"}`)
	f.sendNow()
	var got struct {
		Round int64 `json:"round"`
	}
	w := f.wait("b")
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if w.Code != 200 || got.Round != f.round {
		t.Fatalf("b got round %d (%d), want its own %d", got.Round, w.Code, f.round)
	}
	if w := f.wait("b"); w.Code != 204 {
		t.Fatalf("b took a's send too: %d %s", w.Code, w.Body)
	}
	w = f.wait("a")
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if w.Code != 200 || got.Round != first {
		t.Fatalf("a got round %d (%d), want %d", got.Round, w.Code, first)
	}
}

// The harness's reports and the agent's replies move a send through the
// HTTP endpoints, and the review shows each move: a claim leaves it sent
// with a token, held marks it, landed puts it with the agent, a working
// reply gives its line, another session's reply is refused with the reason,
// and an answer settles it, which the turn's end leaves as it is.
func TestReportsAndRepliesMoveASend(t *testing.T) {
	f := newFixture(t)
	f.srv.waitUnit = 20 * time.Millisecond
	f.presence("a", "pi · a")
	f.presence("b", "pi · b")
	f.do("POST", "/api/agent/review", `{"session":"a"}`)
	f.sendNow()
	w := f.wait("a")
	var got struct {
		ID     int64  `json:"id"`
		Report string `json:"report"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil || w.Code != 200 || len(got.Report) != 16 {
		t.Fatalf("wait %d %s", w.Code, w.Body)
	}
	last := func() recentSend {
		t.Helper()
		rs := f.review().RecentSends
		if len(rs) != 1 || rs[0].ID != got.ID {
			t.Fatalf("recent sends %+v", rs)
		}
		return rs[0]
	}
	if s := last(); s.State != msgstate.Sent || !s.Claimed || s.Held {
		t.Fatalf("claimed %+v", s)
	}
	post := func(path, body string, code int) {
		t.Helper()
		if w := f.do("POST", path, body); w.Code != code {
			t.Fatalf("%s %s: %d %s", path, body, w.Code, w.Body)
		}
	}
	post("/api/agent/report", fmt.Sprintf(`{"session":"a","event":"held","token":%q}`, got.Report), 204)
	if s := last(); s.State != msgstate.Sent || !s.Held {
		t.Fatalf("held %+v", s)
	}
	post("/api/agent/report", fmt.Sprintf(`{"session":"a","event":"landed","token":%q}`, got.Report), 204)
	if s := last(); s.State != msgstate.WithAgent || s.Held {
		t.Fatalf("landed %+v", s)
	}
	post("/api/agent/reply", fmt.Sprintf(`{"session":"a","ids":[%d],"state":"working","text":"applying"}`, got.ID), 204)
	if s := last(); s.State != msgstate.Working || s.Reply != "applying" {
		t.Fatalf("working %+v", s)
	}
	if w := f.do("POST", "/api/agent/reply", fmt.Sprintf(`{"session":"b","ids":[%d],"state":"answered","text":"x"}`, got.ID)); w.Code != 409 || !strings.Contains(w.Body.String(), "another session") {
		t.Fatalf("b's reply %d %s", w.Code, w.Body)
	}
	post("/api/agent/reply", fmt.Sprintf(`{"session":"a","ids":[%d],"state":"answered","text":"applied: one PR"}`, got.ID), 204)
	post("/api/agent/report", `{"session":"a","event":"turn_ended"}`, 204)
	if s := last(); s.State != msgstate.Answered || s.Reply != "applied: one PR" {
		t.Fatalf("answered %+v", s)
	}
}

// Status lists the sends not yet settled.
func TestStatusListsTheUnsettledSends(t *testing.T) {
	f := newFixture(t)
	f.sendNow()
	w := f.do("GET", "/api/agent/status", "")
	var got struct {
		Unsettled []struct {
			ID    int64          `json:"id"`
			State msgstate.State `json:"state"`
		} `json:"unsettled"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil || len(got.Unsettled) != 1 || got.Unsettled[0].State != msgstate.Sent {
		t.Fatalf("%s", w.Body)
	}
}
