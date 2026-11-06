package projcli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/subaud-tools/tackle/internal/proj"
	tools "github.com/subaud-tools/tools-common"
)

// Seams: the live sessions and the platform (Ghostty exists only on macOS).
var (
	liveSessions = proj.LiveSessions
	restoreGOOS  = runtime.GOOS
)

// liveJSON is one running session as `proj live --json` prints it and
// `proj restore --from` reads it: what RestoreSession needs to reopen it.
type liveJSON struct {
	Socket       string `json:"socket"`
	Name         string `json:"name"`
	Dir          string `json:"dir"`
	Agent        string `json:"agent"`
	Conversation string `json:"conversation"`
}

func newLiveFlags() (*flag.FlagSet, *listFlags) {
	fs := flag.NewFlagSet("live", flag.ContinueOnError)
	f := &listFlags{json: fs.Bool("json", false, "emit JSON")}
	tools.SetUsage(fs, "proj live --json", liveHelp)
	return fs, f
}

const liveHelp = "Prints the running sessions as a JSON array of {socket,name,dir,agent,conversation}, " +
	"the shape `proj restore --from` reads. conversation is the agent's session id (announced " +
	"in the session, else proj's record), empty when unknown."

// cmdLive implements `proj live --json`.
func cmdLive(args []string, out, errw io.Writer) error {
	fs, f := newLiveFlags()
	if err := tools.ParseFlags(fs, args, out); err != nil {
		return err
	}
	if !*f.json {
		return tools.UsageError{Msg: "only --json is supported"}
	}
	rec, _ := proj.LoadRecord()
	live := []liveJSON{}
	for _, s := range liveSessions() {
		conv := s.Conversation
		if conv == "" {
			conv = rec.Sessions[s.Name].Conversation
		}
		live = append(live, liveJSON{Socket: s.Socket, Name: s.Name, Dir: s.Dir, Agent: s.Agent, Conversation: conv})
	}
	return writeJSON(out, live)
}

type restoreFlags struct {
	from *string
	json *bool
}

func newRestoreFlags() (*flag.FlagSet, *restoreFlags) {
	fs := flag.NewFlagSet("restore", flag.ContinueOnError)
	f := &restoreFlags{
		from: fs.String("from", "", "a file `proj live --json` wrote"),
		json: fs.Bool("json", false, "emit JSON"),
	}
	tools.SetUsage(fs, "proj restore --from <file> [--json]", restoreHelp)
	return fs, f
}

const restoreHelp = "Reopens each session in the file (as `proj live --json` printed it) that is not " +
	"running, resuming its agent's conversation (pi --session, claude --resume) when its transcript " +
	"exists. On macOS it then opens their Ghostty windows. --json prints " +
	"{restored,skipped,failed:[{name,error}]}. Exits 0 unless the file can't be read."

type restoreFailure struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}

// restoreJSON is what `proj restore --json` prints.
type restoreJSON struct {
	Restored []string         `json:"restored"`
	Skipped  []string         `json:"skipped"`
	Failed   []restoreFailure `json:"failed"`
}

// cmdRestore implements `proj restore --from <file>`, the non-interactive
// restore a workspace runs on boot for the sessions live when it stopped.
func cmdRestore(args []string, out, errw io.Writer) error {
	fs, f := newRestoreFlags()
	if err := tools.ParseFlags(fs, args, out); err != nil {
		return err
	}
	if *f.from == "" {
		return tools.UsageError{Msg: "--from is required"}
	}
	b, err := os.ReadFile(*f.from)
	if err != nil {
		return tools.Exitf(1, "%v", err)
	}
	var sessions []liveJSON
	if err := json.Unmarshal(b, &sessions); err != nil {
		return tools.Exitf(1, "%s: %v", *f.from, err)
	}

	res := restoreJSON{Restored: []string{}, Skipped: []string{}, Failed: []restoreFailure{}}
	running := liveState().Running
	transcripts := proj.DefaultTranscripts()
	cfg := proj.LoadConfig()
	var opened []proj.SavedSession
	for _, s := range sessions {
		if running[s.Name] {
			res.Skipped = append(res.Skipped, s.Name)
			continue
		}
		agent := s.Agent
		if agent == "shell" {
			agent = "" // a plain shell: nothing to launch or pin
		}
		conv := s.Conversation
		if _, ok := transcripts.Locate(conv); !ok {
			conv = ""
		}
		dir := proj.ResolveDir(s.Socket, s.Name, s.Dir)
		created, err := restoreSession(s.Socket, s.Name, dir, agent, conv)
		switch {
		case err != nil:
			res.Failed = append(res.Failed, restoreFailure{Name: s.Name, Error: err.Error()})
			continue
		case !created:
			res.Skipped = append(res.Skipped, s.Name)
			continue
		}
		res.Restored = append(res.Restored, s.Name)
		if cfg.SidebarFor(proj.ProjectFromSocket(s.Socket)) {
			spawnSidebar(s.Socket, s.Name, dir)
		}
		opened = append(opened, proj.SavedSession{Name: s.Name,
			Entry: proj.Entry{Socket: s.Socket, Dir: dir, Agent: agent, Conversation: conv}})
	}

	if restoreGOOS == "darwin" && len(opened) > 0 {
		if err := ghosttyOpen(proj.PlanGhostty(opened), errw); err != nil {
			_, _ = fmt.Fprintf(errw, "ghostty: %v; attach with: proj\n", err)
		}
	}
	if *f.json {
		return writeJSON(out, res)
	}
	for _, n := range res.Restored {
		_, _ = fmt.Fprintf(out, "restored %s\n", n)
	}
	for _, n := range res.Skipped {
		_, _ = fmt.Fprintf(out, "exists %s\n", n)
	}
	for _, x := range res.Failed {
		_, _ = fmt.Fprintf(out, "failed %s: %s\n", x.Name, x.Error)
	}
	return nil
}
