package projtui

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"

	"github.com/subaud-tools/tackle/internal/proj"
)

// openedFixture is p/a and p/c opened, p/b not opened (started two days
// before fixedNow), listed in that interleaved order.
func openedFixture() []Row {
	created := fixedNow.Add(-50 * time.Hour).Unix()
	return []Row{
		{Kind: RowSession, Label: "p/a", Name: "p/a", Project: "p", Agent: "pi", State: "idle", Opened: true},
		{Kind: RowSession, Label: "p/b", Name: "p/b", Project: "p", Agent: "pi", State: "waiting", Created: created},
		{Kind: RowSession, Label: "p/c", Name: "p/c", Project: "p", Agent: "pi", State: "working", Opened: true},
	}
}

var fixedNow = time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC)

// pinClock makes ages count from fixedNow.
func pinClock(t *testing.T) {
	t.Helper()
	orig := clock
	clock = func() time.Time { return fixedNow }
	t.Cleanup(func() { clock = orig })
}

// labels is each visible row's kind-tagged label: "# " for a header.
func labels(m Model) []string {
	var out []string
	for _, r := range m.visibleRows() {
		if r.Kind == RowHeader {
			out = append(out, "# "+r.Label)
			continue
		}
		out = append(out, r.Label)
	}
	return out
}

func TestBuildRowsCarriesOpenedAndCreated(t *testing.T) {
	sessions, _ := buildRows(proj.Roots{}, []proj.Session{{Name: "p/a", Opened: true, Created: 42}})
	if len(sessions) != 1 || !sessions[0].Opened || sessions[0].Created != 42 {
		t.Fatalf("rows = %+v, want Opened and Created 42", sessions)
	}
}

func TestProjectViewGroupsOpenedFirst(t *testing.T) {
	m := newProjectModel("p", openedFixture())
	want := []string{"+ new work…", "# opened · 2", "p/a", "p/c", "# not opened · 1", "p/b"}
	if got := labels(m); strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("rows = %q\nwant   %q", got, want)
	}
}

func TestSessionsScopeGroupsOpenedFirst(t *testing.T) {
	m := sessionsScope(newTestModel(openedFixture()))
	want := []string{"# opened · 2", "p/a", "p/c", "# not opened · 1", "p/b"}
	if got := labels(m); strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("rows = %q\nwant   %q", got, want)
	}
	if m.cursor != 1 {
		t.Fatalf("cursor = %d, want 1 (off the leading header)", m.cursor)
	}
}

func TestEmptyGroupHasNoHeader(t *testing.T) {
	rows := openedFixture()
	m := sessionsScope(newTestModel([]Row{rows[0], rows[2]}))
	want := []string{"# opened · 2", "p/a", "p/c"}
	if got := labels(m); strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("rows = %q\nwant   %q", got, want)
	}
}

func TestFilterAppliesBeforeGrouping(t *testing.T) {
	m := typeString(newProjectModel("p", openedFixture()), "p/b")
	want := []string{"+ new work…", "# not opened · 1", "p/b"}
	if got := labels(m); strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("rows = %q\nwant   %q", got, want)
	}
}

func TestCursorSkipsGroupHeaders(t *testing.T) {
	m := newProjectModel("p", openedFixture())
	var got []string
	for i := 0; i < 4; i++ {
		got = append(got, m.visibleRows()[m.cursor].Label)
		m = press(m, "down")
	}
	got = append(got, m.visibleRows()[m.cursor].Label)
	want := []string{"+ new work…", "p/a", "p/c", "p/b", "p/b"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("cursor walk = %q\nwant         %q", got, want)
	}
	for i := 0; i < 4; i++ {
		m = press(m, "up")
	}
	if m.visibleRows()[m.cursor].Kind != RowNewWork {
		t.Fatalf("cursor up to %+v, want new work", m.visibleRows()[m.cursor])
	}
}

func TestNotOpenedRowShowsHollowDotAndAge(t *testing.T) {
	pinClock(t)
	m := sized(sessionsScope(newTestModel(openedFixture())))
	var opened, notOpened string
	for _, ln := range screen(m) {
		switch {
		case strings.Contains(ln, "p/a"):
			opened = ln
		case strings.Contains(ln, "p/b"):
			notOpened = ln
		}
	}
	if !strings.Contains(notOpened, "○  p/b  pi·waiting · 2d") {
		t.Fatalf("not-opened row = %q, want ○ and pi·waiting · 2d", notOpened)
	}
	if !strings.Contains(opened, "●  p/a  pi·idle") || strings.Contains(opened, "idle · ") {
		t.Fatalf("opened row = %q, want ● and no age", opened)
	}
}

func TestNotOpenedPreview(t *testing.T) {
	pinClock(t)
	m := sized(sessionsScope(newTestModel(openedFixture())))
	m = press(m, "down")
	m = press(m, "down") // p/a → p/c → p/b
	if got := m.visibleRows()[m.cursor].Label; got != "p/b" {
		t.Fatalf("cursor on %q, want p/b", got)
	}
	if v := ansi.Strip(m.View()); !strings.Contains(v, "not opened · started 2d ago") {
		t.Fatalf("preview missing the not-opened line:\n%s", v)
	}
}

// A session started under a minute ago says "just now", not "now ago".
func TestNotOpenedPreviewJustStarted(t *testing.T) {
	pinClock(t)
	r := openedFixture()[1]
	r.Created = fixedNow.Unix()
	m := sized(sessionsScope(newTestModel([]Row{r})))
	if v := ansi.Strip(m.View()); !strings.Contains(v, "not opened · started just now") {
		t.Fatalf("preview missing the just-started line:\n%s", v)
	}
}

func TestAge(t *testing.T) {
	for _, c := range []struct {
		ago  time.Duration
		want string
	}{
		{30 * time.Second, "now"},
		{5 * time.Minute, "5m"},
		{5*time.Hour + 20*time.Minute, "5h"},
		{50 * time.Hour, "2d"},
	} {
		if got := age(fixedNow.Add(-c.ago).Unix(), fixedNow); got != c.want {
			t.Errorf("age(%v) = %q, want %q", c.ago, got, c.want)
		}
	}
}
