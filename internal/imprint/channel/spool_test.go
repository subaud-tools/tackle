package channel

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// watcher runs Watch for session in dir and hands back each line it prints.
type watcher struct {
	lines chan string
	stop  func()
	done  chan error
}

func startWatch(t *testing.T, dir, session string) *watcher {
	t.Helper()
	r, w := io.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	wt := &watcher{lines: make(chan string, 100), done: make(chan error, 1)}
	go func() {
		err := Watch(ctx, dir, session, w, io.Discard)
		_ = w.Close()
		wt.done <- err
	}()
	go func() {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			wt.lines <- sc.Text()
		}
	}()
	wt.stop = cancel
	t.Cleanup(cancel)
	return wt
}

// next waits for the next printed line.
func (wt *watcher) next(t *testing.T) string {
	t.Helper()
	select {
	case l := <-wt.lines:
		return l
	case <-time.After(5 * time.Second):
		t.Fatal("watch printed no line")
		return ""
	}
}

// none fails if watch prints anything within d.
func (wt *watcher) none(t *testing.T, d time.Duration) {
	t.Helper()
	select {
	case l := <-wt.lines:
		t.Fatalf("watch printed an extra line: %q", l)
	case <-time.After(d):
	}
}

func content(t *testing.T, line string) string {
	t.Helper()
	var ev struct {
		Content string            `json:"content"`
		Meta    map[string]string `json:"meta"`
	}
	if err := json.Unmarshal([]byte(line), &ev); err != nil {
		t.Fatalf("not one JSON line: %q (%v)", line, err)
	}
	return ev.Content
}

// What the channel spooled before watch started is printed, from the start
// of the file, then what it spools after; nothing twice.
func TestWatchPrintsTheSpoolFromItsStartThenFollowsIt(t *testing.T) {
	dir := t.TempDir()
	sp, err := OpenSpool(dir, "s-1")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sp.Close() }()
	_ = sp.Notify("one", map[string]string{"op": "imprint_login"})
	_ = sp.Notify("two", nil)

	wt := startWatch(t, dir, "s-1")
	if got := content(t, wt.next(t)); got != "one" {
		t.Fatalf("first = %q", got)
	}
	if got := content(t, wt.next(t)); got != "two" {
		t.Fatalf("second = %q", got)
	}
	_ = sp.Notify("three", nil)
	if got := content(t, wt.next(t)); got != "three" {
		t.Fatalf("third = %q", got)
	}
	wt.none(t, 500*time.Millisecond)
}

// Watch may start before the channel has made its file.
func TestWatchWaitsForASpoolThatDoesNotExistYet(t *testing.T) {
	dir := t.TempDir()
	wt := startWatch(t, dir, "s-1")
	time.Sleep(300 * time.Millisecond)
	sp, err := OpenSpool(dir, "s-1")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sp.Close() }()
	_ = sp.Notify("late", nil)
	if got := content(t, wt.next(t)); got != "late" {
		t.Fatalf("got %q", got)
	}
}

// A line still being written is not printed until its newline arrives.
func TestWatchNeverPrintsHalfALine(t *testing.T) {
	dir := t.TempDir()
	path, err := SpoolFile(dir, "s-1")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	whole := `{"content":"half","meta":{}}`
	_, _ = f.WriteString(whole[:10])
	wt := startWatch(t, dir, "s-1")
	wt.none(t, 500*time.Millisecond)
	_, _ = f.WriteString(whole[10:] + "\n")
	if got := wt.next(t); got != whole {
		t.Fatalf("got %q, want %q", got, whole)
	}
}

// A channel that restarts in the same session removes its old file and makes
// a new one; watch prints the rest of the old file, then the new one.
func TestWatchFollowsTheNewSpoolAfterTheChannelRestarts(t *testing.T) {
	dir := t.TempDir()
	wt := startWatch(t, dir, "s-1")
	first, err := OpenSpool(dir, "s-1")
	if err != nil {
		t.Fatal(err)
	}
	_ = first.Notify("before", nil)
	if got := content(t, wt.next(t)); got != "before" {
		t.Fatalf("got %q", got)
	}
	_ = first.Notify("last words", nil)
	_ = first.Close()
	second, err := OpenSpool(dir, "s-1")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = second.Close() }()
	_ = second.Notify("after", nil)
	got := []string{content(t, wt.next(t)), content(t, wt.next(t))}
	if strings.Join(got, ",") != "last words,after" {
		t.Fatalf("got %v", got)
	}
	wt.none(t, 500*time.Millisecond)
}

func TestWatchReturnsWhenCancelled(t *testing.T) {
	wt := startWatch(t, t.TempDir(), "s-1")
	wt.stop()
	select {
	case err := <-wt.done:
		if err != nil {
			t.Fatalf("err = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("watch kept running after its context ended")
	}
}

// The channel removes its file when it exits, and on start prunes files
// that crashed sessions left behind days ago, but not a live session's.
func TestSpoolIsRemovedOnCloseAndOldSpoolsArePruned(t *testing.T) {
	dir := t.TempDir()
	stale, _ := SpoolFile(dir, "crashed")
	fresh, _ := SpoolFile(dir, "running")
	if err := os.MkdirAll(filepath.Dir(stale), 0o700); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{stale, fresh} {
		if err := os.WriteFile(p, []byte("{}\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	old := time.Now().Add(-SpoolMaxAge - time.Hour)
	if err := os.Chtimes(stale, old, old); err != nil {
		t.Fatal(err)
	}

	sp, err := OpenSpool(dir, "s-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Errorf("the old spool was not pruned: %v", err)
	}
	if _, err := os.Stat(fresh); err != nil {
		t.Errorf("a recent spool was pruned: %v", err)
	}
	mine, _ := SpoolFile(dir, "s-1")
	if _, err := os.Stat(mine); err != nil {
		t.Fatalf("no spool file: %v", err)
	}
	if err := sp.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(mine); !os.IsNotExist(err) {
		t.Errorf("the spool outlived its channel: %v", err)
	}
}

// A session id becomes a file name, so one that could leave the directory is
// refused.
func TestSpoolFileRefusesAPathInTheSessionID(t *testing.T) {
	for _, s := range []string{"", "../x", "a/b", "."} {
		if _, err := SpoolFile(t.TempDir(), s); err == nil {
			t.Errorf("SpoolFile(%q) took it", s)
		}
	}
}

// Two channels can run in one session (the plugin's and a bare
// ~/.claude.json server, until the bare one is removed): they share the
// spool, neither loses the other's lines, and the file goes only when the
// last one exits. A file no live channel holds is a crashed session's
// leftover: it is replaced, and pruned once old, but a live one never is.
func TestChannelsInOneSessionShareTheSpool(t *testing.T) {
	dir := t.TempDir()
	path, _ := SpoolFile(dir, "s-1")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"content":"left by a crash","meta":{}}`+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	a, err := OpenSpool(dir, "s-1")
	if err != nil {
		t.Fatal(err)
	}
	b, err := OpenSpool(dir, "s-1")
	if err != nil {
		t.Fatal(err)
	}
	wt := startWatch(t, dir, "s-1")
	_ = a.Notify("from a", nil)
	_ = b.Notify("from b", nil)
	got := []string{content(t, wt.next(t)), content(t, wt.next(t))}
	if strings.Join(got, ",") != "from a,from b" {
		t.Fatalf("got %v", got)
	}

	// b is live, so a later channel's prune leaves its file however old.
	old := time.Now().Add(-SpoolMaxAge - time.Hour)
	_ = os.Chtimes(path, old, old)
	other, err := OpenSpool(dir, "s-2")
	if err != nil {
		t.Fatal(err)
	}
	_ = other.Close()

	_ = a.Close()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("the spool went while b still runs: %v", err)
	}
	_ = b.Notify("b again", nil)
	if got := content(t, wt.next(t)); got != "b again" {
		t.Fatalf("got %q", got)
	}
	_ = b.Close()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("the spool outlived its last channel: %v", err)
	}
	wt.none(t, 300*time.Millisecond)
}
