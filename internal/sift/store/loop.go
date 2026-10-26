package store

// The review loop's half of the store: the agent's proposals merged into a
// round, the review's owner, Send and its delivery, and what apply did.

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/subaud-tools/tackle/internal/sift/row"
	"github.com/subaud-tools/tools-common/channelmcp"
	"github.com/subaud-tools/tools-common/msgstate"
)

// AddResult is what AddRows did.
type AddResult struct {
	Updated int // proposals merged into the round's rows
	Added   int // intake rows new to the round
	Cleared int // decisions dropped because the proposal they answered changed
}

// AddRows merges the agent's rows into a backlog or intake round (PerItem),
// all or nothing; an audit round's findings are answered per file
// (Propose). A row the
// round has gets its proposal (verdict, title, destination, text, reason)
// replaced; everything a check found stays as it was. An intake row the
// round lacks is added after the others, never certain. Anything else is an
// error: no id, an unknown verdict, an id not in the round, a decision
// (decisions are the user's), an intake row with no source. A decision on a
// row whose proposal, source or passage changed is dropped, sent or not: it
// answered another row. So is every merge into a row whose source or
// passage changed, or that was not in the round until now (dropMergesInto).
func (s *Store) AddRows(ctx context.Context, roundID int64, in []row.Row) (AddResult, error) {
	var res AddResult
	err := s.db.Tx(ctx, func(tx *sql.Tx) error {
		res = AddResult{}
		var kind string
		if err := tx.QueryRowContext(ctx, `SELECT kind FROM rounds WHERE id = ?`, roundID).Scan(&kind); errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no round %d", roundID)
		} else if err != nil {
			return err
		}
		if !PerItem(kind) {
			return fmt.Errorf("round %d is an audit round: its findings are answered one file at a time, with sift propose (sift next gives the file)", roundID)
		}
		var seq int
		if err := tx.QueryRowContext(ctx, `SELECT coalesce(max(seq), -1) + 1 FROM rows WHERE round_id = ?`, roundID).Scan(&seq); err != nil {
			return err
		}
		for i, r := range in {
			where := fmt.Sprintf("row %d", i+1)
			if r.ID != "" {
				where += " (" + r.ID + ")"
			}
			switch {
			case r.ID == "":
				return fmt.Errorf("%s: no id", where)
			case r.Decision != nil:
				return fmt.Errorf("%s: carries a decision; decisions are made on the review page", where)
			case r.Verdict != "" && !row.ValidVerdict(r.Verdict):
				return fmt.Errorf("%s: %q is not a verdict", where, r.Verdict)
			}
			var body string
			err := tx.QueryRowContext(ctx, `SELECT body FROM rows WHERE round_id = ? AND row_id = ?`, roundID, r.ID).Scan(&body)
			if errors.Is(err, sql.ErrNoRows) {
				if r.Check != "intake" {
					return fmt.Errorf("%s: not in round %d (only intake rows are new)", where, roundID)
				}
				if r.Source.Entry == "" && r.Source.File == "" {
					return fmt.Errorf("%s: an intake row needs a source (entry or file)", where)
				}
				r.Certain = false
				r.Source.Canon = ""
				if r.Source.FromDisk() {
					r.Source.Canon = row.Resolve(r.Source.File)
				}
				if err := putRow(ctx, tx, roundID, seq, r, true); err != nil {
					return err
				}
				seq++
				res.Added++
				// A merge into this id was approved while it was not here.
				k, err := dropMergesInto(ctx, tx, roundID, r.ID)
				if err != nil {
					return err
				}
				res.Cleared += k
				continue
			}
			if err != nil {
				return err
			}
			var cur row.Row
			if err := json.Unmarshal([]byte(body), &cur); err != nil {
				return err
			}
			next := cur
			if r.Check == "intake" && cur.Check == "intake" {
				r.Certain = false
				r.Source.Canon = ""
				if r.Source.FromDisk() {
					r.Source.Canon = row.Resolve(r.Source.File)
				}
				next = r // an intake row is the agent's own: all of it is replaced
			}
			next.Verdict, next.Title, next.Destination, next.Text, next.Reason = r.Verdict, r.Title, r.Destination, r.Text, r.Reason
			changed := actionChanged(cur, next)
			if err := putRow(ctx, tx, roundID, 0, next, false); err != nil {
				return err
			}
			res.Updated++
			if changed {
				d, err := tx.ExecContext(ctx, `DELETE FROM decisions WHERE round_id = ? AND row_id = ?`, roundID, r.ID)
				if err != nil {
					return err
				}
				if k, _ := d.RowsAffected(); k > 0 {
					res.Cleared++
				}
			}
			if cur.Source != next.Source || cur.Passage != next.Passage {
				// A merge into this row was approved against its old source
				// and passage.
				k, err := dropMergesInto(ctx, tx, roundID, r.ID)
				if err != nil {
					return err
				}
				res.Cleared += k
			}
		}
		return bump(ctx, tx, roundID)
	})
	if err != nil {
		return AddResult{}, err
	}
	return res, nil
}

// dropMergesInto drops every decision in the round whose verdict (the
// decision's own, or the proposal it accepted) is merge:id, and says how
// many.
func dropMergesInto(ctx context.Context, tx *sql.Tx, roundID int64, id string) (int, error) {
	q, err := tx.QueryContext(ctx, `SELECT d.row_id, d.verdict, r.body FROM decisions d
		JOIN rows r ON r.round_id = d.round_id AND r.row_id = d.row_id WHERE d.round_id = ?`, roundID)
	if err != nil {
		return 0, err
	}
	var drop []string
	for q.Next() {
		var rowID, verdict, body string
		if err := q.Scan(&rowID, &verdict, &body); err != nil {
			_ = q.Close()
			return 0, err
		}
		if verdict == "" {
			var r row.Row
			if err := json.Unmarshal([]byte(body), &r); err != nil {
				_ = q.Close()
				return 0, err
			}
			verdict = r.Verdict
		}
		if row.MergeTarget(verdict) == id {
			drop = append(drop, rowID)
		}
	}
	if err := q.Err(); err != nil {
		_ = q.Close()
		return 0, err
	}
	if err := q.Close(); err != nil {
		return 0, err
	}
	for _, rowID := range drop {
		if _, err := tx.ExecContext(ctx, `DELETE FROM decisions WHERE round_id = ? AND row_id = ?`, roundID, rowID); err != nil {
			return 0, err
		}
	}
	return len(drop), nil
}

// actionChanged reports whether what applying a row would do changed: its
// proposal (verdict, title, destination, text) or what it is about (its
// source and passage). A decision answered the old row, so it is dropped.
func actionChanged(a, b row.Row) bool {
	return a.Verdict != b.Verdict || a.Title != b.Title || a.Destination != b.Destination || a.Text != b.Text ||
		a.Source != b.Source || a.Passage != b.Passage
}

func putRow(ctx context.Context, tx *sql.Tx, roundID int64, seq int, r row.Row, insert bool) error {
	r.Decision, r.Fingerprint = nil, ""
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	if insert {
		_, err = tx.ExecContext(ctx, `INSERT INTO rows(round_id, seq, row_id, check_, body) VALUES (?,?,?,?,?)`, roundID, seq, r.ID, r.Check, string(b))
	} else {
		_, err = tx.ExecContext(ctx, `UPDATE rows SET body = ? WHERE round_id = ? AND row_id = ?`, string(b), roundID, r.ID)
	}
	return err
}

// SetOwner makes session (shown as label) the one the round's review goes to.
func (s *Store) SetOwner(ctx context.Context, roundID int64, session, label string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE rounds SET owner_session = ?, owner_label = ? WHERE id = ?`, session, label, roundID)
	return err
}

// Counts is what one Send carried.
type Counts struct {
	Accept int `json:"accept"`
	Edit   int `json:"edit"`
	Reject int `json:"reject"`
}

// Total is every answer counted.
func (c Counts) Total() int { return c.Accept + c.Edit + c.Reject }

// Note is the user's note on one row or file, with where it is.
type Note struct {
	Row  string `json:"row"` // file:line · check, or the file
	Note string `json:"note"`
	// Again: the user disagrees that the file needs no change; the agent
	// recommends it again.
	Again bool `json:"again,omitempty"`
}

// Send is one press of Send. ID is 0 when nothing was new (nothing is
// recorded then).
type Send struct {
	ID        int64
	Round     int64
	CreatedAt time.Time
	Counts    Counts
	Notes     []Note
	Owner     string // the review's owner session at send time ("" if none)
	// ClaimedBy is the session whose wait took the send ("" until one did),
	// at ClaimedAt, with Claim, the report token the harness names it by.
	ClaimedBy string
	ClaimedAt time.Time
	Claim     string
	// State is where the send stands (msgstate); Held, that the harness
	// holds it while the agent is mid-turn; Reply, the agent's working line
	// or its reply.
	State msgstate.State
	Held  bool
	Reply string
	// Files and Rows name the decisions this Send marked sent (file keys,
	// row ids); not recorded.
	Files, Rows []string
}

// Shown is the decisions the page showed unsent when Send was pressed, by
// file key and by row id: each by its decision id and the item's print as
// the page showed it (for an edited merge, the print covers the target's).
// A Send given it sends each only while both still match; a decision made
// since, even an identical one, has another id and is left for the next.
type Shown struct {
	Files map[string]Seen
	Rows  map[string]Seen
}

// Send marks the round's unsent decisions sent, on files and on rows, and
// records one send for them, in one transaction. Given what the page
// showed (shown, non-nil), it sends only those decisions, each only while
// its id and its item's print are still the ones shown: a decision made
// since Send was pressed waits for the next, and so does one whose item
// changed under it. nil sends every unsent decision. With nothing to
// send it returns a zero Send and records nothing. The Send names the
// decisions it marked (Files, Rows).
func (s *Store) Send(ctx context.Context, roundID int64, owner string, shown *Shown) (Send, error) {
	now := time.Now()
	sd := Send{Round: roundID, CreatedAt: now, Owner: owner}
	err := s.db.Tx(ctx, func(tx *sql.Tx) error {
		// The first statement is a write (the main connection begins deferred).
		if _, err := tx.ExecContext(ctx, `UPDATE rounds SET rev = rev WHERE id = ?`, roundID); err != nil {
			return err
		}
		count := func(action string) {
			switch action {
			case "accept":
				sd.Counts.Accept++
			case "edit":
				sd.Counts.Edit++
			case "reject":
				sd.Counts.Reject++
			}
		}
		items, err := filesIn(ctx, tx, roundID)
		if err != nil {
			return err
		}
		for _, it := range items {
			d := it.Decision
			if d == nil || d.Sent || !shown.hasFile(it.Key, d.ID, it.Fingerprint) {
				continue
			}
			sd.Files = append(sd.Files, it.Key)
			count(d.Action)
			if d.Note != "" {
				sd.Notes = append(sd.Notes, Note{Row: it.Source.File, Note: d.Note, Again: it.Unchanged && d.Action == "reject"})
			}
		}
		rows, err := rowsOf(ctx, tx, roundID)
		if err != nil {
			return err
		}
		for _, r := range rows {
			d := r.Decision
			if d == nil || d.Sent || !shown.hasRow(r.ID, d.ID, r.Fingerprint) {
				continue
			}
			sd.Rows = append(sd.Rows, r.ID)
			count(d.Action)
			if d.Note != "" {
				sd.Notes = append(sd.Notes, Note{Row: where(r), Note: d.Note})
			}
		}
		if len(sd.Files)+len(sd.Rows) == 0 {
			sd.Counts, sd.Notes = Counts{}, nil
			return nil
		}
		for _, k := range sd.Files {
			if _, err := tx.ExecContext(ctx, `UPDATE file_decisions SET sent_at = ? WHERE round_id = ? AND key = ?`, now.UnixMilli(), roundID, k); err != nil {
				return err
			}
		}
		for _, id := range sd.Rows {
			if _, err := tx.ExecContext(ctx, `UPDATE decisions SET sent_at = ? WHERE round_id = ? AND row_id = ?`, now.UnixMilli(), roundID, id); err != nil {
				return err
			}
		}
		counts, _ := json.Marshal(sd.Counts)
		notes, _ := json.Marshal(append([]Note{}, sd.Notes...))
		res, err := tx.ExecContext(ctx, `INSERT INTO sends(round_id, created_at, counts, notes, owner) VALUES (?,?,?,?,?)`,
			roundID, now.UnixMilli(), string(counts), string(notes), owner)
		if err != nil {
			return err
		}
		if sd.ID, err = res.LastInsertId(); err != nil {
			return err
		}
		return bump(ctx, tx, roundID)
	})
	if err != nil {
		return Send{}, err
	}
	return sd, nil
}

// hasFile reports whether the page showed decision id on file key, with the
// file's print as it is now (a nil Shown shows everything).
func (sh *Shown) hasFile(key, id, print string) bool {
	return sh == nil || saw(sh.Files[key], id, print)
}

// hasRow is hasFile for a row.
func (sh *Shown) hasRow(rowID, id, print string) bool {
	return sh == nil || saw(sh.Rows[rowID], id, print)
}

// saw reports whether w is decision id on an item whose print is print.
func saw(w Seen, id, print string) bool {
	return id != "" && w.Decision == id && w.Fingerprint == print
}

// where names a row for a note: file:line · check.
func where(r row.Row) string {
	at := r.Source.File
	if at == "" {
		at = r.Source.Entry
	}
	if r.Source.Start > 0 {
		at = fmt.Sprintf("%s:%d", at, r.Source.Start)
	}
	return at + " · " + r.Check
}

const sendCols = `id, round_id, created_at, counts, notes, owner, delivered_to, delivered_at, claim, state, held, reply`

type scanner interface{ Scan(...any) error }

func scanSend(r scanner) (Send, error) {
	var (
		sd           Send
		at           int64
		counts, note string
		to           sql.NullString
		dat          sql.NullInt64
	)
	if err := r.Scan(&sd.ID, &sd.Round, &at, &counts, &note, &sd.Owner, &to, &dat, &sd.Claim, &sd.State, &sd.Held, &sd.Reply); err != nil {
		return Send{}, err
	}
	sd.CreatedAt, sd.ClaimedBy = time.UnixMilli(at), to.String
	if dat.Valid {
		sd.ClaimedAt = time.UnixMilli(dat.Int64)
	}
	if err := json.Unmarshal([]byte(counts), &sd.Counts); err != nil {
		return Send{}, err
	}
	if err := json.Unmarshal([]byte(note), &sd.Notes); err != nil {
		return Send{}, err
	}
	return sd, nil
}

// ClaimSend hands session the oldest send it may take: one whose recorded
// owner is session, has no owner, or whose owner is not in present (the
// sessions serve sees connected). A send is claimable while no session has
// claimed it, or while it is still sent and its claimer is not in present:
// a session that died before the send landed does not strand it. A send
// that reached the agent is never handed out again. The claim records the
// session, the time and a new report token, and leaves the send sent: the
// harness's landed report moves it. Picking and claiming it is one
// statement, so two claimers never both get it. ok is false when there is
// none.
func (s *Store) ClaimSend(ctx context.Context, session string, present []string) (Send, bool, error) {
	var sd Send
	found := true
	// notIn is "col is not one of present"; with nobody present, every
	// owner and claimer is gone.
	notIn := func(col string) string {
		if len(present) == 0 {
			return "1"
		}
		return col + " NOT IN (?" + strings.Repeat(",?", len(present)-1) + ")"
	}
	claimable := `(delivered_to IS NULL OR (state = 'sent' AND ` + notIn("delivered_to") + `))`
	args := []any{session, time.Now().UnixMilli(), reportToken()}
	addPresent := func() {
		for _, p := range present {
			args = append(args, p)
		}
	}
	addPresent()
	addPresent()
	args = append(args, session)
	addPresent()
	err := s.db.Tx(ctx, func(tx *sql.Tx) error {
		var err error
		sd, err = scanSend(tx.QueryRowContext(ctx, `UPDATE sends SET delivered_to = ?, delivered_at = ?, claim = ?
			WHERE `+claimable+` AND id = (SELECT id FROM sends WHERE `+claimable+` AND (owner = ? OR owner = '' OR `+notIn("owner")+`) ORDER BY id LIMIT 1)
			RETURNING `+sendCols, args...))
		if errors.Is(err, sql.ErrNoRows) {
			found = false
			return nil
		}
		return err
	})
	if err != nil {
		return Send{}, false, err
	}
	return sd, found, nil
}

// reportToken is a claim's report token: 16 random hex characters, which
// the Claude Code mod passes as an argument, so alphanumeric.
func reportToken() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b) // crypto/rand.Read never fails
	return hex.EncodeToString(b)
}

// Sends counts the round's sends.
func (s *Store) Sends(ctx context.Context, roundID int64) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT count(*) FROM sends WHERE round_id = ?`, roundID).Scan(&n)
	return n, err
}

// Unclaimed lists the sends no session has claimed, oldest first.
func (s *Store) Unclaimed(ctx context.Context) ([]Send, error) {
	return listSends(ctx, s.db, `WHERE delivered_to IS NULL ORDER BY id`)
}

// Unsettled lists the sends not answered, declined or failed, oldest first.
func (s *Store) Unsettled(ctx context.Context) ([]Send, error) {
	return listSends(ctx, s.db, `WHERE state NOT IN ('answered', 'declined', 'failed') ORDER BY id`)
}

// RecentSends lists the round's newest n sends, newest first.
func (s *Store) RecentSends(ctx context.Context, roundID int64, n int) ([]Send, error) {
	return listSends(ctx, s.db, `WHERE round_id = ? ORDER BY id DESC LIMIT ?`, roundID, n)
}

func listSends(ctx context.Context, db querier, where string, args ...any) ([]Send, error) {
	rs, err := db.QueryContext(ctx, `SELECT `+sendCols+` FROM sends `+where, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rs.Close() }()
	var out []Send
	for rs.Next() {
		sd, err := scanSend(rs)
		if err != nil {
			return nil, err
		}
		out = append(out, sd)
	}
	return out, rs.Err()
}

// ReportSends applies a harness's report for session: held marks the
// token's send held while it is still sent; landed moves it to with the
// agent and clears held; turn_ended moves each send session claimed that is
// with the agent or working to no reply. Each move is one msgstate.Allowed
// draws, and a move it refuses is skipped, so a late landed report leaves a
// working send working. Each round a report changes has its revision
// bumped, so the page redraws.
func (s *Store) ReportSends(ctx context.Context, session string, event channelmcp.ReportEvent, token string) error {
	switch event {
	case channelmcp.ReportHeld, channelmcp.ReportLanded, channelmcp.ReportTurnEnded:
	default:
		return fmt.Errorf("event %q: want held, landed or turn_ended", event)
	}
	return s.db.Tx(ctx, func(tx *sql.Tx) error {
		// The first statement is a write (the main connection begins deferred).
		if _, err := tx.ExecContext(ctx, `UPDATE sends SET held = held WHERE delivered_to = ?`, session); err != nil {
			return err
		}
		where, args := `WHERE delivered_to = ? AND claim = ?`, []any{session, token}
		if event == channelmcp.ReportTurnEnded {
			where, args = `WHERE delivered_to = ?`, []any{session}
		}
		sends, err := listSends(ctx, tx, where, args...)
		if err != nil {
			return err
		}
		for _, sd := range sends {
			state, held := sd.State, sd.Held
			switch event {
			case channelmcp.ReportHeld:
				held = held || sd.State == msgstate.Sent
			case channelmcp.ReportLanded:
				if msgstate.Allowed(sd.State, msgstate.WithAgent) {
					state = msgstate.WithAgent
				}
				held = false
			case channelmcp.ReportTurnEnded:
				if msgstate.Allowed(sd.State, msgstate.NoReply) {
					state = msgstate.NoReply
				}
			}
			if state == sd.State && held == sd.Held {
				continue
			}
			if _, err := tx.ExecContext(ctx, `UPDATE sends SET state = ?, held = ? WHERE id = ?`, state, held, sd.ID); err != nil {
				return err
			}
			if err := bump(ctx, tx, sd.Round); err != nil {
				return err
			}
		}
		return nil
	})
}

// Refused is a reply the lifecycle refuses, with the reason.
type Refused struct{ Reason string }

func (r Refused) Error() string { return r.Reason }

// ReplySends is the agent's acknowledgement of sends ids, all or nothing:
// working with its line, or answered, declined or failed with its reply.
// Every id must be a send session claimed that msgstate.Acknowledge lets
// move to state; one already in the settled state asked for is left as it
// is, and one already working takes the new line. Anything else refuses
// the whole call (Refused) and changes nothing.
func (s *Store) ReplySends(ctx context.Context, session string, ids []int64, state msgstate.State, text string) error {
	if state != msgstate.Working && !msgstate.Final(state) {
		return Refused{fmt.Sprintf("state %q: want working, answered, declined or failed", state)}
	}
	if len(ids) == 0 {
		return Refused{"ids: name the sends this reply is for"}
	}
	text = strings.TrimSpace(text)
	return s.db.Tx(ctx, func(tx *sql.Tx) error {
		// The first statement is a write (the main connection begins deferred).
		if _, err := tx.ExecContext(ctx, `UPDATE sends SET held = held WHERE delivered_to = ?`, session); err != nil {
			return err
		}
		var move []Send
		for _, id := range ids {
			sd, err := scanSend(tx.QueryRowContext(ctx, `SELECT `+sendCols+` FROM sends WHERE id = ?`, id))
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return Refused{fmt.Sprintf("no send %d", id)}
			case err != nil:
				return err
			case sd.ClaimedBy == "":
				return Refused{fmt.Sprintf("send %d: no session has it yet", id)}
			case sd.ClaimedBy != session:
				return Refused{fmt.Sprintf("send %d: another session has it", id)}
			case sd.State == state && msgstate.Final(state):
				continue
			case sd.State == state:
				// working again: only the line changes
			case !msgstate.Acknowledge(sd.State, state):
				return Refused{fmt.Sprintf("send %d: is %s and cannot be marked %s", id, sd.State.Label(), state.Label())}
			}
			move = append(move, sd)
		}
		for _, sd := range move {
			if _, err := tx.ExecContext(ctx, `UPDATE sends SET state = ?, held = 0, reply = ? WHERE id = ?`, state, text, sd.ID); err != nil {
				return err
			}
			if err := bump(ctx, tx, sd.Round); err != nil {
				return err
			}
		}
		return nil
	})
}

// Apply is what sift apply did in one repo for a round.
type Apply struct {
	Round  int64
	Repo   string // the primary clone
	Base   string // the ref the branch was cut from
	Branch string
	PR     string // the pull request's URL, when one was opened
	// State is pr (branch pushed, PR opened), branch (committed, left
	// local), held (not touched; Detail says why) or failed.
	State  string
	Detail string
	Rows   []string // the row ids applied
	At     time.Time
	// Last is a later attempt that did not replace a success (pr or branch):
	// a success is permanent. Its State is "" when there was none.
	Last Attempt
}

// Attempt is one apply that is not the record's own.
type Attempt struct {
	State  string
	Detail string
	At     time.Time
}

// Succeeded reports whether apply wrote the repo's branch: a pr or branch
// record.
func (a Apply) Succeeded() bool { return a.State == "pr" || a.State == "branch" }

// RecordApply stores what apply did in one repo. A success (pr or branch)
// is permanent: a later attempt after it is kept as Last and changes
// nothing else.
func (s *Store) RecordApply(ctx context.Context, a Apply) error {
	if a.At.IsZero() {
		a.At = time.Now()
	}
	ids, _ := json.Marshal(append([]string{}, a.Rows...))
	return s.db.Tx(ctx, func(tx *sql.Tx) error {
		var state string
		err := tx.QueryRowContext(ctx, `SELECT state FROM applies WHERE round_id = ? AND repo = ?`, a.Round, a.Repo).Scan(&state)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			if _, err := tx.ExecContext(ctx, `INSERT INTO applies(round_id, repo, base, branch, pr, state, detail, rows, at) VALUES (?,?,?,?,?,?,?,?,?)`,
				a.Round, a.Repo, a.Base, a.Branch, a.PR, a.State, a.Detail, string(ids), a.At.UnixMilli()); err != nil {
				return err
			}
		case err != nil:
			return err
		case (Apply{State: state}).Succeeded():
			if _, err := tx.ExecContext(ctx, `UPDATE applies SET last_state = ?, last_detail = ?, last_at = ? WHERE round_id = ? AND repo = ?`,
				a.State, a.Detail, a.At.UnixMilli(), a.Round, a.Repo); err != nil {
				return err
			}
		default:
			if _, err := tx.ExecContext(ctx, `UPDATE applies SET base = ?, branch = ?, pr = ?, state = ?, detail = ?, rows = ?, at = ?,
				last_state = '', last_detail = '', last_at = 0 WHERE round_id = ? AND repo = ?`,
				a.Base, a.Branch, a.PR, a.State, a.Detail, string(ids), a.At.UnixMilli(), a.Round, a.Repo); err != nil {
				return err
			}
		}
		return bump(ctx, tx, a.Round)
	})
}

// Applies lists what apply did for a round, by repo.
func (s *Store) Applies(ctx context.Context, roundID int64) ([]Apply, error) {
	rs, err := s.db.QueryContext(ctx, `SELECT round_id, repo, base, branch, pr, state, detail, rows, at, last_state, last_detail, last_at
		FROM applies WHERE round_id = ? ORDER BY repo`, roundID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rs.Close() }()
	var out []Apply
	for rs.Next() {
		var a Apply
		var ids string
		var at, lastAt int64
		if err := rs.Scan(&a.Round, &a.Repo, &a.Base, &a.Branch, &a.PR, &a.State, &a.Detail, &ids, &at,
			&a.Last.State, &a.Last.Detail, &lastAt); err != nil {
			return nil, err
		}
		a.At = time.UnixMilli(at)
		if lastAt != 0 {
			a.Last.At = time.UnixMilli(lastAt)
		}
		if err := json.Unmarshal([]byte(ids), &a.Rows); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rs.Err()
}
