package proj

import "time"

// Session is a live tmux session discovered across the per-project servers.
// Opened is false for a session nobody has attached to yet (one started in
// the background); Created is its session_created epoch.
type Session struct {
	Name, Socket, Dir, Agent, State string
	Unread, ActionRequired          int
	Opened                          bool
	Created                         int64
}

// LiveSessions returns every live session across all Servers(), enriched with
// its working directory, agent kind/state, and muster attention counts. Each
// server is read with one ScanServer call.
func LiveSessions() []Session {
	counts := MusterCounts()
	device := MusterDevice()
	now := time.Now()
	var out []Session
	for _, sock := range Servers() {
		scanned, err := ScanServer(sock)
		if err != nil {
			continue
		}
		for _, s := range scanned {
			agent := classifyAgent(s.Command)
			state := agentState(agent, s.Bell, s.Activity, now)
			// A recorded @proj_agent wins over pane process detection: the pane's
			// foreground process may be a launcher (e.g. hail running the shim),
			// not the agent itself. State still comes from the pane.
			if s.Agent != "" {
				agent = s.Agent
			}
			a := AttentionFor(counts, device, AliasFor(sock, s.Label))
			out = append(out, Session{
				Name:           s.Name,
				Socket:         sock,
				Dir:            s.Dir,
				Agent:          agent,
				State:          state,
				Unread:         a.Unread,
				ActionRequired: a.ActionRequired,
				Opened:         s.LastAttached > 0,
				Created:        s.Created,
			})
		}
	}
	return out
}
