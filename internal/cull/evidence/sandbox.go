// Package evidence measures which tests catch small planted bugs: it
// copies a project into sandboxes (the working tree is never modified),
// plants mutants in one sandbox at a time, runs the tests exercising each
// mutated module, and records the per-module kill matrix.
package evidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	tools "github.com/schuettc/tools-common"
)

// Sandbox is one private copy of the project, owned by at most one worker
// at a time.
type Sandbox struct{ Dir string }

// heavyDirs are dependency directories symlinked into a sandbox rather
// than copied. vendor is only heavy when git tracks nothing under it.
var heavyDirs = map[string]bool{"node_modules": true, "venv": true, ".venv": true, "vendor": true}

// Prepare returns n sandboxes for root under
// tools.CacheDir("cull")/evidence/<sha256(root)[:12]>/w<N>. Each holds a
// copy of every file `git ls-files -co --exclude-standard` lists, with
// heavy dependency directories (node_modules, venv, .venv, and vendor when
// untracked) symlinked to root's instead; one that is itself a symlink in
// root is linked to the same target (a relative target resolved against
// root). Existing sandboxes are
// refreshed: changed files are rewritten, and anything not in the listing
// (stale files, build output) is removed.
//
// Unreadable untracked files are skipped (logged to os.Stderr). Prepare
// holds the project's sandbox lock (see Run) while it works and fails
// fast if another Prepare or Run holds it.
func Prepare(root string, n int) ([]Sandbox, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	base, unlock, err := lockBase(abs)
	if err != nil {
		return nil, err
	}
	defer unlock()
	return prepare(abs, base, n, func(format string, args ...any) {
		fmt.Fprintf(os.Stderr, "cull evidence: "+format+"\n", args...)
	})
}

// sandboxBase is the directory holding abs's sandboxes.
func sandboxBase(abs string) string {
	sum := sha256.Sum256([]byte(abs))
	return filepath.Join(tools.CacheDir("cull"), "evidence", hex.EncodeToString(sum[:])[:12])
}

// lockBase takes an exclusive, non-blocking lock on <base>/.lock for the
// project at abs, so two runs never share w0..wN. It returns the base
// directory and the unlock function.
func lockBase(abs string) (string, func(), error) {
	base := sandboxBase(abs)
	if err := tools.EnsureDir(base); err != nil {
		return "", nil, err
	}
	lockPath := filepath.Join(base, ".lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return "", nil, fmt.Errorf("evidence: open lock %s: %w", lockPath, err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return "", nil, fmt.Errorf("evidence: another cull evidence run on %s holds the lock %s: %w", abs, lockPath, err)
	}
	return base, func() {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
	}, nil
}

// prepare is Prepare for an absolute root whose lock the caller holds.
func prepare(abs, base string, n int, logf func(string, ...any)) ([]Sandbox, error) {
	if n < 1 {
		n = 1
	}
	plan, err := planCopy(abs, logf)
	if err != nil {
		return nil, err
	}
	out := make([]Sandbox, 0, n)
	for i := 0; i < n; i++ {
		dir := filepath.Join(base, fmt.Sprintf("w%d", i))
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return nil, err
		}
		if err := plan.sync(abs, dir); err != nil {
			return nil, fmt.Errorf("evidence: sandbox %s: %w", dir, err)
		}
		out = append(out, Sandbox{Dir: dir})
	}
	return out, nil
}

// copyPlan is the set of root-relative ('/'-separated) paths a sandbox
// holds: files copied, heavy directories symlinked (heavy maps each to its
// link target), and every directory that is an ancestor of one of them.
type copyPlan struct {
	files map[string]bool
	heavy map[string]string
	dirs  map[string]bool
}

func gitList(root string, args ...string) ([]string, error) {
	cmd := exec.Command("git", append([]string{"ls-files", "-z"}, args...)...)
	cmd.Dir = root
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("evidence: git ls-files %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	var paths []string
	for _, p := range strings.Split(string(out), "\x00") {
		if p != "" {
			paths = append(paths, p)
		}
	}
	return paths, nil
}

func planCopy(root string, logf func(string, ...any)) (*copyPlan, error) {
	listed, err := gitList(root, "-co", "--exclude-standard")
	if err != nil {
		return nil, err
	}
	tracked, err := gitList(root)
	if err != nil {
		return nil, err
	}
	ignoredDirs, err := gitList(root, "-oi", "--exclude-standard", "--directory")
	if err != nil {
		return nil, err
	}
	trackedDirs := map[string]bool{}
	isTracked := map[string]bool{}
	for _, p := range tracked {
		isTracked[p] = true
		for d := path.Dir(p); d != "."; d = path.Dir(d) {
			trackedDirs[d] = true
		}
	}
	// heavyPrefix returns the outermost heavy directory containing p (or p
	// itself), or "" if there is none.
	heavyPrefix := func(p string) string {
		parts := strings.Split(strings.TrimSuffix(p, "/"), "/")
		for i, c := range parts {
			if !heavyDirs[c] {
				continue
			}
			d := strings.Join(parts[:i+1], "/")
			if c == "vendor" && trackedDirs[d] {
				continue
			}
			return d
		}
		return ""
	}

	plan := &copyPlan{files: map[string]bool{}, heavy: map[string]string{}, dirs: map[string]bool{}}
	addDirs := func(p string) {
		for d := path.Dir(p); d != "."; d = path.Dir(d) {
			plan.dirs[d] = true
		}
	}
	addHeavy := func(d string) {
		src := filepath.Join(root, filepath.FromSlash(d))
		fi, err := os.Lstat(src)
		if err != nil {
			return
		}
		target := src
		if fi.Mode()&fs.ModeSymlink != 0 {
			t, err := os.Readlink(src)
			if err != nil {
				return
			}
			if !filepath.IsAbs(t) {
				t = filepath.Join(filepath.Dir(src), t)
			}
			target = t
		}
		plan.heavy[d] = target
		addDirs(d)
	}
	for _, p := range listed {
		if h := heavyPrefix(p); h != "" {
			addHeavy(h)
			continue
		}
		if plan.files[p] {
			continue
		}
		fi, err := os.Lstat(filepath.Join(root, filepath.FromSlash(p)))
		if err != nil || fi.IsDir() {
			continue // deleted but still tracked, or a submodule
		}
		if !isTracked[p] && fi.Mode().IsRegular() {
			f, err := os.Open(filepath.Join(root, filepath.FromSlash(p)))
			if err != nil {
				logf("%s: %v; not copied", p, err)
				continue
			}
			f.Close()
		}
		plan.files[p] = true
		addDirs(p)
	}
	for _, p := range ignoredDirs {
		if !strings.HasSuffix(p, "/") {
			// An ignored symlink is listed as a file: keep it only if it
			// is a heavy directory itself.
			if fi, err := os.Lstat(filepath.Join(root, filepath.FromSlash(p))); err != nil ||
				fi.Mode()&fs.ModeSymlink == 0 || heavyPrefix(p) != p {
				continue
			}
		}
		if h := heavyPrefix(p); h != "" {
			addHeavy(h)
		}
	}
	return plan, nil
}

// sync makes dir match the plan: prune anything the plan doesn't hold,
// then copy changed files and (re)link heavy directories.
func (p *copyPlan) sync(root, dir string) error {
	err := filepath.WalkDir(dir, func(abs string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if abs == dir {
			return nil
		}
		rel, err := filepath.Rel(dir, abs)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		keep := false
		switch {
		case p.heavy[rel] != "":
			keep = d.Type()&fs.ModeSymlink != 0
		case p.files[rel]:
			keep = !d.IsDir()
		case p.dirs[rel]:
			keep = d.IsDir()
		}
		if keep {
			return nil
		}
		if err := os.RemoveAll(abs); err != nil {
			return err
		}
		if d.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return err
	}

	for rel := range p.files {
		if err := syncFile(filepath.Join(root, filepath.FromSlash(rel)), filepath.Join(dir, filepath.FromSlash(rel))); err != nil {
			return err
		}
	}
	for rel, src := range p.heavy {
		dst := filepath.Join(dir, filepath.FromSlash(rel))
		if t, err := os.Readlink(dst); err == nil && t == src {
			continue
		}
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.Symlink(src, dst); err != nil {
			return err
		}
	}
	return nil
}

// syncFile copies src to dst unless dst already has identical content (a
// symlink is recreated with the same target, never followed).
func syncFile(src, dst string) error {
	fi, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if fi.Mode()&fs.ModeSymlink != 0 {
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		if t, err := os.Readlink(dst); err == nil && t == target {
			return nil
		}
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		return os.Symlink(target, dst)
	}
	want, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if dfi, err := os.Lstat(dst); err == nil && dfi.Mode().IsRegular() && dfi.Size() == int64(len(want)) {
		if have, err := os.ReadFile(dst); err == nil && bytes.Equal(have, want) {
			if dfi.Mode().Perm() != sandboxPerm(fi) {
				return os.Chmod(dst, sandboxPerm(fi))
			}
			return nil
		}
	}
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, want, sandboxPerm(fi))
}

// sandboxPerm is src's permission bits made owner-writable, so a mutant
// can always be planted in (and restored to) a read-only module's copy.
func sandboxPerm(src fs.FileInfo) fs.FileMode {
	return src.Mode().Perm() | 0o200
}

// copiedFile reports whether rel is a regular file copied into sandbox dir
// (no path component is a symlink), i.e. safe to mutate there.
func copiedFile(dir, rel string) bool {
	clean := path.Clean(filepath.ToSlash(rel))
	if clean == "." || strings.HasPrefix(clean, "../") || clean == ".." || path.IsAbs(clean) {
		return false
	}
	cur := dir
	parts := strings.Split(clean, "/")
	for i, c := range parts {
		cur = filepath.Join(cur, c)
		fi, err := os.Lstat(cur)
		if err != nil || fi.Mode()&fs.ModeSymlink != 0 {
			return false
		}
		if i == len(parts)-1 {
			return fi.Mode().IsRegular()
		}
	}
	return false
}
