package projcli

import (
	"flag"
	"fmt"
	"io"

	"github.com/subaud-tools/tackle/internal/proj"
	"github.com/subaud-tools/tools-common"
)

// splitTab is proj.Split; a seam for tests.
var splitTab = proj.Split

// splitFlags holds `proj split`'s parsed flags.
type splitFlags struct {
	socket  *string
	session *string
	right   *string
	json    *bool
}

// newSplitFlags is proj split's side-effect-free flag constructor.
func newSplitFlags() (*flag.FlagSet, *splitFlags) {
	fs := flag.NewFlagSet("split", flag.ContinueOnError)
	f := &splitFlags{
		socket:  fs.String("socket", "", "the tmux server's socket path, as `proj new --json` prints it"),
		session: fs.String("session", "", "the session's name"),
		right:   fs.String("right", "", "one of our pages, on 127.0.0.1 or this machine's public page address, to open as a tab beside the terminal"),
		json:    fs.Bool("json", false, "emit JSON"),
	}
	tools.SetUsage(fs, "proj split --socket <path> --session <name> --right <url> [--json]", splitHelp)
	return fs, f
}

const splitHelp = "Opens --right as a tab of the session's split page, its terminal on the left and its tabs on " +
	"the right, through proj serve (started when it is not running), with no Chrome. Prints the split page's " +
	"address, {split} with --json: proj's public address when PAGE_PUBLIC_DOMAIN is set, else its 127.0.0.1 one. " +
	"--right is a page on 127.0.0.1, or https://<tool>.<PAGE_PUBLIC_DOMAIN>. A session that does not exist, or " +
	"a page that is refused, exits 1 with the reason."

// cmdSplit implements `proj split`.
func cmdSplit(args []string, out, errw io.Writer) error {
	fs, f := newSplitFlags()
	if err := tools.ParseFlags(fs, args, out); err != nil {
		return err
	}
	if *f.socket == "" || *f.session == "" || *f.right == "" {
		return tools.UsageError{Msg: "--socket, --session and --right are required"}
	}
	page, err := splitTab(*f.socket, *f.session, *f.right)
	if err != nil {
		return tools.Exitf(1, "%v", err)
	}
	if *f.json {
		return writeJSON(out, proj.SplitResult{Split: page})
	}
	_, err = fmt.Fprintln(out, page)
	return err
}
