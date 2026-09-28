package runner

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// golangRunner runs `go test -json` per package (grouped from testFiles'
// directories), filtered to the top-level Test funcs those files declare.
type golangRunner struct{}

// goTestFuncRe finds top-level Go test function declarations: good enough
// to build a `-run` filter without a full AST parse (extraction already
// did that; this only needs names).
var goTestFuncRe = regexp.MustCompile(`(?m)^func\s+(Test[A-Za-z0-9_]*)\s*\(`)

// goTestEvent is one line of `go test -json` output (the fields this
// package needs; others are ignored by encoding/json).
type goTestEvent struct {
	Action string
	Test   string
}

func (golangRunner) Run(ctx context.Context, dir string, testFiles []string, index map[string]string, timeout time.Duration) Outcome {
	out := newOutcome()

	groups := map[string]map[string]bool{} // package dir, relative to dir ("." = dir itself) -> top-level test names
	for _, f := range testFiles {
		slash := filepath.ToSlash(f)
		sub := path.Dir(slash)
		names, err := goTopLevelTestNames(filepath.Join(dir, filepath.FromSlash(f)))
		if err != nil {
			continue
		}
		if groups[sub] == nil {
			groups[sub] = map[string]bool{}
		}
		for _, n := range names {
			groups[sub][n] = true
		}
	}

	subs := make([]string, 0, len(groups))
	for s := range groups {
		subs = append(subs, s)
	}
	sort.Strings(subs)

	for _, sub := range subs {
		names := make([]string, 0, len(groups[sub]))
		for n := range groups[sub] {
			names = append(names, n)
		}
		if len(names) == 0 {
			continue
		}
		sort.Strings(names)

		pkgSpec := "./"
		if sub != "." {
			pkgSpec = "./" + sub
		}
		runRe := "^(" + strings.Join(names, "|") + ")$"
		argv := []string{"go", "test", "-json", "-run", runRe, pkgSpec}

		res := runProcess(ctx, dir, argv, timeout)
		if res.StartErr != nil {
			out.Err = res.StartErr.Error()
			continue
		}
		if res.TimedOut {
			out.TimedOut = true
			continue
		}
		parseGoEvents(res.Stdout, index, &out)
	}
	return out
}

// goTopLevelTestNames returns the top-level `func TestXxx(t *testing.T)`
// names declared in file (TestMain excluded: it's not a test).
func goTopLevelTestNames(file string) ([]string, error) {
	src, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	matches := goTestFuncRe.FindAllSubmatch(src, -1)
	names := make([]string, 0, len(matches))
	for _, m := range matches {
		name := string(m[1])
		if name == "TestMain" {
			continue
		}
		names = append(names, name)
	}
	return names, nil
}

// parseGoEvents scans `go test -json` stdout and attributes each pass/fail
// Test event to index. go test -json reports an event for every level of
// a subtest tree, including pure containers (a top-level Test whose body
// is only t.Run calls, like TestParse, has no id of its own in
// extraction); a container's own event is skipped rather than counted in
// Unmapped when it doesn't map directly, since its real leaves are
// reported (and attributed) separately.
func parseGoEvents(stdout []byte, index map[string]string, out *Outcome) {
	type result struct {
		test   string
		passed bool
	}
	var results []result
	names := map[string]bool{}

	scanner := bufio.NewScanner(bytes.NewReader(stdout))
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var ev goTestEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			continue
		}
		if ev.Test == "" {
			continue
		}
		names[ev.Test] = true
		switch ev.Action {
		case "pass":
			results = append(results, result{ev.Test, true})
		case "fail":
			results = append(results, result{ev.Test, false})
		}
	}

	isContainer := func(name string) bool {
		prefix := name + "/"
		for n := range names {
			if strings.HasPrefix(n, prefix) {
				return true
			}
		}
		return false
	}
	for _, r := range results {
		attributeGo(index, r.test, r.passed, isContainer, out)
	}
}

// attributeGo maps a go test -json Test path (e.g. "TestTable/one") to a
// cull id via index, falling back to the nearest ancestor ("TestTable")
// that index does map when the leaf itself doesn't: the leaf rule for a
// loop-generated subtest, which extraction treats as one unit with its
// parent. A name that still doesn't map once fully stripped to its
// top-level test is counted in Unmapped, unless isContainer says it's a
// pure container (its own event is then simply dropped: its real leaves,
// reported separately, are what matter).
func attributeGo(index map[string]string, native string, passed bool, isContainer func(string) bool, out *Outcome) {
	for {
		if id, ok := index[native]; ok {
			markResult(out, id, passed)
			return
		}
		i := strings.LastIndex(native, "/")
		if i < 0 {
			if !isContainer(native) {
				out.Unmapped++
			}
			return
		}
		native = native[:i]
	}
}
