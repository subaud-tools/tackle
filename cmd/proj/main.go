package main

import (
	"os"
	"time"

	"github.com/subaud-tools/tackle/internal/proj"
	"github.com/subaud-tools/tackle/internal/projcli"
)

func main() {
	// tmux servers proj starts inherit this limit; at the default 256 a busy
	// server can't start its hook jobs. proj itself works at a lower limit,
	// so a failure is ignored.
	_ = proj.RaiseFileLimit(8192)
	code := projcli.Dispatch(os.Args[1:], os.Stdout, os.Stderr)
	// A session this command created starts muster and hail in the background;
	// give those starts time to launch before the process ends.
	proj.WaitStarts(2 * time.Second)
	os.Exit(code)
}
