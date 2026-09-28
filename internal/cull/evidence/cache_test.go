package evidence

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func cachePkgTest(log, pkg, name, body string) string {
	return fmt.Sprintf("package %s\n\nimport (\n\t\"os\"\n\t\"testing\"\n)\n\nfunc %s(t *testing.T) {\n\tf, _ := os.OpenFile(%q, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)\n\tf.WriteString(%q)\n\tf.Close()\n\t%s\n}\n",
		pkg, name, log, name+"\n", body)
}

// A rerun on an unchanged project reuses every module's cached result
// (no test runs); editing one module's test file redoes only that
// module; NoCache forces a fresh run. Cache files are 0600.
func TestCacheReuse(t *testing.T) {
	needTools(t)
	log := filepath.Join(t.TempDir(), "runs")
	root := newRepo(t, map[string]string{
		"go.mod":      "module fixture\n\ngo 1.22\n",
		"a/a.go":      "package a\n\nfunc Add(x, y int) int { return x + y }\n",
		"a/a_test.go": cachePkgTest(log, "a", "TestAdd", "if Add(2, 3) != 5 {\n\t\tt.Fatal(\"add\")\n\t}"),
		"b/b.go":      "package b\n\nfunc Pos(n int) bool { return n > 0 }\n",
		"b/b_test.go": cachePkgTest(log, "b", "TestPos", "if !Pos(1) || Pos(0) {\n\t\tt.Fatal(\"pos\")\n\t}"),
	})
	mods := []Module{
		{Rel: "a/a.go", Lang: "go", TestFiles: []string{"a/a_test.go"}, Tests: []string{"go:a/a_test.go:TestAdd"}},
		{Rel: "b/b.go", Lang: "go", TestFiles: []string{"b/b_test.go"}, Tests: []string{"go:b/b_test.go:TestPos"}},
	}
	ran := func() []string {
		b, _ := os.ReadFile(log)
		_ = os.Remove(log)
		set := map[string]bool{}
		for _, l := range strings.Fields(string(b)) {
			set[l] = true
		}
		return sortedKeys(set)
	}
	run := func(o Options) (Matrix, string) {
		var buf bytes.Buffer
		o.Stderr = &buf
		m, err := RunMatrix(context.Background(), root, mods, o)
		if err != nil {
			t.Fatal(err)
		}
		return m, buf.String()
	}
	cached := opts(2)
	cached.NoCache = false
	first, _ := run(cached)
	if got := ran(); !reflect.DeepEqual(got, []string{"TestAdd", "TestPos"}) {
		t.Fatalf("first run ran %v", got)
	}
	second, stderr := run(cached)
	if got := ran(); len(got) != 0 {
		t.Errorf("unchanged rerun ran %v\n%s", got, stderr)
	}
	if !reflect.DeepEqual(first, second) {
		t.Errorf("cached matrix differs:\n%+v\n%+v", first, second)
	}
	ents, _ := filepath.Glob(filepath.Join(sandboxBase(mustAbs(t, root)), "cache", "*"))
	if len(ents) != 2 {
		t.Errorf("cache entries %v", ents)
	}
	for _, e := range ents {
		if perm(t, e) != 0o600 {
			t.Errorf("%s: %v", e, perm(t, e))
		}
	}

	write(t, root, "b/b_test.go", cachePkgTest(log, "b", "TestPos", "if !Pos(1) || Pos(0) {\n\t\tt.Fatal(\"pos!\")\n\t}"))
	third, _ := run(cached)
	if got := ran(); !reflect.DeepEqual(got, []string{"TestPos"}) {
		t.Errorf("after editing b's test ran %v", got)
	}
	if !reflect.DeepEqual(first, third) {
		t.Errorf("matrix differs after edit:\n%+v\n%+v", first, third)
	}

	o := opts(2)
	o.NoCache = true
	run(o)
	if got := ran(); !reflect.DeepEqual(got, []string{"TestAdd", "TestPos"}) {
		t.Errorf("NoCache ran %v", got)
	}
}
