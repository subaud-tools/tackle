// Command imprint guides a coding agent through npm package work: it logs
// in, publishes, sets trusted publishers, deprecates and confirms CI
// releases, opening each approval for the user and checking every step on
// the registry.
package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/subaud-tools/tackle/internal/integrationsetup"

	"github.com/subaud-tools/tackle/internal/imprint/cli"
)

//go:embed kempt.toml
var declaration []byte

func main() {
	source, err := integrationsetup.Source("imprint", declaration)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(cli.Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, &source))
}
