// Package channel is `cull channel`: the MCP server each agent session starts
// (Claude Code through the cull plugin, pi through channels.tools). It starts
// cull serve when it isn't running, gives the agent cull's tools and, under
// pi, registers the session's presence, holds a long poll open for Court's
// answers and injects each Send as a notifications/claude/channel event.
// Under Claude Code `cull watch` runs that loop instead and prints each Send
// as a JSON line for the plugin's mod. The channel, its client and its loop
// are tools-common's pageserve; this package is cull's part of it: its tools,
// its instructions and how it reads serve's answers.
package channel

import (
	"encoding/json"
	"strconv"

	"github.com/subaud-tools/tackle/internal/cull/discover"
	"github.com/subaud-tools/tools-common/channelmcp"
	"github.com/subaud-tools/tools-common/pageserve"
)

// Instructions are the standing rules the agent gets with the channel.
const Instructions = `cull judges automated tests with Jev and removes the ones that don't earn their place. Court reviews only what Jev is unsure about, on cull's review page.
- After you write or change tests, run cull_check with base set to the merge base of your branch (or no base for the whole suite when Court asks for it).
- Then run cull_apply with no ids: it removes every cut, both Jev's confident ones and Court's, runs the tests before and after, and puts everything back if they break. Never commit on cull's behalf; report what was removed.
- If cull_check reports to_review above 0, run cull_review once to open the page for Court, tell him how many items wait, and carry on with other work. Never ask Court about a test anywhere but the page, and never open the page repeatedly.
- Court's answers arrive as an event from cull: <channel source="cull" project="…" send_id="…"> in pi, or a turn whose first line is cull event project=… send_id=… in Claude Code. On a send that takes more than a moment, first call cull_reply with its send_id, state working and one line on what you are doing. Follow its steps: cull_check, then cull_apply for his cuts; rewrite each group he chose to merge as one table test that keeps every row the event lists, then run cull_check_group on it and fix what it reports. When done, cull_reply with answered and a short report, or declined or failed with why.
- cull_status tells you what is open, answered and sent for this project, and which sends are not settled yet.
- When your session runs in a folder that holds several repositories (a workspace), pass the repository's path to every cull tool. The event from Court names the repository's path; use it.
- You never answer for Court. If cull reports an error, say what failed; don't work around it by editing tests by hand.
- When Court asks why tests are slow, or checks take long, run cull_time and fix what it and cull_check's speed findings point at; then run cull_time again and report before and after.`

// Channel is a running cull channel: pageserve's, with cull's Spec.
type Channel = pageserve.Channel

// Spec is cull's channel: its tools, its instructions, the project root
// discover finds, and serve's answers as events with source "cull".
func Spec(version string) pageserve.Spec {
	return pageserve.Spec{
		Tool: "cull", Version: version, Instructions: Instructions, Tools: Tools(), Call: Call,
		Root: discover.Root, Wake: wake, Answers: "Court's answers",
	}
}

// wake reads one answer of GET /api/agent/wait into the event pushed to the
// agent: the Send's text, with its project, send id and report token.
func wake(body json.RawMessage) (string, map[string]string, bool) {
	var got struct {
		ID     int64  `json:"id"`
		Root   string `json:"root"`
		Text   string `json:"text"`
		Report string `json:"report"`
	}
	if json.Unmarshal(body, &got) != nil || got.ID == 0 {
		return "", nil, false
	}
	return got.Text, map[string]string{"source": "cull", "project": got.Root, "send_id": strconv.FormatInt(got.ID, 10),
		channelmcp.MetaReport: got.Report}, true
}
