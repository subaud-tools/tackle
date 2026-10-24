package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/internal/sift/channel"
	"github.com/subaud-tools/tackle/internal/sift/serve"
)

// Inside pi's claude-bridge the hosting pi session's channel delivers
// the user's decisions, so `sift watch` prints one stop line and returns at once
// without reaching sift serve: the mod then never restarts it.
func TestWatchStopsInABridgeChild(t *testing.T) {
	t.Setenv("AGENT_SESSION_CHILD", "1")
	t.Setenv("AGENT_SESSION_ID", "x")
	var mu sync.Mutex
	finds := 0
	cl := channel.NewClient(nil)
	cl.Find = func() (serve.Advert, error) {
		mu.Lock()
		finds++
		mu.Unlock()
		return serve.Advert{Base: "http://127.0.0.1:1", Token: "tok"}, nil
	}
	var out bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	done := make(chan error, 1)
	go func() { done <- watch(ctx, &out, io.Discard, cl) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("watch = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("sift watch started its loop in a bridge child")
	}
	mu.Lock()
	n := finds
	mu.Unlock()
	if n != 0 {
		t.Fatalf("watch looked for sift serve %d times in a bridge child", n)
	}
	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	var stop struct {
		Stop string `json:"stop"`
	}
	if len(lines) != 1 || json.Unmarshal([]byte(lines[0]), &stop) != nil || stop.Stop == "" {
		t.Fatalf("want exactly one stop line, got %q", out.String())
	}
}

// Outside the bridge `sift watch` prints each Send with `text`, the whole
// prompt the plugin's mod submits as the turn.
func TestWatchLineCarriesText(t *testing.T) {
	t.Setenv("AGENT_SESSION_CHILD", "")
	t.Setenv("AGENT_SESSION_ID", "")
	t.Setenv("CLAUDE_CODE_SESSION_ID", "s1")
	var mu sync.Mutex
	served := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/agent/presence":
			_, _ = io.WriteString(w, "{}")
		case "/api/agent/wait":
			mu.Lock()
			first := !served && r.URL.Query().Get("session") == "s1"
			served = served || first
			mu.Unlock()
			if first {
				_, _ = io.WriteString(w, `{"id":1,"round":2,"text":"The user sent their decisions."}`)
				return
			}
			time.Sleep(20 * time.Millisecond)
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	cl := channel.NewClient(nil)
	cl.Find = func() (serve.Advert, error) { return serve.Advert{Base: srv.URL, Token: "tok"}, nil }

	r, w := io.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go func() { _ = watch(ctx, w, io.Discard, cl) }()

	lines := make(chan string, 1)
	go func() {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1<<20), 1<<24)
		if sc.Scan() {
			lines <- sc.Text()
		}
	}()
	var line string
	select {
	case line = <-lines:
	case <-time.After(5 * time.Second):
		t.Fatal("no line printed for the Send")
	}
	var ev struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(line), &ev); err != nil {
		t.Fatalf("not one JSON line: %q (%v)", line, err)
	}
	if !strings.HasPrefix(ev.Text, "sift event") {
		t.Fatalf("text = %q, want it to start with \"sift event\"", ev.Text)
	}
}
