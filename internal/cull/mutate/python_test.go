package mutate

import (
	"bytes"
	"os/exec"
	"regexp"
	"testing"
)

// requirePython3 skips the test unless python3 is on PATH.
func requirePython3(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not found in PATH, skipping")
	}
}

const pyOperatorsFixture = `"""Module docstring, never mutated."""
import os


def helper(a, b, ok, done):
    """Helper docstring, never mutated."""
    if a == b:
        pass
    if a != b:
        pass
    if a < b:
        pass
    if a > b:
        pass
    if a <= b:
        pass
    if a >= b:
        pass
    if a is b:
        pass
    if a is not b:
        pass
    if a in (b,):
        pass
    if a not in (b,):
        pass
    if ok and not done:
        pass
    if ok or not done:
        pass
    if not ok:
        pass
    s = a + b
    d = a - b
    p = a * b
    n = 41
    f = 3.5
    greeting = "hi"
    flag_t = True
    flag_f = False
    return s


def other():
    return 1 + 2
`

var pyDescSuffix = regexp.MustCompile(`^line \d+: (.+)$`)

func TestPythonSitesOperators(t *testing.T) {
	requirePython3(t)
	m := For("python")
	if m == nil {
		t.Fatal(`For("python") = nil`)
	}
	sites, err := m.Sites("", "fixture.py", []byte(pyOperatorsFixture))
	if err != nil {
		t.Fatalf("Sites: %v", err)
	}

	want := map[string]int{
		"flip ==":       1,
		"flip !=":       1,
		"flip <":        1,
		"flip >":        1,
		"flip <=":       1,
		"flip >=":       1,
		"flip is":       1,
		"flip is not":   1,
		"flip in":       1,
		"flip not in":   1,
		"flip and":      1,
		"flip or":       1,
		"drop not":      3,
		"flip +":        2,
		"flip -":        1,
		"flip *":        1,
		"+1 int":        3,
		"+1 float":      1,
		"mutate string": 1,
		"flip True":     1,
		"flip False":    1,
		"return None":   2,
	}
	got := map[string]int{}
	total := 0
	for _, s := range sites {
		sub := pyDescSuffix.FindStringSubmatch(s.Desc)
		if sub == nil {
			t.Errorf("Desc %q does not match expected format", s.Desc)
			continue
		}
		got[sub[1]]++
		total++
	}
	if total != 28 {
		t.Errorf("len(sites) = %d, want 28; got desc counts %+v", total, got)
	}
	for k, wantN := range want {
		if got[k] != wantN {
			t.Errorf("desc %q: got %d, want %d", k, got[k], wantN)
		}
	}

	// Every single-site mutation of the fixture must still parse as
	// Python.
	src := []byte(pyOperatorsFixture)
	for _, s := range sites {
		mutated, err := Apply(src, s)
		if err != nil {
			t.Fatalf("Apply(%+v): %v", s, err)
		}
		parsePythonFixture(t, mutated, s.Desc)
	}
}

func TestPythonSitesSkipDocstringsAndImports(t *testing.T) {
	requirePython3(t)
	m := For("python")

	const withDocstringsAndImports = `"""Module docstring."""
import os
import sys as system


class C:
    """Class docstring."""

    def method(self):
        """Method docstring."""
        return os.sep == system.path[0]
`
	sites, err := m.Sites("", "fixture.py", []byte(withDocstringsAndImports))
	if err != nil {
		t.Fatalf("Sites: %v", err)
	}
	src := []byte(withDocstringsAndImports)
	for _, s := range sites {
		text := string(src[s.Start:s.End])
		if text == `"""Module docstring."""` || text == `"""Class docstring."""` || text == `"""Method docstring."""` {
			t.Errorf("mutated a docstring: %+v", s)
		}
		if text == "os" || text == "sys" {
			t.Errorf("mutated an import name: %+v", s)
		}
	}
}

func TestPythonSitesUnicodeOffsets(t *testing.T) {
	requirePython3(t)
	const src = `greeting = "héllo wörld"


def f(a, b):
    return a == b
`
	m := For("python")
	sites, err := m.Sites("", "fixture.py", []byte(src))
	if err != nil {
		t.Fatalf("Sites: %v", err)
	}
	var eq *Site
	for i := range sites {
		if sites[i].Desc == "line 5: flip ==" {
			eq = &sites[i]
		}
	}
	if eq == nil {
		t.Fatalf("no comparison site found among %+v", sites)
	}
	out, err := Apply([]byte(src), *eq)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	want := `greeting = "héllo wörld"


def f(a, b):
    return a != b
`
	if string(out) != want {
		t.Fatalf("Apply with unicode-preceded offsets:\ngot:  %q\nwant: %q", out, want)
	}
}

func TestPythonSitesMissingInterpreter(t *testing.T) {
	t.Setenv("PATH", "")
	m := For("python")
	sites, err := m.Sites("", "fixture.py", []byte(pyOperatorsFixture))
	if err != nil {
		t.Fatalf("Sites: %v", err)
	}
	if sites != nil {
		t.Fatalf("Sites with no python3 on PATH = %+v, want nil", sites)
	}
}

// parsePythonFixture runs `python3 -c "import ast,sys; ast.parse(...)"` on
// src, failing t if it isn't valid Python. desc identifies the mutation
// under test in failure messages.
func parsePythonFixture(t *testing.T, src []byte, desc string) {
	t.Helper()
	cmd := exec.Command("python3", "-c", "import ast, sys; ast.parse(sys.stdin.buffer.read().decode('utf-8'))")
	cmd.Stdin = bytes.NewReader(src)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("parse failed for mutation %q:\n%s\nsrc:\n%s", desc, out, src)
	}
}
