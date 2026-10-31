package projcli

import (
	"context"
	"flag"
	"io"
	"os/signal"
	"syscall"

	"github.com/subaud-tools/tackle/internal/proj/pages"
	"github.com/subaud-tools/tools-common"
)

// newPagesFlags is proj pages' side-effect-free flag constructor.
func newPagesFlags() *flag.FlagSet {
	fs := flag.NewFlagSet("pages", flag.ContinueOnError)
	tools.SetUsage(fs, "proj pages",
		"Serves the local page that shows any session's terminal, until killed.")
	return fs
}

// cmdPages implements `proj pages`.
func cmdPages(args []string, out, errw io.Writer) error {
	if err := tools.ParseFlags(newPagesFlags(), args, out); err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()
	if err := pages.Run(ctx, out); err != nil {
		return tools.Exitf(1, "%v", err)
	}
	return nil
}
