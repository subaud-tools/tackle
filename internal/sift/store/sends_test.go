package store

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/subaud-tools/tackle/internal/sift/row"
	"github.com/subaud-tools/tools-common/channelmcp"
	"github.com/subaud-tools/tools-common/msgstate"
	"github.com/subaud-tools/tools-common/sqlitedb"
)

// sendOne records a round, decides one row and presses Send with owner:
// one send.
func sendOne(t *testing.T, s *Store, owner string) Send {
	t.Helper()
	id, rows := round(t, s)
	if err := s.Decide(ctx, id, rows[0].ID, row.Decision{Action: "accept"}); err != nil {
		t.Fatal(err)
	}
	sd, err := s.Send(ctx, id, owner, nil)
	if err != nil || sd.ID == 0 {
		t.Fatalf("send %+v %v", sd, err)
	}
	return sd
}

// claimed sends one send and has session claim it, with a and b present.
func claimed(t *testing.T, s *Store, session string) Send {
	t.Helper()
	sendOne(t, s, "")
	sd, ok, err := s.ClaimSend(ctx, session, []string{"a", "b"})
	if err != nil || !ok {
		t.Fatalf("claim %+v %v %v", sd, ok, err)
	}
	return sd
}

func sendByID(t *testing.T, s *Store, id int64) Send {
	t.Helper()
	sd, err := scanSend(s.db.QueryRowContext(ctx, `SELECT `+sendCols+` FROM sends WHERE id = ?`, id))
	if err != nil {
		t.Fatal(err)
	}
	return sd
}

func report(t *testing.T, s *Store, session string, event channelmcp.ReportEvent, token string) {
	t.Helper()
	if err := s.ReportSends(ctx, session, event, token); err != nil {
		t.Fatalf("report %s: %v", event, err)
	}
}

func reply(t *testing.T, s *Store, session string, ids []int64, state msgstate.State, text string) {
	t.Helper()
	if err := s.ReplySends(ctx, session, ids, state, text); err != nil {
		t.Fatalf("reply %s: %v", state, err)
	}
}

// A database at schema 6 with one delivered and one undelivered send opens
// at 7: the delivered one is with the agent, the other still sent.
func TestMigrateGivesEachSendItsState(t *testing.T) {
	p := filepath.Join(t.TempDir(), "sift.db")
	v6, err := sqlitedb.Open(ctx, p, sqlitedb.Options{Migrations: Migrations[:6]})
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range []string{
		`INSERT INTO sends(round_id, created_at, counts, notes, owner, delivered_to, delivered_at) VALUES (1, 1, '{}', '[]', '', 'sess-1', 2)`,
		`INSERT INTO sends(round_id, created_at, counts, notes, owner) VALUES (1, 3, '{}', '[]', '')`,
	} {
		if _, err := v6.ExecContext(ctx, q); err != nil {
			t.Fatal(err)
		}
	}
	_ = v6.Close()
	s, err := Open(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()
	if v, _ := s.db.Version(ctx); v != 7 {
		t.Fatalf("version %d", v)
	}
	if a := sendByID(t, s, 1); a.State != msgstate.WithAgent || a.ClaimedBy != "sess-1" {
		t.Errorf("delivered send %+v", a)
	}
	if b := sendByID(t, s, 2); b.State != msgstate.Sent || b.ClaimedBy != "" {
		t.Errorf("undelivered send %+v", b)
	}
}

// A claim records who took the send, with a fresh report token, and leaves
// it sent: only the harness's landed report says it reached the agent. The
// round's revision moves, so the page stops saying no session has it.
func TestClaimLeavesTheSendSentWithAToken(t *testing.T) {
	s, _ := open(t)
	sendOne(t, s, "")
	_, rev0, _ := s.Latest(ctx)
	sd, ok, err := s.ClaimSend(ctx, "a", []string{"a"})
	if err != nil || !ok {
		t.Fatalf("%v %v", ok, err)
	}
	if sd.State != msgstate.Sent || sd.ClaimedBy != "a" || sd.ClaimedAt.IsZero() {
		t.Errorf("claimed %+v", sd)
	}
	if !regexp.MustCompile(`^[0-9a-f]{16}$`).MatchString(sd.Claim) {
		t.Errorf("token %q", sd.Claim)
	}
	if got := sendByID(t, s, sd.ID); got.Claim != sd.Claim || got.State != msgstate.Sent {
		t.Errorf("stored %+v", got)
	}
	if _, rev, _ := s.Latest(ctx); rev <= rev0 {
		t.Errorf("rev %d not past %d", rev, rev0)
	}
}

// A session that claimed a send and died before it landed does not strand
// it: once the claimer is gone, another session takes it, with a new token.
func TestASentSendWhoseClaimerIsGoneIsClaimedAgain(t *testing.T) {
	s, _ := open(t)
	first := claimed(t, s, "a")
	if _, ok, _ := s.ClaimSend(ctx, "b", []string{"a", "b"}); ok {
		t.Fatal("b took a send its present claimer holds")
	}
	sd, ok, err := s.ClaimSend(ctx, "b", []string{"b"})
	if err != nil || !ok || sd.ID != first.ID || sd.ClaimedBy != "b" || sd.Claim == first.Claim {
		t.Fatalf("reclaim %+v %v %v", sd, ok, err)
	}
}

// A send that reached the agent is never handed out again, even when its
// claimer is gone.
func TestAWithAgentSendIsNeverClaimedAgain(t *testing.T) {
	s, _ := open(t)
	sd := claimed(t, s, "a")
	report(t, s, "a", channelmcp.ReportLanded, sd.Claim)
	if got, ok, _ := s.ClaimSend(ctx, "b", []string{"b"}); ok {
		t.Fatalf("claimed a send with the agent: %+v", got)
	}
}

// landed moves the token's send to with the agent and clears held, and the
// round's revision moves so the page redraws.
func TestLandedMovesTheSendToWithAgent(t *testing.T) {
	s, _ := open(t)
	sd := claimed(t, s, "a")
	report(t, s, "a", channelmcp.ReportHeld, sd.Claim)
	_, rev0, _ := s.Latest(ctx)
	report(t, s, "a", channelmcp.ReportLanded, sd.Claim)
	if got := sendByID(t, s, sd.ID); got.State != msgstate.WithAgent || got.Held {
		t.Errorf("after landed %+v", got)
	}
	if _, rev, _ := s.Latest(ctx); rev <= rev0 {
		t.Errorf("rev %d not past %d", rev, rev0)
	}
}

// A late landed report leaves a send the agent already marked working as it
// is.
func TestLandedAfterWorkingLeavesItWorking(t *testing.T) {
	s, _ := open(t)
	sd := claimed(t, s, "a")
	reply(t, s, "a", []int64{sd.ID}, msgstate.Working, "applying")
	report(t, s, "a", channelmcp.ReportLanded, sd.Claim)
	if got := sendByID(t, s, sd.ID); got.State != msgstate.Working || got.Reply != "applying" {
		t.Errorf("%+v", got)
	}
}

// held marks a send that is still sent, and only that.
func TestHeldMarksOnlyASentSend(t *testing.T) {
	s, _ := open(t)
	sd := claimed(t, s, "a")
	report(t, s, "a", channelmcp.ReportHeld, sd.Claim)
	if got := sendByID(t, s, sd.ID); got.State != msgstate.Sent || !got.Held {
		t.Errorf("held %+v", got)
	}
	other := claimed(t, s, "a")
	report(t, s, "a", channelmcp.ReportLanded, other.Claim)
	report(t, s, "a", channelmcp.ReportHeld, other.Claim)
	if got := sendByID(t, s, other.ID); got.Held {
		t.Errorf("held a landed send %+v", got)
	}
}

// When the session's turn ends, its sends with the agent or working have
// no reply; another session's, and an answered one, stay as they are.
func TestTurnEndedLeavesTheSessionsOpenSendsWithNoReply(t *testing.T) {
	s, _ := open(t)
	with := claimed(t, s, "a")
	report(t, s, "a", channelmcp.ReportLanded, with.Claim)
	working := claimed(t, s, "a")
	reply(t, s, "a", []int64{working.ID}, msgstate.Working, "on it")
	done := claimed(t, s, "a")
	reply(t, s, "a", []int64{done.ID}, msgstate.Answered, "applied")
	theirs := claimed(t, s, "b")
	report(t, s, "b", channelmcp.ReportLanded, theirs.Claim)
	report(t, s, "a", channelmcp.ReportTurnEnded, "")
	for _, c := range []struct {
		id   int64
		want msgstate.State
	}{{with.ID, msgstate.NoReply}, {working.ID, msgstate.NoReply}, {done.ID, msgstate.Answered}, {theirs.ID, msgstate.WithAgent}} {
		if got := sendByID(t, s, c.id); got.State != c.want {
			t.Errorf("send %d: %s, want %s", c.id, got.State, c.want)
		}
	}
}

// A reply is all or nothing: one id it may not settle refuses the whole
// call, with the reason, and nothing changes.
func TestAReplyIsRefusedWhole(t *testing.T) {
	s, _ := open(t)
	mine := claimed(t, s, "a")
	declined := claimed(t, s, "a")
	reply(t, s, "a", []int64{declined.ID}, msgstate.Declined, "not now")
	theirs := claimed(t, s, "b")
	for name, c := range map[string]struct {
		ids   []int64
		state msgstate.State
		want  string
	}{
		"another session's": {[]int64{mine.ID, theirs.ID}, msgstate.Answered, "another session"},
		"after declined":    {[]int64{mine.ID, declined.ID}, msgstate.Answered, "declined"},
		"an unknown id":     {[]int64{mine.ID, 999}, msgstate.Answered, "no send 999"},
	} {
		err := s.ReplySends(ctx, "a", c.ids, c.state, "done")
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: err %v, want %q", name, err, c.want)
		}
		if got := sendByID(t, s, mine.ID); got.State != msgstate.Sent || got.Reply != "" {
			t.Errorf("%s: a refused reply moved %+v", name, got)
		}
	}
}

// A reply after the turn ended still settles the send.
func TestALateAnswerSettlesANoReplySend(t *testing.T) {
	s, _ := open(t)
	sd := claimed(t, s, "a")
	report(t, s, "a", channelmcp.ReportLanded, sd.Claim)
	report(t, s, "a", channelmcp.ReportTurnEnded, "")
	reply(t, s, "a", []int64{sd.ID}, msgstate.Answered, "applied after all")
	if got := sendByID(t, s, sd.ID); got.State != msgstate.Answered || got.Reply != "applied after all" {
		t.Errorf("%+v", got)
	}
}

// The page lists a round's newest sends first.
func TestRecentSendsAreNewestFirst(t *testing.T) {
	s, _ := open(t)
	id, rows := round(t, s)
	var sent []int64
	for _, r := range rows {
		if err := s.Decide(ctx, id, r.ID, row.Decision{Action: "accept"}); err != nil {
			t.Fatal(err)
		}
		sd, err := s.Send(ctx, id, "", nil)
		if err != nil {
			t.Fatal(err)
		}
		sent = append(sent, sd.ID)
	}
	got, err := s.RecentSends(ctx, id, 2)
	if err != nil || len(got) != 2 || got[0].ID != sent[2] || got[1].ID != sent[1] {
		t.Fatalf("%+v %v", got, err)
	}
}
