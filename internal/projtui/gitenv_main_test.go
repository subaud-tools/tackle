package projtui

import (
	"os"
	"testing"

	"github.com/subaud-tools/tackle/internal/sift/gitenv"
)

// TestMain clears the variables that point git at another repository. A git
// hook exports them (the pre-push hook from a linked worktree sets GIT_DIR),
// and without this the tests' git calls would write into that repository.
// It also points XDG_CONFIG_HOME at a temp dir, so a test that saves a picker
// choice (^o) never writes the real proj config.
func TestMain(m *testing.M) {
	for _, k := range gitenv.Selectors {
		_ = os.Unsetenv(k)
	}
	dir, err := os.MkdirTemp("", "projtui-config-")
	if err != nil {
		panic(err)
	}
	_ = os.Setenv("XDG_CONFIG_HOME", dir)
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}
