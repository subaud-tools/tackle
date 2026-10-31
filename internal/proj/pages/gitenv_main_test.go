package pages

import (
	"os"
	"testing"

	"github.com/subaud-tools/tackle/internal/sift/gitenv"
)

// TestMain clears the variables that point git at another repository, which a
// git hook exports when the pre-push gate runs these tests.
func TestMain(m *testing.M) {
	for _, k := range gitenv.Selectors {
		_ = os.Unsetenv(k)
	}
	os.Exit(m.Run())
}
