package evidence

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCanarySource(t *testing.T) {
	cases := []struct{ lang, src, want string }{
		{"python", "\"\"\"Doc.\n\nMore.\n\"\"\"\n# c\nfrom __future__ import (\n    annotations,\n)\nfrom __future__ import division\n\nimport os\n",
			"\"\"\"Doc.\n\nMore.\n\"\"\"\n# c\nfrom __future__ import (\n    annotations,\n)\nfrom __future__ import division\n\nraise RuntimeError(\"cull canary\")\nimport os\n"},
		{"python", "#!/usr/bin/env python\nr'''one line'''\nX = 1", "#!/usr/bin/env python\nr'''one line'''\nraise RuntimeError(\"cull canary\")\nX = 1"},
		{"python", "X = 1\n", "raise RuntimeError(\"cull canary\")\nX = 1\n"},
		{"python", "\"\"\"only doc\"\"\"", "\"\"\"only doc\"\"\"\nraise RuntimeError(\"cull canary\")\n"},
		{"go", "package a\n\nfunc F() {}", "package a\n\nfunc F() {}\nvar _ = cullCanaryUndefined\n"},
		{"typescript", "// hi\n'use strict';\nimport {\n  a,\n  b,\n} from './x'\nimport type { T } from \"./t\";\n/* block\n */\nimport './side';\nexport const c = 1;\n",
			"// hi\n'use strict';\nimport {\n  a,\n  b,\n} from './x'\nimport type { T } from \"./t\";\n/* block\n */\nimport './side';\nthrow new Error(\"cull canary\");\nexport const c = 1;\n"},
		{"typescript", "const a = import('./x');\n", "throw new Error(\"cull canary\");\nconst a = import('./x');\n"},
	}
	for _, c := range cases {
		got, ok := canarySource(c.lang, []byte(c.src))
		if !ok || string(got) != c.want {
			t.Errorf("%s %q:\n got %q (%v)\nwant %q", c.lang, c.src, got, ok, c.want)
		}
	}
	if _, ok := canarySource("rust", []byte("fn main() {}")); ok {
		t.Error("canary for an unknown language")
	}
}

// A Go canary that breaks the build is an Invalid run: the module is
// loaded, and its mutants run as usual.
func TestGoCanaryCompileFailureLoaded(t *testing.T) {
	needTools(t)
	root := calcRepo(t)
	var log bytes.Buffer
	o := Options{Workers: 1, MaxMutants: 30, RunTimeout: 60 * time.Second, Stderr: &log}
	m, err := RunMatrix(context.Background(), root, []Module{calcModule()}, o)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(log.String(), "canary calc.go: loaded (invalid:") {
		t.Errorf("no loaded canary in log:\n%s", log.String())
	}
	if len(m.Shadowed) != 0 || len(m.Mutants) == 0 || len(m.Baseline) != 3 {
		t.Errorf("shadowed %v, %d mutants, baseline %v", m.Shadowed, len(m.Mutants), m.Baseline)
	}
	if b, _ := os.ReadFile(filepath.Join(sandboxBase(root), "w0", "calc.go")); string(b) != calcSrc {
		t.Errorf("canary not restored: %q", b)
	}
}

const pyCalc = "\"\"\"Calc.\"\"\"\nfrom __future__ import annotations\n\n\ndef add(a, b):\n    return a + b\n"

// A module whose tests import another copy (here a conftest path hack
// pointing outside the project) is unmeasured, never "catches nothing".
func TestShadowedPythonUnmeasured(t *testing.T) {
	needPython(t)
	outside := t.TempDir()
	write(t, outside, "calc.py", pyCalc)
	root := newRepo(t, map[string]string{
		"calc.py":            pyCalc,
		"conftest.py":        "import sys\nsys.path.insert(0, " + pyQuote(outside) + ")\n",
		"tests/test_calc.py": "from calc import add\n\n\ndef test_add():\n    assert add(2, 3) == 5\n",
	})
	writeLast(t, root, []string{"tests/test_calc.py"}, map[string]verdict{"test_add": {"review", "review_band"}})
	res, err := Run(context.Background(), root, settleOpts())
	if err != nil {
		t.Fatal(err)
	}
	id := "py:tests/test_calc.py:test_add"
	if len(res.Settled) != 0 {
		t.Errorf("Settled = %v", res.Settled)
	}
	if res.Reasons[id] != "tests don't load this copy of the module (an installed package or path setting shadows it)" {
		t.Errorf("reason = %q (modules %d, mutants %d)", res.Reasons[id], res.Modules, res.Mutants)
	}
	if res.Mutants != 0 {
		t.Errorf("ran %d mutants of a shadowed module", res.Mutants)
	}
	_, by := loadLast(t, root)
	if by["test_add"].Verdict != "review" || by["test_add"].Rule != RuleUnmeasured {
		t.Errorf("test_add: %s/%s", by["test_add"].Verdict, by["test_add"].Rule)
	}
}

func pyQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "\\'") + "'" }

// A src/ layout package is imported from the sandbox, even with another
// copy on PYTHONPATH (as an editable install of another checkout would be).
func TestSrcLayoutSandboxFirst(t *testing.T) {
	needPython(t)
	for _, name := range []string{"plain", "outside copy on PYTHONPATH"} {
		t.Run(name, func(t *testing.T) {
			if name != "plain" {
				outside := t.TempDir()
				write(t, outside, "pkg/__init__.py", "")
				write(t, outside, "pkg/calc.py", pyCalc)
				t.Setenv("PYTHONPATH", outside)
			}
			root := newRepo(t, map[string]string{
				"src/pkg/__init__.py": "",
				"src/pkg/calc.py":     pyCalc,
				"tests/test_calc.py":  "from pkg.calc import add\n\n\ndef test_add():\n    assert add(2, 3) == 5\n",
			})
			var log bytes.Buffer
			mod := Module{Rel: "src/pkg/calc.py", Lang: "python", TestFiles: []string{"tests/test_calc.py"},
				Tests: []string{"py:tests/test_calc.py:test_add"}}
			m, err := RunMatrix(context.Background(), root, []Module{mod},
				Options{Workers: 1, MaxMutants: 30, RunTimeout: 60 * time.Second, Stderr: &log})
			if err != nil {
				t.Fatal(err)
			}
			killed := 0
			for _, mu := range m.Mutants {
				if len(mu.KilledBy) > 0 {
					killed++
				}
			}
			if len(m.Shadowed) != 0 || !m.Baseline["py:tests/test_calc.py:test_add"] || killed == 0 {
				t.Errorf("shadowed %v, baseline %v, %d/%d killed\n%s", m.Shadowed, m.Baseline, killed, len(m.Mutants), log.String())
			}
		})
	}
}

// Dependency dirs that are symlinks in the project (ignored or not, with
// absolute or relative targets) are linked to the same target.
func TestSymlinkedDepDirsLinked(t *testing.T) {
	needTools(t)
	venv, dotVenv := t.TempDir(), t.TempDir()
	write(t, venv, "marker", "v\n")
	root := newRepo(t, map[string]string{".gitignore": "venv\n.venv/\n", "app.py": "X = 1\n"})
	nm := filepath.Join(filepath.Dir(root), "nm-shared")
	write(t, nm, "pkg/i.js", "module.exports = 1\n")
	for link, target := range map[string]string{"venv": venv, ".venv": dotVenv, "node_modules": "../nm-shared"} {
		if err := os.Symlink(target, filepath.Join(root, link)); err != nil {
			t.Fatal(err)
		}
	}
	sbs, err := Prepare(root, 1)
	if err != nil {
		t.Fatal(err)
	}
	for link, want := range map[string]string{"venv": venv, ".venv": dotVenv, "node_modules": nm} {
		got, err := os.Readlink(filepath.Join(sbs[0].Dir, link))
		if err != nil || got != want {
			t.Errorf("%s -> %q (%v), want %q", link, got, err, want)
		}
	}
	if _, err := os.Stat(filepath.Join(sbs[0].Dir, "venv", "marker")); err != nil {
		t.Error(err)
	}
	// Refresh keeps them.
	if sbs, err = Prepare(root, 1); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.Readlink(filepath.Join(sbs[0].Dir, "venv")); got != venv {
		t.Errorf("after refresh venv -> %q", got)
	}
}
