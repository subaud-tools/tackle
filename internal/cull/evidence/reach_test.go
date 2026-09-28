package evidence

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/schuettc/tackle/internal/cull/mutate"
)

// mutantsBefore fails t if evidence.json holds a mutant of rel at or past
// the first byte of marker in src: code the targets never reach.
func mutantsBefore(t *testing.T, root, rel, src, marker string) {
	t.Helper()
	cut := strings.Index(src, marker)
	if cut < 0 {
		t.Fatalf("marker %q not in source", marker)
	}
	data, err := os.ReadFile(filepath.Join(root, ".cull", "evidence.json"))
	if err != nil {
		t.Fatal(err)
	}
	var f evidenceFile
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatal(err)
	}
	for _, mu := range f.Mutants {
		if mu.Module == rel && mu.Site.Start >= cut {
			t.Errorf("mutant planted in unreached code: %+v", mu.Site)
		}
	}
}

// The nfl-dk false cut: a target calls build, which calls the module
// helper _opt_str and constructs a dataclass whose __post_init__
// validates. The target's only unique catch is in the helper (or in
// __post_init__): both are reached from the callee, so it is kept; the
// unreferenced function is never mutated.
func TestPythonReachedDefinitionsKeep(t *testing.T) {
	needPython(t)
	const helperMod = `from dataclasses import dataclass


def _opt_str(x):
    return str(x + 1)


@dataclass
class Cause:
    text: str

    def __post_init__(self):
        if self.text == "":
            raise ValueError("empty")


def build(x):
    if x < 0:
        return None
    return Cause(_opt_str(x))


def unrelated(x):
    return x * 3
`
	const postMod = `from dataclasses import dataclass


@dataclass
class Cause:
    text: str

    def __post_init__(self):
        if self.text == "":
            raise ValueError("empty")


def build(s):
    return Cause(s.strip())


def unrelated(x):
    return x * 3
`
	for _, c := range []struct {
		name, mod, tests, target string
	}{
		{"catch in helper", helperMod,
			"from calc import build\n\n\ndef test_neg():\n    assert build(-1) is None\n    assert build(1) is not None\n\n\n" +
				"def test_helper():\n    assert build(1).text == \"2\"\n",
			"test_helper"},
		{"catch in __post_init__", postMod,
			"import pytest\n\nfrom calc import build\n\n\ndef test_other():\n    assert build(\"a\").text == \"a\"\n\n\n" +
				"def test_post():\n    with pytest.raises(ValueError, match=\"^empty$\"):\n        build(\"  \")\n",
			"test_post"},
	} {
		t.Run(c.name, func(t *testing.T) {
			root := newRepo(t, map[string]string{"calc.py": c.mod, "test_calc.py": c.tests})
			writeLast(t, root, []string{"test_calc.py"}, map[string]verdict{c.target: {"review", "review_band"}})
			o := settleOpts()
			var log bytes.Buffer
			o.Stderr = &log
			if _, err := Run(context.Background(), root, o); err != nil {
				t.Fatal(err)
			}
			_, by := loadLast(t, root)
			if tr := by[c.target]; tr.Verdict != "keep" || tr.Rule != RuleSoleGuard {
				t.Errorf("%s: %s/%s %v\n%s", c.target, tr.Verdict, tr.Rule, tr.Reasons, log.String())
			}
			mutantsBefore(t, root, "calc.py", c.mod, "def unrelated")
		})
	}
}

// Go: the target calls Make, which calls helper and formats a
// constructed Box (its String method runs only through fmt). Each
// target's unique catch is in one of them; the unreferenced type and
// function are never mutated.
func TestGoReachedDefinitionsKeep(t *testing.T) {
	needTools(t)
	const src = `package calc

import (
	"fmt"
	"strconv"
)

type Box struct{ v int }

func (b Box) String() string { return strconv.Itoa(b.v) + "#box" }

func helper(x int) int { return x + 1 }

func Make(x int) string {
	if x < 0 {
		return "neg"
	}
	return fmt.Sprint(Box{v: helper(x)})
}

type Other struct{}

func (Other) Twice(x int) int { return x * 2 }

func Unrelated(x int) int { return x - 7 }
`
	root := newRepo(t, map[string]string{
		"go.mod":  "module fixture\n\ngo 1.22\n",
		"calc.go": src,
		"calc_test.go": "package calc\n\nimport (\n\t\"strings\"\n\t\"testing\"\n)\n\n" +
			"func TestNeg(t *testing.T) {\n\tif Make(-1) != \"neg\" || Make(1) == \"neg\" {\n\t\tt.Fatal(\"neg\")\n\t}\n}\n\n" +
			"func TestHelper(t *testing.T) {\n\tif !strings.HasPrefix(Make(1), \"2#\") {\n\t\tt.Fatal(\"helper\")\n\t}\n}\n\n" +
			"func TestString(t *testing.T) {\n\tif !strings.HasSuffix(Make(1), \"#box\") {\n\t\tt.Fatal(\"string\")\n\t}\n}\n",
	})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"calc_test.go"}, map[string]verdict{"TestHelper": band, "TestString": band})
	o := settleOpts()
	var log bytes.Buffer
	o.Stderr = &log
	if _, err := Run(context.Background(), root, o); err != nil {
		t.Fatal(err)
	}
	_, by := loadLast(t, root)
	for _, id := range []string{"TestHelper", "TestString"} {
		if tr := by[id]; tr.Verdict != "keep" || tr.Rule != RuleSoleGuard {
			t.Errorf("%s: %s/%s %v", id, tr.Verdict, tr.Rule, tr.Reasons)
		}
	}
	mutantsBefore(t, root, "calc.go", src, "type Other")
	if t.Failed() {
		t.Log(log.String())
	}
}

// reachNames is reach from the defs named seeds, as sorted qualified names.
func reachNames(defs []mutate.Def, seeds ...string) []string {
	var idx []int
	for i, d := range defs {
		if slices.Contains(seeds, qualName(d)) {
			idx = append(idx, i)
		}
	}
	var out []string
	for _, i := range reach(defs, idx) {
		out = append(out, qualName(defs[i]))
	}
	sort.Strings(out)
	return out
}

// Reach follows calls, constructed classes (all their methods), methods
// by name and their owner classes, transitively, and nothing else.
func TestReachDoesNotLeak(t *testing.T) {
	defs := []mutate.Def{
		{Name: "f", Kind: "function", Refs: []string{"x", "g", "C"}},
		{Name: "g", Kind: "function", Refs: []string{"self", "run"}},
		{Name: "C", Kind: "class", Refs: []string{"dataclass"}},
		{Name: "__post_init__", Owner: "C", Kind: "method", Refs: []string{"self", "check"}},
		{Name: "check", Owner: "C", Kind: "method"},
		{Name: "run", Owner: "R", Kind: "method", Refs: []string{"h"}},
		{Name: "h", Kind: "function"},
		{Name: "R", Kind: "class"},
		{Name: "other", Owner: "R", Kind: "method"},
		{Name: "unrelated", Kind: "function", Refs: []string{"f", "D"}},
		{Name: "D", Kind: "class"},
		{Name: "__init__", Owner: "D", Kind: "method"},
	}
	got := reachNames(defs, "f")
	want := []string{"C", "C.__post_init__", "C.check", "R", "R.other", "R.run", "f", "g", "h"}
	if !slices.Equal(got, want) {
		t.Errorf("reach(f) = %v, want %v", got, want)
	}
}

// TypeScript definitions from the helper's --defs mode: make reaches its
// helper and every member of the class it constructs, not the
// unreferenced class or function.
func TestTSReach(t *testing.T) {
	if os.Getenv("CULL_TS") == "" {
		t.Skip("CULL_TS not set")
	}
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not on PATH")
	}
	const src = `class Box {
  constructor(public v: number) {
    this.check();
  }
  check(): void {
    if (this.v < 0) throw new Error("neg");
  }
  get label(): string {
    return "#" + this.v;
  }
}

function helper(x: number): number {
  return x + 1;
}

export function make(x: number): string {
  return new Box(helper(x)).label;
}

class Other {
  twice(x: number): number {
    return x * 2;
  }
}

export function unrelated(x: number): number {
  return new Other().twice(x) - 7;
}
`
	defs, err := mutate.Defs("typescript", t.TempDir(), "calc.ts", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	got := reachNames(defs, "make")
	want := []string{"Box", "Box.check", "Box.constructor", "Box.label", "helper", "make"}
	if !slices.Equal(got, want) {
		t.Errorf("reach(make) = %v, want %v", got, want)
	}
	if d := defs[slices.IndexFunc(defs, func(d mutate.Def) bool { return d.Name == "helper" })]; !strings.HasPrefix(src[d.Start:d.End], "function helper") {
		t.Errorf("helper span = %q", src[d.Start:d.End])
	}
}

// Go init funcs are reached whenever anything in the file is.
func TestReachGoInit(t *testing.T) {
	defs := []mutate.Def{
		{Name: "f", Kind: "function"},
		{Name: "init", Kind: "init", Refs: []string{"g"}},
		{Name: "g", Kind: "function"},
		{Name: "init", Kind: "init"},
		{Name: "h", Kind: "function"},
	}
	if got := reach(defs, []int{0}); !slices.Equal(got, []int{0, 1, 2, 3}) {
		t.Errorf("reach(f) = %v", got)
	}
	if got := reach(defs, nil); len(got) != 0 {
		t.Errorf("reach() = %v", got)
	}
}

// TypeScript: #private members are refs, and namespace members are
// definitions (reached by name, and through the namespace).
func TestTSReachPrivateAndNamespace(t *testing.T) {
	if os.Getenv("CULL_TS") == "" {
		t.Skip("CULL_TS not set")
	}
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not on PATH")
	}
	const src = `class Box {
  #scale(x: number): number {
    return x * 2;
  }
  run(x: number): number {
    return this.#scale(x);
  }
}

namespace Util {
  export function inc(x: number): number {
    return x + 1;
  }
  export namespace Deep {
    export function dec(x: number): number {
      return x - 1;
    }
  }
}

export function make(b: Box, x: number): number {
  return b.run(Util.inc(x));
}

export function other(x: number): number {
  return Util.Deep.dec(x);
}
`
	defs, err := mutate.Defs("typescript", t.TempDir(), "calc.ts", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range defs {
		t.Logf("%s %s %v", qualName(d), d.Kind, d.Refs)
	}
	for _, name := range []string{"Box.#scale", "Util", "inc", "Deep", "dec"} {
		if !slices.ContainsFunc(defs, func(d mutate.Def) bool { return qualName(d) == name }) {
			t.Errorf("no def %s", name)
		}
	}
	if i := slices.IndexFunc(defs, func(d mutate.Def) bool { return qualName(d) == "Box.run" }); i < 0 || !slices.Contains(defs[i].Refs, "#scale") {
		t.Errorf("Box.run does not ref #scale")
	}
	got := reachNames(defs, "make")
	for _, name := range []string{"Box.#scale", "inc"} {
		if !slices.Contains(got, name) {
			t.Errorf("reach(make) = %v, missing %s", got, name)
		}
	}
	if got := reachNames(defs, "other"); !slices.Contains(got, "dec") {
		t.Errorf("reach(other) = %v, missing dec", got)
	}
}

// A reached def reaches the defs nested in its span, named or not.
func TestReachNested(t *testing.T) {
	defs := []mutate.Def{
		{Name: "f", Kind: "function", Start: 0, End: 100},
		{Name: "g", Kind: "function", Start: 10, End: 50},
		{Name: "h", Kind: "function", Start: 100, End: 150},
	}
	if got := reach(defs, []int{0}); !slices.Equal(got, []int{0, 1}) {
		t.Errorf("reach(f) = %v", got)
	}
}
