package cli

import (
	"path/filepath"
	"testing"

	"github.com/subaud-tools/tools-common/claudeplugin"
)

// The imprint plugin is rendered by tools-common, not written by hand: its
// manifest, the .mcp.json that runs `imprint channel` for the tools, and the
// mod that runs `imprint watch` and starts a turn with each event. Regenerate
// with CLAUDEPLUGIN_UPDATE=1 go test -run '^TestPluginIsCurrent$' ./internal/imprint/cli
func TestPluginIsCurrent(t *testing.T) {
	claudeplugin.AssertCurrent(t, filepath.Join("..", "..", "..", "cmd", "imprint", "plugin"), claudeplugin.Spec{
		Name:        "imprint",
		DisplayName: "imprint",
		Bin:         "imprint",
		ServerName:  "imprint",
		Description: "npm publishing, trusted publishers and deprecation, one step at a time; each step's progress lands in the session that started it.",
		Author:      "Court Schuett",
		AuthorURL:   "https://tackle.tools",
		Homepage:    "https://tackle.tools",
		Repository:  "https://github.com/subaud-tools/tackle",
		License:     "MIT",
		Keywords:    []string{"npm", "publishing", "trusted-publishing"},
	})
}
