// Package servers is what proj's servers scope and servers page show and do:
// every family page server on the machine, up or down, the left-over
// processes no server holds, and the actions on them. It reads and acts only
// through pageserve, so it names no tool: a server is listed by existing.
package servers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/subaud-tools/tools-common/pageserve"
)

// NewerBuild is the note on a server whose binary was replaced since it
// started: a restart runs the newer build.
const NewerBuild = "newer build installed"

// LogLines is how much of a server's serve.log Log shows.
const LogLines = 40

// Seams: tests list fixed servers and stop nothing.
var (
	listServers   = pageserve.Servers
	listLeftovers = pageserve.Leftovers
	startTool     = pageserve.StartTool
	stopTool      = pageserve.StopTool
	restartTool   = pageserve.RestartTool
	stopLeftover  = pageserve.StopLeftover
	now           = time.Now
)

// Row is one page server.
type Row struct {
	Tool    string `json:"tool"`
	Up      bool   `json:"up"`
	Address string `json:"address,omitempty"` // host:port
	URL     string `json:"url,omitempty"`     // its page, to open
	Version string `json:"version,omitempty"`
	Uptime  string `json:"uptime,omitempty"`
	Note    string `json:"note,omitempty"`
}

// Leftover is one process of a server's tool that no server holds.
type Leftover struct {
	Tool    string `json:"tool"`
	PID     int    `json:"pid"`
	Age     string `json:"age"`
	Command string `json:"command"`
}

// List is the three groups both views show.
type List struct {
	Up   []Row      `json:"up"`
	Down []Row      `json:"down"`
	Left []Leftover `json:"left"`
}

// Read lists the servers and the left-over processes.
func Read() (List, error) {
	ss, err := listServers()
	if err != nil {
		return List{}, err
	}
	l := List{Up: []Row{}, Down: []Row{}, Left: []Leftover{}}
	for _, s := range ss {
		r := Row{Tool: s.Tool, Up: s.Up}
		if !s.Up {
			l.Down = append(l.Down, r)
			continue
		}
		e := s.Entry
		r.URL, r.Version = e.URL, e.Version
		r.Address = host(e.Base)
		if r.Address == "" {
			r.Address = host(e.URL)
		}
		if !e.StartedAt.IsZero() {
			r.Uptime = Duration(now().Sub(e.StartedAt))
		}
		if e.Stale {
			r.Note = NewerBuild
		}
		l.Up = append(l.Up, r)
	}
	lo, err := listLeftovers(ss)
	if err != nil {
		return l, err
	}
	for _, p := range lo {
		l.Left = append(l.Left, Leftover{Tool: p.Tool, PID: p.PID, Age: Duration(p.Age), Command: p.Command})
	}
	return l, nil
}

func host(u string) string {
	p, err := url.Parse(u)
	if err != nil {
		return ""
	}
	return p.Host
}

// Duration is d in its two largest units: 45s, 12m, 3h12m, 1d2h.
func Duration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	s := int(d / time.Second)
	switch days, h, m := s/86400, s/3600%24, s/60%60; {
	case days > 0:
		return fmt.Sprintf("%dd%dh", days, h)
	case h > 0:
		return fmt.Sprintf("%dh%dm", h, m)
	case m > 0:
		return fmt.Sprintf("%dm", m)
	default:
		return fmt.Sprintf("%ds", s)
	}
}

// server finds tool among the listed servers. Every action names a listed
// server, so a name from a page or key never runs another program.
func server(tool string) (pageserve.Server, error) {
	ss, err := listServers()
	if err != nil {
		return pageserve.Server{}, err
	}
	for _, s := range ss {
		if s.Tool == tool {
			return s, nil
		}
	}
	return pageserve.Server{}, fmt.Errorf("%q is not a page server", tool)
}

// Log is the last LogLines lines of tool's serve.log.
func Log(tool string) (string, error) {
	s, err := server(tool)
	if err != nil {
		return "", err
	}
	f, err := os.Open(s.LogPath)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	// The tail is in the last 64 KiB unless lines run past 1.6 KiB each.
	const window = 64 << 10
	if st, err := f.Stat(); err == nil && st.Size() > window {
		if _, err := f.Seek(-window, io.SeekEnd); err != nil {
			return "", err
		}
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	if len(lines) > LogLines {
		lines = lines[len(lines)-LogLines:]
	}
	return strings.Join(lines, "\n") + "\n", nil
}

// Start ensures tool's serve.
func Start(tool string) error {
	if _, err := server(tool); err != nil {
		return err
	}
	_, err := startTool(tool)
	return err
}

// Stop asks tool's serve to exit; Start undoes it.
func Stop(ctx context.Context, tool string) error {
	if _, err := server(tool); err != nil {
		return err
	}
	return stopTool(ctx, tool)
}

// Restart stops tool's serve, when it runs, then starts it.
func Restart(ctx context.Context, tool string) error {
	if _, err := server(tool); err != nil {
		return err
	}
	_, err := restartTool(ctx, tool)
	return err
}

// StopLeftover ends pid, only when Leftovers lists it now: a pid that ended
// may have been reused by another process.
func StopLeftover(pid int) error {
	ss, err := listServers()
	if err != nil {
		return err
	}
	lo, err := listLeftovers(ss)
	if err != nil {
		return err
	}
	for _, p := range lo {
		if p.PID == pid {
			return stopLeftover(pid)
		}
	}
	return fmt.Errorf("pid %d is not a left-over process", pid)
}

// StopLeftovers ends every process Leftovers lists now, and says how many.
func StopLeftovers() (int, error) {
	ss, err := listServers()
	if err != nil {
		return 0, err
	}
	lo, err := listLeftovers(ss)
	if err != nil {
		return 0, err
	}
	var errs []error
	n := 0
	for _, p := range lo {
		if err := stopLeftover(p.PID); err != nil {
			errs = append(errs, fmt.Errorf("pid %d: %w", p.PID, err))
			continue
		}
		n++
	}
	return n, errors.Join(errs...)
}
