package proj

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/subaud-tools/tackle/internal/sift/gitenv"
)

// TestMain clears the variables that point git at another repository. A git
// hook exports them (the pre-push hook from a linked worktree sets GIT_DIR),
// and without this the tests' git calls would write into that repository.
func TestMain(m *testing.M) {
	for _, k := range gitenv.Selectors {
		_ = os.Unsetenv(k)
	}
	// Fixtures isolate HOME to prevent login shells/tmux loading the operator's
	// config. Keep builds of the hook probe on the existing Go caches instead
	// of downloading read-only modules into testing.TempDir's cleanup tree.
	keys := []string{"GOPATH", "GOCACHE", "GOMODCACHE"}
	out, err := exec.Command("go", append([]string{"env"}, keys...)...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "resolve test Go caches:", err)
		os.Exit(1)
	}
	for i, value := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if i < len(keys) {
			_ = os.Setenv(keys[i], value)
		}
	}
	os.Exit(m.Run())
}
