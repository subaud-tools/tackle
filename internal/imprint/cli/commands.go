package cli

import (
	"io"

	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/channelmcp"
)

func commands(stdin io.Reader) []tools.Command {
	return []tools.Command{
		{
			Name: "skills", Group: "setup", Synopsis: "install [--agent NAME]…", Subcommands: []string{"install"},
			Summary:  "install the imprint skill for each agent harness, the version this binary was built with",
			NewFlags: skillsFlags, Run: runSkills,
		},
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
