package projcli

import (
	"context"
	"io"

	"github.com/subaud-tools/tackle/internal/proj/pages"
	"github.com/subaud-tools/tackle/internal/version"
	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/pageserve"
)

// serveCommand is `proj serve`: pageserve's, over proj's page server.
func serveCommand() tools.Command {
	return pageserve.Commands{
		Spec:    pageserve.Spec{Tool: "proj", Version: version.Number()},
		Process: pages.Process(),
		Serve:   serveForeground,
	}.ServeCommand(serveHelp)
}

var serveHelp = pageserve.Help{
	Synopsis: "[--no-open] [--port N] [--stop]",
	Summary:  "serve the terminal pages of Chrome session windows",
	Long: "Starts proj's page server in the background (one per machine), which shows each session's\n" +
		"terminal page for its Chrome window, and opens it; when it is already running, opens the\n" +
		"running one. proj window starts it when it is not running. --stop stops it.",
}

// serveForeground is `proj serve --foreground`: the server in this process
// until ctx ends.
func serveForeground(ctx context.Context, port int, errw io.Writer) error {
	return pages.Run(ctx, pages.Options{Port: port, Version: version.Number(), Log: errw})
}
