package cli

import (
	"path/filepath"
	"testing"

	"github.com/subaud-tools/tools-common/claudeplugin"
)

// The cull plugin is rendered by tools-common, not written by hand: its
// manifest, the .mcp.json that runs `cull channel` for the tools, and the mod
// that runs `cull watch` and starts a turn with each event. Regenerate with
// CLAUDEPLUGIN_UPDATE=1 go test -run '^TestPluginIsCurrent$' ./internal/cull/cli
func TestPluginIsCurrent(t *testing.T) {
	claudeplugin.AssertCurrent(t, filepath.Join("..", "..", "..", "cmd", "cull", "plugin"), claudeplugin.Spec{
		Name:        "cull",
		DisplayName: "cull",
		Bin:         "cull",
		ServerName:  "cull",
		Description: "Judge your tests with Jev; Court's answers from cull's review page land in the session that asked.",
		Author:      "Court Schuett",
		AuthorURL:   "https://tackle.tools",
		Homepage:    "https://tackle.tools",
		Repository:  "https://github.com/subaud-tools/tackle",
		License:     "MIT",
		Keywords:    []string{"tests", "testing", "review"},
	})
}
