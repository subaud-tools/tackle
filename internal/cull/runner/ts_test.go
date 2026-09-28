package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestForTSNoPerTestReporterError(t *testing.T) {
	root := t.TempDir()
	write(t, root, "web/package.json", `{"scripts": {"test": "jest"}}`)
	write(t, root, "web/foo.test.ts", "// placeholder\n")
	_, err := For(root, "typescript", "web/foo.test.ts")
	if err == nil {
		t.Fatal("For: want error for a TS project with no per-test reporter script")
	}
	if err.Error() != "no per-test reporter for this TypeScript runner" {
		t.Errorf("err = %q, want %q", err.Error(), "no per-test reporter for this TypeScript runner")
	}
}

func TestForUnknownLangErrors(t *testing.T) {
	if _, err := For(t.TempDir(), "ruby", "x_test.rb"); err == nil {
		t.Fatal("For: want error for unknown language")
	}
}

func write(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func requireCullTS(t *testing.T) {
	t.Helper()
	if os.Getenv("CULL_TS") == "" {
		t.Skip("CULL_TS not set; skipping TypeScript runner fixtures")
	}
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not found")
	}
}

func TestTSRunNodeTestPassFailAndDescribe(t *testing.T) {
	requireCullTS(t)
	root := t.TempDir()
	write(t, root, "package.json", `{"scripts": {"test": "node --test"}}`)
	write(t, root, "a.test.js", `const { test, describe } = require('node:test');
const assert = require('node:assert');

describe('outer', () => {
  test('passes', () => {
    assert.ok(true);
  });
  test('fails', () => {
    assert.ok(false);
  });
});

test('top level', () => {
  assert.ok(true);
});

test('env clean', () => {
  if (process.env.TYPESAFE_API_KEY) {
    throw new Error('leaked TYPESAFE_API_KEY');
  }
});
`)
	r, err := For(root, "typescript", "a.test.js")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	index := map[string]string{
		"outer > passes": "ts:a.test.js:outer > passes",
		"outer > fails":  "ts:a.test.js:outer > fails",
		"top level":      "ts:a.test.js:top level",
		"env clean":      "ts:a.test.js:env clean",
	}
	out := r.Run(context.Background(), root, []string{"a.test.js"}, index, 30*time.Second)
	if out.TimedOut {
		t.Fatalf("TimedOut = true, want false: %+v", out)
	}
	if !out.Passed["ts:a.test.js:outer > passes"] {
		t.Errorf("Passed[outer > passes] = false, want true; out=%+v", out)
	}
	if !out.Failed["ts:a.test.js:outer > fails"] {
		t.Errorf("Failed[outer > fails] = false, want true; out=%+v", out)
	}
	if !out.Passed["ts:a.test.js:top level"] {
		t.Errorf("Passed[top level] = false, want true; out=%+v", out)
	}
	if !out.Passed["ts:a.test.js:env clean"] {
		t.Errorf("Passed[env clean] = false, want true (TYPESAFE_API_KEY must not leak); out=%+v", out)
	}
}

func TestTSRunNodeTestScriptErrorsBeforeReportingIsInvalid(t *testing.T) {
	requireCullTS(t)
	root := t.TempDir()
	write(t, root, "package.json", `{"scripts": {"test": "node --test"}}`)
	// a.test.js is intentionally never created: node exits ("Could not
	// find 'a.test.js'") before its test harness ever opens the
	// --test-reporter-destination file, so no report exists to parse.
	r, err := For(root, "typescript", "a.test.js")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	index := map[string]string{}
	out := r.Run(context.Background(), root, []string{"a.test.js"}, index, 30*time.Second)
	if !out.Invalid {
		t.Fatalf("Invalid = false, want true (script errored before reporting): out=%+v", out)
	}
	if out.InvalidReason == "" {
		t.Errorf("InvalidReason = %q, want a non-empty reason", out.InvalidReason)
	}
}

func TestTSRunVitestPassFailAndUnmapped(t *testing.T) {
	requireCullTS(t)
	if _, err := exec.LookPath("npx"); err != nil {
		t.Skip("npx not found")
	}
	probe := exec.Command("npx", "--offline", "--no-install", "vitest", "--version")
	if out, err := probe.CombinedOutput(); err != nil {
		t.Skipf("vitest not resolvable offline via npx: %v: %s", err, out)
	}

	root := t.TempDir()
	write(t, root, "package.json", `{"scripts": {"test": "vitest run"}}`)
	write(t, root, "a.test.ts", `import { describe, it, expect } from 'vitest';

describe('outer', () => {
  it('passes', () => {
    expect(1).toBe(1);
  });
  it('fails', () => {
    expect(1).toBe(2);
  });
});

it('top level', () => {
  expect(true).toBe(true);
});
`)
	r, err := For(root, "typescript", "a.test.ts")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	index := map[string]string{
		"outer > passes": "ts:a.test.ts:outer > passes",
		"outer > fails":  "ts:a.test.ts:outer > fails",
		// "top level" deliberately omitted -> Unmapped.
	}
	out := r.Run(context.Background(), root, []string{"a.test.ts"}, index, 60*time.Second)
	if out.TimedOut {
		t.Fatalf("TimedOut = true, want false: %+v", out)
	}
	if !out.Passed["ts:a.test.ts:outer > passes"] {
		t.Errorf("Passed[outer > passes] = false, want true; out=%+v", out)
	}
	if !out.Failed["ts:a.test.ts:outer > fails"] {
		t.Errorf("Failed[outer > fails] = false, want true; out=%+v", out)
	}
	if out.Unmapped != 1 {
		t.Errorf("Unmapped = %d, want 1 (top level not in index): out=%+v", out.Unmapped, out)
	}
}

// A node:test file that throws while loading is reported as one failed
// testcase named after the file: the run didn't execute its tests, so
// it is Invalid (like a pytest collection error).
func TestTSRunNodeFileLoadErrorIsInvalid(t *testing.T) {
	requireCullTS(t)
	root := t.TempDir()
	write(t, root, "package.json", `{"scripts": {"test": "node --test"}}`)
	write(t, root, "a.test.js", "throw new Error('boom');\nrequire('node:test').test('t', () => {});\n")
	r, err := For(root, "typescript", "a.test.js")
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	out := r.Run(context.Background(), root, []string{"a.test.js"}, map[string]string{"t": "ts:a.test.js:t"}, 30*time.Second)
	if !out.Invalid || !strings.Contains(out.InvalidReason, "a.test.js") {
		t.Fatalf("out = %+v", out)
	}
}
