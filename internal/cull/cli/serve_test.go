package cli

// Starting serve in the background and opening the browser are pageserve's
// ServeCommand, tested there; this tests cull's part: the path resolves to
// its project's review, and serve's flags reach the running serve's entry.

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/subaud-tools/tackle/internal/cull/serve"
	"github.com/subaud-tools/tackle/internal/cull/store"
	"github.com/subaud-tools/tools-common/pageserve"
)

func serveRig(t *testing.T) (root string) {
	t.Helper()
	t.Setenv("CULL_HOME", t.TempDir())
	return t.TempDir()
}

func writeEntry(t *testing.T, e pageserve.Entry) {
	t.Helper()
	e.PID = os.Getpid()
	if err := serve.Live().Advertise(e); err != nil {
		t.Fatal(err)
	}
}

func TestServePageOpensTheProjectsReview(t *testing.T) {
	root := serveRig(t)
	fragment, note, err := servePage(context.Background(), []string{root})
	if err != nil {
		t.Fatal(err)
	}
	if want := "#/p/" + projectID(t, root); fragment != want {
		t.Errorf("fragment %q, want %q", fragment, want)
	}
	if !strings.Contains(note, "no review for") || !strings.Contains(note, "run cull check") {
		t.Errorf("want no-review note: %q", note)
	}
}

// projectID is the store's id for the project rooted at root.
func projectID(t *testing.T, root string) string {
	t.Helper()
	ctx := context.Background()
	st, err := store.Open(ctx, store.Path())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	p, err := st.Project(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	return strconv.FormatInt(p.ID, 10)
}

func TestServeSubdirOpensSameProject(t *testing.T) {
	root := serveRig(t)
	if out, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Skipf("git init: %v %s", err, out)
	}
	sub := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	real, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	want := "#/p/" + projectID(t, real)
	for _, p := range []string{root, sub} {
		if fragment, _, err := servePage(context.Background(), []string{p}); err != nil || fragment != want {
			t.Errorf("%s: %q %v, want %q", p, fragment, err, want)
		}
	}
}

func TestServeReusesRunning(t *testing.T) {
	root := serveRig(t)
	writeEntry(t, pageserve.Entry{URL: "http://127.0.0.1:2/?t=x", Base: "http://127.0.0.1:2"})
	code, out, errw := run(t, "", "serve", root, "--no-open")
	if code != 0 {
		t.Fatalf("code %d: %s", code, errw)
	}
	if !strings.Contains(out, "cull serve is running: http://127.0.0.1:2") || !strings.Contains(out, "no review for") {
		t.Errorf("out=%q", out)
	}
}

func TestServeStop(t *testing.T) {
	serveRig(t)
	code, out, _ := run(t, "", "serve", "--stop")
	if code != 0 || !strings.Contains(out, "cull serve is not running") {
		t.Errorf("code %d out %q", code, out)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = os.Remove(serve.Live().ServePath())
	}))
	defer ts.Close()
	writeEntry(t, pageserve.Entry{URL: ts.URL, Base: ts.URL})
	code, out, errw := run(t, "", "serve", "--stop")
	if code != 0 || !strings.Contains(out, "cull serve stopped") {
		t.Errorf("code %d out %q err %q", code, out, errw)
	}
}

func TestServeMissingPathExits2(t *testing.T) {
	serveRig(t)
	code, _, errw := run(t, "", "serve", filepath.Join(t.TempDir(), "nope"))
	if code != 2 || !strings.Contains(errw, "nope") {
		t.Errorf("code %d err %q", code, errw)
	}
}
