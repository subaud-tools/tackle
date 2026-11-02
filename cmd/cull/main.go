// Command cull judges whether tests are worth keeping, using TypeSafe's Jev
// classifier, and measures its rubrics against a labeled corpus.
package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/subaud-tools/tackle/internal/integrationsetup"

	"github.com/subaud-tools/tackle/internal/cull/cli"
)

//go:embed kempt.toml
var declaration []byte

func main() {
	source, err := integrationsetup.Source("cull", declaration)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(cli.Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, &source))
}
