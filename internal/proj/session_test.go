package proj

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestEnsureSessionFromHomeSetsLabel(t *testing.T) {
	requireTmux(t)
	sock := "proj-phase1-sess"
	defer func() { _, _ = Run(sock, "kill-server") }()
	dir := t.TempDir()
	// tmux canonicalizes pane paths; resolve symlinks (macOS /var → /private/var)
	// so the pane_current_path contract holds cross-platform.
	if r, err := filepath.EvalSymlinks(dir); err == nil {
		dir = r
	}
	if _, err := EnsureSession(sock, "proj-phase1-sess/w", dir, "none", "", nil); err != nil {
		t.Fatalf("EnsureSession: %v", err)
	}
	// pane dir is the requested dir...
	if got := Query(sock, "proj-phase1-sess/w", "#{pane_current_path}"); got != dir {
		t.Fatalf("pane dir = %q want %q", got, dir)
	}
	// ...but the SERVER's start dir is $HOME (cwd-poison defense): the session's
	// own start path is $HOME because new-session ran from there.
	if got := Query(sock, "proj-phase1-sess/w", "#{session_path}"); got != os.Getenv("HOME") {
		t.Fatalf("session_path = %q want $HOME", got)
	}
	if got := Query(sock, "proj-phase1-sess/w", "#{@claude_task}"); got != "w" {
		t.Fatalf("@claude_task = %q want w", got)
	}
	// idempotent — reused, so created=false
	if created, err := EnsureSession(sock, "proj-phase1-sess/w", dir, "none", "", nil); err != nil {
		t.Fatalf("second EnsureSession: %v", err)
	} else if created {
		t.Fatal("second EnsureSession reported created=true; want false (reuse)")
	}
}

// Creating a session starts muster and hail; reusing one starts nothing.
func TestCreatingASessionStartsMusterAndHail(t *testing.T) {
	withFakeRun(t)
	var started []string
	orig := startServers
	t.Cleanup(func() { startServers = orig })
	startServers = func(tools ...string) { started = append(started, tools...) }
	if created, err := EnsureSession("proj-x", "x/w", "/tmp/dir", "none", "", nil); err != nil || !created {
		t.Fatalf("created=%v err=%v", created, err)
	}
	if want := []string{"muster", "hail"}; !slices.Equal(started, want) {
		t.Fatalf("started %v, want %v", started, want)
	}
	started = nil
	runner = func(string, ...string) (string, error) { return "", nil } // has-session succeeds
	if created, err := RestoreSession("proj-x", "x/w", "/tmp/dir", "none", ""); err != nil || created {
		t.Fatalf("reuse: created=%v err=%v", created, err)
	}
	if len(started) != 0 {
		t.Fatalf("reusing a session started %v", started)
	}
}
