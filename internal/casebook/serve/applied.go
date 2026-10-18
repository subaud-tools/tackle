package serve

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/subaud-tools/tackle/internal/casebook/config"
	"github.com/subaud-tools/tackle/internal/casebook/observe"
)

// doneUntilObserved is each verified item still to keep out of review: one
// casebook verified done after GitHub was last read for its owner (read), so
// the index, built from that read, still shows it waiting. An owner never
// read goes by the index's build time.
func doneUntilObserved(verified map[string]time.Time, read func(owner string) (time.Time, bool), built time.Time) map[string]bool {
	out := map[string]bool{}
	for key, at := range verified {
		seen, ok := read(ownerOf(key))
		if !ok {
			seen = built
		}
		if at.After(seen) {
			out[key] = true
		}
	}
	return out
}

// ownerOf is a key's owner, lower case: "repo:o/r", "pr:o/r#4",
// "branch:o/r@b" → "o".
func ownerOf(key string) string {
	rest := key[strings.IndexByte(key, ':')+1:]
	if i := strings.IndexByte(rest, '/'); i >= 0 {
		rest = rest[:i]
	}
	return strings.ToLower(rest)
}

// ownerReads is when GitHub was last read for each owner, from the GitHub
// cache, read again only when the file changes.
type ownerReads struct {
	mu    sync.Mutex
	mod   time.Time
	reads map[string]time.Time
}

func (o *ownerReads) get() map[string]time.Time {
	o.mu.Lock()
	defer o.mu.Unlock()
	path := config.CachePath()
	fi, err := os.Stat(path)
	if err != nil {
		return nil
	}
	if o.reads != nil && fi.ModTime().Equal(o.mod) {
		return o.reads
	}
	gh, err := observe.LoadGitHub(path)
	if err != nil || gh == nil {
		return o.reads
	}
	m := map[string]time.Time{}
	for login, ow := range gh.Owners {
		if !ow.FetchedAt.IsZero() {
			m[strings.ToLower(login)] = ow.FetchedAt
		}
	}
	o.mod, o.reads = fi.ModTime(), m
	return m
}

// applied is Index.Applied: what casebook verified done that the index's
// GitHub read hasn't seen yet.
func (s *Server) applied(built time.Time) map[string]bool {
	reads := s.reads.get()
	oldest := built
	for _, at := range reads {
		if at.Before(oldest) {
			oldest = at
		}
	}
	verified, err := s.Apply.VerifiedSince(context.Background(), oldest)
	if err != nil {
		return nil
	}
	return doneUntilObserved(verified, func(owner string) (time.Time, bool) {
		at, ok := reads[owner]
		return at, ok
	}, built)
}
