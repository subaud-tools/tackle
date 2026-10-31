package proj

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/sessiontheme"
)

// Client joins a native terminal to a tmux session on an explicit server.
type Client struct {
	Socket, SessionID, TTY string
	generation             string
}

var themeNow = time.Now
var themeSockets = Servers

var themeSessionOptions = []string{"@proj_theme_bg", "@proj_theme_fg", "@proj_theme_panel", "@proj_theme_selected_bg", "@proj_theme_muted", "@proj_theme_border", "@proj_theme_accent", "@proj_theme_accent_text", "@proj_theme_success", "@proj_theme_warning", "@proj_theme_error", "status-style"}
var themeWindowOptions = []string{"window-status-style", "window-status-current-style", "pane-border-style", "pane-active-border-style"}

type themeSession struct {
	socket, id, name, stable string
	entry                    Entry
	options                  []string
}
type themeWindow struct {
	socket, session, id string
	options             []string
}
type themeObservation struct {
	terminals []GhosttyTerminal
	clients   []Client
	sessions  []themeSession
	windows   []themeWindow
}
type themeSyncState struct {
	ObservedAt time.Time         `json:"observed_at"`
	Applied    map[string]string `json:"applied"`
}

func optionFormat(keys []string) string {
	var f []string
	for _, key := range keys {
		f = append(f, "#{"+key+"}")
	}
	return strings.Join(f, "\t")
}

// observeThemes finishes every read before the record is touched. All fields
// needed from a server are batched, never queried once per palette role.
func observeThemes() (themeObservation, error) {
	terminals, err := ObserveGhostty()
	if err != nil {
		return themeObservation{}, err
	}
	o := themeObservation{terminals: terminals}
	for _, sock := range themeSockets() {
		// Reuse ScanServer's main-pane selection and its batched field layout.
		panes, err := runner(sock, "list-panes", "-a", "-F", scanFormat())
		if err != nil {
			return o, err
		}
		scanned := map[string]ScannedSession{}
		for _, s := range parseScan(panes) {
			scanned[s.Name] = s
		}
		out, err := runner(sock, "list-sessions", "-F", "#{session_id}\t#{session_name}\t#{@proj_session_id}\t"+optionFormat(themeSessionOptions))
		if err != nil {
			return o, err
		}
		for _, line := range splitLines(out) {
			f := strings.Split(line, "\t")
			if len(f) != 3+len(themeSessionOptions) {
				return o, fmt.Errorf("proj: malformed session observation")
			}
			s, ok := scanned[f[1]]
			if !ok {
				return o, fmt.Errorf("proj: session changed during theme observation: %s", f[1])
			}
			o.sessions = append(o.sessions, themeSession{socket: sock, id: f[0], name: f[1], stable: f[2], entry: Entry{ID: f[2], Socket: sock, Project: ProjectFromSocket(sock), Dir: s.Dir, Conversation: s.Conversation, Agent: s.Agent}, options: f[3:]})
		}
		out, err = runner(sock, "list-clients", "-F", "#{session_id}\t#{client_tty}\t#{client_created}:#{client_pid}")
		if err != nil {
			return o, err
		}
		for _, line := range splitLines(out) {
			f := strings.Split(line, "\t")
			if len(f) != 3 {
				return o, fmt.Errorf("proj: malformed client observation")
			}
			o.clients = append(o.clients, Client{Socket: sock, SessionID: f[0], TTY: f[1], generation: f[2]})
		}
		out, err = runner(sock, "list-windows", "-a", "-F", "#{session_id}\t#{window_id}\t"+optionFormat(themeWindowOptions))
		if err != nil {
			return o, err
		}
		for _, line := range splitLines(out) {
			f := strings.Split(line, "\t")
			if len(f) != 2+len(themeWindowOptions) {
				return o, fmt.Errorf("proj: malformed window observation")
			}
			o.windows = append(o.windows, themeWindow{socket: sock, session: f[0], id: f[1], options: f[2:]})
		}
	}
	return o, nil
}

func (o themeObservation) matched() map[string]themeSession {
	sessions := map[string]themeSession{}
	for _, s := range o.sessions {
		sessions[s.socket+"\x00"+s.id] = s
	}
	byTTY := map[string]themeSession{}
	for _, c := range o.clients {
		if s, ok := sessions[c.Socket+"\x00"+c.SessionID]; ok {
			byTTY[c.TTY] = s
		}
	}
	return byTTY
}

func (o themeObservation) save() (Record, error) {
	var saved Record
	byTTY := o.matched()
	err := UpdateRecord(func(r *Record) bool {
		before, _ := json.Marshal(r)
		var live []LiveEntry
		for _, s := range o.sessions {
			live = append(live, LiveEntry{Name: s.name, Entry: s.entry})
		}
		mergeLive(r, live, themeNow())
		var native []string
		groups := map[string][]string{}
		for _, t := range o.terminals {
			if s, ok := byTTY[t.TTY]; ok {
				if _, seen := groups[t.WindowID]; !seen {
					native = append(native, t.WindowID)
				}
				if !containsName(groups[t.WindowID], s.name) {
					groups[t.WindowID] = append(groups[t.WindowID], s.name)
				}
			}
		}
		var cur [][]string
		for _, id := range native {
			names := groups[id]
			assignWindowGroup(r, id, names)
			cur = append(cur, names)
		}
		oldWindows, oldGroupIDs := r.Layout.Windows, r.Layout.GroupIDs
		wins := mergeLayout(oldWindows, cur)
		if !reflect.DeepEqual(r.Layout.Windows, wins) || r.Layout.SavedAt.IsZero() {
			r.Layout.Windows = wins
			r.Layout.GroupIDs = nil
			r.Layout.SavedAt = themeNow()
		}
		// Rebuild layout bindings after moves, including absent saved work.
		r.Layout.GroupIDs = make([]string, len(wins))
		for i, w := range wins {
			for _, name := range w {
				if id := r.Sessions[name].GroupID; id != "" {
					r.Layout.GroupIDs[i] = id
					break
				}
			}
			// Absent windows may have no session entry left to carry their
			// identity. Match the saved window, not its pre-merge index.
			if r.Layout.GroupIDs[i] == "" {
				for j, old := range oldWindows {
					if j < len(oldGroupIDs) && reflect.DeepEqual(w, old) {
						r.Layout.GroupIDs[i] = oldGroupIDs[j]
						break
					}
				}
			}
		}
		saved = *r
		after, _ := json.Marshal(r)
		return string(before) != string(after)
	})
	return saved, err
}

func containsName(names []string, name string) bool {
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}

// ReconcileThemes coalesces status refreshes across processes. Assignment and
// successful terminal delivery are separate: a failed write is never cached.
func ReconcileThemes(force bool) error {
	if !SessionThemesEnabled() {
		if force {
			_, _, err := SaveLayout(nil)
			return err
		}
		return nil
	}
	return withThemeSyncLock(func() error {
		path := filepath.Join(tools.StateDir("proj"), "theme-sync.json")
		state := themeSyncState{Applied: map[string]string{}}
		if b, err := os.ReadFile(path); err == nil {
			if err = json.Unmarshal(b, &state); err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			return err
		}
		now := themeNow()
		if !force && !state.ObservedAt.IsZero() && now.Sub(state.ObservedAt) >= 0 && now.Sub(state.ObservedAt) < 5*time.Second {
			return nil
		}
		o, err := observeThemes()
		if err != nil {
			return err
		}
		r, err := o.save()
		if err != nil {
			return err
		}
		next := themeSyncState{ObservedAt: now, Applied: map[string]string{}}
		byTTY := o.matched()
		appliedSessions := map[string]bool{}
		var applyErr error
		// Pi watches these stable paths; refreshing them never types a command
		// into, or restarts, an existing pane. Keep failures independent of the
		// successful-terminal ledger so its retry semantics remain unchanged.
		for _, s := range o.sessions {
			e := r.Sessions[s.name]
			if e.Agent == "pi" && e.GroupID != "" {
				if _, err := preparePiTheme(e, r); err != nil {
					applyErr = fmt.Errorf("proj: refresh pi theme for %s: %w", s.name, err)
				}
			}
		}
		for _, t := range o.terminals {
			s, ok := byTTY[t.TTY]
			if !ok {
				continue
			}
			e := r.Sessions[s.name]
			p, err := sessiontheme.Resolve(sessiontheme.Assignment{Hue: r.Groups[e.GroupID].Hue, Shade: e.Shade})
			if err != nil {
				return err
			}
			sessionKey := s.socket + "\x00" + s.id
			if !appliedSessions[sessionKey] {
				if err := applyTmuxTheme(s, o.windows, e.ID, p); err != nil {
					applyErr = fmt.Errorf("proj: apply tmux theme on %s/%s: %w", s.socket, s.name, err)
				}
				appliedSessions[sessionKey] = true
			}
			generation := ""
			for _, c := range o.clients {
				if c.TTY == t.TTY && c.Socket == s.socket {
					generation = c.generation
					break
				}
			}
			key := s.socket + "\x00" + t.TTY
			signature := e.ID + "\x00" + generation + "\x00" + string(TerminalSequences(p))
			if state.Applied[key] != signature {
				if err := writeTerminal(t.TTY, TerminalSequences(p)); err != nil {
					applyErr = err
					continue
				}
			}
			next.Applied[key] = signature
		}
		b, err := json.Marshal(next)
		if err != nil {
			return err
		}
		if err := tools.WriteFileAtomic(path, b, 0600); err != nil {
			return err
		}
		return applyErr
	})
}

func applyTmuxTheme(s themeSession, windows []themeWindow, stable string, p sessiontheme.Palette) error {
	values := []string{p.Background, p.Foreground, p.Panel, p.SelectedBackground, p.Muted, p.Border, p.Accent, p.AccentText, p.Success, p.Warning, p.Error, "bg=" + p.Background + ",fg=" + p.Foreground}
	var args []string
	add := func(cmd, target, key, value string) {
		if len(args) > 0 {
			args = append(args, ";")
		}
		args = append(args, cmd, "-t", target, key, value)
	}
	if s.stable != stable {
		add("set-option", s.id, "@proj_session_id", stable)
	}
	for i, key := range themeSessionOptions {
		if s.options[i] != values[i] {
			add("set-option", s.id, key, values[i])
		}
	}
	styles := []string{"bg=" + p.Panel + ",fg=" + p.Foreground, "bg=" + p.SelectedBackground + ",fg=" + p.Foreground, "fg=" + p.Border, "fg=" + p.Accent}
	for _, w := range windows {
		if w.socket != s.socket || w.session != s.id {
			continue
		}
		for i, key := range themeWindowOptions {
			if w.options[i] != styles[i] {
				add("set-window-option", w.id, key, styles[i])
			}
		}
	}
	if len(args) == 0 {
		return nil
	}
	_, err := runner(s.socket, args...)
	return err
}
