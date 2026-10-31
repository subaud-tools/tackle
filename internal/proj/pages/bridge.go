package pages

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/coder/websocket"
	"github.com/creack/pty"
)

// closeNoSession is the close code when the session does not exist.
const closeNoSession websocket.StatusCode = 4404

// maxInput bounds one message from the page: a large paste, not a stream.
const maxInput = 16 << 20

// bridge joins websockets to tmux clients in ptys.
type bridge struct {
	mu     sync.Mutex
	synced map[string]bool // sockets told the page supports synchronized output
}

func newBridge() *bridge { return &bridge{synced: map[string]bool{}} }

func (b *bridge) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("GET /api/term/ws", b.serveTerm)
	return mux
}

// serveTerm upgrades, then runs one tmux client for the socket's lifetime.
// localweb's Origin check skips GETs, so this upgrade relies on Accept's own:
// it refuses an Origin whose host is not the request's.
func (b *bridge) serveTerm(w http.ResponseWriter, r *http.Request) {
	socket, session := r.URL.Query().Get("socket"), r.URL.Query().Get("session")
	if socket == "" || session == "" {
		http.Error(w, "socket and session are required", http.StatusBadRequest)
		return
	}
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer func() { _ = c.CloseNow() }()
	c.SetReadLimit(maxInput)
	target := "=" + session
	if tmuxCmd(socket, "has-session", "-t", target).Run() != nil {
		_ = c.Close(closeNoSession, "no session "+session)
		return
	}
	b.syncOnce(socket)

	cmd := tmuxCmd(socket, "-u", "attach", "-t", target)
	cmd.Env = attachEnv()
	tty, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 24, Cols: 80})
	if err != nil {
		_ = c.Close(websocket.StatusInternalError, "attach: "+err.Error())
		return
	}
	defer func() {
		_ = tty.Close()
		_ = cmd.Process.Kill() // the client only: the session keeps running
		_ = cmd.Wait()
	}()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	go func() {
		defer cancel()
		buf := make([]byte, 64<<10)
		for {
			n, err := tty.Read(buf)
			if n > 0 && c.Write(ctx, websocket.MessageBinary, buf[:n]) != nil {
				return
			}
			if err != nil {
				_ = c.Close(websocket.StatusNormalClosure, "detached")
				return
			}
		}
	}()
	for {
		kind, data, err := c.Read(ctx)
		if err != nil {
			return
		}
		if kind == websocket.MessageText {
			applyResize(tty, data)
			continue
		}
		if _, err := tty.Write(data); err != nil {
			return
		}
	}
}

// syncOnce tells the session's server, once per socket, that this terminal
// supports synchronized output, so fast output does not tear.
func (b *bridge) syncOnce(socket string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.synced[socket] {
		return
	}
	if tmuxCmd(socket, "set", "-as", "terminal-features", "xterm-256color:Sync").Run() == nil {
		b.synced[socket] = true
	}
}

// applyResize sets the pty's size from a {"type":"resize"} message; tmux
// follows the pty. Anything else is ignored.
func applyResize(tty *os.File, data []byte) {
	var m struct {
		Type string `json:"type"`
		Cols uint16 `json:"cols"`
		Rows uint16 `json:"rows"`
	}
	if json.Unmarshal(data, &m) != nil || m.Type != "resize" || m.Cols == 0 || m.Rows == 0 {
		return
	}
	_ = pty.Setsize(tty, &pty.Winsize{Rows: m.Rows, Cols: m.Cols})
}

// tmuxCmd is tmux against socket: a full path (as $TMUX gives it) or a server
// name. The args may lead with a client flag such as -u.
func tmuxCmd(socket string, args ...string) *exec.Cmd {
	flag := "-L"
	if strings.Contains(socket, "/") {
		flag = "-S"
	}
	return exec.Command("tmux", append([]string{flag, socket}, args...)...) //nolint:gosec // G204: tmux argv, no shell; the socket comes from an authenticated request
}

// attachEnv is this process's environment for a tmux client: the server itself
// runs inside tmux, and a client refuses to start with $TMUX set.
func attachEnv() []string {
	var env []string
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "TMUX=") && !strings.HasPrefix(kv, "TMUX_PANE=") && !strings.HasPrefix(kv, "TERM=") {
			env = append(env, kv)
		}
	}
	return append(env, "TERM=xterm-256color")
}
