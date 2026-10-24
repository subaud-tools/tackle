package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/internal/imprint/channel"
)

// Inside pi's claude-bridge the hosting pi session's channel delivers
// imprint's messages, so `imprint watch` prints one stop line and returns at
// once without following a spool: the mod then never restarts it.
func TestWatchStopsInABridgeChild(t *testing.T) {
	t.Setenv("AGENT_SESSION_CHILD", "1")
	t.Setenv("AGENT_SESSION_ID", "x")
	var out bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	done := make(chan error, 1)
	go func() { done <- watch(ctx, t.TempDir(), &out, io.Discard) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("watch = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("imprint watch followed the spool in a bridge child")
	}
	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	var stop struct {
		Stop string `json:"stop"`
	}
	if len(lines) != 1 || json.Unmarshal([]byte(lines[0]), &stop) != nil || stop.Stop == "" {
		t.Fatalf("want exactly one stop line, got %q", out.String())
	}
}

// Outside the bridge `imprint watch` prints each message the channel spools
// with `text`, the whole prompt the plugin's mod submits as the turn.
func TestWatchLineCarriesText(t *testing.T) {
	t.Setenv("AGENT_SESSION_CHILD", "")
	t.Setenv("AGENT_SESSION_ID", "")
	t.Setenv("CLAUDE_CODE_SESSION_ID", "s1")
	dir := t.TempDir()
	spool, err := channel.OpenSpool(dir, "s1")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = spool.Close() })
	if err := spool.Notify("logged in as demo-user", map[string]string{"op": "imprint_login"}); err != nil {
		t.Fatal(err)
	}

	r, w := io.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { _ = watch(ctx, dir, w, io.Discard); close(done) }()
	// watch saves its cursor under dir: it must be gone before dir is removed.
	t.Cleanup(func() { cancel(); _ = r.Close(); <-done })

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
		t.Fatal("no line printed for the spooled message")
	}
	var ev struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(line), &ev); err != nil {
		t.Fatalf("not one JSON line: %q (%v)", line, err)
	}
	if !strings.HasPrefix(ev.Text, "imprint event") {
		t.Fatalf("text = %q, want it to start with \"imprint event\"", ev.Text)
	}
}
