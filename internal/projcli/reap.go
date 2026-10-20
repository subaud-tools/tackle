package projcli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/subaud-tools/tackle/internal/proj"
	tools "github.com/subaud-tools/tools-common"
)

// newReapFlags is proj reap's side-effect-free flag constructor, shared by
// cmdReap (actual parsing) and the registry's NewFlags (-h/help and man).
func newReapFlags() (*flag.FlagSet, *string) {
	fs := flag.NewFlagSet("reap", flag.ContinueOnError)
	socket := fs.String("socket", "", "tmux socket name")
	tools.SetUsage(fs, "reap [<session>] [--socket S]",
		"Forgets a session from proj's saved record and kills its tmux session if it is running, "+
			"as the picker's ^x ^x does. It may reap the session it runs in. Socket resolves from "+
			"--socket, then the ambient server, then a server search by session name. A bare "+
			"`proj reap` opens the picker to reap or forget sessions (^x ^x); it never opens one.")
	return fs, socket
}

// cmdReap implements `proj reap <session> [--socket S]`; Dispatch routes a
// bare `proj reap` to the reap picker before it gets here.
func cmdReap(args []string, out, errw io.Writer) error {
	fs, socketFlag := newReapFlags()
	var positional []string
	rest := args
	for {
		// --help prints usage and returns flag.ErrHelp, which Dispatch maps to a
		// clean exit.
		if err := tools.ParseFlags(fs, rest, os.Stdout); err != nil {
			return err
		}
		rest = fs.Args()
		if len(rest) == 0 {
			break
		}
		positional = append(positional, rest[0])
		rest = rest[1:]
	}
	if len(positional) != 1 {
		return tools.UsageError{Msg: "usage: proj reap <session> [--socket S]"}
	}
	name := positional[0]

	socket := *socketFlag
	if socket == "" {
		socket = proj.CurrentServer()
	}
	if socket == "" {
		if s, ok := proj.FindServer(name); ok {
			socket = s
		}
	}
	live := socket != "" && hasSession(socket, name)
	// The ambient server may not be the one hosting name: search the others.
	if !live && *socketFlag == "" {
		if s, ok := proj.FindServer(name); ok {
			socket, live = s, true
		}
	}

	rec, err := proj.LoadRecord()
	if err != nil {
		return tools.Exitf(1, "%v", err)
	}
	if _, saved := rec.Sessions[name]; !saved && !live {
		return tools.Exitf(1, "no session %q", name)
	}

	// Forget first: run from a pane inside this session, the kill ends this process.
	if err := proj.Forget(name); err != nil {
		return tools.Exitf(1, "forget: %v", err)
	}
	if !live {
		_, _ = fmt.Fprintln(out, "forgot "+name)
		return nil
	}
	_, _ = fmt.Fprintln(out, "reaped "+name)
	if err := proj.KillSession(socket, name); err != nil {
		return tools.Exitf(1, "reap: %v", err)
	}
	return nil
}

// hasSession reports whether socket hosts a session named exactly name.
func hasSession(socket, name string) bool {
	_, err := proj.Run(socket, "has-session", "-t", "="+name)
	return err == nil
}
