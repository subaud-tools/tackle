package evidence

import (
	"bytes"
	"context"
	"testing"
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
