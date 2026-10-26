package store

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/subaud-tools/tools-common/channelmcp"
	"github.com/subaud-tools/tools-common/msgstate"
	"github.com/subaud-tools/tools-common/sqlitedb"
)

// sender records a project with a run and returns a function that answers
// the next item and presses Send: one send each call.
func sender(t *testing.T, s *Store) (int64, func() Send) {
	t.Helper()
	p, err := s.Project(ctx, "/x")
	if err != nil {
		t.Fatal(err)
	}
	run, err := s.RecordRun(ctx, Run{ProjectID: p.ID}, items(12))
	if err != nil {
		t.Fatal(err)
	}
	i := 0
	return p.ID, func() Send {
		t.Helper()
		if err := s.SaveAnswers(ctx, p.ID, run, []Answer{answerOf(i, "cut", "")}); err != nil {
			t.Fatal(err)
		}
		i++
		sd, err := s.Send(ctx, p.ID, "")
		if err != nil || sd.ID == 0 {
			t.Fatalf("send %+v %v", sd, err)
		}
		return sd
	}
}

// claimed sends one send and has session claim it, with a and b present.
func claimed(t *testing.T, s *Store, project int64, send func() Send, session string) Send {
	t.Helper()
	send()
	sd, ok, err := s.ClaimSend(ctx, project, session, []string{"a", "b"})
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

func report(t *testing.T, s *Store, session string, event channelmcp.ReportEvent, token string) []int64 {
	t.Helper()
	ps, err := s.ReportSends(ctx, session, event, token)
	if err != nil {
		t.Fatalf("report %s: %v", event, err)
	}
	return ps
}

func reply(t *testing.T, s *Store, session string, ids []int64, state msgstate.State, text string) {
	t.Helper()
	if _, err := s.ReplySends(ctx, session, ids, state, text); err != nil {
		t.Fatalf("reply %s: %v", state, err)
	}
}

// A database at schema 2 with one delivered and one undelivered send opens
// at 3: the delivered one is with the agent, the other still sent.
func TestMigrateGivesEachSendItsState(t *testing.T) {
	p := filepath.Join(t.TempDir(), "cull.db")
	v2, err := sqlitedb.Open(ctx, p, sqlitedb.Options{Migrations: Migrations[:2]})
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range []string{
		`INSERT INTO projects(id, root) VALUES (1, '/x')`,
		`INSERT INTO sends(project_id, created_at, counts, notes, owner, delivered_to, delivered_at) VALUES (1, 1, '{}', '[]', '', 'sess-1', 2)`,
		`INSERT INTO sends(project_id, created_at, counts, notes, owner) VALUES (1, 3, '{}', '[]', '')`,
	} {
		if _, err := v2.ExecContext(ctx, q); err != nil {
			t.Fatal(err)
		}
	}
	_ = v2.Close()
	s, err := Open(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()
	if v, _ := s.db.Version(ctx); v != 3 || SchemaVersion != 3 {
		t.Fatalf("version %d, SchemaVersion %d", v, SchemaVersion)
	}
	if a := sendByID(t, s, 1); a.State != msgstate.WithAgent || a.ClaimedBy != "sess-1" {
		t.Errorf("delivered send %+v", a)
	}
	if b := sendByID(t, s, 2); b.State != msgstate.Sent || b.ClaimedBy != "" {
		t.Errorf("undelivered send %+v", b)
	}
}

// A claim records who took the send, with a fresh report token, and leaves
// it sent: only the harness's landed report says it reached the agent.
func TestClaimLeavesTheSendSentWithAToken(t *testing.T) {
	s, _ := open(t)
	project, send := sender(t, s)
	send()
	sd, ok, err := s.ClaimSend(ctx, project, "a", []string{"a"})
	if err != nil || !ok {
		t.Fatalf("%v %v", ok, err)
	}
	if sd.State != msgstate.Sent || sd.ClaimedBy != "a" || sd.ClaimedAt.IsZero() {
		t.Errorf("claimed %+v", sd)
	}
	if !regexp.MustCompile(`^[0-9a-f]{16}$`).MatchString(sd.Claim) {
		t.Errorf("token %q", sd.Claim)
	}
}

// A session that claimed a send and died before it landed does not strand
// it: once the claimer is gone, another session takes it, with a new token.
func TestASentSendWhoseClaimerIsGoneIsClaimedAgain(t *testing.T) {
	s, _ := open(t)
	project, send := sender(t, s)
	first := claimed(t, s, project, send, "a")
	if _, ok, _ := s.ClaimSend(ctx, project, "b", []string{"a", "b"}); ok {
		t.Fatal("b took a send its present claimer holds")
	}
	sd, ok, err := s.ClaimSend(ctx, project, "b", []string{"b"})
	if err != nil || !ok || sd.ID != first.ID || sd.ClaimedBy != "b" || sd.Claim == first.Claim {
		t.Fatalf("reclaim %+v %v %v", sd, ok, err)
	}
}

// A send that reached the agent is never handed out again, even when its
// claimer is gone: cull_apply must not run twice on the same answers.
func TestAWithAgentSendIsNeverClaimedAgain(t *testing.T) {
	s, _ := open(t)
	project, send := sender(t, s)
	sd := claimed(t, s, project, send, "a")
	report(t, s, "a", channelmcp.ReportLanded, sd.Claim)
	if got, ok, _ := s.ClaimSend(ctx, project, "b", []string{"b"}); ok {
		t.Fatalf("claimed a send with the agent: %+v", got)
	}
}

// landed moves the token's send to with the agent and clears held, and says
// which project changed, so the page redraws.
func TestLandedMovesTheSendToWithAgent(t *testing.T) {
	s, _ := open(t)
	project, send := sender(t, s)
	sd := claimed(t, s, project, send, "a")
	report(t, s, "a", channelmcp.ReportHeld, sd.Claim)
	changed := report(t, s, "a", channelmcp.ReportLanded, sd.Claim)
	if got := sendByID(t, s, sd.ID); got.State != msgstate.WithAgent || got.Held {
		t.Errorf("after landed %+v", got)
	}
	if !slices.Equal(changed, []int64{project}) {
		t.Errorf("changed %v", changed)
	}
}

// A late landed report leaves a send the agent already marked working as it
// is.
func TestLandedAfterWorkingLeavesItWorking(t *testing.T) {
	s, _ := open(t)
	project, send := sender(t, s)
	sd := claimed(t, s, project, send, "a")
	reply(t, s, "a", []int64{sd.ID}, msgstate.Working, "cutting")
	report(t, s, "a", channelmcp.ReportLanded, sd.Claim)
	if got := sendByID(t, s, sd.ID); got.State != msgstate.Working || got.Reply != "cutting" {
		t.Errorf("%+v", got)
	}
}

// held marks a send that is still sent, and only that.
func TestHeldMarksOnlyASentSend(t *testing.T) {
	s, _ := open(t)
	project, send := sender(t, s)
	sd := claimed(t, s, project, send, "a")
	report(t, s, "a", channelmcp.ReportHeld, sd.Claim)
	if got := sendByID(t, s, sd.ID); got.State != msgstate.Sent || !got.Held {
		t.Errorf("held %+v", got)
	}
	other := claimed(t, s, project, send, "a")
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
	project, send := sender(t, s)
	with := claimed(t, s, project, send, "a")
	report(t, s, "a", channelmcp.ReportLanded, with.Claim)
	working := claimed(t, s, project, send, "a")
	reply(t, s, "a", []int64{working.ID}, msgstate.Working, "on it")
	done := claimed(t, s, project, send, "a")
	reply(t, s, "a", []int64{done.ID}, msgstate.Answered, "applied")
	theirs := claimed(t, s, project, send, "b")
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
	project, send := sender(t, s)
	mine := claimed(t, s, project, send, "a")
	declined := claimed(t, s, project, send, "a")
	reply(t, s, "a", []int64{declined.ID}, msgstate.Declined, "not now")
	theirs := claimed(t, s, project, send, "b")
	for name, c := range map[string]struct {
		ids  []int64
		want string
	}{
		"another session's": {[]int64{mine.ID, theirs.ID}, "another session"},
		"after declined":    {[]int64{mine.ID, declined.ID}, "declined"},
		"an unknown id":     {[]int64{mine.ID, 999}, "no send 999"},
	} {
		_, err := s.ReplySends(ctx, "a", c.ids, msgstate.Answered, "done")
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
	project, send := sender(t, s)
	sd := claimed(t, s, project, send, "a")
	report(t, s, "a", channelmcp.ReportLanded, sd.Claim)
	report(t, s, "a", channelmcp.ReportTurnEnded, "")
	reply(t, s, "a", []int64{sd.ID}, msgstate.Answered, "applied after all")
	if got := sendByID(t, s, sd.ID); got.State != msgstate.Answered || got.Reply != "applied after all" {
		t.Errorf("%+v", got)
	}
}

// The page lists a project's newest sends first.
func TestRecentSendsAreNewestFirst(t *testing.T) {
	s, _ := open(t)
	project, send := sender(t, s)
	a, b, c := send(), send(), send()
	_ = a
	got, err := s.RecentSends(ctx, project, 2)
	if err != nil || len(got) != 2 || got[0].ID != c.ID || got[1].ID != b.ID {
		t.Fatalf("%+v %v", got, err)
	}
}
