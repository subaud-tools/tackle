package projtui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/subaud-tools/tackle/internal/proj"
)

// footerOf renders m's footer wide enough that nothing is cut.
func footerOf(m Model) string { return m.footerView(200) }

func namingWork(t *testing.T) Model {
	t.Helper()
	m := newProjectModel("tools-workspace", nil)
	return press(m, "enter") // open the "+ new work…" input
}

func TestCtrlOFlipsTheChoiceWhileNamingNewWork(t *testing.T) {
	m := namingWork(t)
	if f := footerOf(m); !strings.Contains(f, "^o open:terminal") {
		t.Fatalf("footer should start on terminal:\n%s", f)
	}
	m = press(m, "ctrl+o")
	if f := footerOf(m); !strings.Contains(f, "^o open:chrome") {
		t.Fatalf("ctrl+o should show chrome:\n%s", f)
	}
	m = press(m, "ctrl+o")
	if f := footerOf(m); !strings.Contains(f, "^o open:terminal") {
		t.Fatalf("a second ctrl+o should show terminal:\n%s", f)
	}
}

func TestCtrlOFlipsTheChoiceWhileBrowsing(t *testing.T) {
	m := sessionsScope(newTestModel([]Row{{Kind: RowSession, Socket: "proj-a", Name: "a/w", Label: "a/w"}}))
	if f := footerOf(m); !strings.Contains(f, "^o open:terminal") {
		t.Fatalf("browse footer should start on terminal:\n%s", f)
	}
	m = press(m, "ctrl+o")
	if f := footerOf(m); !strings.Contains(f, "^o open:chrome") {
		t.Fatalf("ctrl+o should show chrome while browsing:\n%s", f)
	}
}

func TestChromeTurnsSidebarOffAndTerminalRestoresIt(t *testing.T) {
	for _, before := range []bool{true, false} {
		m := namingWork(t)
		m.sidebarChoice = before
		m = press(m, "ctrl+o")
		if m.sidebarChoice {
			t.Fatalf("before=%v: chrome should turn the sidebar off", before)
		}
		m = press(m, "ctrl+s")
		if !m.sidebarChoice {
			t.Fatalf("before=%v: ctrl+s should still turn the sidebar on under chrome", before)
		}
		m = press(m, "ctrl+o")
		if m.sidebarChoice != before {
			t.Fatalf("back on terminal the sidebar = %v, want the earlier %v", m.sidebarChoice, before)
		}
	}
}

func TestEnterOnSessionCarriesTheChoice(t *testing.T) {
	rows := []Row{{Kind: RowSession, Socket: "proj-a", Name: "a/w", Label: "a/w"}}
	m := press(sessionsScope(newTestModel(rows)), "enter")
	if m.Result.Kind != "jump" || m.Result.Chrome {
		t.Fatalf("terminal: Result = %+v, want a jump without Chrome", m.Result)
	}
	m = press(press(sessionsScope(newTestModel(rows)), "ctrl+o"), "enter")
	if m.Result.Kind != "jump" || !m.Result.Chrome || m.Result.Socket != "proj-a" || m.Result.Name != "a/w" {
		t.Fatalf("chrome: Result = %+v, want a jump to a/w asking for Chrome", m.Result)
	}
}

func TestEnterOnNewWorkCarriesTheChoice(t *testing.T) {
	m := press(typeString(namingWork(t), "fix"), "enter")
	if m.Result.Kind != "new" || m.Result.Chrome || !m.Result.Sidebar {
		t.Fatalf("terminal: Result = %+v, want new work with the sidebar and no Chrome", m.Result)
	}
	m = press(typeString(press(namingWork(t), "ctrl+o"), "fix"), "enter")
	if m.Result.Kind != "new" || !m.Result.Chrome || m.Result.Sidebar {
		t.Fatalf("chrome: Result = %+v, want new work in Chrome without the sidebar", m.Result)
	}
}

func TestCtrlOSavesTheChoice(t *testing.T) {
	m := namingWork(t)
	var saved []string
	m.saveOpenIn = func(where string) error { saved = append(saved, where); return nil }
	m = press(press(m, "ctrl+o"), "ctrl+o")
	if strings.Join(saved, ",") != proj.OpenChrome+","+proj.OpenTerminal {
		t.Fatalf("saved = %v", saved)
	}
	m.saveOpenIn = func(string) error { return errors.New("disk full") }
	m = press(m, "ctrl+o")
	if !strings.Contains(m.footerHint, "disk full") {
		t.Fatalf("a failed save should say so, hint = %q", m.footerHint)
	}
}

func TestSavedChromeStartsTheNextPickerOnChromeWithoutSidebar(t *testing.T) {
	m := newProjectModel("tools-workspace", nil).withOpenIn(proj.OpenChrome)
	if f := footerOf(m); !strings.Contains(f, "^o open:chrome") {
		t.Fatalf("browse footer should read the saved choice:\n%s", f)
	}
	m = press(m, "enter")
	if m.sidebarChoice {
		t.Fatal("new work under a saved chrome choice should have no sidebar")
	}
	m = press(m, "ctrl+o")
	if !m.sidebarChoice {
		t.Fatal("switching back to terminal should restore the configured sidebar")
	}
}

// Restore stays terminal-only and a bare `proj reap` opens nothing, so neither
// offers the toggle.
func TestCtrlONotOfferedOnSavedTabOrInReap(t *testing.T) {
	saved := sessionsScope(newTestModel(nil))
	saved.scope = scopeSaved
	reap := newTestModel(nil).reapMode()
	for name, m := range map[string]Model{"saved": saved, "reap": reap} {
		if f := footerOf(m); strings.Contains(f, "open:") {
			t.Fatalf("%s footer offers the toggle:\n%s", name, f)
		}
		next, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlO})
		if next.(Model).openChrome {
			t.Fatalf("%s: ctrl+o flipped the choice", name)
		}
	}
}
