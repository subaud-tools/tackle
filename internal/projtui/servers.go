package projtui

import (
	"context"
	"fmt"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/subaud-tools/tools-common/localweb"
	"github.com/subaud-tools/tools-common/pageserve"

	"github.com/subaud-tools/tackle/internal/proj/pages"
	"github.com/subaud-tools/tackle/internal/proj/servers"
)

// The servers scope lists every family page server, up then down, and the
// left-over processes no server holds (shown as "also running"), refreshed every serversInterval while
// it is on screen. Its keys act on the highlighted row through the servers
// package; nothing it does waits inside Update.

// serversInterval is how often the servers scope rereads the list while it
// is on screen. A var so tests can shorten it.
var serversInterval = 2 * time.Second

// actionTimeout bounds a stop or restart, which waits for the server to go.
const actionTimeout = 20 * time.Second

// serverActions are the servers scope's side effects. newModel leaves them
// nil so tests act on nothing; New wires the servers package.
type serverActions struct {
	list          func() (servers.List, error)
	start         func(tool string) error
	stop          func(tool string) error
	restart       func(tool string) error
	log           func(tool string) (string, error)
	stopLeftover  func(pid int) error
	stopLeftovers func() (int, error)
	open          func(url string) error
	page          func() (string, error) // the servers page's URL
}

// withServers wires the servers scope to the real servers and browser.
func (m Model) withServers() Model {
	ctxDo := func(f func(context.Context, string) error) func(string) error {
		return func(tool string) error {
			ctx, cancel := context.WithTimeout(context.Background(), actionTimeout)
			defer cancel()
			return f(ctx, tool)
		}
	}
	m.srv = serverActions{
		list:          servers.Read,
		start:         servers.Start,
		stop:          ctxDo(servers.Stop),
		restart:       ctxDo(servers.Restart),
		log:           servers.Log,
		stopLeftover:  servers.StopLeftover,
		stopLeftovers: servers.StopLeftovers,
		open:          localweb.OpenBrowser,
		page: func() (string, error) {
			e, err := pages.Process().Ensure(pageserve.Executable())
			return e.URL, err
		},
	}
	return m
}

// serversTickMsg rereads the list while the scope is on screen.
type serversTickMsg struct{}

// serversLoadedMsg carries a read of the list.
type serversLoadedMsg struct {
	list servers.List
	err  error
}

// serverDoneMsg reports an action: hint is the footer line.
type serverDoneMsg struct{ hint string }

// serverLogMsg carries a server's log tail for the preview pane.
type serverLogMsg struct{ tool, text string }

func serversTick() tea.Cmd {
	return tea.Tick(serversInterval, func(time.Time) tea.Msg { return serversTickMsg{} })
}

func (m Model) onServers() bool { return m.view == viewEntrance && m.scope == scopeServers }

// enterServers reads the list now and starts the refresh tick, unless one
// is already running.
func (m Model) enterServers() (Model, tea.Cmd) {
	cmds := []tea.Cmd{m.loadServers()}
	if !m.srvTicking {
		m.srvTicking = true
		cmds = append(cmds, serversTick())
	}
	return m, tea.Batch(cmds...)
}

func (m Model) loadServers() tea.Cmd {
	list := m.srv.list
	if list == nil {
		return nil
	}
	return func() tea.Msg {
		l, err := list()
		return serversLoadedMsg{l, err}
	}
}

// serversTicked rereads the list and re-arms while the scope is on screen;
// off screen the tick stops.
func (m Model) serversTicked() (tea.Model, tea.Cmd) {
	if !m.onServers() {
		m.srvTicking = false
		return m, nil
	}
	return m, tea.Batch(m.loadServers(), serversTick())
}

func (m Model) applyServers(msg serversLoadedMsg) Model {
	m.srvLoaded = true
	if msg.err != nil {
		m.footerHint = "servers: " + msg.err.Error()
	}
	m.srvList = msg.list
	if m.onServers() {
		m = m.rebuildPreserving()
	}
	return m
}

// serverRows are the list's rows; visibleRows adds the group headers.
func serverRows(l servers.List) []Row {
	var rows []Row
	for _, g := range [][]servers.Row{l.Up, l.Down} {
		for _, s := range g {
			rows = append(rows, Row{Kind: RowServer, Label: s.Tool, Name: s.Tool, Server: s})
		}
	}
	for _, p := range l.Left {
		rows = append(rows, Row{Kind: RowLeftover, Label: p.Tool, Name: "pid " + strconv.Itoa(p.PID), Leftover: p})
	}
	return rows
}

// groupServers puts rows under "up · N", "down · N" and "also running · N"
// (the left-over processes).
func groupServers(rows []Row) []Row {
	var up, down, left []Row
	for _, r := range rows {
		switch {
		case r.Kind == RowLeftover:
			left = append(left, r)
		case r.Server.Up:
			up = append(up, r)
		default:
			down = append(down, r)
		}
	}
	var out []Row
	for _, g := range []struct {
		label string
		rows  []Row
	}{{"up", up}, {"down", down}, {"also running", left}} {
		if len(g.rows) > 0 {
			out = append(out, Row{Kind: RowHeader, Label: fmt.Sprintf("%s · %d", g.label, len(g.rows))})
			out = append(out, g.rows...)
		}
	}
	return out
}

// act runs f off the update loop; its serverDoneMsg rereads the list.
func (m Model) act(f func() string) tea.Cmd {
	return func() tea.Msg { return serverDoneMsg{f()} }
}

func done(verb, what string, err error) string {
	if err != nil {
		return verb + " " + what + ": " + err.Error()
	}
	return verb + " " + what
}

// updateServers handles the servers scope's keys. ok=false hands the key on
// to the shared list handling (navigation, tab, ?, esc). Letters are this
// scope's commands, so typing does not filter here.
func (m Model) updateServers(msg tea.KeyMsg) (Model, tea.Cmd, bool) {
	if m.srvConfirm {
		m.srvConfirm = false
		m.footerHint = ""
		if msg.String() == "y" && m.srv.stopLeftovers != nil {
			stop := m.srv.stopLeftovers
			m.footerHint = "stopping other processes…"
			return m, m.act(func() string {
				n, err := stop()
				return done("stopped", fmt.Sprintf("%d other processes", n), err)
			}), true
		}
		return m, nil, true
	}
	var row Row
	if vis := m.visibleRows(); m.cursor < len(vis) {
		row = vis[m.cursor]
	}
	tool := row.Server.Tool
	switch msg.String() {
	case "enter":
		switch {
		case row.Kind != RowServer:
		case !row.Server.Up:
			m.footerHint = tool + " is not running: s starts it"
		case m.srv.open != nil:
			if err := m.srv.open(row.Server.URL); err != nil {
				m.footerHint = "open: " + err.Error()
			}
		}
		return m, nil, true
	case "r":
		if row.Kind == RowServer && m.srv.restart != nil {
			f := m.srv.restart
			m.footerHint = "restarting " + tool + "…"
			return m, m.act(func() string { return done("restarted", tool, f(tool)) }), true
		}
		return m, nil, true
	case "x":
		switch {
		case row.Kind == RowLeftover && m.srv.stopLeftover != nil:
			f, pid := m.srv.stopLeftover, row.Leftover.PID
			return m, m.act(func() string { return done("stopped", "pid "+strconv.Itoa(pid), f(pid)) }), true
		case row.Kind == RowServer && !row.Server.Up:
			m.footerHint = tool + " is not running"
		case row.Kind == RowServer && m.srv.stop != nil:
			f := m.srv.stop
			m.footerHint = "stopping " + tool + "…"
			return m, m.act(func() string { return done("stopped", tool, f(tool)) }), true
		}
		return m, nil, true
	case "s":
		switch {
		case row.Kind != RowServer:
		case row.Server.Up:
			m.footerHint = tool + " is already up"
		case m.srv.start != nil:
			f := m.srv.start
			m.footerHint = "starting " + tool + "…"
			return m, m.act(func() string { return done("started", tool, f(tool)) }), true
		}
		return m, nil, true
	case "l":
		if row.Kind == RowServer && m.srv.log != nil {
			f := m.srv.log
			return m, func() tea.Msg {
				text, err := f(tool)
				if err != nil {
					text = err.Error()
				}
				return serverLogMsg{tool, text}
			}, true
		}
		return m, nil, true
	case "X":
		if n := len(m.srvList.Left); n == 0 {
			m.footerHint = "nothing else running"
		} else {
			m.srvConfirm = true
			m.footerHint = fmt.Sprintf("stop %d other processes? y/n", n)
		}
		return m, nil, true
	case "w":
		if m.srv.page != nil && m.srv.open != nil {
			page, open := m.srv.page, m.srv.open
			return m, func() tea.Msg {
				u, err := page()
				if err == nil {
					err = open(u)
				}
				if err != nil {
					return serverDoneMsg{"servers page: " + err.Error()}
				}
				return serverDoneMsg{"opened the servers page"}
			}, true
		}
		return m, nil, true
	}
	if runeText(msg) != "" {
		return m, nil, true
	}
	return m, nil, false
}
