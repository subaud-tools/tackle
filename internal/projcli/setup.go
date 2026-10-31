package projcli

import (
	_ "embed"

	"github.com/subaud-tools/tackle/internal/version"
	"github.com/subaud-tools/tackle/shell"
	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/setupbundle"
)

//go:embed setup.toml
var setupManifest []byte

func setupCommand() tools.Command {
	return setupbundle.BundleCommand("proj", version.Number(), setupbundle.Source{
		Manifest: setupManifest,
		Files:    shell.Proj(),
	})
}
