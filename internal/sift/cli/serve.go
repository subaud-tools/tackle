package cli

import (
	"context"
	"io"
	"time"

	"github.com/subaud-tools/tackle/internal/sift/channel"
	"github.com/subaud-tools/tackle/internal/sift/serve"
	"github.com/subaud-tools/tackle/internal/sift/store"
	"github.com/subaud-tools/tackle/internal/version"
	"github.com/subaud-tools/tackle/skills"
	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/pageserve"
)

// shell is sift's serve, channel, watch and report commands: pageserve's,
// with sift's channel, its serve, its review page and its skill.
func shell() pageserve.Commands {
	return pageserve.Commands{
		// The channel's tools find gh through lookPath, the cli's test seam.
		Spec:    channel.Spec(version.Number(), &channel.Env{LookPath: func(file string) (string, error) { return lookPath(file) }}),
		Process: pageserve.Process{Live: serve.Live(), Args: []string{"serve", "--foreground", "--no-open"}},
		Group:   "review",
		Serve:   serveForeground,
		Page:    servePage,
		Skill:   skill,
	}
}

// skill is the sift skill this binary was built with.
func skill() ([]byte, error) { return skills.FS.ReadFile("sift/SKILL.md") }

// client finds sift serve, else starts it from this binary.
func client() *pageserve.Client {
	p := shell().Process
	return pageserve.NewClient(p.Live, func() (pageserve.Entry, error) { return p.Ensure(pageserve.Executable()) })
}

var serveHelp = pageserve.Help{
	Synopsis: "[--no-open] [--port N] [--stop]",
	Summary:  "open the review page (starts the local server if needed)",
	Long: "Starts the review page's server in the background (one per machine) and opens the latest\n" +
		"round's review. When it is already running, opens the running one. Needs no network.",
}

var channelHelp = pageserve.Help{
	Summary: "MCP server for an agent session: check, next, propose, review, apply, status, and Send as events",
	Long: "Runs sift's MCP server on stdin/stdout for an agent session (Claude Code or pi starts it).\n" +
		"It starts sift serve if needed and gives the agent sift_check, sift_review, sift_apply and\n" +
		"sift_status. Under pi it also registers the session and delivers the user's decisions as\n" +
		"events; under Claude Code `sift watch`, which the sift plugin runs, does that.\n" +
		"stdout carries MCP frames only; logs go to stderr.",
}

var watchHelp = pageserve.Help{
	Summary: "print the user's Sends as JSON lines: the Claude Code plugin's feed",
	Long: "Registers the session with sift serve and prints each of the user's Sends for it as one\n" +
		"JSON line {\"text\": ..., \"content\": ..., \"meta\": {...}} on stdout. The sift plugin's mod\n" +
		"runs it for a Claude Code session and starts a turn with each line's text. Logs go to\n" +
		"stderr. Runs until signalled. In a session pi's claude-bridge runs, it prints\n" +
		"{\"stop\": ...} and exits 0.",
}

// serveForeground is `sift serve --foreground`: the server over the store, in
// this process until ctx ends.
func serveForeground(ctx context.Context, port int, errw io.Writer) error {
	st, err := store.Open(ctx, store.Path())
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()
	return serve.Run(ctx, st, serve.Options{Port: port, Idle: 8 * time.Hour, Version: version.Number(), Log: errw,
		Ready: func(string) {}})
}

// servePage is `sift serve`'s page: the latest round's review, which is the
// page's own start, and a note when there is no round yet.
func servePage(ctx context.Context, args []string) (fragment, note string, err error) {
	if len(args) > 0 {
		return "", "", tools.UsageError{Msg: "serve takes no arguments"}
	}
	if st, err := store.Open(ctx, store.Path()); err == nil {
		if _, _, err := st.Latest(ctx); err != nil {
			note = "no round yet: run sift check"
		}
		_ = st.Close()
	}
	return "", note, nil
}
