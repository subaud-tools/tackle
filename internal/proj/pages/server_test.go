package pages

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/subaud-tools/tackle/internal/proj/servers"
	"github.com/subaud-tools/tools-common/localweb"
	"github.com/subaud-tools/tools-common/splitpage"
)

// A restart keeps the token, so a window's terminal page keeps its address.
func TestTheTokenSurvivesARestart(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	ready := make(chan struct{})
	go func() { done <- Run(ctx, Options{Version: "test", Ready: func(string) { close(ready) }}) }()
	<-ready
	first, err := Live().Running()
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if second := runServer(t); second.Token != first.Token {
		t.Fatalf("the token changed across a restart: %q, then %q", first.Token, second.Token)
	}
}

func TestTheAPINeedsTheToken(t *testing.T) {
	srv := startServer(t)
	for _, path := range []string{"/api/term/ws?socket=/tmp/a/s&session=p/w", "/api/split/tabs?session=p/w"} {
		if code, _ := get(t, "http://"+srv.Addr()+path, nil); code != http.StatusUnauthorized {
			t.Errorf("%s without the token = %d, want 401", path, code)
		}
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

func get(t *testing.T, u string, header http.Header) (int, string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header = header
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, string(b)
}

func TestTermServesTheSplitPageThatTakes(t *testing.T) {
	srv := startServer(t)
	q := url.Values{"socket": {"/tmp/a/s"}, "session": {"p/w"}}
	code, body := get(t, "http://"+srv.Addr()+"/term?"+q.Encode(), nil)
	if code != http.StatusOK || !strings.Contains(body, `data-take="true"`) || !strings.Contains(body, `data-socket="/api/term/ws"`) {
		t.Fatalf("/term = %d:\n%s", code, body)
	}
	if code, _ := get(t, "http://"+srv.Addr()+"/split.js", nil); code != http.StatusOK {
		t.Fatalf("/split.js = %d", code)
	}
}

func TestATabOpenedForOneSessionIsNotListedForAnother(t *testing.T) {
	srv := startServer(t)
	base := "http://" + srv.Addr()
	if _, err := splitpage.Open(base, srv.Token, "p/w", "http://127.0.0.1:4321/plan", "plan"); err != nil {
		t.Fatal(err)
	}
	auth := http.Header{localweb.TokenHeader: {srv.Token}}
	var other, own splitpage.TabList
	for _, c := range []struct {
		session string
		into    *splitpage.TabList
	}{{"p/x", &other}, {"p/w", &own}} {
		code, body := get(t, base+"/api/split/tabs?"+url.Values{"session": {c.session}}.Encode(), auth)
		if code != http.StatusOK {
			t.Fatalf("tabs for %s = %d: %s", c.session, code, body)
		}
		if err := json.Unmarshal([]byte(body), c.into); err != nil {
			t.Fatal(err)
		}
	}
	if len(other.Tabs) != 0 || len(own.Tabs) != 1 || own.Tabs[0].Address != "http://127.0.0.1:4321/plan" {
		t.Fatalf("other session's tabs %+v, own %+v", other, own)
	}
}

// GET /api/servers needs the token and lists the groups: here proj serve
// itself, up, and a tool whose live folder holds no entry, down.
func TestTheServersAPINeedsTheTokenAndListsTheGroups(t *testing.T) {
	srv := startServer(t)
	if err := os.MkdirAll(filepath.Join(os.Getenv("XDG_STATE_HOME"), "zz", "live"), 0o700); err != nil {
		t.Fatal(err)
	}
	u := "http://" + srv.Addr() + "/api/servers"
	if code, _ := get(t, u, nil); code != http.StatusUnauthorized {
		t.Fatalf("without the token = %d, want 401", code)
	}
	code, body := get(t, u, http.Header{localweb.TokenHeader: {srv.Token}})
	if code != http.StatusOK {
		t.Fatalf("= %d: %s", code, body)
	}
	var l servers.List
	if err := json.Unmarshal([]byte(body), &l); err != nil {
		t.Fatal(err)
	}
	if len(l.Up) != 1 || l.Up[0].Tool != "proj" || len(l.Down) != 1 || l.Down[0].Tool != "zz" || l.Left == nil {
		t.Fatalf("groups %s", body)
	}
}

// A left-over stop on a pid Leftovers doesn't list is refused, and the
// process lives.
func TestALeftOverStopOnAPidNotListedIsRefused(t *testing.T) {
	srv := startServer(t)
	cmd := exec.Command("sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
	pid := cmd.Process.Pid
	code, body := post(t, srv, "/api/leftovers/"+strconv.Itoa(pid)+"/stop")
	if code != http.StatusConflict {
		t.Fatalf("stop pid %d = %d: %s, want 409", pid, code, body)
	}
	if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
		t.Fatalf("pid %d was stopped: %v", pid, err)
	}
}

func post(t *testing.T, srv served, path string) (int, string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, "http://"+srv.Addr()+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(localweb.TokenHeader, srv.Token)
	req.Header.Set("Origin", "http://"+srv.Addr())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}
