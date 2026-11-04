package proj

import (
	"sync/atomic"
	"testing"
	"time"
)

// A start running in the background when proj is about to exit gets to finish,
// so a short proj command does not end before muster and hail are launched.
func TestWaitStartsLetsARunningStartFinish(t *testing.T) {
	var done atomic.Bool
	goStart(func() { time.Sleep(50 * time.Millisecond); done.Store(true) })
	WaitStarts(time.Second)
	if !done.Load() {
		t.Fatal("WaitStarts returned before the start finished")
	}
}

// A start that hangs never holds proj's exit past the limit.
func TestWaitStartsGivesUpAtTheLimit(t *testing.T) {
	release := make(chan struct{})
	t.Cleanup(func() { close(release) })
	goStart(func() { <-release })
	began := time.Now()
	WaitStarts(100 * time.Millisecond)
	if d := time.Since(began); d > time.Second {
		t.Fatalf("WaitStarts waited %v", d)
	}
}
