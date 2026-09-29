package runner

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
// package needs; others are ignored by encoding/json). FailedBuild is set
// on the package-level "fail" event when the package didn't build at all
// (no Test events follow); Output is scanned as a fallback for older Go
// versions that only print "[build failed]" rather than setting
// FailedBuild.
type goTestEvent struct {
	Action      string
	Test        string
	Output      string
	FailedBuild string
}

func (golangRunner) Run(ctx context.Context, dir string, testFiles []string, index map[string]string, timeout time.Duration) Outcome {
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
	return runGoGroups(ctx, dir, groups, index, timeout)
}

// runGoGroups runs, per package directory, the named top-level tests.
func runGoGroups(ctx context.Context, dir string, groups map[string]map[string]bool, index map[string]string, timeout time.Duration) Outcome {
	out := newOutcome()
	subs := make([]string, 0, len(groups))
	for s := range groups {
		subs = append(subs, s)
	}
	sort.Strings(subs)

	// One deadline for the whole Run, shared across every package group:
	// otherwise a hang in an early package would let each later package
	// consume its own full timeout too (N x timeout for N packages).
	// verify.RunProcess derives its own per-call context.WithTimeout(ctx,
	// timeout), which takes the earlier of that and deadlineCtx's already-
	// running deadline, so passing the unmodified timeout on every call
	// here is correct: deadlineCtx's absolute deadline is what actually
	// cuts a later call off early.
	deadlineCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for _, sub := range subs {
		if err := ctx.Err(); err != nil {
			out.Err = err.Error() // cancelled: start nothing more
			break
		}
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
		// -count=1: never replay a cached result. Evidence runs the same
		// unmutated tests twice to detect flakiness; a cache hit would
		// replay the first run and hide it.
		argv := []string{"go", "test", "-json", "-count=1", "-run", runRe, pkgSpec}

		res := runProcess(deadlineCtx, dir, argv, timeout)
		if res.StartErr != nil {
			out.Err = res.StartErr.Error()
			continue
		}
		if res.TimedOut {
			out.TimedOut = true
			// One shared deadline: stop entirely rather than letting the
			// next package group start a fresh timeout window of its own.
			break
		}
		sawTest, buildFailed := parseGoEvents(res.Stdout, index, &out)
		if buildFailed || (!sawTest && res.ExitCode != 0) {
			out.Invalid = true
			if out.InvalidReason == "" {
				out.InvalidReason = fmt.Sprintf("go: package %s failed to build", pkgSpec)
			}
		}
	}
	return out
}

// RunEach builds each package's test binary once (go test -c, which
// vets like go test) and runs it once per test file, with that file's
// top-level tests (see EachRunner).
func (golangRunner) RunEach(ctx context.Context, dir string, testFiles []string, index map[string]string, timeout time.Duration) Outcome {
	files := append([]string(nil), testFiles...)
	sort.Strings(files)
	tmp, err := os.MkdirTemp("", "cull-gorun-*")
	if err != nil {
		o := newOutcome()
		o.Err = err.Error()
		return o
	}
	defer func() { _ = os.RemoveAll(tmp) }()
	bins := map[string]string{} // package dir -> its test binary
	return runEach(ctx, files, timeout, func(ctx context.Context, f string) Outcome {
		out := newOutcome()
		names, err := goTopLevelTestNames(filepath.Join(dir, filepath.FromSlash(f)))
		if err != nil || len(names) == 0 {
			return out
		}
		sub := path.Dir(filepath.ToSlash(f))
		pkg := "./"
		if sub != "." {
			pkg = "./" + sub
		}
		bin, ok := bins[sub]
		if !ok {
			bin = filepath.Join(tmp, fmt.Sprintf("p%d.test", len(bins)))
			bins[sub] = bin
			res := runProcess(ctx, dir, []string{"go", "test", "-c", "-o", bin, pkg}, timeout)
			switch {
			case res.StartErr != nil:
				out.Err = res.StartErr.Error()
				return out
			case res.TimedOut:
				out.TimedOut = true
				return out
			case res.ExitCode != 0:
				out.Invalid = true
				out.InvalidReason = fmt.Sprintf("go: package %s failed to build", pkg)
				return out
			}
		}
		argv := []string{"go", "tool", "test2json", "-t", "-p", pkg, bin, "-test.v=test2json",
			"-test.count=1", "-test.run", "^(" + strings.Join(names, "|") + ")$"}
		res := runProcess(ctx, filepath.Join(dir, filepath.FromSlash(sub)), argv, timeout)
		switch {
		case res.StartErr != nil:
			out.Err = res.StartErr.Error()
			return out
		case res.TimedOut:
			out.TimedOut = true
			return out
		}
		if sawTest, _ := parseGoEvents(res.Stdout, index, &out); !sawTest && res.ExitCode != 0 {
			out.Invalid = true
			out.InvalidReason = fmt.Sprintf("go: %s's test binary exited %d before any test ran", pkg, res.ExitCode)
		}
		return out
	})
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
//
// It returns sawTest (whether any Test event was reported at all — false
// means the package never got as far as running tests) and buildFailed
// (a package-level build failure was reported: FailedBuild set, a
// "build-fail" action, or "[build failed]" in package output — the
// several shapes go test -json has used across Go versions).
func parseGoEvents(stdout []byte, index map[string]string, out *Outcome) (sawTest, buildFailed bool) {
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
		if ev.Action == "build-fail" || ev.FailedBuild != "" || strings.Contains(ev.Output, "[build failed]") {
			buildFailed = true
		}
		if ev.Test == "" {
			continue
		}
		sawTest = true
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
	return sawTest, buildFailed
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
