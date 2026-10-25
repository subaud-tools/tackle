package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/subaud-tools/tackle/internal/sift/channel"
	"github.com/subaud-tools/tackle/internal/sift/serve"
	"github.com/subaud-tools/tackle/internal/version"
	"github.com/subaud-tools/tackle/skills"
	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/channelmcp"
	commonskills "github.com/subaud-tools/tools-common/skills"
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
		refreshSkill(errw)
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		client := channel.NewClient(func() (serve.Advert, error) { return serve.Start(executable(), 0) })
		ch := channel.ForSession(client, version.Number())
		ch.Log = errw
		ch.LookPath = lookPath
		return ch.Run(ctx, stdin, out)
	}
}

// refreshSkill rewrites each installed copy of the sift skill that differs
// from the one this binary was built with, so the skill an agent reads
// matches the binary it calls. It logs each path it rewrote; an error is
// logged and the channel starts anyway.
func refreshSkill(errw io.Writer) {
	body, err := skills.FS.ReadFile("sift/SKILL.md")
	if err == nil {
		var paths []string
		paths, err = commonskills.Refresh("sift", body)
		for _, p := range paths {
			_, _ = fmt.Fprintf(errw, "sift channel: refreshed the sift skill at %s\n", p)
		}
	}
	if err != nil {
		_, _ = fmt.Fprintf(errw, "sift channel: refreshing the sift skill: %v\n", err)
	}
}

var watchFlags = flags("watch", "sift watch",
	"Registers the session with sift serve and prints each of the user's Sends for it as one\n"+
		"JSON line {\"text\": ..., \"content\": ..., \"meta\": {...}} on stdout. The sift plugin's mod\n"+
		"runs it for a Claude Code session and starts a turn with each line's text. Logs go to\n"+
		"stderr. Runs until signalled. In a session pi's claude-bridge runs, it prints\n"+
		"{\"stop\": ...} and exits 0.",
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
	return watch(ctx, out, errw, client)
}

// watch is runWatch after its flags, until ctx ends, with the client to
// serve passed in so a test can run it.
func watch(ctx context.Context, out, errw io.Writer, client *channel.Client) error {
	// In a bridge child this prints the stop line and the mod never restarts
	// watch: the hosting pi session's channel delivers the user's decisions.
	lines, run, err := channelmcp.StartWatch(out, "sift")
	if !run || err != nil {
		return err
	}
	ch := channel.New(channel.FromEnv(), client, version.Number())
	ch.Log = errw
	ch.Push = lines
	return ch.Watch(ctx)
}
