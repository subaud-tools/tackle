package projcli

import (
	"flag"
	"fmt"
	"io"

	"github.com/subaud-tools/tackle/internal/proj"
	"github.com/subaud-tools/tools-common"
)

// windowFlags holds `proj window`'s parsed flags.
type windowFlags struct {
	socket *string
	pane   *string
	right  *string
	keep   *bool
	json   *bool
}

// newWindowFlags is proj window's side-effect-free flag constructor.
func newWindowFlags() (*flag.FlagSet, *windowFlags) {
	fs := flag.NewFlagSet("window", flag.ContinueOnError)
	f := &windowFlags{
		socket: fs.String("socket", "", "the tmux server's socket path (the part of $TMUX before the first comma)"),
		pane:   fs.String("pane", "", "a pane of the session ($TMUX_PANE)"),
		right:  fs.String("right", "", "one of our pages on 127.0.0.1 to open as a tab beside the terminal"),
		keep:   fs.Bool("keep-view", false, "leave a terminal page another terminal detached as it is"),
		json:   fs.Bool("json", false, "emit JSON"),
	}
	tools.SetUsage(fs, "proj window --socket <path> --pane <%id> [--right <url>] [--keep-view] [--json]", windowHelp)
	return fs, f
}

const windowHelp = "Moves the session to its Chrome window, detaching any other terminal, and brings the window " +
	"to the front, or opens a new one on the session's split page: its terminal on the left and its tabs " +
	"on the right. Prints {window,id}: window is opened or raised. --right opens the page as a tab of the " +
	"session's split page, or shows the tab already open on it; only pages on 127.0.0.1 open there. " +
	"With --keep-view a page another terminal detached stays detached. A failure names the step and exits 1."

// cmdWindow implements `proj window`.
func cmdWindow(args []string, out, errw io.Writer) error {
	fs, f := newWindowFlags()
	if err := tools.ParseFlags(fs, args, out); err != nil {
		return err
	}
	if *f.socket == "" || *f.pane == "" {
		return tools.UsageError{Msg: "--socket and --pane are required"}
	}
	show := proj.Window
	if *f.keep {
		show = proj.WindowKeepingView
	}
	res, err := show(*f.socket, *f.pane, *f.right)
	if err != nil {
		return tools.Exitf(1, "%v", err)
	}
	if *f.json {
		return writeJSON(out, res)
	}
	_, err = fmt.Fprintf(out, "%s Chrome window %s\n", res.Window, res.ID)
	return err
}
