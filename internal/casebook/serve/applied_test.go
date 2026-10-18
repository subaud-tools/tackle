package serve

import (
	"testing"
	"time"
)

// TestDoneUntilObserved (Court, 2026-10-17): an item casebook verified done
// stays out of review until GitHub has been read for its owner after it;
// the index's build time doesn't count (a restart rebuilds from the cache).
func TestDoneUntilObserved(t *testing.T) {
	t0 := time.Date(2026, 10, 17, 23, 29, 40, 0, time.UTC)
	verified := map[string]time.Time{
		"repo:schuettc/a":     t0.Add(2 * time.Minute), // after schuettc was read
		"repo:other/b":        t0.Add(2 * time.Minute), // other was read since
		"branch:schuettc/c@x": t0.Add(-time.Hour),      // before the read: it shows the truth
		"pr:Schuettc/d#4":     t0.Add(time.Minute),     // owner in any case
		"repo:unknown/e":      t0.Add(time.Minute),     // owner never read: the build time
	}
	read := map[string]time.Time{"schuettc": t0, "other": t0.Add(5 * time.Minute)}
	built := t0.Add(10 * time.Minute)
	got := doneUntilObserved(verified, func(owner string) (time.Time, bool) {
		at, ok := read[owner]
		return at, ok
	}, built)
	want := map[string]bool{"repo:schuettc/a": true, "pr:Schuettc/d#4": true}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for k := range want {
		if !got[k] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
