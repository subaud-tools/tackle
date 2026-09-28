package evidence

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/schuettc/tackle/internal/cull/cases"
	"github.com/schuettc/tackle/internal/cull/mutate"
)

// pyKeep runs evidence on files and fails t unless target (a review-band
// test in test_calc.py) is kept as a sole guard.
func pyKeep(t *testing.T, files map[string]string, target string) string {
	t.Helper()
	needPython(t)
	root := newRepo(t, files)
	writeLast(t, root, []string{"test_calc.py"}, map[string]verdict{target: {"review", "review_band"}})
	o := settleOpts()
	var log bytes.Buffer
	o.Stderr = &log
	if _, err := Run(context.Background(), root, o); err != nil {
		t.Fatal(err)
	}
	_, by := loadLast(t, root)
	if tr := by[target]; tr.Verdict != "keep" || tr.Rule != RuleSoleGuard {
		t.Errorf("%s: %s/%s %v\n%s", target, tr.Verdict, tr.Rule, tr.Reasons, log.String())
	}
	return root
}

const postTests = "import pytest\n\nfrom calc import build\nfrom mk import kind\n\n\n" +
	"def test_other():\n    assert build(kind(), \"a\").text == \"a\"\n\n\n" +
	"def test_post():\n    with pytest.raises(ValueError, match=\"^empty$\"):\n        build(kind(), \"  \")\n"

// A reached method reaches its owner class even when it constructs it
// without naming it: cls(...) in a classmethod reached through a
// parameter, and dataclasses.replace(self, ...). Base.__post_init__ holds
// the target's unique catch.
func TestOwnerClassReached(t *testing.T) {
	const clsMod = `from dataclasses import dataclass


@dataclass
class Base:
    text: str

    def __post_init__(self):
        if self.text == "":
            raise ValueError("empty")


@dataclass
class Cause(Base):
    @classmethod
    def parse(cls, s):
        return cls(s.strip())


def build(k, s):
    return k.parse(s)
`
	const replaceMod = `from dataclasses import dataclass, replace


@dataclass
class Cause:
    text: str

    def __post_init__(self):
        if self.text == "":
            raise ValueError("empty")

    def with_text(self, s):
        return replace(self, text=s.strip())


def build(c, s):
    return c.with_text(s)
`
	t.Run("cls", func(t *testing.T) {
		pyKeep(t, map[string]string{
			"calc.py": clsMod, "test_calc.py": postTests,
			"mk.py": "import calc\n\n\ndef kind():\n    return calc.Cause\n",
		}, "test_post")
	})
	t.Run("replace", func(t *testing.T) {
		pyKeep(t, map[string]string{
			"calc.py": replaceMod, "test_calc.py": postTests,
			"mk.py": "import calc\n\n\ndef kind():\n    return calc.Cause(\"x\")\n",
		}, "test_post")
	})
}

// Python definitions inside top-level try/if/with blocks are definitions:
// _compat, defined in a try block, holds the target's unique catch.
func TestPythonBlockDefinitionsReached(t *testing.T) {
	const mod = `import sys
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from typing import List

try:
    import json

    def _compat(x):
        return x + "!"
except ImportError:
    def _compat(x):
        return x
finally:
    pass

if sys.version_info >= (3, 0):
    def _tail(x):
        return _compat(x)
else:
    def _tail(x):
        return x


def build(s):
    return _tail(s.strip())
`
	pyKeep(t, map[string]string{
		"calc.py": mod,
		"test_calc.py": "from calc import build\n\n\ndef test_other():\n    assert build(\" a \").startswith(\"a\")\n\n\n" +
			"def test_compat():\n    assert build(\"a\").endswith(\"!\")\n",
	}, "test_compat")
}

// Go: every init in the file is reached whenever anything in it is:
// the second init fills the table Make reads, and holds TestTable's
// unique catch.
func TestGoInitReached(t *testing.T) {
	needTools(t)
	const src = `package calc

var table map[string]int

func init() {
	_ = 0
}

func Make(k string) int {
	if k == "" {
		return -1
	}
	return table[k]
}

func init() {
	table = map[string]int{"a": 1 + 1}
}
`
	root := newRepo(t, map[string]string{
		"go.mod":  "module fixture\n\ngo 1.22\n",
		"calc.go": src,
		"calc_test.go": "package calc\n\nimport \"testing\"\n\n" +
			"func TestNeg(t *testing.T) {\n\tif Make(\"\") != -1 || Make(\"a\") == -1 {\n\t\tt.Fatal(\"neg\")\n\t}\n}\n\n" +
			"func TestTable(t *testing.T) {\n\tif Make(\"a\") != 2 {\n\t\tt.Fatal(\"table\")\n\t}\n}\n",
	})
	writeLast(t, root, []string{"calc_test.go"}, map[string]verdict{"TestTable": {"review", "review_band"}})
	o := settleOpts()
	var log bytes.Buffer
	o.Stderr = &log
	if _, err := Run(context.Background(), root, o); err != nil {
		t.Fatal(err)
	}
	_, by := loadLast(t, root)
	if tr := by["TestTable"]; tr.Verdict != "keep" || tr.Rule != RuleSoleGuard {
		t.Errorf("TestTable: %s/%s %v\n%s", tr.Verdict, tr.Rule, tr.Reasons, log.String())
	}
}

// Each reached definition samples only its own sites: a class's field
// sites are always planted however many sites its methods hold, and each
// method gets its own sample.
func TestClassOwnSample(t *testing.T) {
	src := []byte("F" + strings.Repeat(".", 99) + "C" + strings.Repeat(".", 899))
	for v := 0; v < 20; v++ {
		cls := fmt.Sprintf("C%d", v)
		defs := []mutate.Def{
			{Name: "f", Kind: "function", Start: 0, End: 100, Refs: []string{cls}},
			{Name: cls, Kind: "class", Start: 100, End: 1000},
			{Name: "m1", Owner: cls, Kind: "method", Start: 200, End: 600},
			{Name: "m2", Owner: cls, Kind: "method", Start: 600, End: 1000},
		}
		sites := []mutate.Site{{Start: 110, End: 111, Desc: "field"}, {Start: 120, End: 121, Desc: "field"}}
		for i := 0; i < 80; i++ {
			sites = append(sites, mutate.Site{Start: 200 + i*10, End: 201 + i*10, Desc: "method"})
		}
		got := targetSites("m.py", src, sites, []cases.Callee{{Symbol: "f", Source: string(src[:100])}}, defs, nolog)
		n := map[string]int{}
		for _, s := range got {
			n[s.Desc]++
		}
		if n["field"] != 2 || n["method"] != 2*perCallee {
			t.Errorf("%s: sampled %v", cls, n)
		}
	}
}

// A nested helper gets its own sample: its site is planted though its
// parent has many more than perCallee other sites.
func TestNestedHelperPlanted(t *testing.T) {
	var b strings.Builder
	b.WriteString("def build(x):\n    def _norm(s):\n        return s + \"!\"\n\n    t = 0\n")
	for i := 0; i < 30; i++ {
		fmt.Fprintf(&b, "    if x > %d:\n        t = t + %d\n", i+100, i)
	}
	b.WriteString("    return _norm(str(t))\n")
	pyKeep(t, map[string]string{
		"calc.py": b.String(),
		"test_calc.py": "from calc import build\n\n\ndef test_other():\n    assert build(0).startswith(\"0\")\n\n\n" +
			"def test_bang():\n    assert build(0).endswith(\"!\")\n",
	}, "test_bang")
}
