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
	json   *bool
}

// newWindowFlags is proj window's side-effect-free flag constructor.
func newWindowFlags() (*flag.FlagSet, *windowFlags) {
	fs := flag.NewFlagSet("window", flag.ContinueOnError)
	f := &windowFlags{
		socket: fs.String("socket", "", "the tmux server's socket path (the part of $TMUX before the first comma)"),
		pane:   fs.String("pane", "", "a pane of the session ($TMUX_PANE)"),
		right:  fs.String("right", "", "the URL for the split's right half when a window is opened"),
		json:   fs.Bool("json", false, "emit JSON"),
	}
	tools.SetUsage(fs, "proj window --socket <path> --pane <%id> [--right <url>] [--json]", windowHelp)
	return fs, f
}

const windowHelp = "Brings the session's Chrome window to the front, or opens a new one with the session's " +
	"terminal on the left of a split and --right (or an empty tab) on the right, and prints " +
	"{window,id}: window is opened or raised. A failure names the step and exits 1; the split " +
	"needs macOS Accessibility permission for the app running proj."

// cmdWindow implements `proj window`.
func cmdWindow(args []string, out, errw io.Writer) error {
	fs, f := newWindowFlags()
	if err := tools.ParseFlags(fs, args, out); err != nil {
		return err
	}
	if *f.socket == "" || *f.pane == "" {
		return tools.UsageError{Msg: "--socket and --pane are required"}
	}
	res, err := proj.Window(*f.socket, *f.pane, *f.right)
	if err != nil {
		return tools.Exitf(1, "%v", err)
	}
	if *f.json {
		return writeJSON(out, res)
	}
	_, err = fmt.Fprintf(out, "%s Chrome window %s\n", res.Window, res.ID)
	return err
}
