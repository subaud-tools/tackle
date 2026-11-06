package proj

import "strings"

// ScannedSession is one session read by ScanServer, with its main pane's fields.
type ScannedSession struct {
	Name         string
	Attached     bool
	Pane         string // main pane id, e.g. "%10"
	Command      string // main pane's pane_current_command
	Dir          string // main pane's pane_current_path
	Bell         bool   // window_bell_flag of the main pane's window
	Activity     int64  // window_activity epoch of the main pane's window, 0 if absent
	Conversation string // @harness_session
	Agent        string // @proj_agent (AgentOption())
	Label        string // LabelOption()
	LastAttached int64  // session_last_attached epoch, 0 if never attached
	Created      int64  // session_created epoch
}

// unitFields splits a line of -F output on the unit separator. tmux 3.4
// prints that byte as the four characters \037 (a later tmux prints the
// byte), so both are read as the separator.
func unitFields(line string) []string {
	return strings.Split(strings.ReplaceAll(line, `\037`, "\x1f"), "\x1f")
}

// scanFormat is the list-panes -a format ScanServer reads, one pane per line.
// Fields are separated by \x1f so an empty one (an unset option, a blank
// command) never shifts the columns; see mainPaneCmd. window_active, last,
// keeps the main pane in the session's current window, as list-panes -t
// <session> did. session_last_attached and session_created come after it,
// so the earlier columns keep their indexes. A function, because the option
// names come from the environment.
func scanFormat() string {
	return strings.Join([]string{
		"#{session_name}",
		"#{session_attached}",
		"#{pane_left}",
		"#{pane_height}",
		"#{pane_id}",
		"#{pane_current_command}",
		"#{pane_current_path}",
		"#{window_bell_flag}",
		"#{window_activity}",
		"#{@harness_session}",
		"#{" + AgentOption() + "}",
		"#{" + LabelOption() + "}",
		"#{window_active}",
		"#{session_last_attached}",
		"#{session_created}",
	}, "\x1f")
}

const scanFields = 15

// ScanServer reads socket's sessions with one `list-panes -a` call.
// A server that is not running returns the tmux error.
func ScanServer(socket string) ([]ScannedSession, error) {
	sessions, _, err := scanServerBusy(socket)
	return sessions, err
}

// scanServerBusy is ScanServer with each session's busy panes, from the same
// list-panes call.
func scanServerBusy(socket string) ([]ScannedSession, map[string][]BusyPane, error) {
	out, err := Run(socket, "list-panes", "-a", "-F", scanFormat())
	if err != nil {
		return nil, nil, err
	}
	sessions := parseScan(out)
	return sessions, parseBusyPanes(out, sessions), nil
}

// BusyPane is a pane, other than its session's main pane, running something
// other than a shell: a long job, a server, a second agent.
type BusyPane struct {
	Pane    string `json:"pane"`
	Command string `json:"command"`
}

// parseBusyPanes lists, per session, the panes in any of its windows other
// than its main pane whose current command is not a shell (a login shell's
// "-zsh" included) or a sidebar tool (yazi, scratch), which sit idle all day.
func parseBusyPanes(out string, sessions []ScannedSession) map[string][]BusyPane {
	main := make(map[string]string, len(sessions))
	for _, s := range sessions {
		main[s.Name] = s.Pane
	}
	busy := map[string][]BusyPane{}
	for _, ln := range splitLines(out) {
		f := unitFields(ln)
		if len(f) < scanFields || f[0] == "" {
			continue
		}
		pane, cmd := f[4], f[5]
		if m, ok := main[f[0]]; !ok || pane == m || cmd == "" ||
			classifyAgent(strings.TrimPrefix(cmd, "-")) == "shell" {
			continue
		}
		busy[f[0]] = append(busy[f[0]], BusyPane{Pane: pane, Command: cmd})
	}
	return busy
}

// parseScan turns list-panes -a output in scanFormat() into sessions, in
// first-seen order, choosing each session's main pane in its current window
// with pickMainPane's rule.
func parseScan(out string) []ScannedSession {
	type best struct{ left, height int }
	var sessions []ScannedSession
	idx := map[string]int{}
	var bests []best
	for _, ln := range splitLines(out) {
		f := unitFields(ln)
		if len(f) < scanFields || f[0] == "" || f[12] != "1" {
			continue
		}
		l, h := atoi(f[2]), atoi(f[3])
		i, seen := idx[f[0]]
		if seen && !betterPane(l, h, bests[i].left, bests[i].height) {
			continue
		}
		s := ScannedSession{
			Name:         f[0],
			Attached:     atoi(f[1]) > 0,
			Pane:         f[4],
			Command:      f[5],
			Dir:          f[6],
			Bell:         f[7] == "1",
			Activity:     int64(atoi(f[8])),
			Conversation: f[9],
			Agent:        f[10],
			Label:        f[11],
			LastAttached: int64(atoi(f[13])),
			Created:      int64(atoi(f[14])),
		}
		if !seen {
			idx[f[0]] = len(sessions)
			sessions = append(sessions, s)
			bests = append(bests, best{l, h})
			continue
		}
		sessions[i] = s
		bests[i] = best{l, h}
	}
	return sessions
}
