package proj

import (
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"testing"
)

// The raised limit must reach a child (tmux), not just proj: Go hands
// children the original limit unless the program set it itself.
func TestRaiseFileLimitReachesChild(t *testing.T) {
	var orig syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &orig); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &orig) })
	low := syscall.Rlimit{Cur: 256, Max: orig.Max}
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &low); err != nil {
		t.Fatal(err)
	}

	if err := RaiseFileLimit(8192); err != nil {
		t.Fatalf("RaiseFileLimit: %v", err)
	}
	out, err := exec.Command("sh", "-c", "ulimit -n").Output()
	if err != nil {
		t.Fatal(err)
	}
	want := uint64(8192)
	if orig.Max < want {
		want = orig.Max
	}
	if got := strings.TrimSpace(string(out)); got != strconv.FormatUint(want, 10) {
		t.Fatalf("child ulimit -n = %s, want %d", got, want)
	}
}
