package runner

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/schuettc/tackle/internal/cull/verify"
)

// tsRunner runs TypeScript tests with whichever per-test reporter its
// package.json test script implies: vitest's JSON reporter, or node:test's
// JUnit reporter. newTSRunner decides which (or errors) once, from the
// package.json owning the representative testFile For was given.
type tsRunner struct {
	kind    string // "vitest" or "node"
	nodeCmd string // node:test's script command text (kind == "node" only), split into argv at Run time
}

// errNoReporter is returned whenever a TypeScript project's test script
// isn't one runner knows how to get per-test results from: those tests
// are unmeasurable, not a bug in this package.
var errNoReporter = fmt.Errorf("no per-test reporter for this TypeScript runner")

func newTSRunner(root, testFile string) (Runner, error) {
	dirRel := path.Dir(filepath.ToSlash(testFile))
	pkgDir, ok := verify.NearestTSTestDir(root, dirRel)
	if !ok {
		return nil, errNoReporter
	}
	absDir := root
	if pkgDir != "." {
		absDir = filepath.Join(root, filepath.FromSlash(pkgDir))
	}
	script, ok := verify.PackageJSONTestScript(absDir)
	if !ok {
		return nil, errNoReporter
	}
	switch {
	case strings.Contains(script, "vitest"):
		return &tsRunner{kind: "vitest"}, nil
	case strings.Contains(script, "--test"):
		return &tsRunner{kind: "node", nodeCmd: script}, nil
	default:
		return nil, errNoReporter
	}
}

func (r *tsRunner) Run(ctx context.Context, dir string, testFiles []string, index map[string]string, timeout time.Duration) Outcome {
	out := newOutcome()
	if len(testFiles) == 0 {
		return out
	}
	if err := ctx.Err(); err != nil {
		out.Err = err.Error()
		return out
	}
	switch r.kind {
	case "vitest":
		return r.runVitest(ctx, dir, testFiles, index, timeout)
	case "node":
		return r.runNode(ctx, dir, testFiles, index, timeout)
	default:
		return out
	}
}

// vitestReport is the shape of vitest's --reporter=json --outputFile this
// package needs.
type vitestReport struct {
	TestResults []struct {
		AssertionResults []struct {
			AncestorTitles []string `json:"ancestorTitles"`
			Title          string   `json:"title"`
			Status         string   `json:"status"`
		} `json:"assertionResults"`
	} `json:"testResults"`
}

func (r *tsRunner) runVitest(ctx context.Context, dir string, testFiles []string, index map[string]string, timeout time.Duration) Outcome {
	out := newOutcome()

	tmp, err := os.CreateTemp("", "cull-vitest-*.json")
	if err != nil {
		out.Err = err.Error()
		return out
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmpPath) }()

	argv := append([]string{"npx", "vitest", "run", "--reporter=json", "--outputFile=" + tmpPath}, testFiles...)
	res := runProcess(ctx, dir, argv, timeout)
	if res.StartErr != nil {
		out.Err = res.StartErr.Error()
		return out
	}
	if res.TimedOut {
		out.TimedOut = true
		return out
	}

	data, err := os.ReadFile(tmpPath)
	if err != nil || len(data) == 0 {
		out.Invalid = true
		out.InvalidReason = "vitest exited without writing a report file"
		return out
	}
	var report vitestReport
	if err := json.Unmarshal(data, &report); err != nil {
		out.Invalid = true
		out.InvalidReason = "vitest report file was not valid JSON: " + err.Error()
		return out
	}
	for _, tr := range report.TestResults {
		for _, ar := range tr.AssertionResults {
			if ar.Status != "passed" && ar.Status != "failed" {
				continue
			}
			native := strings.Join(append(append([]string{}, ar.AncestorTitles...), ar.Title), " > ")
			attributeTS(index, native, ar.Status == "passed", &out)
		}
	}
	return out
}

// nodeJUnitSuites is the root <testsuites> element node --test-reporter=junit
// writes: top-level testcases (tests outside any describe) alongside
// nested testsuite elements (one per describe level).
type nodeJUnitSuites struct {
	Suites []nodeJUnitSuite `xml:"testsuite"`
	Cases  []nodeJUnitCase  `xml:"testcase"`
}

type nodeJUnitSuite struct {
	Name   string           `xml:"name,attr"`
	Suites []nodeJUnitSuite `xml:"testsuite"`
	Cases  []nodeJUnitCase  `xml:"testcase"`
}

type nodeJUnitCase struct {
	Name    string    `xml:"name,attr"`
	Failure *struct{} `xml:"failure"`
	Error   *struct{} `xml:"error"`
	Skipped *struct{} `xml:"skipped"`
}

func (r *tsRunner) runNode(ctx context.Context, dir string, testFiles []string, index map[string]string, timeout time.Duration) Outcome {
	out := newOutcome()

	tmp, err := os.CreateTemp("", "cull-nodetest-*.xml")
	if err != nil {
		out.Err = err.Error()
		return out
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmpPath) }()

	fields := strings.Fields(r.nodeCmd)
	argv := append(append([]string{}, fields...), "--test-reporter=junit", "--test-reporter-destination="+tmpPath)
	argv = append(argv, testFiles...)

	res := runProcess(ctx, dir, argv, timeout)
	if res.StartErr != nil {
		out.Err = res.StartErr.Error()
		return out
	}
	if res.TimedOut {
		out.TimedOut = true
		return out
	}
	readNodeJUnit(tmpPath, testFiles, index, &out)
	return out
}

// readNodeJUnit adds the results of node --test's junit report at path
// to out (Invalid when it's missing or unreadable, or a file failed to
// load).
func readNodeJUnit(path string, testFiles []string, index map[string]string, out *Outcome) {
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		out.Invalid = true
		out.InvalidReason = "node --test exited without writing a report file"
		return
	}
	var suites nodeJUnitSuites
	if err := xml.Unmarshal(data, &suites); err != nil {
		out.Invalid = true
		out.InvalidReason = "node --test report file was not valid XML: " + err.Error()
		return
	}
	for _, c := range suites.Cases {
		if c.Skipped != nil {
			continue
		}
		failed := c.Failure != nil || c.Error != nil
		if f, ok := fileCase(c.Name, testFiles); ok && failed {
			// node reports a file that failed to load (or crashed
			// outside any test) as one failed testcase named after it.
			out.Invalid = true
			if out.InvalidReason == "" {
				out.InvalidReason = "node --test: " + f + " failed outside its tests (load error)"
			}
			continue
		}
		attributeTS(index, c.Name, !failed, out)
	}
	walkNodeSuites(suites.Suites, nil, index, out)
}

// fileCase reports which of testFiles a top-level testcase name is (node
// names a file-level result by the path it was given).
func fileCase(name string, testFiles []string) (string, bool) {
	for _, f := range testFiles {
		if name == f || strings.HasSuffix(filepath.ToSlash(name), "/"+filepath.ToSlash(f)) {
			return f, true
		}
	}
	return "", false
}

func walkNodeSuites(suites []nodeJUnitSuite, ancestors []string, index map[string]string, out *Outcome) {
	for _, s := range suites {
		next := append(append([]string{}, ancestors...), s.Name)
		for _, c := range s.Cases {
			if c.Skipped != nil {
				continue
			}
			native := strings.Join(append(append([]string{}, next...), c.Name), " > ")
			attributeTS(index, native, c.Failure == nil && c.Error == nil, out)
		}
		walkNodeSuites(s.Suites, next, index, out)
	}
}

// attributeTS maps a " > "-joined ancestor/title chain to a cull id via
// index, falling back to the nearest ancestor chain that index does map
// when the leaf itself doesn't (containers are never units themselves,
// but this mirrors the Go/Python leaf-rule fallback for consistency).
func attributeTS(index map[string]string, native string, passed bool, out *Outcome) {
	for {
		if id, ok := index[native]; ok {
			markResult(out, id, passed)
			return
		}
		i := strings.LastIndex(native, " > ")
		if i < 0 {
			out.Unmapped++
			return
		}
		native = native[:i]
	}
}
