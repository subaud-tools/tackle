// Package channel is `sift channel`: the MCP server each agent session
// starts (Claude Code through the sift plugin, pi through channels.tools). It
// starts sift serve when it isn't running, gives the agent sift's tools and,
// under pi, registers the session's presence, holds a long poll open for the
// user's decisions and injects each Send as a notifications/claude/channel
// event. Under Claude Code `sift watch` runs that loop instead and prints
// each Send as a JSON line for the plugin's mod. The channel, its client and
// its loop are tools-common's pageserve; this package is sift's part of it:
// its tools, its instructions and how it reads serve's answers.
package channel

import (
	"context"
	"encoding/json"
	"os/exec"
	"strconv"

	"github.com/subaud-tools/tackle/internal/sift/config"
	"github.com/subaud-tools/tools-common/channelmcp"
	"github.com/subaud-tools/tools-common/pageserve"
)

// Instructions are the standing guidance the agent gets with the channel.
const Instructions = `sift audits the instruction files coding agents read (global files, repo files, skills). For each file with findings you recommend one revised version, and the user accepts, edits or rejects each file on a review page.
- Run sift_check when the user asks for an audit. It records a round and says how many files need a recommendation.
- Recommend every file: sift_next gives one file with its content, findings and the guidance for the rewrite; sift_propose stores the whole revised file. Repeat until sift_next says done. Recommend files that move text between them in one sift_propose call, each linking the other.
- Then run sift_review once to open the page for the user, tell them how many files wait, and carry on with other work. The page opens only once every file has a recommendation.
- The user's decisions arrive as an event from sift: <channel source="sift" round="…" send_id="…"> in pi, or a turn whose first line is sift event round=… send_id=… in Claude Code. On a send that takes more than a moment, first call sift_reply with its send_id, state working and one line on what you are doing. Follow its steps: sift_apply writes a branch per repo; then run ` + "`sift reconcile`" + ` and review each branch. When done, sift_reply with answered and a short report, or declined or failed with why.
- The pull requests are the user's to merge. Once they are merged or closed, run sift_clean to remove the branches and worktrees apply made.
- sift_status tells you the round's state, what is open, decided and sent, and which sends are not settled yet.
- When sift reports an error, say what failed and leave the files as they are.`

// Channel is a running sift channel: pageserve's, with sift's Spec.
type Channel = pageserve.Channel

// Env is what sift's tools reach outside serve: LoadConfig reads sift's
// config, LookPath finds gh. A nil field is the real one; tests replace them.
type Env struct {
	LoadConfig func() (config.Config, error)
	LookPath   func(string) (string, error)
}

// Spec is sift's channel: its tools over env, its instructions, no project
// (sift's rounds are the machine's), and serve's answers as events with
// source "sift". A nil env is the real one.
func Spec(version string, env *Env) pageserve.Spec {
	if env == nil {
		env = &Env{}
	}
	if env.LoadConfig == nil {
		env.LoadConfig = func() (config.Config, error) { return config.Load(config.Path()) }
	}
	if env.LookPath == nil {
		env.LookPath = exec.LookPath
	}
	return pageserve.Spec{
		Tool: "sift", Version: version, Instructions: Instructions, Tools: Tools(),
		Call: func(ctx context.Context, ch *Channel, name string, args json.RawMessage) (string, error) {
			return calls{ch, env}.Call(ctx, name, args)
		},
		Wake: wake, Answers: "the user's decisions",
	}
}

// wake reads one answer of GET /api/agent/wait into the event pushed to the
// agent: the Send's text, with its round, send id and report token.
func wake(body json.RawMessage) (string, map[string]string, bool) {
	var got struct {
		ID     int64  `json:"id"`
		Round  int64  `json:"round"`
		Text   string `json:"text"`
		Report string `json:"report"`
	}
	if json.Unmarshal(body, &got) != nil || got.ID == 0 {
		return "", nil, false
	}
	return got.Text, map[string]string{"source": "sift", "round": strconv.FormatInt(got.Round, 10), "send_id": strconv.FormatInt(got.ID, 10),
		channelmcp.MetaReport: got.Report}, true
}
