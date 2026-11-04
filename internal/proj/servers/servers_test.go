package servers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/subaud-tools/tools-common/pageserve"
)

var t0 = time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC)

// fixed makes Servers and Leftovers return these, and the clock read t0.
func fixed(t *testing.T, s []pageserve.Server, l []pageserve.Leftover) {
	t.Helper()
	ls, ll, ln := listServers, listLeftovers, now
	t.Cleanup(func() { listServers, listLeftovers, now = ls, ll, ln })
	listServers = func() ([]pageserve.Server, error) { return s, nil }
	listLeftovers = func([]pageserve.Server) ([]pageserve.Leftover, error) { return l, nil }
	now = func() time.Time { return t0 }
}

func TestRowsGroupUpDownAndLeftOver(t *testing.T) {
	fixed(t, []pageserve.Server{
		{Tool: "cull"},
		{Tool: "diazo", Up: true, Entry: pageserve.Entry{Base: "http://127.0.0.1:4100", URL: "http://127.0.0.1:4100/?t=x", PID: 7, Version: "1.2.0", StartedAt: t0.Add(-3*time.Hour - 12*time.Minute)}},
	}, []pageserve.Leftover{{Tool: "galley", PID: 99, Age: 26 * time.Hour, Command: "galley serve"}})
	l, err := Read()
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Up) != 1 || len(l.Down) != 1 || len(l.Left) != 1 {
		t.Fatalf("up %v, down %v, left %v", l.Up, l.Down, l.Left)
	}
	up := l.Up[0]
	if up.Tool != "diazo" || up.Address != "127.0.0.1:4100" || up.URL != "http://127.0.0.1:4100/?t=x" || up.Version != "1.2.0" || up.Uptime != "3h12m" {
		t.Errorf("up row %+v", up)
	}
	if l.Down[0].Tool != "cull" || l.Down[0].Up {
		t.Errorf("down row %+v", l.Down[0])
	}
	if lo := l.Left[0]; lo.Tool != "galley" || lo.PID != 99 || lo.Age != "1d2h" || lo.Command != "galley serve" {
		t.Errorf("left-over row %+v", lo)
	}
}

func TestAStaleEntryCarriesTheNewerBuildNote(t *testing.T) {
	fixed(t, []pageserve.Server{{Tool: "diazo", Up: true, Entry: pageserve.Entry{Base: "http://127.0.0.1:1", Stale: true, StartedAt: t0}}}, nil)
	l, err := Read()
	if err != nil {
		t.Fatal(err)
	}
	if l.Up[0].Note != NewerBuild {
		t.Fatalf("note = %q, want %q", l.Up[0].Note, NewerBuild)
	}
}

func TestTheLogIsTheLast40Lines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "serve.log")
	var b strings.Builder
	for i := 1; i <= 100; i++ {
		fmt.Fprintf(&b, "line %d\n", i)
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	fixed(t, []pageserve.Server{{Tool: "cull", LogPath: path}}, nil)
	got, err := Log("cull")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(lines) != 40 || lines[0] != "line 61" || lines[39] != "line 100" {
		t.Fatalf("%d lines, %q … %q", len(lines), lines[0], lines[len(lines)-1])
	}
}

func TestALeftOverStopRefusesAPidNoLongerListed(t *testing.T) {
	fixed(t, []pageserve.Server{{Tool: "galley"}}, []pageserve.Leftover{{Tool: "galley", PID: 99}})
	var stopped []int
	o := stopLeftover
	t.Cleanup(func() { stopLeftover = o })
	stopLeftover = func(pid int) error { stopped = append(stopped, pid); return nil }
	if err := StopLeftover(98); err == nil {
		t.Fatal("stopping pid 98, which Leftovers doesn't list, was not refused")
	}
	if err := StopLeftover(99); err != nil {
		t.Fatal(err)
	}
	if len(stopped) != 1 || stopped[0] != 99 {
		t.Fatalf("stopped %v, want [99]", stopped)
	}
}
