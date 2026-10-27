// Package serve is `cull serve`: the review page's API over the store, served
// by pageserve's one local server per machine, which streams live events and
// writes the entry other processes find (and stop) it by.
package serve

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"time"

	"github.com/subaud-tools/tackle/internal/cull/store"
	"github.com/subaud-tools/tools-common/localweb/page"
	"github.com/subaud-tools/tools-common/pageserve"
)

//go:embed assets
var assets embed.FS

// Live is where cull serve's entry is: StateDir/live/serve.json.
func Live() pageserve.Live { return pageserve.LiveFor("cull") }

// Running returns the live serve's entry, or pageserve.ErrNotRunning.
func Running() (pageserve.Entry, error) { return Live().Running() }

// Server is cull's API over the store.
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
	PollEvery time.Duration // run-change poll; default 2s
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

// Watch polls the store every `every` for new runs and emits a run event for
// each project whose latest run id changed (the first read only records).
func (s *Server) Watch(ctx context.Context, every time.Duration) {
	var last map[int64]int64
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		cur, err := s.st.LatestRunIDs(ctx)
		if err == nil {
			if last != nil {
				for p, id := range cur {
					if last[p] != id {
						s.events.Emit("run", map[string]int64{"project": p, "run": id})
					}
				}
			}
			last = cur
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}
