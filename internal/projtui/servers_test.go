package projtui

import (
	"slices"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/subaud-tools/tackle/internal/proj/servers"
)

type serverFakes struct {
	calls []string // "restart cull", "stop-left-over", ...
}

// newServersModel is an entrance model on the servers scope, listing diazo
// up, cull down and one left-over galley, with every action recorded.
func newServersModel(t *testing.T) (Model, *serverFakes) {
	t.Helper()
	f := &serverFakes{}
	old := serversInterval
	t.Cleanup(func() { serversInterval = old })
	serversInterval = time.Millisecond
	m := newTestModel(nil)
	m.srv = serverActions{
		list: func() (servers.List, error) {
			return servers.List{
				Up:   []servers.Row{{Tool: "diazo", Up: true, Address: "127.0.0.1:4100", URL: "http://127.0.0.1:4100/"}},
				Down: []servers.Row{{Tool: "cull"}},
				Left: []servers.Leftover{{Tool: "galley", PID: 99, Age: "3d", Command: "galley serve"}},
			}, nil
		},
		start:         func(tool string) error { f.calls = append(f.calls, "start "+tool); return nil },
		stop:          func(tool string) error { f.calls = append(f.calls, "stop "+tool); return nil },
		restart:       func(tool string) error { f.calls = append(f.calls, "restart "+tool); return nil },
		log:           func(tool string) (string, error) { return "log of " + tool, nil },
		stopLeftovers: func() (int, error) { f.calls = append(f.calls, "stop left over"); return 1, nil },
	}
	for m.scope != scopeServers {
		m = run(t, m, keyMsg("tab"))
	}
	if len(m.visibleRows()) == 0 {
		t.Fatal("the servers scope listed nothing")
	}
	return m, f
}

func keyMsg(k string) tea.Msg {
	if k == "tab" {
		return tea.KeyMsg{Type: tea.KeyTab}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
}

// run sends msg and then runs the commands it returns, and theirs, except the
// refresh ticks, as Bubble Tea would.
func run(t *testing.T, m Model, msg tea.Msg) Model {
	t.Helper()
	queue := []tea.Msg{msg}
	for len(queue) > 0 {
		next, cmd := m.Update(queue[0])
		m, queue = next.(Model), queue[1:]
		queue = append(queue, drain(cmd)...)
	}
	return m
}

func drain(cmd tea.Cmd) []tea.Msg {
	if cmd == nil {
		return nil
	}
	switch msg := cmd().(type) {
	case tea.BatchMsg:
		var out []tea.Msg
		for _, c := range msg {
			out = append(out, drain(c)...)
		}
		return out
	case serversTickMsg, tickMsg, nil:
		return nil
	default:
		return []tea.Msg{msg}
	}
}

// selectTool puts the cursor on tool's row.
func selectTool(t *testing.T, m Model, tool string) Model {
	t.Helper()
	for i, r := range m.visibleRows() {
		if r.Kind != RowHeader && r.Name == tool {
			m.cursor = i
			return m
		}
	}
	t.Fatalf("no row for %s in %v", tool, m.visibleRows())
	return m
}

func TestServerKeysActOnTheHighlightedRow(t *testing.T) {
	m, f := newServersModel(t)
	m = run(t, selectTool(t, m, "diazo"), keyMsg("r"))
	m = run(t, selectTool(t, m, "diazo"), keyMsg("x"))
	run(t, selectTool(t, m, "cull"), keyMsg("s"))
	if want := []string{"restart diazo", "stop diazo", "start cull"}; !slices.Equal(f.calls, want) {
		t.Fatalf("calls %v, want %v", f.calls, want)
	}
}

func TestStopAllLeftOverAsksAndActsOnlyOnY(t *testing.T) {
	m, f := newServersModel(t)
	m = run(t, m, keyMsg("X"))
	if len(f.calls) != 0 || m.footerHint == "" {
		t.Fatalf("X acted without asking: calls %v, hint %q", f.calls, m.footerHint)
	}
	m = run(t, m, keyMsg("n"))
	if len(f.calls) != 0 {
		t.Fatalf("n stopped them: %v", f.calls)
	}
	m = run(t, run(t, m, keyMsg("X")), keyMsg("y"))
	if want := []string{"stop left over"}; !slices.Equal(f.calls, want) {
		t.Fatalf("calls %v, want %v", f.calls, want)
	}
}
