package evidence

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"
)

func needPython(t *testing.T) {
	t.Helper()
	needTools(t)
	if err := exec.Command("python3", "-m", "pytest", "--version").Run(); err != nil {
		t.Skip("python3 with pytest not available")
	}
}

// Same-size mutants planted in one sandbox within one second must not run
// each other's cached bytecode. plantHook pins every planted mutant's mtime
// to one instant, so CPython's (mtime, size) pyc check can't tell them apart.
func TestNoStaleBytecode(t *testing.T) {
	needPython(t)
	root := newRepo(t, map[string]string{
		"calc.py": "def f(a, b):\n    s = a + b\n    ok = a < b\n    return (s, ok)\n",
		"test_calc.py": "from calc import f\n\n\ndef test_add():\n    assert f(2, 3)[0] == 5\n\n\n" +
			"def test_lt():\n    assert f(1, 2)[1]\n    assert not f(2, 1)[1]\n    assert not f(2, 2)[1]\n",
	})
	pinned := time.Now().Truncate(time.Second)
	plantHook = func(file string) { _ = os.Chtimes(file, pinned, pinned) }
	defer func() { plantHook = nil }()
	mod := Module{Rel: "calc.py", Lang: "python", TestFiles: []string{"test_calc.py"},
		Tests: []string{"py:test_calc.py:test_add", "py:test_calc.py:test_lt"}}
	m, err := Run(context.Background(), root, []Module{mod}, opts(1))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]string{
		"line 2: flip +":      {"py:test_calc.py:test_add"},
		"line 3: flip <":      {"py:test_calc.py:test_lt"},
		"line 4: return None": {"py:test_calc.py:test_add", "py:test_calc.py:test_lt"},
	}
	got := map[string][]string{}
	for _, mu := range m.Mutants {
		got[mu.Site.Desc] = mu.KilledBy
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("matrix:\n got %v\nwant %v", got, want)
	}
}

// A test shared by two modules that flakes only in the second module's
// baseline is never credited with a kill in the first.
func TestFlakyElsewhereNotCredited(t *testing.T) {
	needTools(t)
	dir := t.TempDir()
	m0, m1 := filepath.Join(dir, "m0"), filepath.Join(dir, "m1")
	// Passes its first two runs (module a's baseline), fails from then on.
	shared := `package calc

import (
	"syscall"
	"testing"
)

func exists(p string) bool {
	fd, err := syscall.Open(p, syscall.O_RDONLY, 0)
	if err == nil {
		syscall.Close(fd)
	}
	return err == nil
}

func touch(p string) {
	if fd, err := syscall.Open(p, syscall.O_CREAT|syscall.O_WRONLY, 0o644); err == nil {
		syscall.Close(fd)
	}
}

func TestShared(t *testing.T) {
	if exists(` + "`" + m1 + "`" + `) {
		t.Fatal("third run on fails")
	}
	if exists(` + "`" + m0 + "`" + `) {
		touch(` + "`" + m1 + "`" + `)
	} else {
		touch(` + "`" + m0 + "`" + `)
	}
	if Add(2, 3) != 5 || Sub(3, 2) != 1 {
		t.Fatal("math")
	}
}
`
	root := newRepo(t, map[string]string{
		"go.mod":         "module fixture\n\ngo 1.22\n",
		"a.go":           "package calc\n\nfunc Add(a, b int) int { return a + b }\n",
		"b.go":           "package calc\n\nfunc Sub(a, b int) int { return a - b }\n",
		"shared_test.go": shared,
	})
	id := "go:shared_test.go:TestShared"
	mods := []Module{
		{Rel: "a.go", Lang: "go", TestFiles: []string{"shared_test.go"}, Tests: []string{id}},
		{Rel: "b.go", Lang: "go", TestFiles: []string{"shared_test.go"}, Tests: []string{id}},
	}
	m, err := Run(context.Background(), root, mods, opts(1))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(m.Flaky, []string{id}) {
		t.Fatalf("flaky = %v", m.Flaky)
	}
	if len(m.Mutants) == 0 {
		t.Fatal("no mutants for a.go")
	}
	for _, mu := range m.Mutants {
		if mu.KilledBy == nil {
			t.Errorf("%s excluded, want non-nil empty", mu.Site.Desc)
		}
		for _, k := range mu.KilledBy {
			if k == id {
				t.Errorf("%s %s credited to flaky %s", mu.Module, mu.Site.Desc, id)
			}
		}
	}
}

func skipRoot(t *testing.T) {
	t.Helper()
	if os.Geteuid() == 0 {
		t.Skip("permission checks don't apply to root")
	}
}

// A read-only module is still mutable in its sandbox.
func TestReadOnlyModuleMutable(t *testing.T) {
	needTools(t)
	skipRoot(t)
	root := calcRepo(t)
	if err := os.Chmod(filepath.Join(root, "calc.go"), 0o444); err != nil {
		t.Fatal(err)
	}
	m, err := Run(context.Background(), root, []Module{calcModule()}, opts(1))
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Mutants) != 3 {
		t.Fatalf("mutants = %v", rows(m))
	}
	fi, err := os.Stat(filepath.Join(root, "calc.go"))
	if err != nil || fi.Mode().Perm() != 0o444 {
		t.Fatalf("working tree mode changed: %v %v", fi.Mode(), err)
	}
}

// An unreadable untracked file is skipped, not fatal.
func TestUnreadableUntrackedSkipped(t *testing.T) {
	needTools(t)
	skipRoot(t)
	root := newRepo(t, map[string]string{"a.txt": "a\n", "secret.txt": "s\n"})
	git(t, root, "add", "a.txt")
	secret := filepath.Join(root, "secret.txt")
	if err := os.Chmod(secret, 0); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(secret, 0o644)
	sbs, err := Prepare(root, 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(sbs[0].Dir, "secret.txt")); err == nil {
		t.Error("unreadable file copied")
	}
	if b, _ := os.ReadFile(filepath.Join(sbs[0].Dir, "a.txt")); string(b) != "a\n" {
		t.Errorf("a.txt = %q", b)
	}
}

// A second concurrent run on the same project fails fast, naming the lock.
func TestConcurrentRunLocked(t *testing.T) {
	needTools(t)
	root := calcRepo(t)
	sbs, err := Prepare(root, 1)
	if err != nil {
		t.Fatal(err)
	}
	lockPath := filepath.Join(filepath.Dir(sbs[0].Dir), ".lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	_, err = Run(context.Background(), root, []Module{calcModule()}, opts(1))
	if err == nil || !strings.Contains(err.Error(), lockPath) {
		t.Fatalf("Run err = %v, want lock error naming %s", err, lockPath)
	}
	if _, err := Prepare(root, 1); err == nil || !strings.Contains(err.Error(), lockPath) {
		t.Fatalf("Prepare err = %v", err)
	}
	if el := time.Since(start); el > 5*time.Second {
		t.Fatalf("lock failure took %v", el)
	}
	_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	if _, err := Run(context.Background(), root, []Module{calcModule()}, opts(1)); err != nil {
		t.Fatalf("after unlock: %v", err)
	}
}
