package proj

import (
	"errors"
	"os/exec"
	"testing"
)

func TestQueryEmptyOnBadSocket(t *testing.T) {
	// A socket that cannot exist yields "" (degradation), never a panic.
	if got := Query("/nonexistent/proj-nope", "=nope", "#{session_name}"); got != "" {
		t.Fatalf("want empty, got %q", got)
	}
}

func TestRunErrorNamesCommandAndTmuxMessage(t *testing.T) {
	orig := execCommand
	t.Cleanup(func() { execCommand = orig })
	execCommand = func(string, ...string) *exec.Cmd {
		return orig("sh", "-c", "echo 'no current client' >&2; exit 1")
	}
	_, err := Run("", "switch-client", "-t", "=a/b")
	if err == nil {
		t.Fatal("want an error")
	}
	if want := "tmux switch-client -t =a/b: no current client"; err.Error() != want {
		t.Fatalf("got %q, want %q", err, want)
	}
	var ee *exec.ExitError
	if !errors.As(err, &ee) {
		t.Fatalf("errors.As does not find the exit error in %v", err)
	}
}

func TestRunErrorFallsBackToExitStatus(t *testing.T) {
	orig := execCommand
	t.Cleanup(func() { execCommand = orig })
	execCommand = func(string, ...string) *exec.Cmd { return orig("sh", "-c", "exit 1") }
	_, err := Run("", "has-session", "-t", "=x")
	if err == nil || err.Error() != "tmux has-session -t =x: exit status 1" {
		t.Fatalf("got %v", err)
	}
}
