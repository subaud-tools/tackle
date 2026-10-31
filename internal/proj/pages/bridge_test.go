package pages

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/subaud-tools/tools-common/localweb"
)

const testSession = "pages-test/w"

// tmuxServer starts a throwaway tmux server on a socket in a new temp dir,
// with one detached session running sh, and returns the socket. $TMUX is set,
// as it is for the real server, which runs inside tmux session `pages`.
func tmuxServer(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("tmux"); err != nil {
		t.Skip("tmux not installed")
	}
	// Kept short: macOS caps unix socket paths at 104 bytes.
	tmp, err := os.MkdirTemp("/tmp", "pjp")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmp) })
	socket := filepath.Join(tmp, "s")
	if !strings.HasPrefix(socket, tmp) {
		t.Fatalf("socket %s is outside the temp dir", socket)
	}
	t.Setenv("TMUX", "")
	out, err := exec.Command("tmux", "-f", "/dev/null", "-S", socket, "new-session", "-d",
		"-s", testSession, "-x", "80", "-y", "24", "/bin/sh").CombinedOutput()
	if err != nil {
		t.Fatalf("new-session: %v: %s", err, out)
	}
	t.Cleanup(func() { _ = exec.Command("tmux", "-S", socket, "kill-server").Run() })
	t.Setenv("TMUX", "/tmp/elsewhere,1,0")
	return socket
}

func tmuxOut(t *testing.T, socket string, args ...string) string {
	t.Helper()
	out, err := exec.Command("tmux", append([]string{"-S", socket}, args...)...).Output()
	if err != nil {
		t.Fatalf("tmux %v: %v", args, err)
	}
	return string(out)
}

// startServer runs the page server until the test ends, its state in a temp dir.
func startServer(t *testing.T) *localweb.Server {
	t.Helper()
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	srv, err := Start(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return srv
}

func wsURL(srv *localweb.Server, socket, session string) string {
	q := url.Values{"socket": {socket}, "session": {session}}
	return "ws://" + srv.Addr() + "/api/term/ws?" + q.Encode()
}

func withToken(srv *localweb.Server) *websocket.DialOptions {
	return &websocket.DialOptions{HTTPHeader: http.Header{localweb.TokenHeader: {srv.Token}}}
}

func dial(t *testing.T, srv *localweb.Server, socket string) *websocket.Conn {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, _, err := websocket.Dial(ctx, wsURL(srv, socket, testSession), withToken(srv)) //nolint:bodyclose // a 101 has no body; the conn is closed on cleanup
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = c.CloseNow() })
	go func() { // keep draining, as the page does, so the bridge never blocks
		for {
			if _, _, err := c.Read(context.Background()); err != nil {
				return
			}
		}
	}()
	return c
}

// waitFor polls until check is true.
func waitFor(t *testing.T, what string, check func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if check() {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

func TestSocketWithoutTokenIsRefused(t *testing.T) {
	socket := tmuxServer(t)
	srv := startServer(t)
	_, resp, err := websocket.Dial(context.Background(), wsURL(srv, socket, testSession), nil) //nolint:bodyclose // Dial's refused response holds only a copy of the body
	if err == nil || resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("got %v, %v; want a 401", resp, err)
	}
}

func TestSocketFromAnotherSiteIsRefused(t *testing.T) {
	socket := tmuxServer(t)
	srv := startServer(t)
	o := withToken(srv)
	o.HTTPHeader.Set("Origin", "https://evil.example")
	_, resp, err := websocket.Dial(context.Background(), wsURL(srv, socket, testSession), o) //nolint:bodyclose // Dial's refused response holds only a copy of the body
	if err == nil || resp == nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("got %v, %v; want a 403", resp, err)
	}
	if got := tmuxOut(t, socket, "list-clients"); got != "" {
		t.Fatalf("clients = %q, want none", got)
	}
}

func TestBytesTypedReachThePane(t *testing.T) {
	socket := tmuxServer(t)
	c := dial(t, startServer(t), socket)
	waitFor(t, "the client to attach", func() bool { return tmuxOut(t, socket, "list-clients") != "" })
	if err := c.Write(context.Background(), websocket.MessageBinary, []byte("echo hi\r")); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "hi in the pane", func() bool {
		for _, l := range strings.Split(tmuxOut(t, socket, "capture-pane", "-p", "-t", "="+testSession+":"), "\n") {
			if strings.TrimSpace(l) == "hi" {
				return true
			}
		}
		return false
	})
}

func TestResizeMessageChangesTheClientSize(t *testing.T) {
	socket := tmuxServer(t)
	c := dial(t, startServer(t), socket)
	waitFor(t, "the client to attach", func() bool { return tmuxOut(t, socket, "list-clients") != "" })
	msg := `{"type":"resize","cols":101,"rows":33}`
	if err := c.Write(context.Background(), websocket.MessageText, []byte(msg)); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "the client's width", func() bool {
		return strings.TrimSpace(tmuxOut(t, socket, "list-clients", "-F", "#{client_width}")) == "101"
	})
}

func TestMissingSessionClosesWithAnErrorFrame(t *testing.T) {
	socket := tmuxServer(t)
	srv := startServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, _, err := websocket.Dial(ctx, wsURL(srv, socket, "no-such-session"), withToken(srv)) //nolint:bodyclose // a 101 has no body
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.CloseNow() }()
	_, _, err = c.Read(ctx)
	var ce websocket.CloseError
	if !errors.As(err, &ce) || ce.Code != closeNoSession {
		t.Fatalf("read = %v, want close code %d", err, closeNoSession)
	}
	if got := tmuxOut(t, socket, "list-clients"); got != "" {
		t.Fatalf("clients = %q, want none", got)
	}
	if got := strings.TrimSpace(tmuxOut(t, socket, "list-sessions", "-F", "#{session_name}")); got != testSession {
		t.Fatalf("sessions = %q, want only %s", got, testSession)
	}
}

func TestAttachTellsTmuxThePageSupportsSync(t *testing.T) {
	socket := tmuxServer(t)
	dial(t, startServer(t), socket)
	waitFor(t, "the client to attach", func() bool { return tmuxOut(t, socket, "list-clients") != "" })
	if got := tmuxOut(t, socket, "show", "-s", "terminal-features"); !strings.Contains(got, "xterm-256color:Sync") {
		t.Fatalf("terminal-features = %q", got)
	}
}
