// Package serve is `sift serve`: the review page's API over the store, served
// by pageserve's one local server per machine, which streams live events and
// writes the entry other processes find (and stop) it by.
package serve

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"time"

	"github.com/subaud-tools/tackle/internal/sift/store"
	"github.com/subaud-tools/tools-common/localweb/page"
	"github.com/subaud-tools/tools-common/pageserve"
)

//go:embed assets
var assets embed.FS

// Live is where sift serve's entry is: StateDir/live/serve.json.
func Live() pageserve.Live { return pageserve.LiveFor("sift") }

// Running returns the live serve's entry, or pageserve.ErrNotRunning.
func Running() (pageserve.Entry, error) { return Live().Running() }

// Server is sift's API over the store.
type Server struct {
	st *store.Store

	events *pageserve.Events

	ag       agents
	now      func() time.Time // presence clock; nil: time.Now
	waitUnit time.Duration    // one "second" of a wait timeout; 0: a second
	pageURL  string           // the page's address with the token, for agents
}

// New returns a Server over st.
func New(st *store.Store) *Server {
	return &Server{st: st, events: pageserve.NewEvents()}
}

// Options configures Run.
type Options struct {
	Port      int           // 0: remembered, else free
	Idle      time.Duration // shut down after this long with no API traffic (0: never)
	Version   string
	Log       io.Writer // one line per start/stop
	Ready     func(url string)
	PollEvery time.Duration // round-change poll; default 2s
}

// Run serves until ctx ends, /api/stop, or the idle limit. pageserve writes
// the entry while running and removes it on the way out.
func Run(ctx context.Context, st *store.Store, o Options) error {
	if o.PollEvery <= 0 {
		o.PollEvery = 2 * time.Second
	}
	assetsFS, err := fs.Sub(assets, "assets")
	if err != nil {
		return err
	}
	s := New(st)
	return pageserve.Run(ctx, pageserve.Options{
		Live: Live(), Assets: page.With(assetsFS), API: s.Handler(), Events: s.events,
		Port: o.Port, Idle: o.Idle, Version: o.Version, Log: o.Log, Ready: o.Ready,
		OnStart: func(ctx context.Context, url string) {
			s.pageURL = url
			go s.Watch(ctx, o.PollEvery)
		},
	})
}

// Watch polls the store every `every` and emits a round event when the
// latest round or its revision changed: a new check, the agent's proposals,
// an apply (the first read only records).
func (s *Server) Watch(ctx context.Context, every time.Duration) {
	var lastID, lastRev int64 = -1, -1
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		id, rev, err := s.st.Latest(ctx)
		if err != nil {
			id, rev = 0, 0
		}
		if lastID >= 0 && (id != lastID || rev != lastRev) {
			s.events.Emit("round", map[string]int64{"round": id, "rev": rev})
		}
		lastID, lastRev = id, rev
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}
