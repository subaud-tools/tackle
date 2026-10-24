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
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/subaud-tools/tools-common/channelmcp"
)

// imprint's messages come from jobs inside `imprint channel`, so under Claude
// Code a separate `imprint watch` (the plugin's mod) cannot wait for them the
// way cull's and sift's watch wait on their serve. Instead the channel writes
// each message as one JSON line to a spool file for its session, and watch
// follows that file from its start.

// Each live channel holds a shared flock on its spool, so a file no channel
// holds is a crashed session's leftover, and two channels of one session (the
// plugin's and a bare ~/.claude.json server) share one file. Whoever asks
// whether a spool is held (a starting or exiting channel, prune, watch) does
// so by trying its exclusive lock, under one gate lock per spool directory,
// so no two probes meet and none makes a channel misread the answer.

// SpoolMaxAge is how old a spool file left by a session that crashed may get
// before a starting channel prunes it. A file a live channel holds is never
// pruned.
const SpoolMaxAge = 3 * 24 * time.Hour

var (
	sessionName = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9._-]*$`)
	// watchPoll is how often watch looks for new lines.
	watchPoll = 100 * time.Millisecond
	// openedHook, when set (tests only), runs right after OpenSpool opens
	// the file, before it looks at the lock: where another starter can
	// interleave.
	openedHook func()
)

// SpoolFile is the session's spool under imprint's state directory dir.
func SpoolFile(dir, session string) (string, error) {
	if !sessionName.MatchString(session) {
		return "", fmt.Errorf("session id %q can't name a spool file", session)
	}
	return filepath.Join(dir, "spool", session+".jsonl"), nil
}

// Spool is a channel's spool file: each Notify appends one whole JSON line,
// {"text": ..., "content": ..., "meta": {...}}, in one write to a file opened
// for append.
type Spool struct {
	path  string
	f     *os.File
	lines *channelmcp.Lines
	once  sync.Once
}

var _ channelmcp.Pusher = (*Spool)(nil)

// OpenSpool prunes leftover spool files older than SpoolMaxAge and opens the
// session's spool. A file another live channel of this session holds is
// shared; one no channel holds is a crash's leftover, replaced rather than
// appended to, so watch never replays it.
func OpenSpool(dir, session string) (*Spool, error) {
	path, err := SpoolFile(dir, session)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	prune(filepath.Dir(path), time.Now().Add(-SpoolMaxAge))
	for range 10 {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		if err != nil {
			return nil, err
		}
		if openedHook != nil {
			openedHook()
		}
		retry, lerr := false, error(nil)
		if err := gate(filepath.Dir(path), func() {
			if alone(f) {
				if replaced(f, path) {
					// Another starter replaced this leftover after we opened
					// it: the path names its live spool now. Start over.
					retry = true
					return
				}
				if info, err := f.Stat(); err == nil && info.Size() > 0 {
					_ = removeSpool(path) // a leftover: start a new file
					retry = true
					return
				}
			}
			lerr = syscall.Flock(int(f.Fd()), syscall.LOCK_SH)
			retry = lerr == nil && replaced(f, path) // the last holder removed it meanwhile
		}); err != nil {
			lerr = err
		}
		if lerr != nil {
			_ = f.Close()
			return nil, lerr
		}
		if retry {
			_ = f.Close()
			continue
		}
		return &Spool{path: path, f: f, lines: channelmcp.NewLinesFor(f, "imprint")}, nil
	}
	return nil, fmt.Errorf("spool %s keeps changing under this channel", path)
}

// alone takes f's exclusive lock when no live channel holds it. Call it
// only under gate.
func alone(f *os.File) bool {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB) == nil
}

// gate runs fn holding the spool directory's gate lock.
func gate(dir string, fn func()) error {
	g, err := os.OpenFile(filepath.Join(dir, ".gate"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return err
	}
	defer func() { _ = g.Close() }() // releases the lock
	if err := syscall.Flock(int(g.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	fn()
	return nil
}

// held says whether a live channel holds f, an open spool.
func held(f *os.File, dir string) bool {
	live := false
	if err := gate(dir, func() {
		if live = !alone(f); !live {
			_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		}
	}); err != nil {
		return false
	}
	return live
}

// Path is the spool file's path.
func (s *Spool) Path() string { return s.path }

// Notify appends one message.
func (s *Spool) Notify(content string, meta map[string]string) error {
	return s.lines.Notify(content, meta)
}

// Close ends the spool and, when no other channel of the session holds it
// any more, removes its file.
func (s *Spool) Close() error {
	var err error
	s.once.Do(func() {
		gerr := gate(filepath.Dir(s.path), func() {
			if alone(s.f) && !replaced(s.f, s.path) {
				err = removeSpool(s.path)
			}
			if cerr := s.f.Close(); err == nil {
				err = cerr
			}
		})
		if gerr != nil {
			err = errors.Join(gerr, s.f.Close())
		}
	})
	return err
}

// prune removes spool files last changed before cutoff that no live channel
// holds, with their cursors, and old cursors whose spool is gone.
func prune(dir string, cutoff time.Time) {
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		info, err := e.Info()
		if e.IsDir() || err != nil || !info.ModTime().Before(cutoff) {
			continue
		}
		p := filepath.Join(dir, e.Name())
		if spool, ok := strings.CutSuffix(p, ".offset"); ok {
			if _, err := os.Stat(spool); errors.Is(err, os.ErrNotExist) {
				_ = os.Remove(p)
			}
			continue
		}
		if !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		_ = gate(dir, func() {
			if alone(f) && !replaced(f, p) {
				_ = removeSpool(p)
			}
		})
		_ = f.Close()
	}
}

// removeSpool removes a spool file and watch's cursor beside it.
func removeSpool(path string) error {
	_ = os.Remove(cursorFile(path))
	return os.Remove(path)
}

// cursorFile is where watch records how far into the spool at path it has
// printed: "<dev> <inode> <offset>", so a replaced spool starts from 0.
func cursorFile(path string) string { return path + ".offset" }

// fileID is the spool file's device and inode, as the cursor records them.
func fileID(f *os.File) (string, bool) {
	info, err := f.Stat()
	if err != nil {
		return "", false
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%d %d", st.Dev, st.Ino), true
}

// resume is the offset in f a previous watch of the session stopped at: its
// cursor's, when the cursor names this file and lies within it; else 0.
func resume(f *os.File, path string) int64 {
	id, ok := fileID(f)
	if !ok {
		return 0
	}
	b, err := os.ReadFile(cursorFile(path))
	if err != nil {
		return 0
	}
	fields := strings.Fields(string(b))
	if len(fields) != 3 || fields[0]+" "+fields[1] != id {
		return 0
	}
	off, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return 0
	}
	if info, err := f.Stat(); err != nil || off < 0 || off > info.Size() {
		return 0
	}
	return off
}

// saveCursor records, atomically and readable only by the user, that watch
// has printed f up to off.
func saveCursor(f *os.File, path string, off int64) error {
	id, ok := fileID(f)
	if !ok {
		return errors.New("no file identity for the cursor")
	}
	cursor := cursorFile(path)
	tmp, err := os.CreateTemp(filepath.Dir(cursor), filepath.Base(cursor)+".tmp*")
	if err != nil {
		return err
	}
	_, werr := fmt.Fprintf(tmp, "%s %d\n", id, off)
	cerr := tmp.Close()
	if werr == nil {
		werr = cerr
	}
	if werr == nil {
		werr = os.Chmod(tmp.Name(), 0o600)
	}
	if werr == nil {
		werr = os.Rename(tmp.Name(), cursor)
	}
	if werr != nil {
		_ = os.Remove(tmp.Name())
	}
	return werr
}

// Watch is `imprint watch`: it prints each line of the session's spool on
// out, from the file's start or from where the session's last watch stopped
// (its cursor), then follows it until ctx ends. It follows only a spool a
// live channel holds: until there is one, a leftover counts as no file. The
// file need not exist yet.
// Only whole lines are printed. A line is printed again when watch is killed
// between printing it and saving the cursor; when the channel replaces the file
// (it restarted), watch finishes the old file and starts on the new one.
func Watch(ctx context.Context, dir, session string, out, log io.Writer) error {
	path, err := SpoolFile(dir, session)
	if err != nil {
		return err
	}
	var (
		f    *os.File
		rest []byte // the start of a line whose newline hasn't been written yet
		off  int64  // where in f the next whole line starts
	)
	printed := func(n int) {
		off += int64(n)
		if err := saveCursor(f, path, off); err != nil {
			_, _ = fmt.Fprintf(log, "imprint watch: cursor: %v\n", err)
		}
	}
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
			} else if !held(f, filepath.Dir(path)) {
				_ = f.Close() // a crash's leftover: wait for a live channel's
				f = nil
			} else if off = resume(f, path); off > 0 {
				if _, err := f.Seek(off, io.SeekStart); err != nil {
					off = 0
				}
			}
		}
		if f != nil {
			if rest, err = drain(f, rest, out, printed); err != nil {
				return err
			}
			if replaced(f, path) {
				// Lines written before the swap are in the old file still.
				if rest, err = drain(f, rest, out, printed); err != nil {
					return err
				}
				if len(rest) > 0 {
					_, _ = fmt.Fprintf(log, "imprint watch: dropped an unfinished line from the old spool: %q\n", rest)
				}
				_ = f.Close()
				f, rest, off = nil, nil, 0
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

// drain reads f to its end and prints each whole line, telling done how many
// bytes each line took, and returns what is left of an unfinished one.
func drain(f *os.File, rest []byte, out io.Writer, done func(int)) ([]byte, error) {
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
			done(i + 1)
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
