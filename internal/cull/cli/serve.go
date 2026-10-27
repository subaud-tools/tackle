package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/subaud-tools/tackle/internal/cull/channel"
	"github.com/subaud-tools/tackle/internal/cull/discover"
	"github.com/subaud-tools/tackle/internal/cull/serve"
	"github.com/subaud-tools/tackle/internal/cull/store"
	"github.com/subaud-tools/tackle/internal/version"
	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/pageserve"
)

// shell is cull's serve, channel, watch and report commands: pageserve's,
// with cull's channel, its serve and its review page.
func shell() pageserve.Commands {
	return pageserve.Commands{
		Spec:    channel.Spec(version.Number()),
		Process: pageserve.Process{Live: serve.Live(), Args: []string{"serve", "--foreground", "--no-open"}},
		Group:   "check",
		Serve:   serveForeground,
		Page:    servePage,
	}
}

// client finds cull serve, else starts it from this binary.
func client() *pageserve.Client {
	p := shell().Process
	return pageserve.NewClient(p.Live, func() (pageserve.Entry, error) { return p.Ensure(pageserve.Executable()) })
}

var serveHelp = pageserve.Help{
	Synopsis: "[path] [--no-open] [--port N] [--stop]",
	Summary:  "open the review page (starts the local server if needed)",
	Long: "Starts the review page's server in the background (one per machine) and opens the\n" +
		"project's review. When it is already running, opens the running one. Needs no egress.",
}

var channelHelp = pageserve.Help{
	Summary: "MCP server for an agent session: check, apply, review and Court's answers as events",
	Long: "Runs cull's MCP server on stdin/stdout for an agent session (Claude Code or pi starts it).\n" +
		"It starts cull serve if needed and gives the agent cull_check, cull_apply, cull_check_group,\n" +
		"cull_review and cull_status. Under pi it also registers the session and delivers Court's\n" +
		"answers as events; under Claude Code `cull watch`, which the cull plugin runs, does that.\n" +
		"stdout carries MCP frames only; logs go to stderr.",
}

var watchHelp = pageserve.Help{
	Summary: "print Court's answers as JSON lines: the Claude Code plugin's feed",
	Long: "Registers the session with cull serve and prints each of Court's Sends for it as one JSON\n" +
		"line {\"text\": ..., \"content\": ..., \"meta\": {...}} on stdout. The cull plugin's mod runs\n" +
		"it for a Claude Code session and starts a turn with each line's text. Logs go to stderr.\n" +
		"Runs until signalled. In a session pi's claude-bridge runs, it prints {\"stop\": ...} and\n" +
		"exits 0.",
}

// serveForeground is `cull serve --foreground`: the server over the store, in
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

// servePage resolves `cull serve [path]` to the project's review: the
// fragment that opens it, and a note when it has no review yet.
func servePage(ctx context.Context, args []string) (fragment, note string, err error) {
	if len(args) > 1 {
		return "", "", tools.UsageError{Msg: "serve takes at most one path"}
	}
	path := "."
	if len(args) == 1 {
		path = args[0]
	}
	if _, err := os.Stat(path); err != nil {
		return "", "", tools.Exitf(2, "%v", err)
	}
	root, err := discover.Root(path)
	if err != nil {
		return "", "", tools.Exitf(2, "%v", err)
	}
	st, err := store.Open(ctx, store.Path())
	if err != nil {
		return "", "", err
	}
	defer func() { _ = st.Close() }()
	proj, err := st.Project(ctx, root)
	if err != nil {
		return "", "", err
	}
	_, _, runErr := st.LatestRun(ctx, proj.ID)
	if runErr != nil && !errors.Is(runErr, sql.ErrNoRows) {
		return "", "", runErr
	}
	if runErr != nil {
		note = fmt.Sprintf("no review for %s yet — run cull check", root)
	}
	return "#/p/" + strconv.FormatInt(proj.ID, 10), note, nil
}
