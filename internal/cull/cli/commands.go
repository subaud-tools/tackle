package cli

import (
	"io"

	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/channelmcp"
)

func commands(stdin io.Reader) []tools.Command {
	return []tools.Command{
		{
			Name: "check", Group: "check", Synopsis: "[path] [--diff base]",
			Summary:  "judge tests and near-duplicate groups (suite or diff), write .cull/last.json",
			NewFlags: checkFlags, Run: runCheck(stdin),
		},
		{
			Name: "apply", Group: "check", Synopsis: "--ids <id>… | --verdict cut",
			Summary:  "remove tests judged cut, tidy imports, verify; roll back on failure",
			NewFlags: applyFlags, Run: runApply,
		},
		{
			Name: "time", Group: "check", Synopsis: "[path] [--timeout D] [--json]",
			Summary:  "run the project's checks with timing and report where the time goes",
			NewFlags: timeFlags, Run: runTime,
		},
		{
			Name: "serve", Group: "check", Synopsis: "[path] [--no-open] [--port N] [--stop]",
			Summary:  "open the review page (starts the local server if needed)",
			NewFlags: serveFlags, Run: runServe,
		},
		{
			Name: "channel", Group: "check", Synopsis: "",
			Summary:  "MCP server for an agent session: check, apply, review and Court's answers as events",
			NewFlags: channelFlags, Run: runChannel(stdin),
		},
		{
			Name: "watch", Group: "check", Synopsis: "",
			Summary:  "print Court's answers as JSON lines: the Claude Code plugin's feed",
			NewFlags: watchFlags, Run: runWatch,
		},
		channelmcp.ReportCommand(deliverReport),
		{
			Name: "wait", Group: "check", Synopsis: "[path] [--timeout D]",
			Summary:  "wait until Court sends his answers for this project, then print what to do (for sessions without the cull channel)",
			NewFlags: waitFlags, Run: runWait,
		},
		{
			Name: "init", Group: "setup", Synopsis: "[path]",
			Summary:  "store your TypeSafe key, allow egress for the project, ignore .cull/ (terminal only)",
			NewFlags: initFlags, Run: runInit(stdin),
		},
		{
			Name: "doctor", Group: "setup", Synopsis: "[path]",
			Summary:  "check what cull needs: key, egress, python3, node, test command, serve, agent setup",
			NewFlags: doctorFlags, Run: runDoctor,
		},
		{
			Name: "judge", Group: "judge", Synopsis: "[file|-]",
			Summary:  "judge tests or duplicate groups with Jev (plumbing; JSONL in, JSONL out)",
			NewFlags: judgeFlags, Run: runJudge(stdin),
		},
	}
}

// deliverReport is `cull report`'s delivery, which the Claude Code plugin's
// mod runs. It does nothing: tackle's tools set no report token
// (channelmcp.MetaReport), so there is nothing to deliver. It is mounted so
// the regenerated mod never calls a command the binary lacks.
func deliverReport(string, channelmcp.ReportEvent, string) error { return nil }
