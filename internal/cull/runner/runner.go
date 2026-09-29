// Package runner runs a set of test files once (for a language: Go,
// Python, or TypeScript) and reports per-test pass/fail results, mapped
// through a caller-supplied index (runner-native test name -> cull
// TestCase id) to cull ids. It is evidence's execution primitive: given
// files with a planted bug and files without, evidence runs the tests
// before/after and compares Outcomes.
package runner

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/schuettc/tackle/internal/cull/verify"
)

// Outcome is one Run call's per-test results, keyed by cull TestCase id
// (the values index maps native names to).
type Outcome struct {
	Passed   map[string]bool
	Failed   map[string]bool
	Unmapped int  // reported native results that matched nothing, even after the leaf-rule parent fallback
	TimedOut bool // the whole run was killed for exceeding its timeout; Passed/Failed reflect only what was reported before then
	Err      string

	// Invalid is true when the run did not execute the requested tests as
	// a whole, so Passed/Failed/Unmapped (whatever exists) are not a
	// reliable signal: a Go package failed to build, pytest hit a
	// collection/import error, or a TypeScript runner exited without a
	// parseable report file. Callers must treat an Invalid Outcome as
	// unusable, not as "nothing selected".
	Invalid       bool
	InvalidReason string
}

// Runner runs testFiles (paths relative to dir) once in dir and reports
// their per-test results. index maps a runner-native test name (its exact
// shape is Runner-specific: e.g. Go's dotted subtest path, Python's
// Class::func, TypeScript's " > "-joined describe/title chain) to the
// cull TestCase id it corresponds to. A reported result whose native name
// isn't in index is attributed to its nearest ancestor that is (the
// leaf-rule fallback: a loop-generated Go subtest or a parametrized
// pytest case collapses to its parent unit), or else counted in Unmapped.
// A cancelled ctx starts no process (Err is set) and kills running ones.
type Runner interface {
	Run(ctx context.Context, dir string, testFiles []string, index map[string]string, timeout time.Duration) Outcome
}

// For returns the Runner for lang ("go", "python", or "typescript"),
// rooted at root. For TypeScript, it inspects the package.json owning
// testFile (root-relative, '/'-separated) to decide vitest vs node:test
// vs "no per-test reporter" (an error, since the tests are then
// unmeasurable per-test). Evidence always runs its own detected test
// command (go test / pytest / the reporter-augmented TS script), never a
// project's custom test_command override, so For takes no testCommand
// parameter.
func For(root, lang, testFile string) (Runner, error) {
	switch lang {
	case "go":
		return golangRunner{}, nil
	case "python":
		return pythonRunner{}, nil
	case "typescript":
		return newTSRunner(root, testFile)
	default:
		return nil, fmt.Errorf("runner: unknown language %q", lang)
	}
}

func newOutcome() Outcome {
	return Outcome{Passed: map[string]bool{}, Failed: map[string]bool{}}
}

// markResult records one native test result for id. A failure always
// wins: an id that ends up marked both failed and passed (e.g. a
// table-driven Go test whose loop-generated subtests all collapse to one
// parent id via the leaf-rule fallback, and only some of them fail) is
// reported Failed only, matching go test's own aggregate result for the
// parent.
func markResult(o *Outcome, id string, passed bool) {
	if passed {
		if !o.Failed[id] {
			o.Passed[id] = true
		}
		return
	}
	o.Failed[id] = true
	delete(o.Passed, id)
}

// procResult is one subprocess run via verify.RunProcess, with stdout and
// stderr captured separately and in full (unlike verify.Result, which
// tails a single combined buffer for a human): the language runners parse
// go test -json's stdout, or a JUnit/vitest-JSON report file, not a tail.
type procResult struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
	TimedOut bool
	StartErr error
}

// runProcess runs argv in dir with verify's shared timeout/process-group/
// env-filtering plumbing (verify.RunProcess): env never contains
// TYPESAFE_API_KEY, and the whole process group is killed (with a bounded
// post-kill wait) if timeout is exceeded.
func runProcess(ctx context.Context, dir string, argv []string, timeout time.Duration) procResult {
	return runProcessEnv(ctx, dir, argv, verify.FilteredEnv(), timeout)
}

// runProcessEnv is runProcess with an explicit environment.
func runProcessEnv(ctx context.Context, dir string, argv, env []string, timeout time.Duration) procResult {
	var stdout, stderr bytes.Buffer
	raw := verify.RunProcess(ctx, dir, argv, env, timeout, &stdout, &stderr)
	return procResult{
		Stdout:   stdout.Bytes(),
		Stderr:   stderr.Bytes(),
		ExitCode: raw.ExitCode,
		TimedOut: raw.TimedOut,
		StartErr: raw.StartErr,
	}
}

// EachRunner is a Runner that can run test files one per process.
type EachRunner interface {
	// RunEach runs each of testFiles in its own process and merges the
	// results per test id: a test file then runs as it does alone,
	// whatever else is selected. One timeout covers every run; an
	// Invalid or timed-out run ends it.
	RunEach(ctx context.Context, dir string, testFiles []string, index map[string]string, timeout time.Duration) Outcome
}

// runEach runs each file with run, under one deadline for all of them,
// merging the Outcomes (see EachRunner).
func runEach(ctx context.Context, files []string, timeout time.Duration, run func(ctx context.Context, file string) Outcome) Outcome {
	out := newOutcome()
	dctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for _, f := range files {
		if err := ctx.Err(); err != nil {
			out.Err = err.Error()
			break
		}
		if dctx.Err() != nil {
			out.TimedOut = true
			break
		}
		mergeOutcome(&out, run(dctx, f))
		if out.TimedOut || out.Invalid || out.Err != "" {
			break
		}
	}
	return out
}

// mergeOutcome adds o's results to dst (a failure wins, as in
// markResult).
func mergeOutcome(dst *Outcome, o Outcome) {
	for id := range o.Passed {
		markResult(dst, id, true)
	}
	for id := range o.Failed {
		markResult(dst, id, false)
	}
	dst.Unmapped += o.Unmapped
	dst.TimedOut = dst.TimedOut || o.TimedOut
	if o.Invalid && !dst.Invalid {
		dst.Invalid, dst.InvalidReason = true, o.InvalidReason
	}
	if dst.Err == "" {
		dst.Err = o.Err
	}
}
