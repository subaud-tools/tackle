package serve

// Run's server, entry, idle exit, stop and token are pageserve's, tested
// there; this tests what sift adds: its page and routes, and its entry where
// sift's has always been.

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/internal/sift/store"
	"github.com/subaud-tools/tools-common/localweb"
	"github.com/subaud-tools/tools-common/pageserve"
)

func TestRunServesSiftsRoutesFromItsEntry(t *testing.T) {
	home := t.TempDir()
	t.Setenv("SIFT_HOME", home)
	st, err := store.Open(context.Background(), filepath.Join(t.TempDir(), "sift.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	ready, done := make(chan string, 1), make(chan error, 1)
	go func() {
		done <- Run(context.Background(), st, Options{Version: "v", Ready: func(u string) { ready <- u }})
	}()
	var url string
	select {
	case url = <-ready:
	case err := <-done:
		t.Fatalf("run exited: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("not ready")
	}
	if _, err := os.Stat(filepath.Join(home, "state", "live", "serve.json")); err != nil {
		t.Fatalf("entry not at SIFT_HOME/state/live/serve.json: %v", err)
	}
	e, err := Running()
	if err != nil || e.PID != os.Getpid() || e.Version != "v" || e.URL != url {
		t.Fatalf("%+v %v", e, err)
	}
	for _, path := range []string{"/api/review", "/api/agent/status", "/api/poll"} {
		req, _ := http.NewRequest("GET", e.Base+path, nil)
		req.Header.Set(localweb.TokenHeader, e.Token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("%s: %d", path, resp.StatusCode)
		}
	}
	if resp, err := http.Get(e.Base + "/"); err != nil || resp.StatusCode != 200 {
		t.Fatalf("index: %v %v", resp, err)
	} else {
		_ = resp.Body.Close()
	}
	if err := pageserve.Stop(context.Background(), Live()); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("did not exit")
	}
	if _, err := Running(); !errors.Is(err, pageserve.ErrNotRunning) {
		t.Fatalf("running after stop: %v", err)
	}
}
