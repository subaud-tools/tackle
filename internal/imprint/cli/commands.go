package cli

import (
	"io"

	"github.com/subaud-tools/tackle/skills"
	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/channelmcp"
	"github.com/subaud-tools/tools-common/pageserve"
)

func commands(stdin io.Reader) []tools.Command {
	return []tools.Command{
		pageserve.SkillsCommand("imprint", func() ([]byte, error) { return skills.FS.ReadFile("imprint/SKILL.md") }),
		{
			Name: "channel", Group: "agent", Synopsis: "",
			Summary:  "MCP server for an agent session: npm login, show, publish, trust, deprecate, expect and status, with progress as events",
			NewFlags: channelFlags, Run: runChannel(stdin),
		},
		{
			Name: "watch", Group: "agent", Synopsis: "",
			Summary:  "print imprint's messages for this Claude Code session as JSON lines: the plugin's feed",
			NewFlags: watchFlags, Run: runWatch,
		},
		channelmcp.ReportCommand(deliverReport),
	}
}

// deliverReport is `imprint report`'s delivery, which the Claude Code plugin's
// mod runs. It does nothing: tackle's tools set no report token
// (channelmcp.MetaReport), so there is nothing to deliver. It is mounted so
// the regenerated mod never calls a command the binary lacks.
func deliverReport(string, channelmcp.ReportEvent, string) error { return nil }
