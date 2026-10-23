// Package channel is `imprint channel`: the MCP server an agent session
// starts. Each tool does one npm step through package ops, returns at once
// with what it started, and sends the step's progress to the session as
// events: under pi as notifications/claude/channel, under Claude Code as JSON
// lines in the session's spool, which `imprint watch` (the imprint plugin's
// mod) prints and the mod turns into turns.
//
// stdout carries MCP frames only; every log line goes to the Log writer
// (stderr in the command).
package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/subaud-tools/tackle/internal/imprint/ops"
	"github.com/subaud-tools/tools-common/channelmcp"
	"github.com/subaud-tools/tools-common/harness"
)

// Instructions are the standing guidance the agent gets with the channel.
const Instructions = `imprint guides you through npm package work: logging in, publishing a tarball, setting a package's trusted publisher, deprecating a package, and confirming that a CI release landed signed by the right repo.
- You decide every step; imprint does the one npm step you ask for and tells you what it means. The user logs in and approves each npm operation in the browser: imprint opens each approval link for them and tells you it is waiting.
- Every tool except imprint_show and imprint_status returns at once with what it started. Its progress arrives as events from imprint: <channel source="imprint" op="…" package="…"> in pi, or a turn whose first line is imprint event op=… package=… in Claude Code. Each says "waiting for your approval: <url>", "approved", "live on npm", or "failed:" with what npm's error usually means and the file holding npm's full output. Relay approval links to the user.
- A step is done only when imprint says the registry shows it. After a publish, don't work out the next version until imprint says it is live.
- imprint_status lists what is waiting on the user's approval or on the registry.
- Read the imprint skill for the procedures: moving a trusted publisher, renaming a package, proving a new trusted publisher with one CI release.`

// Channel is a running imprint channel.
type Channel struct {
	Ops    *ops.Ops
	Server *channelmcp.Server
	Log    io.Writer

	// Push is where each job's reports go: Server under pi; under Claude
	// Code the session's spool, which `imprint watch` follows.
	Push  channelmcp.Pusher
	spool *Spool
	// claude is a Claude Code session; spoolErr why its spool didn't open.
	claude   bool
	spoolErr error

	mu     sync.Mutex
	jobs   map[int]*job
	nextID int
	runCtx context.Context
}

// job is one operation in flight.
type job struct {
	tool, pkg, what string
	last            string
}

// New builds a channel running o's operations.
func New(o *ops.Ops, version string) *Channel {
	ch := &Channel{Ops: o, Log: io.Discard, jobs: map[int]*job{}}
	ch.Server = channelmcp.New(channelmcp.Handler{
		Name: "imprint", Version: version, Instructions: Instructions, Tools: Tools(),
		Call: func(name string, args json.RawMessage) (string, error) {
			return ch.Call(name, args)
		},
	})
	ch.Push = ch.Server
	return ch
}

// ForSession is the channel `imprint channel` runs. Under Claude Code its
// messages go to the session's spool under imprint's state directory dir
// (see Spool); under pi to Server, as always.
func ForSession(o *ops.Ops, version, dir string) *Channel {
	ch := New(o, version)
	c := harness.FromEnv()
	if !c.IsClaude() {
		return ch
	}
	ch.claude = true
	sp, err := OpenSpool(dir, c.SessionID)
	if err != nil {
		ch.spoolErr = err
		return ch
	}
	ch.spool, ch.Push = sp, sp
	return ch
}

// Run serves MCP on r/w until r ends or ctx does. Operations still running end with
// it: an npm call stops with the session. The session's spool is removed
// when it ends.
func (ch *Channel) Run(ctx context.Context, r io.Reader, w io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	switch {
	case ch.spool != nil:
		defer func() { _ = ch.spool.Close() }()
		_, _ = fmt.Fprintf(ch.Log, "imprint channel: Claude Code session: messages go to %s for the imprint plugin's `imprint watch`\n", ch.spool.Path())
	case ch.spoolErr != nil:
		_, _ = fmt.Fprintf(ch.Log, "imprint channel: Claude Code session, but no spool (%v): messages go out as channel notifications, which Claude Code drops unless the channel flag is on\n", ch.spoolErr)
	}
	ch.mu.Lock()
	ch.runCtx = ctx // set before Server.Run starts reading
	ch.mu.Unlock()
	// Return on a signal too, not only when r ends: Claude Code signals the
	// channel while its stdin may still be open, and the spool must go.
	served := make(chan error, 1)
	go func() { served <- ch.Server.Run(r, w) }()
	select {
	case err := <-served:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (ch *Channel) ctx() context.Context {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	if ch.runCtx == nil {
		return context.Background()
	}
	return ch.runCtx
}

// start runs fn in the background as a job of tool on pkg, sending each
// of its reports to the session, and returns what the tool says at once.
func (ch *Channel) start(tool, pkg, what string, fn func(ctx context.Context, report ops.Report) error) string {
	ch.mu.Lock()
	ch.nextID++
	id := ch.nextID
	j := &job{tool: tool, pkg: pkg, what: what, last: "started"}
	ch.jobs[id] = j
	ch.mu.Unlock()
	meta := map[string]string{"op": tool}
	if pkg != "" {
		meta["package"] = pkg
	}
	report := func(text string) {
		ch.mu.Lock()
		j.last = text
		ch.mu.Unlock()
		if err := ch.Push.Notify(text, meta); err != nil {
			_, _ = fmt.Fprintf(ch.Log, "imprint channel: notify: %v\n", err)
		}
	}
	ctx := ch.ctx()
	go func() {
		defer func() {
			ch.mu.Lock()
			delete(ch.jobs, id)
			ch.mu.Unlock()
		}()
		_ = fn(ctx, report) // every operation reports its own outcome
	}()
	return "started: " + what + "; progress arrives as messages from imprint"
}

// status lists the operations in flight with their last message.
func (ch *Channel) status() string {
	s := ch.jobsStatus()
	if ch.claude {
		s += "\n\nThis is a Claude Code session: imprint's messages reach it through `imprint watch`, which the imprint plugin starts, " +
			"not through this process. If one never arrived, check that the imprint plugin is installed and enabled (/plugin)."
	}
	return s
}

func (ch *Channel) jobsStatus() string {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	if len(ch.jobs) == 0 {
		return "nothing is running: nothing waits on the user's approval or on the registry"
	}
	ids := make([]int, 0, len(ch.jobs))
	for id := range ch.jobs {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	var b strings.Builder
	for _, id := range ids {
		j := ch.jobs[id]
		fmt.Fprintf(&b, "%s: %s; last: %s\n", j.tool, j.what, j.last)
	}
	return strings.TrimRight(b.String(), "\n")
}
