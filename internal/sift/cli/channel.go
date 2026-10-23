package cli

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/subaud-tools/tackle/internal/sift/channel"
	"github.com/subaud-tools/tackle/internal/sift/serve"
	"github.com/subaud-tools/tackle/internal/version"
	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/channelmcp"
)

var channelFlags = flags("channel", "sift channel",
	"Runs sift's MCP server on stdin/stdout for an agent session (Claude Code or pi starts it).\n"+
		"It starts sift serve if needed and gives the agent sift_check, sift_review, sift_apply and\n"+
		"sift_status. Under pi it also registers the session and delivers the user's decisions as\n"+
		"events; under Claude Code `sift watch`, which the sift plugin runs, does that.\n"+
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
		ch.LookPath = lookPath
		return ch.Run(ctx, stdin, out)
	}
}

var watchFlags = flags("watch", "sift watch",
	"Registers the session with sift serve and prints each of the user's Sends for it as one\n"+
		"JSON line {\"content\": ..., \"meta\": {...}} on stdout. The sift plugin's mod runs it for\n"+
		"a Claude Code session and starts a turn with each line. Logs go to stderr. Runs until\n"+
		"signalled.",
	nil)

// runWatch is `sift watch`: the channel's presence and wait loop, printing
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
	ch := channel.New(channel.FromEnv(), client, version.Number())
	ch.Log = errw
	ch.Push = channelmcp.NewLines(out)
	return ch.Watch(ctx)
}
