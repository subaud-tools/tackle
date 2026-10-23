package cli

import (
	"path/filepath"
	"testing"

	"github.com/subaud-tools/tools-common/claudeplugin"
)

// The sift plugin is rendered by tools-common, not written by hand: its
// manifest, the .mcp.json that runs `sift channel` for the tools, and the mod
// that runs `sift watch` and starts a turn with each event. Regenerate with
// CLAUDEPLUGIN_UPDATE=1 go test -run '^TestPluginIsCurrent$' ./internal/sift/cli
func TestPluginIsCurrent(t *testing.T) {
	claudeplugin.AssertCurrent(t, filepath.Join("..", "..", "..", "cmd", "sift", "plugin"), claudeplugin.Spec{
		Name:        "sift",
		DisplayName: "sift",
		Bin:         "sift",
		ServerName:  "sift",
		Description: "Audit agent instruction files; the user's decisions from sift's review page land in the session that asked.",
		Author:      "Court Schuett",
		AuthorURL:   "https://tackle.tools",
		Homepage:    "https://tackle.tools",
		Repository:  "https://github.com/subaud-tools/tackle",
		License:     "MIT",
		Keywords:    []string{"instructions", "agents", "review"},
	})
}
