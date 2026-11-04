package pages

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/subaud-tools/tackle/internal/proj/servers"
	"github.com/subaud-tools/tools-common/pageserve"
)

// The servers page is proj serve's /: every family page server, the
// left-over processes, and the actions on them, through the servers package.

// actionTimeout bounds a stop or restart, which waits for the server to go.
const actionTimeout = 20 * time.Second

// serversAPI is the servers page's API, under the token like every /api/
// route: GET /api/servers, POST /api/servers/<tool>/{start,stop,restart},
// GET /api/servers/<tool>/log, POST /api/leftovers/<pid>/stop and POST
// /api/leftovers/stop.
func serversAPI(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/servers", func(w http.ResponseWriter, _ *http.Request) {
		l, err := servers.Read()
		if err != nil {
			fail(w, err)
			return
		}
		reply(w, l)
	})
	mux.HandleFunc("POST /api/servers/{tool}/{action}", func(w http.ResponseWriter, r *http.Request) {
		tool := r.PathValue("tool")
		ctx, cancel := context.WithTimeout(r.Context(), actionTimeout)
		defer cancel()
		var err error
		switch r.PathValue("action") {
		case "start":
			err = servers.Start(tool)
		case "stop":
			err = servers.Stop(ctx, tool)
		case "restart":
			err = servers.Restart(ctx, tool)
		default:
			http.NotFound(w, r)
			return
		}
		if err != nil {
			fail(w, err)
			return
		}
		reply(w, map[string]bool{"ok": true})
	})
	mux.HandleFunc("GET /api/servers/{tool}/log", func(w http.ResponseWriter, r *http.Request) {
		text, err := servers.Log(r.PathValue("tool"))
		if err != nil {
			fail(w, err)
			return
		}
		reply(w, map[string]string{"log": text})
	})
	mux.HandleFunc("POST /api/leftovers/{pid}/stop", func(w http.ResponseWriter, r *http.Request) {
		pid, err := strconv.Atoi(r.PathValue("pid"))
		if err != nil {
			http.Error(w, "bad pid", http.StatusBadRequest)
			return
		}
		if err := servers.StopLeftover(pid); err != nil {
			fail(w, err)
			return
		}
		reply(w, map[string]bool{"ok": true})
	})
	mux.HandleFunc("POST /api/leftovers/stop", func(w http.ResponseWriter, _ *http.Request) {
		n, err := servers.StopLeftovers()
		if err != nil {
			fail(w, err)
			return
		}
		reply(w, map[string]int{"stopped": n})
	})
}

func reply(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// fail answers 409 for a tool or pid not listed now, else 500.
func fail(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	if errors.Is(err, servers.ErrNotListed) {
		code = http.StatusConflict
	}
	http.Error(w, err.Error(), code)
}

// watch emits a "servers" event whenever the list changes, so an open page
// reloads it. It reads only while a page is open, since each read runs ps;
// the first read after none was open emits too, since the list may have
// changed meanwhile.
func watch(ctx context.Context, events *pageserve.Events, every time.Duration) {
	var last servers.List
	fresh := false
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
		if events.Streams() == 0 {
			fresh = false
			continue
		}
		now, err := servers.Read()
		if err != nil || (fresh && reflect.DeepEqual(now, last)) {
			continue
		}
		last, fresh = now, true
		events.Emit("servers", nil)
	}
}
