package cli

import (
	"io"

	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/pageserve"
)

func commands(stdin io.Reader) []tools.Command {
	return []tools.Command{
		{
			Name: "init", Group: "setup", Synopsis: "[--root DIR]… [--yes] [--force]",
			Summary:  "detect the agent harnesses you use and write sift's config (profiles, roots, budgets)",
			NewFlags: initFlags, Run: runInit(stdin),
		},
		pageserve.SkillsCommand("sift", skill),
		{
			Name: "doctor", Group: "setup", Synopsis: "",
			Summary:  "check what sift needs: config, profiles, roots, state; gh, the page server and channels are optional",
			NewFlags: doctorFlags, Run: runDoctor,
		},
		{
			Name: "check", Group: "audit", Synopsis: "[--json]",
			Summary:  "audit the instruction files and record the round's findings as rows",
			NewFlags: checkFlags, Run: runCheck,
		},
		{
			Name: "rows", Group: "audit", Synopsis: "add [--round N] < rows.jsonl", Subcommands: []string{"add"},
			Summary:  "merge the agent's backlog and intake rows (JSON lines on stdin) into their round",
			NewFlags: rowsFlags, Run: runRows(stdin),
		},
		{
			Name: "next", Group: "audit", Synopsis: "[--round N] [--json]",
			Summary:  "the next file to recommend: its content at the audit, its findings and the guidance",
			NewFlags: nextFlags, Run: runNext,
		},
		{
			Name: "propose", Group: "audit", Synopsis: "[--round N] < recommendation.json",
			Summary:  "store file recommendations (JSON on stdin), checked all or nothing",
			NewFlags: proposeFlags, Run: runPropose(stdin),
		},
		shell().ServeCommand(serveHelp),
		{
			Name: "wait", Group: "review", Synopsis: "[--timeout D]",
			Summary:  "wait until the user presses Send, then print what to do (for sessions without the sift channel)",
			NewFlags: waitFlags, Run: runWait,
		},
		shell().ChannelCommand(stdin, channelHelp),
		shell().WatchCommand(watchHelp),
		shell().ReportCommand(),
		{
			Name: "apply", Group: "review", Synopsis: "[--round N] [--dry-run] [--json]",
			Summary:  "write each approved file whole, a branch per repo; a pull request with gh",
			NewFlags: applyFlags, Run: runApply,
		},
		{
			Name: "clean", Group: "review", Synopsis: "[--dry-run] [--json]",
			Summary:  "remove the branches and worktrees apply recorded, once their pull requests are merged or closed",
			NewFlags: cleanFlags, Run: runClean,
		},
		{
			Name: "reconcile", Group: "review", Synopsis: "[--round N] [--json]",
			Summary:  "check each apply branch: every file as approved, nothing else changed",
			NewFlags: reconcileFlags, Run: runReconcile,
		},
	}
}
