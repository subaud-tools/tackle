package projcli

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/internal/version"
	tools "github.com/subaud-tools/tools-common"
)

func TestServeWritesItsEntryWithTheVersion(t *testing.T) {
	t.Setenv("PROJ_HOME", t.TempDir())
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- serveForeground(ctx, 0, io.Discard) }()
	t.Cleanup(func() { cancel(); <-done })
	path := filepath.Join(tools.StateDir("proj"), "live", "serve.json")
	var b []byte
	for range 100 {
		var err error
		if b, err = os.ReadFile(path); err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	var e struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(b, &e); err != nil {
		t.Fatalf("%s = %q: %v", path, b, err)
	}
	if e.Version != version.Number() {
		t.Fatalf("entry version = %q, want %q", e.Version, version.Number())
	}
}
