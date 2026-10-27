package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/internal/sift/apply"
	"github.com/subaud-tools/tackle/internal/sift/rec"
	"github.com/subaud-tools/tackle/internal/sift/row"
	"github.com/subaud-tools/tackle/internal/sift/serve"
	st "github.com/subaud-tools/tackle/internal/sift/sifttest"
	"github.com/subaud-tools/tackle/internal/sift/store"
	"github.com/subaud-tools/tools-common/localweb"
	"github.com/subaud-tools/tools-common/msgstate"
	"github.com/subaud-tools/tools-common/pageserve"
)

// ---- serve ------------------------------------------------------------------

// serveRig isolates sift's state. Starting serve in the background and
// opening the browser are pageserve's ServeCommand, tested there; these test
// sift's part: the note, no arguments, and the running serve's entry.
func serveRig(t *testing.T) {
	t.Helper()
	siftEnv(t)
}

func writeEntry(t *testing.T, e pageserve.Entry) {
	t.Helper()
	e.PID = os.Getpid()
	if err := serve.Live().Advertise(e); err != nil {
		t.Fatal(err)
	}
}

func TestServePageNotesThereIsNoRound(t *testing.T) {
	serveRig(t)
	fragment, note, err := servePage(context.Background(), nil)
	if err != nil || fragment != "" || note != "no round yet: run sift check" {
		t.Fatalf("%q %q %v", fragment, note, err)
	}
	seedRound(t)
	if _, note, err := servePage(context.Background(), nil); err != nil || note != "" {
		t.Fatalf("with a round: %q %v", note, err)
	}
	if code, _, errw := run(t, "", "serve", "x"); code != 2 || !strings.Contains(errw, "serve takes no arguments") {
		t.Fatalf("an argument: %d %q", code, errw)
	}
}

func TestServeNoOpenAndRunning(t *testing.T) {
	serveRig(t)
	seedRound(t)
	writeEntry(t, pageserve.Entry{URL: "http://127.0.0.1:2/?t=x", Base: "http://127.0.0.1:2"})
	code, out, _ := run(t, "", "serve", "--no-open")
	if code != 0 || !strings.Contains(out, "is running: http://127.0.0.1:2") || strings.Contains(out, "no round") {
		t.Fatalf("code %d out %q", code, out)
	}
}

func TestServeStopWhenNotRunning(t *testing.T) {
	serveRig(t)
	if code, out, _ := run(t, "", "serve", "--stop"); code != 0 || !strings.Contains(out, "not running") {
		t.Fatalf("%d %q", code, out)
	}
}

// ---- wait (cull's wait_test shape) ------------------------------------------

// waitRig runs serve in-process.
func waitRig(t *testing.T) *store.Store {
	t.Helper()
	siftEnv(t)
	s, err := store.Open(context.Background(), store.Path())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	ctx, cancel := context.WithCancel(context.Background())
	ready := make(chan struct{})
	done := make(chan error, 1)
	go func() { done <- serve.Run(ctx, s, serve.Options{Ready: func(string) { close(ready) }}) }()
	select {
	case <-ready:
	case err := <-done:
		t.Fatalf("serve: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("serve did not start")
	}
	t.Cleanup(func() {
		cancel()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Error("serve did not stop")
		}
	})
	return s
}

// decideAndPressSend records a round of one row, rejects it with a note and
// presses Send.
func decideAndPressSend(t *testing.T, s *store.Store) {
	t.Helper()
	ctx := context.Background()
	r := row.Row{ID: "r1", Check: "stale-status", Source: row.Source{File: "/w/a/CLAUDE.md", Start: 3, End: 3}, Passage: "- Never push."}
	id, err := s.RecordRound(ctx, store.Round{Kind: "on-demand"}, []row.Row{r})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Decide(ctx, id, "r1", row.Decision{Action: "reject", Note: "it is fine"}); err != nil {
		t.Fatal(err)
	}
	_, rows, err := s.Round(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	shown, _ := json.Marshal(map[string]store.Seen{"r1": {Fingerprint: rows[0].Fingerprint, Decision: rows[0].Decision.ID}})
	adv, err := serve.Running()
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("POST", adv.Base+"/api/send", strings.NewReader(fmt.Sprintf(`{"round":%d,"files":{},"rows":%s}`, id, shown)))
	req.Header.Set(localweb.TokenHeader, adv.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("send: %d %s", resp.StatusCode, b)
	}
}

type waitResult struct {
	code      int
	out, errw string
}

func startWait(args ...string) <-chan waitResult {
	ch := make(chan waitResult, 1)
	go func() {
		var o, e strings.Builder
		code := Main(append([]string{"wait"}, args...), strings.NewReader(""), &o, &e)
		ch <- waitResult{code, o.String(), e.String()}
	}()
	return ch
}

func TestWaitPrintsTheSendAndExits0(t *testing.T) {
	s := waitRig(t)
	res := startWait()
	time.Sleep(300 * time.Millisecond)
	decideAndPressSend(t, s)
	select {
	case r := <-res:
		if r.code != 0 || !strings.HasPrefix(r.out, "The user sent their decisions for sift round 1: 0 accepted, 0 edited, 1 rejected") ||
			!strings.Contains(r.out, "- /w/a/CLAUDE.md:3 · stale-status: it is fine\n") {
			t.Errorf("code %d\nout %q\nerr %q", r.code, r.out, r.errw)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("wait never returned")
	}
}

// sift wait reports the send it printed landed: it is with the agent.
func TestWaitReportsTheSendLanded(t *testing.T) {
	s := waitRig(t)
	res := startWait()
	time.Sleep(300 * time.Millisecond)
	decideAndPressSend(t, s)
	select {
	case r := <-res:
		if r.code != 0 {
			t.Fatalf("code %d err %q", r.code, r.errw)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("wait never returned")
	}
	sends, err := s.RecentSends(context.Background(), 1, 1)
	if err != nil || len(sends) != 1 || sends[0].State != msgstate.WithAgent {
		t.Fatalf("%+v %v", sends, err)
	}
}

func TestWaitTimesOutWithExit3(t *testing.T) {
	waitRig(t)
	res := <-startWait("--timeout", "1s")
	if res.code != 3 || res.out != "" || !strings.Contains(res.errw, "nothing sent") {
		t.Errorf("code %d out %q err %q", res.code, res.out, res.errw)
	}
}

func TestWaitBadTimeoutIsUsage(t *testing.T) {
	siftEnv(t)
	code, _, errw := run(t, "", "wait", "--timeout", "0s")
	if code != 2 || !strings.Contains(errw, "--timeout") {
		t.Errorf("code %d err %q", code, errw)
	}
}

// ---- apply and reconcile ------------------------------------------------------

// applyRig is a published repo with a round of one file, its
// recommendation accepted and sent.
func applyRig(t *testing.T) string {
	t.Helper()
	siftEnv(t)
	repo := st.Repo(t, filepath.Join(t.TempDir(), "app"), map[string]string{"CLAUDE.md": "# App\n\n- Never push to main.\n"})
	st.Publish(t, repo)
	s, err := store.Open(context.Background(), store.Path())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()
	f := rec.NewFile(row.Source{File: filepath.Join(repo, "CLAUDE.md"), Repo: repo, Ref: "origin/main", Path: "CLAUDE.md"}, "repo", 6000, "# App\n\n- Never push to main.\n")
	f.Commit, f.Rows = st.Git(t, repo, "rev-parse", "origin/main"), []string{"neg"}
	r := row.Row{ID: "neg", Check: "stale-status", Passage: "- Never push to main.", Source: row.Source{File: f.Source.File, Start: 3, End: 3}}
	id, err := s.RecordAudit(ctx, store.Round{Kind: "on-demand"}, []row.Row{r}, []rec.File{f})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Propose(ctx, id, []rec.Rec{{File: f.Key, Base: f.Base, Content: "# App\n\n- Push to a branch.\n", Summary: "Guidance.",
		Findings: []rec.Account{{Row: "neg", Did: "fixed", How: "guidance"}}}}); err != nil {
		t.Fatal(err)
	}
	items, _ := s.Files(ctx, id)
	if _, err := s.DecideFile(ctx, id, f.Key, rec.Decision{Action: "accept"}, map[string]store.Seen{f.Key: {Fingerprint: items[0].Fingerprint}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Send(context.Background(), id, "", nil); err != nil {
		t.Fatal(err)
	}
	return repo
}

func TestApplyWithoutGhLeavesTheBranch(t *testing.T) {
	repo := applyRig(t)
	code, out, errw := run(t, "", "apply")
	if code != 0 || !strings.Contains(out, "branch sift/round-1 (1 applied, 0 skipped)") || !strings.Contains(out, "no gh") || !strings.Contains(out, "written  accept") {
		t.Fatalf("code %d out %q err %q", code, out, errw)
	}
	if got := st.Git(t, repo, "show", "sift/round-1:CLAUDE.md"); got != "# App\n\n- Push to a branch." {
		t.Errorf("branch content %q", got)
	}
	code, out, _ = run(t, "", "reconcile")
	if code != 0 || !strings.Contains(out, "1 file(s), 0 problem(s)") {
		t.Fatalf("reconcile %d %q", code, out)
	}
	// Again: the branch exists, so the repo is held (exit 1).
	code, out, _ = run(t, "", "apply")
	if code != 1 || !strings.Contains(out, "held") {
		t.Fatalf("again: %d %q", code, out)
	}
}

func TestApplyWithGhOpensAPullRequest(t *testing.T) {
	repo := applyRig(t)
	bare := st.Git(t, repo, "remote", "get-url", "origin")
	st.Git(t, repo, "remote", "set-url", "origin", "https://github.com/owner/app.git")
	st.Git(t, repo, "config", "url."+bare+".insteadOf", "https://github.com/owner/app.git")
	lookPath = func(string) (string, error) { return "/usr/bin/gh", nil }
	orig := ghRunner
	ghRunner = func(_ context.Context, _ string, args ...string) ([]byte, error) {
		return []byte("https://github.com/owner/app/pull/3\n"), nil
	}
	t.Cleanup(func() { ghRunner = orig })
	code, out, errw := run(t, "", "apply", "--json")
	var res apply.Result
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("code %d out %q err %q", code, out, errw)
	}
	if res.Repos[0].State != "pr" || res.Repos[0].PR != "https://github.com/owner/app/pull/3" {
		t.Fatalf("%+v", res.Repos[0])
	}
}

func TestReconcileWithNothingApplied(t *testing.T) {
	siftEnv(t)
	seedRound(t)
	code, out, _ := run(t, "", "reconcile")
	if code != 0 || !strings.Contains(out, "no branch to reconcile") {
		t.Fatalf("%d %q", code, out)
	}
}

// sift clean prints its plan first; --dry-run stops there.
func TestCleanPrintsThePlan(t *testing.T) {
	siftEnv(t)
	code, out, errw := run(t, "", "clean", "--dry-run")
	if code != 0 || !strings.Contains(out, "nothing sift created is left to clean") {
		t.Fatalf("code %d: %s %s", code, out, errw)
	}
}
