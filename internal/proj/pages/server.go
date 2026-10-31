// Package pages is `proj pages`: one long-running local server that shows any
// session's terminal in a browser page, xterm.js over a websocket to a tmux
// client in a pty.
package pages

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/localweb"
)

const tool = "proj"

//go:embed assets
var assets embed.FS

// Start serves the page until ctx ends. The token is kept in StateDir, so
// `proj window` can read it and a restarted server keeps open tabs working.
func Start(ctx context.Context) (*localweb.Server, error) {
	token, err := loadToken()
	if err != nil {
		return nil, err
	}
	sub, err := fs.Sub(assets, "assets")
	if err != nil {
		return nil, err
	}
	return localweb.Start(ctx, localweb.Config{Tool: tool, Assets: sub, API: newBridge().handler(), Token: token})
}

// Run serves until ctx ends, printing the address once it listens.
func Run(ctx context.Context, out io.Writer) error {
	srv, err := Start(ctx)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(out, "proj pages on http://"+srv.Addr())
	return srv.Wait()
}

// Address says where a running server is: its host:port and token. ok is false
// when nothing that holds our token answers on the remembered port.
func Address() (addr, token string, ok bool) {
	dir := tools.StateDir(tool)
	port, err := os.ReadFile(filepath.Join(dir, "port"))
	if err != nil {
		return "", "", false
	}
	if _, err := strconv.Atoi(strings.TrimSpace(string(port))); err != nil {
		return "", "", false
	}
	tok, err := os.ReadFile(tokenPath())
	if err != nil {
		return "", "", false
	}
	addr, token = "127.0.0.1:"+strings.TrimSpace(string(port)), strings.TrimSpace(string(tok))
	req, err := http.NewRequest(http.MethodGet, "http://"+addr+"/api/health", http.NoBody)
	if err != nil {
		return "", "", false
	}
	req.Header.Set(localweb.TokenHeader, token)
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Do(req)
	if err != nil {
		return "", "", false
	}
	_ = resp.Body.Close()
	return addr, token, resp.StatusCode == http.StatusOK
}

// TermPath is the terminal page's path and query for a session, as the browser
// shows it once the token is exchanged for the cookie.
func TermPath(socket, session string) string {
	return "/term?" + url.Values{"socket": {socket}, "session": {session}}.Encode()
}

// TermURL is the address to open: TermPath on addr, with the token.
func TermURL(addr, token, socket, session string) string {
	return "http://" + addr + TermPath(socket, session) + "&t=" + token
}

func tokenPath() string { return filepath.Join(tools.StateDir(tool), "pages-token") }

func loadToken() (string, error) {
	b, err := os.ReadFile(tokenPath())
	if err == nil {
		return strings.TrimSpace(string(b)), nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	tok := make([]byte, 32)
	if _, err := rand.Read(tok); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tok)
	return token, tools.WriteFileAtomic(tokenPath(), []byte(token+"\n"), 0o600)
}
