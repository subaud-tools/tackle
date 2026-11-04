// Package pages is `proj serve`: one local server per machine that shows any
// session's split page, the family's splitpage: its terminal on the left,
// xterm.js over a websocket to a tmux client in a pty, and its tabs on the right.
// Its / is the servers page: every family page server, to see and control.
// It runs on pageserve, so it starts, is found, goes stale and stops as the
// family's other page servers do.
package pages

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/subaud-tools/tools-common/localweb/page"
	"github.com/subaud-tools/tools-common/pageserve"
	"github.com/subaud-tools/tools-common/splitpage"
)

const tool = "proj"

// assets is the server's own servers page at /; everything else is the split
// page's.
//
//go:embed assets
var assets embed.FS

// termPage is the split page every session's /term serves; the page reads its
// session and socket from its own query.
var termPage = splitpage.Page(splitpage.Config{Title: "proj", Socket: "/api/term/ws", Take: true})

// Live is where proj serve's entry is: StateDir/live/serve.json.
func Live() pageserve.Live { return pageserve.LiveFor(tool) }

// Process is how proj serve starts in the background, from proj window or
// `proj serve`.
func Process() pageserve.Process {
	return pageserve.Process{Live: Live(), Args: []string{"serve", "--foreground", "--no-open"}}
}

// Client finds proj serve, else starts it from this binary.
func Client() *pageserve.Client {
	p := Process()
	return pageserve.NewClient(p.Live, func() (pageserve.Entry, error) { return p.Ensure(pageserve.Executable()) })
}

// Options configures Run.
type Options struct {
	Port    int // 0: remembered, else free
	Version string
	Log     io.Writer // one line per start/stop
	Ready   func(url string)
	Every   time.Duration // how often an open servers page's list is checked; default 2s
}

// Run serves until ctx ends or /api/stop. pageserve writes the entry while
// running, keeps the token across restarts and the port in StateDir/port.
func Run(ctx context.Context, o Options) error {
	sub, err := fs.Sub(assets, "assets")
	if err != nil {
		return err
	}
	if o.Every <= 0 {
		o.Every = 2 * time.Second
	}
	api := http.NewServeMux()
	api.Handle("/api/term/", newBridge().handler())
	api.Handle("/api/split/", splitpage.NewTabs().Handler())
	serversAPI(api)
	events := pageserve.NewEvents()
	return pageserve.Run(ctx, pageserve.Options{
		Live: Live(), Assets: page.With(siteFS{own: sub}), API: api, Events: events,
		Port: o.Port, Version: o.Version, Log: o.Log, Ready: o.Ready,
		OnStart: func(ctx context.Context, _ string) { go watch(ctx, events, o.Every) },
	})
}

// siteFS serves the split page at term, the server's own files (the servers page),
// and the split page's script and styles beside them.
type siteFS struct{ own fs.FS }

func (s siteFS) Open(name string) (fs.File, error) {
	if name == "term" {
		return termPage.Open("index.html")
	}
	if f, err := s.own.Open(name); err == nil {
		return f, nil
	}
	return termPage.Open(name)
}

// TermURL is the address to open for a session's split page, with the token,
// to take the session from any other client; the page drops both on load.
func TermURL(addr, token, inst, socket, session string) string {
	return "http://" + addr + "/term?" + url.Values{"socket": {socket}, "session": {session}, "i": {inst}, "t": {token}, "take": {"1"}}.Encode()
}

// Instance names the page server a terminal page came from, in its address,
// where it stays after the page drops its token: its token and its pid. The
// token and port are kept across restarts, so a page left dead by one has the
// current address; the pid tells it from a page of the running server. It is a
// hash, so the address never carries the token itself.
func Instance(token string, pid int) string {
	h := sha256.Sum256([]byte("proj pages instance " + token + " " + strconv.Itoa(pid)))
	return hex.EncodeToString(h[:6])
}

// ViewURL is TermURL without the take: a page that attaches only when no other
// client has the session, else offers to show it there.
func ViewURL(addr, token, inst, socket, session string) string {
	return "http://" + addr + "/term?" + url.Values{"socket": {socket}, "session": {session}, "i": {inst}, "t": {token}}.Encode()
}

// TermMatches says whether pageURL is the terminal page of session on socket:
// a loopback /term page with exactly those two values. The port is left out,
// since a restarted server may land on another.
func TermMatches(pageURL, socket, session string) bool {
	u, err := url.Parse(pageURL)
	if err != nil || u.Scheme != "http" || u.Hostname() != "127.0.0.1" || u.Path != "/term" {
		return false
	}
	q := u.Query()
	return q.Get("socket") == socket && q.Get("session") == session
}
