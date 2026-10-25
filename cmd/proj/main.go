package main

import (
	"os"

	"github.com/subaud-tools/tackle/internal/proj"
	"github.com/subaud-tools/tackle/internal/projcli"
)

func main() {
	// tmux servers proj starts inherit this limit; at the default 256 a busy
	// server can't start its hook jobs. proj itself works at a lower limit,
	// so a failure is ignored.
	_ = proj.RaiseFileLimit(8192)
	os.Exit(projcli.Dispatch(os.Args[1:], os.Stdout, os.Stderr))
}
