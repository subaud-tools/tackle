package evidence

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

// TestMain isolates the sandbox cache (tools.CacheDir("cull")) from the
// developer's real one.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "cull-evidence-cache")
	if err != nil {
		panic(err)
	}
	os.Setenv("CULL_HOME", dir)
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

func needTools(t *testing.T) {
	t.Helper()
	for _, bin := range []string{"git", "go"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("%s not on PATH", bin)
		}
	}
}

func write(t *testing.T, root, rel, body string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func newRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	git(t, root, "init", "-q")
	for rel, body := range files {
		write(t, root, rel, body)
	}
	return root
}

const calcSrc = `package calc

func Add(a, b int) int { return a + b }

func Pos(n int) bool { return n > 0 }
`

const calcTest = `package calc

import "testing"

func TestAddSole(t *testing.T) {
	if Add(2, 3) != 5 {
		t.Fatal("add")
	}
}

func TestPosA(t *testing.T) {
	if !Pos(1) || Pos(0) {
		t.Fatal("pos")
	}
}

func TestPosB(t *testing.T) {
	got := []bool{Pos(1), Pos(0)}
	if got[0] != true || got[1] != false {
		t.Fatal("pos b")
	}
}
`

func calcRepo(t *testing.T) string {
	return newRepo(t, map[string]string{
		"go.mod":       "module fixture\n\ngo 1.22\n",
		"calc.go":      calcSrc,
		"calc_test.go": calcTest,
	})
}

func calcModule() Module {
	return Module{
		Rel:       "calc.go",
		Lang:      "go",
		TestFiles: []string{"calc_test.go"},
		Tests: []string{
			"go:calc_test.go:TestAddSole",
			"go:calc_test.go:TestPosA",
			"go:calc_test.go:TestPosB",
		},
	}
}

func opts(workers int) Options {
	// NoCache: tests rerunning a root expect real runs (cache_test.go
	// covers the cache).
	return Options{Workers: workers, MaxMutants: 30, RunTimeout: 60 * time.Second, Stderr: &bytes.Buffer{}, NoCache: true}
}

type row struct {
	Desc     string
	KilledBy []string
}

func rows(m Matrix) []row {
	var out []row
	for _, mu := range m.Mutants {
		out = append(out, row{mu.Site.Desc, mu.KilledBy})
	}
	return out
}

func TestExactMatrix(t *testing.T) {
	needTools(t)
	root := calcRepo(t)
	m, err := RunMatrix(context.Background(), root, []Module{calcModule()}, opts(2))
	if err != nil {
		t.Fatal(err)
	}
	want := []row{
		{"line 3: flip +", []string{"go:calc_test.go:TestAddSole"}},
		{"line 5: flip >", []string{"go:calc_test.go:TestPosA", "go:calc_test.go:TestPosB"}},
		{"line 5: +1 int", []string{"go:calc_test.go:TestPosA", "go:calc_test.go:TestPosB"}},
	}
	if got := rows(m); !reflect.DeepEqual(got, want) {
		t.Fatalf("matrix:\n got %#v\nwant %#v", got, want)
	}
	for _, mu := range m.Mutants {
		if mu.Module != "calc.go" {
			t.Errorf("mutant module = %q", mu.Module)
		}
	}
	wantBase := map[string]bool{
		"go:calc_test.go:TestAddSole": true,
		"go:calc_test.go:TestPosA":    true,
		"go:calc_test.go:TestPosB":    true,
	}
	if !reflect.DeepEqual(m.Baseline, wantBase) {
		t.Fatalf("baseline = %v", m.Baseline)
	}
	if len(m.Flaky) != 0 {
		t.Fatalf("flaky = %v", m.Flaky)
	}
}

func hashTree(t *testing.T, root string) map[string][32]byte {
	t.Helper()
	out := map[string][32]byte{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		out[p] = sha256.Sum256(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestWorkingTreeUntouched(t *testing.T) {
	needTools(t)
	root := calcRepo(t)
	before := hashTree(t, root)
	if _, err := RunMatrix(context.Background(), root, []Module{calcModule()}, opts(3)); err != nil {
		t.Fatal(err)
	}
	if after := hashTree(t, root); !reflect.DeepEqual(before, after) {
		t.Fatalf("working tree changed:\nbefore %v\nafter  %v", before, after)
	}
}

func TestParallelMatchesSerial(t *testing.T) {
	needTools(t)
	root := calcRepo(t)
	serial, err := RunMatrix(context.Background(), root, []Module{calcModule()}, opts(1))
	if err != nil {
		t.Fatal(err)
	}
	par, err := RunMatrix(context.Background(), root, []Module{calcModule()}, opts(4))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(serial.Mutants, par.Mutants) {
		t.Fatalf("serial %v\nparallel %v", serial.Mutants, par.Mutants)
	}
	if !reflect.DeepEqual(serial.Baseline, par.Baseline) {
		t.Fatalf("baseline serial %v parallel %v", serial.Baseline, par.Baseline)
	}
}

func TestFlakyTestUnmeasured(t *testing.T) {
	needTools(t)
	marker := filepath.Join(t.TempDir(), "ran-once")
	// syscall access bypasses go test's file tracking, so a cached replay of
	// the first baseline run would hide this flake.
	flaky := `package calc

import (
	"syscall"
	"testing"
)

func TestFlaky(t *testing.T) {
	fd, err := syscall.Open(` + "`" + marker + "`" + `, syscall.O_RDONLY, 0)
	if err == nil {
		syscall.Close(fd)
		t.Fatal("second run fails")
	}
	fd, err = syscall.Open(` + "`" + marker + "`" + `, syscall.O_CREAT|syscall.O_WRONLY, 0o644)
	if err == nil {
		syscall.Close(fd)
	}
	if !Pos(1) {
		t.Fatal("pos")
	}
}
`
	root := newRepo(t, map[string]string{
		"go.mod":        "module fixture\n\ngo 1.22\n",
		"calc.go":       calcSrc,
		"calc_test.go":  calcTest,
		"flaky_test.go": flaky,
	})
	mod := calcModule()
	mod.TestFiles = append(mod.TestFiles, "flaky_test.go")
	mod.Tests = append(mod.Tests, "go:flaky_test.go:TestFlaky")
	m, err := RunMatrix(context.Background(), root, []Module{mod}, opts(2))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(m.Flaky, []string{"go:flaky_test.go:TestFlaky"}) {
		t.Fatalf("flaky = %v", m.Flaky)
	}
	if m.Baseline["go:flaky_test.go:TestFlaky"] {
		t.Fatal("flaky test measured")
	}
	if !m.Baseline["go:calc_test.go:TestPosA"] {
		t.Fatal("stable test not measured")
	}
	for _, mu := range m.Mutants {
		for _, id := range mu.KilledBy {
			if id == "go:flaky_test.go:TestFlaky" {
				t.Fatalf("flaky test credited with a kill: %v", mu)
			}
		}
	}
}

func TestHangingMutantExcluded(t *testing.T) {
	needTools(t)
	root := newRepo(t, map[string]string{
		"go.mod": "module fixture\n\ngo 1.22\n",
		"halve.go": `package halve

func Halve(x int) int {
	n := 0
	for x != 0 {
		x -= 2
		n++
	}
	return n
}
`,
		"halve_test.go": `package halve

import "testing"

func TestHalve(t *testing.T) {
	if Halve(4) != 2 {
		t.Fatal("halve")
	}
}
`,
	})
	mod := Module{Rel: "halve.go", Lang: "go", TestFiles: []string{"halve_test.go"}, Tests: []string{"go:halve_test.go:TestHalve"}}
	o := opts(2)
	o.RunTimeout = 8 * time.Second
	start := time.Now()
	m, err := RunMatrix(context.Background(), root, []Module{mod}, o)
	if err != nil {
		t.Fatal(err)
	}
	if el := time.Since(start); el > 60*time.Second {
		t.Fatalf("run took %v", el)
	}
	killer := []string{"go:halve_test.go:TestHalve"}
	want := []row{
		{"line 4: +1 int", killer},
		{"line 5: flip !=", killer},
		{"line 5: +1 int", nil},
		{"line 6: +1 int", nil},
	}
	if got := rows(m); !reflect.DeepEqual(got, want) {
		t.Fatalf("matrix:\n got %#v\nwant %#v", got, want)
	}
}

func TestSymlinkedDepsNotCopied(t *testing.T) {
	needTools(t)
	root := newRepo(t, map[string]string{
		".gitignore":            "node_modules/\n.venv/\n",
		"app.py":                "X = 1\n",
		"node_modules/pkg/i.js": "module.exports = 1\n",
		".venv/lib/site.py":     "Y = 2\n",
		"vendor/dep/dep.py":     "Z = 3\n",
		"lib/vendor/tracked.py": "W = 4\n",
		"ignored.log":           "",
	})
	write(t, root, ".gitignore", "node_modules/\n.venv/\n*.log\n")
	git(t, root, "add", "lib/vendor/tracked.py")

	sbs, err := Prepare(root, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(sbs) != 2 || sbs[0].Dir == sbs[1].Dir {
		t.Fatalf("sandboxes = %v", sbs)
	}
	for _, sb := range sbs {
		if !strings.Contains(sb.Dir, filepath.Join("evidence")) {
			t.Errorf("sandbox %s not under the evidence cache", sb.Dir)
		}
		for _, heavy := range []string{"node_modules", ".venv", "vendor"} {
			fi, err := os.Lstat(filepath.Join(sb.Dir, heavy))
			if err != nil {
				t.Fatalf("%s: %v", heavy, err)
			}
			if fi.Mode()&os.ModeSymlink == 0 {
				t.Errorf("%s copied, want symlink", heavy)
			}
			if target, _ := os.Readlink(filepath.Join(sb.Dir, heavy)); target != filepath.Join(root, heavy) {
				t.Errorf("%s -> %s", heavy, target)
			}
		}
		for _, rel := range []string{"app.py", "lib/vendor/tracked.py", ".gitignore"} {
			fi, err := os.Lstat(filepath.Join(sb.Dir, filepath.FromSlash(rel)))
			if err != nil || !fi.Mode().IsRegular() {
				t.Errorf("%s not copied as a regular file: %v", rel, err)
			}
		}
		if fi, err := os.Lstat(filepath.Join(sb.Dir, "lib", "vendor")); err != nil || !fi.IsDir() {
			t.Errorf("tracked vendor dir not copied: %v", err)
		}
		if _, err := os.Lstat(filepath.Join(sb.Dir, "ignored.log")); err == nil {
			t.Error("ignored file copied")
		}
	}

	// Mutants never target symlinked content.
	m, err := RunMatrix(context.Background(), root, []Module{{Rel: "vendor/dep/dep.py", Lang: "python", TestFiles: []string{"app.py"}}}, opts(1))
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Mutants) != 0 {
		t.Fatalf("mutated symlinked content: %v", m.Mutants)
	}
}

func TestPrepareRefreshes(t *testing.T) {
	needTools(t)
	root := newRepo(t, map[string]string{"a.txt": "one\n", "b.txt": "gone soon\n"})
	sbs, err := Prepare(root, 1)
	if err != nil {
		t.Fatal(err)
	}
	write(t, root, "a.txt", "two\n")
	if err := os.Remove(filepath.Join(root, "b.txt")); err != nil {
		t.Fatal(err)
	}
	write(t, sbs[0].Dir, "stray.txt", "leftover\n")
	sbs, err = Prepare(root, 1)
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(sbs[0].Dir, "a.txt")); string(b) != "two\n" {
		t.Fatalf("a.txt = %q", b)
	}
	var left []string
	entries, _ := os.ReadDir(sbs[0].Dir)
	for _, e := range entries {
		left = append(left, e.Name())
	}
	sort.Strings(left)
	if !reflect.DeepEqual(left, []string{"a.txt"}) {
		t.Fatalf("sandbox entries = %v", left)
	}
}
