// Package pages is `proj pages`: one long-running local server that shows any
// session's terminal in a browser page, xterm.js over a websocket to a tmux
// client in a pty.
package pages

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/hex"
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

// Start serves the page until ctx ends, on a new token each launch, which it
// leaves in StateDir for `proj window` once the server listens.
func Start(ctx context.Context) (*localweb.Server, error) {
	tok := make([]byte, 32)
	if _, err := rand.Read(tok); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(tok)
	sub, err := fs.Sub(assets, "assets")
	if err != nil {
		return nil, err
	}
	ctx, stop := context.WithCancel(ctx)
	srv, err := localweb.Start(ctx, localweb.Config{Tool: tool, Assets: healthFS{FS: sub, token: token}, API: newBridge().handler(), Token: token})
	if err != nil {
		stop()
		return nil, err
	}
	context.AfterFunc(ctx, stop)
	if err := tools.WriteFileAtomic(tokenPath(), []byte(token+"\n"), 0o600); err != nil {
		stop()
		return nil, err
	}
	return srv, nil
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
// when nothing on the remembered port proves it holds our token. The probe
// never sends the token: whatever listens there may be another program.
func Address() (addr, token string, ok bool) {
	port, err := os.ReadFile(filepath.Join(tools.StateDir(tool), "port"))
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
	n := make([]byte, 32)
	if _, err := rand.Read(n); err != nil {
		return "", "", false
	}
	nonce := hex.EncodeToString(n)
	client := &http.Client{
		Timeout:       2 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	resp, err := client.Get("http://" + addr + "/" + healthDir + "/" + nonce)
	if err != nil {
		return "", "", false
	}
	defer func() { _ = resp.Body.Close() }()
	proof, err := io.ReadAll(io.LimitReader(resp.Body, 256))
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", "", false
	}
	return addr, token, hmac.Equal(proof, []byte(healthProof(token, nonce)))
}

// healthDir is where the server answers health challenges. It is an asset path
// because localweb asks for the token on everything under /api/.
const healthDir = "health"

// healthProof shows the token is held without revealing it.
func healthProof(token, nonce string) string {
	m := hmac.New(sha256.New, []byte(token))
	m.Write([]byte(nonce))
	return hex.EncodeToString(m.Sum(nil))
}

// healthFS is the page's assets, plus health/<nonce>, a file whose content is
// the proof for that nonce. A nonce is 64 hex characters.
type healthFS struct {
	fs.FS
	token string
}

func (h healthFS) Open(name string) (fs.File, error) {
	nonce, ok := strings.CutPrefix(name, healthDir+"/")
	if !ok {
		return h.FS.Open(name)
	}
	if b, err := hex.DecodeString(nonce); err != nil || len(b) != 32 {
		return nil, fs.ErrNotExist
	}
	proof := healthProof(h.token, nonce)
	return proofFile{Reader: strings.NewReader(proof), size: int64(len(proof))}, nil
}

// proofFile is an in-memory file, seekable as http.FileServer needs.
type proofFile struct {
	*strings.Reader
	size int64
}

func (f proofFile) Close() error               { return nil }
func (f proofFile) Stat() (fs.FileInfo, error) { return f, nil }
func (f proofFile) Name() string               { return "proof" }
func (f proofFile) Size() int64                { return f.size }
func (f proofFile) Mode() fs.FileMode          { return 0o444 }
func (f proofFile) ModTime() time.Time         { return time.Time{} }
func (f proofFile) IsDir() bool                { return false }
func (f proofFile) Sys() any                   { return nil }

// TermURL is the address to open for a session's terminal page, with the token,
// to take the session from any other client; the page drops both on load.
func TermURL(addr, token, socket, session string) string {
	return "http://" + addr + "/term?" + url.Values{"socket": {socket}, "session": {session}, "t": {token}, "take": {"1"}}.Encode()
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

func tokenPath() string { return filepath.Join(tools.StateDir(tool), "pages-token") }
