package pages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/creack/pty"
)

// closeNoSession is the close code when the session does not exist.
const closeNoSession websocket.StatusCode = 4404

// closeDetached is the close code when another terminal has the session.
const closeDetached websocket.StatusCode = 4409

// ClientOption is the tmux session option holding the pid of the page's latest
// client, which tells it from a Ghostty terminal's.
const ClientOption = "@chrome_client"

// maxInput bounds one message from the page: a large paste, not a stream.
const maxInput = 16 << 20

// bridge joins websockets to tmux clients in ptys.
type bridge struct {
	mu     sync.Mutex
	synced map[string]bool // sockets told the page supports synchronized output

	// Seams, so tests never drive the real Ghostty.
	closeTerminal func(tty string) error
	ttyBusy       func(tty string) (bool, error)
	closeWait     time.Duration // how long a pushed-out tty gets to empty
}

func newBridge() *bridge {
	return &bridge{synced: map[string]bool{}, closeTerminal: CloseTerminal, ttyBusy: ttyHasProcess, closeWait: 2 * time.Second}
}

// CloseTerminal closes the Ghostty terminal with the given tty, if there is one.
// proj sets it: it holds the Ghostty scripts, and it imports this package.
var CloseTerminal = func(string) error { return nil }

// pushedOut lists the ttys a take of target will detach, leaving out the page's
// own client: its pty goes back to the system, and a new terminal may get it.
func (b *bridge) pushedOut(socket, target string) []string {
	out, err := tmuxCmd(socket, "list-clients", "-t", target, "-F", "#{client_pid} #{client_tty}").Output()
	if err != nil {
		return nil
	}
	page, _ := tmuxCmd(socket, "show-options", "-v", "-q", "-t", target+":", ClientOption).Output()
	var ttys []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		pid, tty, ok := strings.Cut(line, " ")
		if ok && pid != strings.TrimSpace(string(page)) {
			ttys = append(ttys, tty)
		}
	}
	return ttys
}

// closePushedOut closes the Ghostty terminal of a tty a take detached, once
// nothing runs on the tty. A tty still running something after closeWait keeps
// its terminal: that is a shell the client dropped back to, not an empty window.
// A failure is only logged: the take has already happened.
func (b *bridge) closePushedOut(tty string) {
	deadline := time.Now().Add(b.closeWait)
	for {
		busy, err := b.ttyBusy(tty)
		if err != nil {
			log.Printf("proj pages: read %s: %v", tty, err)
			return
		}
		if !busy {
			break
		}
		if time.Now().After(deadline) {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err := b.closeTerminal(tty); err != nil {
		log.Printf("proj pages: close the terminal on %s: %v", tty, err)
	}
}

// ttyHasProcess says whether a process other than a zombie is on tty. ps exits 1
// with no processes, and also on a tty it can't read, which only the latter
// explains on stderr.
func ttyHasProcess(tty string) (bool, error) {
	out, err := exec.Command("ps", "-t", tty, "-o", "stat=").Output() //nolint:gosec // G204: ps argv, no shell; the tty comes from tmux
	var ee *exec.ExitError
	if err != nil && (!errors.As(err, &ee) || len(ee.Stderr) > 0) {
		return false, fmt.Errorf("ps: %w", err)
	}
	for _, stat := range strings.Fields(string(out)) {
		if !strings.HasPrefix(stat, "Z") {
			return true, nil
		}
	}
	return false, nil
}

func (b *bridge) handler() http.Handler {
	mux := http.NewServeMux()
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
	// Only take=1 takes the session: a page Chrome restored or reloaded must leave
	// another terminal's view alone.
	attach := []string{"-u", "attach", "-t", target}
	var pushedOut []string
	if r.URL.Query().Get("take") == "1" {
		attach = []string{"-u", "attach", "-d", "-t", target}
		pushedOut = b.pushedOut(socket, target)
	} else if out, err := tmuxCmd(socket, "list-clients", "-t", target).Output(); err == nil && len(out) > 0 {
		_ = c.Close(closeDetached, "shown elsewhere")
		return
	}
	b.syncOnce(socket)

	cmd := tmuxCmd(socket, attach...)
	cmd.Env = attachEnv()
	tty, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 24, Cols: 80})
	if err != nil {
		_ = c.Close(websocket.StatusInternalError, "attach: "+err.Error())
		return
	}
	for _, t := range pushedOut {
		go b.closePushedOut(t)
	}
	// A failure only leaves proj unable to say the page shows the session.
	_ = tmuxCmd(socket, "set-option", "-t", target+":", ClientOption, strconv.Itoa(cmd.Process.Pid)).Run()
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
				if tmuxCmd(socket, "has-session", "-t", target).Run() == nil {
					_ = c.Close(closeDetached, "detached")
				} else {
					_ = c.Close(websocket.StatusNormalClosure, "session ended")
				}
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
