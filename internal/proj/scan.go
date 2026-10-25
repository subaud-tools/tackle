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
}

// scanFormat is the list-panes -a format ScanServer reads, one pane per line.
// Fields are separated by \x1f so an empty one (an unset option, a blank
// command) never shifts the columns; see mainPaneCmd. A function, because the
// option names come from the environment.
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
	}, "\x1f")
}

const scanFields = 12

// ScanServer reads socket's sessions with one `list-panes -a` call.
// A server that is not running returns the tmux error.
func ScanServer(socket string) ([]ScannedSession, error) {
	out, err := Run(socket, "list-panes", "-a", "-F", scanFormat())
	if err != nil {
		return nil, err
	}
	return parseScan(out), nil
}

// parseScan turns list-panes -a output in scanFormat() into sessions, in
// first-seen order, choosing each session's main pane with pickMainPane's rule.
func parseScan(out string) []ScannedSession {
	type best struct{ left, height int }
	var sessions []ScannedSession
	idx := map[string]int{}
	var bests []best
	for _, ln := range splitLines(out) {
		f := strings.Split(ln, "\x1f")
		if len(f) < scanFields || f[0] == "" {
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
