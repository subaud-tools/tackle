package pages

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	tools "github.com/subaud-tools/tools-common"
)

const fakeToken = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

// remember writes the state `proj pages` leaves: srv's port and fakeToken.
func remember(t *testing.T, srv *httptest.Server) {
	t.Helper()
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if err := tools.WriteFileAtomic(filepath.Join(tools.StateDir(tool), "port"), []byte(u.Port()+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := tools.WriteFileAtomic(tokenPath(), []byte(fakeToken+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

// recordingServer answers with handler and keeps each request's full text.
func recordingServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, func() string) {
	t.Helper()
	var mu sync.Mutex
	var seen strings.Builder
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true)
		mu.Lock()
		seen.Write(dump)
		mu.Unlock()
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	return srv, func() string { mu.Lock(); defer mu.Unlock(); return seen.String() }
}

func TestAddressRejectsAnotherServerAndKeepsTheTokenFromIt(t *testing.T) {
	srv, seen := recordingServer(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	remember(t, srv)
	if _, _, ok := Address(); ok {
		t.Fatal("a server that only answers 200 was accepted")
	}
	if got := seen(); got == "" || strings.Contains(got, fakeToken) {
		t.Fatalf("the probe was missing or held the token:\n%s", got)
	}
}

func TestAddressDoesNotFollowRedirects(t *testing.T) {
	target, seen := recordingServer(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	srv, _ := recordingServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	})
	remember(t, srv)
	if _, _, ok := Address(); ok {
		t.Fatal("a redirect was accepted")
	}
	if got := seen(); got != "" {
		t.Fatalf("the redirect was followed:\n%s", got)
	}
}

func TestAddressFindsTheRunningServer(t *testing.T) {
	srv := startServer(t)
	addr, token, ok := Address()
	if !ok || addr != srv.Addr() || token != srv.Token {
		t.Fatalf("Address = %q, %q, %v; want %q and its token", addr, token, ok, srv.Addr())
	}
}

func TestEachStartHasANewToken(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	first, err := Start(ctx)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Start(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if first.Token == second.Token {
		t.Fatal("two starts share a token")
	}
	b, err := os.ReadFile(tokenPath())
	if err != nil || strings.TrimSpace(string(b)) != second.Token {
		t.Fatalf("token file = %q, %v; want the latest start's", b, err)
	}
	fi, err := os.Stat(tokenPath())
	if err != nil || fi.Mode().Perm() != 0o600 {
		t.Fatalf("token file mode = %v, %v", fi, err)
	}
}

func TestTermMatches(t *testing.T) {
	const s, session = "/tmp/a/s", "p/w"
	page := func(socket, name string) string {
		return "http://127.0.0.1:9/term?" + url.Values{"socket": {socket}, "session": {name}}.Encode()
	}
	for _, c := range []struct {
		name, url string
		want      bool
	}{
		{"this session's page", page(s, session), true},
		{"the same name on another socket", page("/tmp/b/s", session), false},
		{"a socket that only starts the same", page(s+"2", session), false},
		{"another session", page(s, "p/x"), false},
		{"another site quoting the page", "https://evil.example/" + page(s, session), false},
		{"another host", strings.Replace(page(s, session), "127.0.0.1", "evil.example", 1), false},
	} {
		if got := TermMatches(c.url, s, session); got != c.want {
			t.Errorf("%s: %q = %v, want %v", c.name, c.url, got, c.want)
		}
	}
}
