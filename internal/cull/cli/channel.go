package cli

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/subaud-tools/tackle/internal/cull/channel"
	"github.com/subaud-tools/tackle/internal/cull/serve"
	"github.com/subaud-tools/tackle/internal/version"
	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/channelmcp"
)

var channelFlags = flags("channel", "cull channel",
	"Runs cull's MCP server on stdin/stdout for an agent session (Claude Code or pi starts it).\n"+
		"It starts cull serve if needed and gives the agent cull_check, cull_apply, cull_check_group,\n"+
		"cull_review and cull_status. Under pi it also registers the session and delivers Court's\n"+
		"answers as events; under Claude Code `cull watch`, which the cull plugin runs, does that.\n"+
		"stdout carries MCP frames only; logs go to stderr.",
	nil)

func runChannel(stdin io.Reader) func(args []string, out, errw io.Writer) error {
	return func(args []string, out, errw io.Writer) error {
		fs := channelFlags()
		pos, err := parse(fs, args, out)
		if err != nil {
			return err
		}
		if len(pos) > 0 {
			return tools.UsageError{Msg: "channel takes no arguments"}
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		client := channel.NewClient(func() (serve.Advert, error) { return serve.Start(executable(), 0) })
		ch := channel.ForSession(client, version.Number())
		ch.Log = errw
		return ch.Run(ctx, stdin, out)
	}
}

var watchFlags = flags("watch", "cull watch",
	"Registers the session with cull serve and prints each of Court's Sends for it as one JSON\n"+
		"line {\"text\": ..., \"content\": ..., \"meta\": {...}} on stdout. The cull plugin's mod runs\n"+
		"it for a Claude Code session and starts a turn with each line's text. Logs go to stderr.\n"+
		"Runs until signalled. In a session pi's claude-bridge runs, it prints {\"stop\": ...} and\n"+
		"exits 0.",
	nil)

// runWatch is `cull watch`: the channel's presence and wait loop, printing
// each Send as one JSON line on stdout for the Claude Code plugin's mod.
func runWatch(args []string, out, errw io.Writer) error {
	fs := watchFlags()
	pos, err := parse(fs, args, out)
	if err != nil {
		return err
	}
	if len(pos) > 0 {
		return tools.UsageError{Msg: "watch takes no arguments"}
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	client := channel.NewClient(func() (serve.Advert, error) { return serve.Start(executable(), 0) })
	return watch(ctx, out, errw, client)
}

// watch is runWatch after its flags, until ctx ends, with the client to
// serve passed in so a test can run it.
func watch(ctx context.Context, out, errw io.Writer, client *channel.Client) error {
	// In a bridge child this prints the stop line and the mod never restarts
	// watch: the hosting pi session's channel delivers Court's answers.
	lines, run, err := channelmcp.StartWatch(out, "cull")
	if !run || err != nil {
		return err
	}
	ch := channel.New(channel.FromEnv(), client, version.Number())
	ch.Log = errw
	ch.Push = lines
	return ch.Watch(ctx)
}
