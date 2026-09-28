package evidence

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// Cancellation stops dispatch at once: no test run starts after the
// cancel, the in-flight (slow) run is killed, every sandbox's module is
// restored, and Run writes nothing to .cull/.
func TestCancelStopsRun(t *testing.T) {
	needPython(t)
	state := t.TempDir()
	runs := filepath.Join(state, "runs")
	slow := filepath.Join(state, "slow")
	src := "def f(x):\n    a = x + 1\n    b = x * 2\n    c = x - 3\n    return a < b and c > 0\n"
	root := newRepo(t, map[string]string{
		"calc.py": src,
		"test_calc.py": "import os, time\nfrom calc import f\n\n" +
			"with open(" + pyQuote(runs) + ", 'a') as fh:\n    fh.write('run\\n')\n\n\n" +
			"def test_f():\n    if os.path.exists(" + pyQuote(slow) + "):\n        time.sleep(30)\n    assert f(10)\n",
	})
	writeLast(t, root, []string{"test_calc.py"}, map[string]verdict{"test_f": {"review", "review_band"}})
	last, _ := os.ReadFile(filepath.Join(root, ".cull", "last.json"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var (
		mu       sync.Mutex
		planted  int
		atCancel int
		canceled time.Time
	)
	count := func() int {
		b, _ := os.ReadFile(runs)
		return strings.Count(string(b), "run\n")
	}
	plantHook = func(string) {
		mu.Lock()
		defer mu.Unlock()
		planted++
		switch planted {
		case 1:
			_ = os.WriteFile(slow, nil, 0o644) // this run hangs until killed
		case 2:
			time.Sleep(time.Second) // let the first slow run start
			atCancel, canceled = count(), time.Now()
			cancel()
		}
	}
	defer func() { plantHook = nil }()
	o := settleOpts()
	o.Workers = 2
	_, err := Run(ctx, root, o)
	elapsed := time.Since(canceled)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
	if elapsed > 8*time.Second {
		t.Errorf("Run returned %s after cancel", elapsed)
	}
	if n := count(); n != atCancel {
		t.Errorf("%d test runs started after cancel", n-atCancel)
	}
	sbs, _ := filepath.Glob(filepath.Join(sandboxBase(mustAbs(t, root)), "w*", "calc.py"))
	if len(sbs) == 0 {
		t.Fatal("no sandboxes")
	}
	for _, p := range sbs {
		if b, _ := os.ReadFile(p); string(b) != src {
			t.Errorf("%s not restored", p)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".cull", "evidence.json")); !os.IsNotExist(err) {
		t.Error("evidence.json written")
	}
	if b, _ := os.ReadFile(filepath.Join(root, ".cull", "last.json")); string(b) != string(last) {
		t.Error("last.json changed")
	}
}

func mustAbs(t *testing.T, p string) string {
	t.Helper()
	a, err := filepath.Abs(p)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

// syncBuf is a goroutine-safe log buffer.
type syncBuf struct {
	mu sync.Mutex
	b  strings.Builder
}

func (s *syncBuf) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuf) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

// Cancelled during a (hanging) baseline, evidence does no further work:
// no test file is isolated or dropped, no baseline is reported, and it
// returns promptly with the cancellation.
func TestCancelDuringBaseline(t *testing.T) {
	needPython(t)
	state := t.TempDir()
	runs := filepath.Join(state, "runs")
	test := func(n string) string {
		return "import time\nfrom calc import f\n\n" +
			"with open(" + pyQuote(runs) + ", 'a') as fh:\n    fh.write('run\\n')\n\n\n" +
			"def test_" + n + "():\n    time.sleep(30)\n    assert f(10)\n"
	}
	root := newRepo(t, map[string]string{
		"calc.py":      "def f(x):\n    return x > 1\n",
		"test_a.py":    test("a"),
		"test_b.py":    test("b"),
		"test_c.py":    test("c"),
		"test_d.py":    test("d"),
		"test_calc.py": test("calc"),
	})
	mod := Module{Rel: "calc.py", Lang: "python", TestFiles: []string{"test_a.py", "test_b.py", "test_c.py", "test_calc.py", "test_d.py"},
		Tests: []string{"py:test_a.py:test_a", "py:test_b.py:test_b", "py:test_c.py:test_c", "py:test_calc.py:test_calc", "py:test_d.py:test_d"}}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var canceled time.Time
	go func() {
		for {
			if b, _ := os.ReadFile(runs); len(b) > 0 {
				time.Sleep(200 * time.Millisecond)
				canceled = time.Now()
				cancel()
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()
	var log syncBuf
	o := opts(1)
	o.Stderr = &log
	_, err := RunMatrix(ctx, root, []Module{mod}, o)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v\n%s", err, log.String())
	}
	if d := time.Since(canceled); d > 8*time.Second {
		t.Errorf("returned %s after cancel", d)
	}
	for _, bad := range []string{"dropped", "baseline calc.py", "canary"} {
		if strings.Contains(log.String(), bad) {
			t.Errorf("work after cancel (%q):\n%s", bad, log.String())
		}
	}
}
