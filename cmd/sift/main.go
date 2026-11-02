// Command sift keeps agent instruction files lean: it audits global, repo and
// skill files for size, duplicates, stale status, dead paths and the like,
// and records each round's findings as rows for an agent and a person to
// decide.
package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/subaud-tools/tackle/internal/integrationsetup"

	"github.com/subaud-tools/tackle/internal/sift/cli"
)

//go:embed kempt.toml
var declaration []byte

func main() {
	source, err := integrationsetup.Source("sift", declaration)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(cli.Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, &source))
}
