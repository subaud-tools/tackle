package channel

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/subaud-tools/tools-common/channelmcp"
)

// imprint's messages come from jobs inside `imprint channel`, so under Claude
// Code a separate `imprint watch` (the plugin's mod) cannot wait for them the
// way cull's and sift's watch wait on their serve. Instead the channel writes
// each message as one JSON line to a spool file for its session, and watch
// follows that file from its start.

// SpoolMaxAge is how old a spool file left by a session that crashed may get
// before a starting channel prunes it. A live channel keeps its own file
// younger than this.
const SpoolMaxAge = 3 * 24 * time.Hour

var (
	sessionName = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9._-]*$`)
	// watchPoll is how often watch looks for new lines; keepFresh how often a
	// channel touches its spool so a prune never takes a live session's.
	watchPoll = 100 * time.Millisecond
	keepFresh = time.Hour
)

// SpoolFile is the session's spool under imprint's state directory dir.
func SpoolFile(dir, session string) (string, error) {
	if !sessionName.MatchString(session) {
		return "", fmt.Errorf("session id %q can't name a spool file", session)
	}
	return filepath.Join(dir, "spool", session+".jsonl"), nil
}

// Spool is a channel's spool file: each Notify appends one whole JSON line,
// {"content": ..., "meta": {...}}, in one write to a file opened for append.
type Spool struct {
	path  string
	f     *os.File
	lines *channelmcp.Lines
	stop  chan struct{}
	once  sync.Once
}

var _ channelmcp.Pusher = (*Spool)(nil)

// OpenSpool prunes spool files older than SpoolMaxAge and starts the
// session's spool afresh: a file a crashed channel of this session left is
// replaced, not appended to.
func OpenSpool(dir, session string) (*Spool, error) {
	path, err := SpoolFile(dir, session)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	prune(filepath.Dir(path), time.Now().Add(-SpoolMaxAge))
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, err
	}
	s := &Spool{path: path, f: f, lines: channelmcp.NewLines(f), stop: make(chan struct{})}
	go s.keepFresh()
	return s, nil
}

// Path is the spool file's path.
func (s *Spool) Path() string { return s.path }

// Notify appends one message.
func (s *Spool) Notify(content string, meta map[string]string) error {
	return s.lines.Notify(content, meta)
}

// Close ends the spool and removes its file, unless a newer channel of the
// same session has already replaced it.
func (s *Spool) Close() error {
	var err error
	s.once.Do(func() {
		close(s.stop)
		mine, _ := s.f.Stat()
		err = s.f.Close()
		if now, serr := os.Stat(s.path); serr == nil && mine != nil && os.SameFile(mine, now) {
			if rerr := os.Remove(s.path); err == nil {
				err = rerr
			}
		}
	})
	return err
}

func (s *Spool) keepFresh() {
	t := time.NewTicker(keepFresh)
	defer t.Stop()
	for {
		select {
		case <-s.stop:
			return
		case <-t.C:
			now := time.Now()
			_ = os.Chtimes(s.path, now, now)
		}
	}
}

// prune removes spool files last changed before cutoff.
func prune(dir string, cutoff time.Time) {
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		if info, err := e.Info(); err == nil && info.ModTime().Before(cutoff) {
			_ = os.Remove(filepath.Join(dir, e.Name()))
		}
	}
}

// Watch is `imprint watch`: it prints each line of the session's spool on
// out, from the file's start, then follows it until ctx ends. The file need
// not exist yet. Only whole lines are printed, each once; when the channel
// replaces the file (it restarted), watch finishes the old file and starts on
// the new one.
func Watch(ctx context.Context, dir, session string, out, log io.Writer) error {
	path, err := SpoolFile(dir, session)
	if err != nil {
		return err
	}
	var (
		f    *os.File
		rest []byte // the start of a line whose newline hasn't been written yet
	)
	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()
	tick := time.NewTicker(watchPoll)
	defer tick.Stop()
	for {
		if f == nil {
			if f, err = os.Open(path); err != nil {
				f = nil
				if !errors.Is(err, os.ErrNotExist) {
					_, _ = fmt.Fprintf(log, "imprint watch: %v\n", err)
				}
			}
		}
		if f != nil {
			if rest, err = drain(f, rest, out); err != nil {
				return err
			}
			if replaced(f, path) {
				// Lines written before the swap are in the old file still.
				if rest, err = drain(f, rest, out); err != nil {
					return err
				}
				if len(rest) > 0 {
					_, _ = fmt.Fprintf(log, "imprint watch: dropped an unfinished line from the old spool: %q\n", rest)
				}
				_ = f.Close()
				f, rest = nil, nil
				continue
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
		}
	}
}

// drain reads f to its end and prints each whole line, returning what is left
// of an unfinished one.
func drain(f *os.File, rest []byte, out io.Writer) ([]byte, error) {
	buf := make([]byte, 32<<10)
	for {
		n, err := f.Read(buf)
		rest = append(rest, buf[:n]...)
		for {
			i := bytes.IndexByte(rest, '\n')
			if i < 0 {
				break
			}
			if line := rest[:i+1]; len(bytes.TrimSpace(line)) > 0 {
				if _, werr := out.Write(line); werr != nil {
					return rest, werr
				}
			}
			rest = rest[i+1:]
		}
		if err == io.EOF || n == 0 {
			return rest, nil
		}
		if err != nil {
			return rest, err
		}
	}
}

// replaced says whether path no longer names the open file f.
func replaced(f *os.File, path string) bool {
	open, err := f.Stat()
	if err != nil {
		return true
	}
	now, err := os.Stat(path)
	return err != nil || !os.SameFile(open, now)
}
