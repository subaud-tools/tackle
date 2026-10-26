// Package channel is `sift channel`: the MCP server each agent session
// starts (Claude Code through the sift plugin, pi through channels.tools). It
// starts sift serve when it isn't running, gives the agent sift's tools and,
// under pi, registers the session's presence, holds a long poll open for the
// user's decisions and injects each Send as a notifications/claude/channel
// event. Under Claude Code `sift watch` runs that loop instead and prints
// each Send as a JSON line for the plugin's mod. The
// presence and wait loops are copied in shape from cull's channel (sift does
// not import it).
//
// stdout carries MCP frames only; every log line goes to the Log writer
// (stderr in the command).
package channel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/subaud-tools/tackle/internal/sift/config"
	"github.com/subaud-tools/tools-common/channelmcp"
	"github.com/subaud-tools/tools-common/harness"
	"github.com/subaud-tools/tools-common/localweb"
)

// Instructions are the standing guidance the agent gets with the channel.
const Instructions = `sift audits the instruction files coding agents read (global files, repo files, skills). For each file with findings you recommend one revised version, and the user accepts, edits or rejects each file on a review page.
- Run sift_check when the user asks for an audit. It records a round and says how many files need a recommendation.
- Recommend every file: sift_next gives one file with its content, findings and the guidance for the rewrite; sift_propose stores the whole revised file. Repeat until sift_next says done. Recommend files that move text between them in one sift_propose call, each linking the other.
- Then run sift_review once to open the page for the user, tell them how many files wait, and carry on with other work. The page opens only once every file has a recommendation.
- The user's decisions arrive as an event from sift: <channel source="sift" round="…" send_id="…"> in pi, or a turn whose first line is sift event round=… send_id=… in Claude Code. On a send that takes more than a moment, first call sift_reply with its send_id, state working and one line on what you are doing. Follow its steps: sift_apply writes a branch per repo; then run ` + "`sift reconcile`" + ` and review each branch. When done, sift_reply with answered and a short report, or declined or failed with why.
- The pull requests are the user's to merge. Once they are merged or closed, run sift_clean to remove the branches and worktrees apply made.
- sift_status tells you the round's state, what is open, decided and sent, and which sends are not settled yet.
- When sift reports an error, say what failed and leave the files as they are.`

// Identity is the session this channel serves.
type Identity struct {
	Session string
	Harness string // "claude" or "pi"
	Label   string
	CWD     string
}

// FromEnv resolves the session by the family rule (tools-common harness).
// Session is "" outside a harness.
func FromEnv() Identity {
	c := harness.FromEnv()
	id := Identity{Session: c.SessionID, CWD: c.CWD, Harness: "pi"}
	if c.IsClaude() {
		id.Harness = "claude"
	}
	if id.CWD == "" {
		id.CWD, _ = os.Getwd()
	}
	id.Label = id.Harness + " · " + filepath.Base(id.CWD)
	return id
}

// Channel is a running sift channel.
type Channel struct {
	ID     Identity
	Client *Client
	Server *channelmcp.Server
	Log    io.Writer
	// Open opens a URL in the user's browser (tests replace it).
	Open func(url string) error
	// LoadConfig reads sift's config; LookPath finds gh (tests replace them).
	LoadConfig func() (config.Config, error)
	LookPath   func(string) (string, error)
	// Retry is how long to wait before trying serve again; Presence how often
	// to refresh presence; Poll the long-poll timeout. Tests shorten them.
	Retry    time.Duration
	Presence time.Duration
	Poll     time.Duration

	// Push is where the wait loop sends each Send: Server for `sift
	// channel`, a JSON line on stdout for `sift watch`. One loop, two outputs.
	Push channelmcp.Pusher
	// ServeOnly is `sift channel` under Claude Code: the plugin's mod runs
	// `sift watch`, which takes the sends, so this process serves the tools
	// and never waits. Two takers would split a session's sends.
	ServeOnly bool

	runCtx  context.Context
	reports chan report // pi's reports, posted to serve in order by one goroutine
}

// report is one report from the harness, waiting to be posted to serve.
type report struct {
	event channelmcp.ReportEvent
	token string
}

// New builds a channel for id.
func New(id Identity, c *Client, version string) *Channel {
	ch := &Channel{ID: id, Client: c, Log: io.Discard, Open: localweb.OpenBrowser,
		LoadConfig: func() (config.Config, error) { return config.Load(config.Path()) }, LookPath: exec.LookPath,
		Retry: 5 * time.Second, Presence: 30 * time.Second, Poll: 55 * time.Second, reports: make(chan report, 64)}
	ch.Server = channelmcp.New(channelmcp.Handler{
		Name: "sift", Version: version, Instructions: Instructions, Tools: Tools(),
		Call: func(name string, args json.RawMessage) (string, error) {
			ctx := ch.runCtx
			if ctx == nil {
				ctx = context.Background()
			}
			return ch.Call(ctx, name, args)
		},
		// Run's read loop hands each report over and goes on: posting it
		// to serve waits on the network, and tool calls must not.
		OnReport: func(event channelmcp.ReportEvent, token string) {
			select {
			case ch.reports <- report{event, token}:
			default:
				_, _ = fmt.Fprintf(ch.Log, "sift channel: dropped report %s %s: too many waiting\n", event, token)
			}
		},
	})
	ch.Push = ch.Server
	return ch
}

// ForSession is the channel `sift channel` runs: this session, and serve-only
// when the session is Claude Code's (see ServeOnly).
func ForSession(c *Client, version string) *Channel {
	ch := New(FromEnv(), c, version)
	ch.ServeOnly = harness.FromEnv().IsClaude()
	return ch
}

// Run serves MCP on r/w and, unless the channel is serve-only, the presence
// and wake loop when the session is known.
func (ch *Channel) Run(ctx context.Context, r io.Reader, w io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ch.runCtx = ctx // set before Server.Run starts reading
	switch {
	case ch.ServeOnly:
		_, _ = fmt.Fprintln(ch.Log, "sift channel: Claude Code session: serving tools only; the sift plugin's `sift watch` delivers the user's decisions")
	case ch.ID.Session == "":
		_, _ = fmt.Fprintf(ch.Log, "sift channel: %v\n", errNoSession)
	default:
		go ch.loop(ctx)
		go ch.postReports(ctx)
	}
	return ch.Server.Run(r, w)
}

var errNoSession = errors.New("no session id (not under Claude Code or pi); the user's decisions will not reach this session")

// postReports posts each report the harness sent to serve, in the order
// sent, until ctx ends. A failed post is logged, not retried.
func (ch *Channel) postReports(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case r := <-ch.reports:
			if err := ch.Client.Report(ctx, ch.ID.Session, r.event, r.token); err != nil && ctx.Err() == nil {
				_, _ = fmt.Fprintf(ch.Log, "sift channel: report %s: %v\n", r.event, err)
			}
		}
	}
}

// Watch is `sift watch`: the presence and wake loop alone, until ctx ends.
func (ch *Channel) Watch(ctx context.Context) error {
	if ch.ID.Session == "" {
		return errNoSession
	}
	ch.loop(ctx)
	return nil
}

func (ch *Channel) presence(ctx context.Context) error {
	_, err := ch.Client.Do(ctx, http.MethodPost, "/api/agent/presence", map[string]any{
		"session": ch.ID.Session, "harness": ch.ID.Harness, "label": ch.ID.Label}, nil)
	return err
}

// loop starts serve if needed, keeps presence fresh and long-polls for sends,
// emitting each claimed send as one channel event. serve claims a send for
// exactly one caller, so an event is never duplicated.
func (ch *Channel) loop(ctx context.Context) {
	var lastPresence time.Time
	var lastErr string
	logOnce := func(err error) {
		if msg := err.Error(); msg != lastErr {
			_, _ = fmt.Fprintf(ch.Log, "sift channel: %s\n", msg)
			lastErr = msg
		}
	}
	for ctx.Err() == nil {
		if time.Since(lastPresence) > ch.Presence {
			if err := ch.presence(ctx); err != nil {
				logOnce(err)
				sleep(ctx, ch.Retry)
				continue
			}
			lastPresence, lastErr = time.Now(), ""
		}
		var got struct {
			ID     int64  `json:"id"`
			Round  int64  `json:"round"`
			Text   string `json:"text"`
			Report string `json:"report"`
		}
		code, err := ch.Client.Do(ctx, http.MethodGet, "/api/agent/wait?"+q(
			"session", ch.ID.Session, "timeout", strconv.Itoa(max(int(ch.Poll.Seconds()), 1))), nil, &got)
		switch {
		case err != nil:
			if ctx.Err() != nil {
				return
			}
			var se *StatusError
			if errors.As(err, &se) {
				lastPresence = time.Time{}
			}
			logOnce(err)
			sleep(ctx, ch.Retry)
		case code == http.StatusOK && got.ID != 0:
			lastErr = ""
			meta := map[string]string{"source": "sift", "round": strconv.FormatInt(got.Round, 10), "send_id": strconv.FormatInt(got.ID, 10),
				channelmcp.MetaReport: got.Report}
			if err := ch.Push.Notify(got.Text, meta); err != nil {
				_, _ = fmt.Fprintf(ch.Log, "sift channel: notify: %v\n", err)
			}
		default:
			lastErr = ""
		}
	}
}

func sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}
