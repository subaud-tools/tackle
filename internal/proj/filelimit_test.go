package proj

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"testing"
)

// fileLimitChildEnv marks the re-executed test binary that plays proj.
const fileLimitChildEnv = "PROJ_FILELIMIT_CHILD"

// The raised limit must reach a child (tmux), not just proj. Go raises its own
// soft limit before main and hands children the original one back unless the
// program calls Setrlimit itself, so this starts a fresh process with a low
// soft limit, as a login shell starts proj, and reads what its child gets.
func TestRaiseFileLimitReachesChild(t *testing.T) {
	var lim syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &lim); err != nil {
		t.Fatal(err)
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("sh", "-c", `ulimit -Sn 256 && exec "$0" -test.run='^TestRaiseFileLimitChild$'`, exe)
	cmd.Env = append(os.Environ(), fileLimitChildEnv+"=1")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("child: %v; out=%s", err, out)
	}
	want := uint64(8192)
	if lim.Max < want {
		want = lim.Max
	}
	got := strings.TrimSpace(strings.SplitN(string(out), "\n", 2)[0])
	if got != strconv.FormatUint(want, 10) {
		t.Fatalf("grandchild ulimit -n = %q, want %d", got, want)
	}
}

// TestRaiseFileLimitChild is the re-executed process: it raises the limit and
// prints what its own child sees. It does nothing in a normal run.
func TestRaiseFileLimitChild(t *testing.T) {
	if os.Getenv(fileLimitChildEnv) != "1" {
		t.Skip("only runs re-executed by TestRaiseFileLimitReachesChild")
	}
	if err := RaiseFileLimit(8192); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command("sh", "-c", "ulimit -n").Output()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stdout.Write(out); err != nil {
		t.Fatal(err)
	}
}
