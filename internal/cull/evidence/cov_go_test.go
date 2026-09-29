package evidence

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const covGoSrc = "package calc\n\nvar limit = 0\n\nfunc Add(a, b int) int { return a + b }\n\nfunc Pos(n int) bool { return n > limit }\n\nfunc Unused(n int) int { return n * 2 }\n"

func covGoTest(log string) string {
	return fmt.Sprintf(`package calc

import (
	"os"
	"testing"
)

func logRun(name string) {
	f, _ := os.OpenFile(%q, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	f.WriteString(name + "\n")
	f.Close()
}

func TestAdd(t *testing.T) {
	logRun("TestAdd")
	if Add(2, 3) != 5 {
		t.Fatal("add")
	}
}
`, log)
}

// covGoPosTest is a second test file of the package.
const covGoPosTest = `package calc

import "testing"

func TestPosA(t *testing.T) {
	logRun("TestPosA")
	if !Pos(1) || Pos(0) {
		t.Fatal("pos")
	}
}

func TestPosB(t *testing.T) {
	logRun("TestPosB")
	t.Run("sub", func(t *testing.T) {
		if Pos(-1) {
			t.Fatal("pos b")
		}
	})
}
`

// Go: coverprofiles select the test files each mutant runs (whole: a
// subtest runs with its top-level test); a package-level var runs every
// test; an unused function's mutants are uncovered. Same matrix as the
// full run.
func TestGoCoverageSelectsTests(t *testing.T) {
	needTools(t)
	l := newRunLog(t, covGoSrc)
	root := newRepo(t, map[string]string{
		"go.mod":       "module fixture\n\ngo 1.22\n",
		"calc.go":      covGoSrc,
		"calc_test.go": covGoTest(l.path),
		"pos_test.go":  covGoPosTest,
	})
	mod := Module{Rel: "calc.go", Lang: "go", TestFiles: []string{"calc_test.go", "pos_test.go"},
		Tests: []string{"go:calc_test.go:TestAdd", "go:pos_test.go:TestPosA", "go:pos_test.go:TestPosB"}}
	var log bytes.Buffer
	o := opts(1)
	o.Stderr = &log
	m, err := RunMatrix(context.Background(), root, []Module{mod}, o)
	if err != nil {
		t.Fatal(err)
	}
	want := map[int][]string{3: {"TestAdd", "TestPosA", "TestPosB"}, 5: {"TestAdd"}, 7: {"TestPosA", "TestPosB"}}
	if got := l.byLine(t); !reflect.DeepEqual(got, want) {
		t.Errorf("tests run per mutated line:\n got %v\nwant %v\n%s", got, want, log.String())
	}
	for _, mu := range m.Mutants {
		if strings.HasPrefix(mu.Site.Desc, "line 9:") != mu.Uncovered {
			t.Errorf("%s: uncovered %v", mu.Site.Desc, mu.Uncovered)
		}
	}
	plantHook = nil
	o.NoCoverage = true
	full, err := RunMatrix(context.Background(), root, []Module{mod}, o)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(rows(m), rows(full)) {
		t.Errorf("selected:\n %v\nfull:\n %v", rows(m), rows(full))
	}
}
