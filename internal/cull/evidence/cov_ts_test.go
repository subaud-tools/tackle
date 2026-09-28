package evidence

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

const covTSSrc = `export const LIMIT: number = 3;

export function f(x: number): number {
  return x + 1;
}

export function g(x: number): boolean {
  return x > LIMIT;
}

export function h(x: number): number {
  return x * 2;
}
`

func covTSTest(log, fn string, tests map[string]string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "import { test } from 'node:test';\nimport assert from 'node:assert';\nimport { appendFileSync } from 'node:fs';\nimport { %s } from './calc.ts';\n\n", fn)
	names := map[string]bool{}
	for n := range tests {
		names[n] = true
	}
	for _, n := range sortedKeys(names) {
		fmt.Fprintf(&b, "test(%q, () => {\n  appendFileSync(%q, %q);\n  %s;\n});\n\n", n, log, n+"\n", tests[n])
	}
	return b.String()
}

// TypeScript (node:test): coverage is per test file, every test in a
// file counting as covering what the file ran; a module-level mutant
// runs every file; an unused function's mutants are uncovered.
func TestTSCoverageSelectsFiles(t *testing.T) {
	needTools(t)
	if os.Getenv("CULL_TS") == "" {
		t.Skip("CULL_TS not set")
	}
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not on PATH")
	}
	l := newRunLog(t, covTSSrc)
	root := newRepo(t, map[string]string{
		"package.json": `{"type": "module", "scripts": {"test": "node --test"}}`,
		"calc.ts":      covTSSrc,
		"a.test.ts":    covTSTest(l.path, "f", map[string]string{"fa": "assert.equal(f(1), 2)", "fb": "f(0)"}),
		"b.test.ts":    covTSTest(l.path, "g", map[string]string{"g1": "assert.ok(g(5))", "g2": "assert.ok(!g(1))"}),
	})
	mod := Module{Rel: "calc.ts", Lang: "typescript", TestFiles: []string{"a.test.ts", "b.test.ts"},
		Tests: []string{"ts:a.test.ts:fa", "ts:a.test.ts:fb", "ts:b.test.ts:g1", "ts:b.test.ts:g2"}}
	var log bytes.Buffer
	o := opts(1)
	o.Stderr = &log
	m, err := RunMatrix(context.Background(), root, []Module{mod}, o)
	if err != nil {
		t.Fatal(err)
	}
	want := map[int][]string{1: {"fa", "fb", "g1", "g2"}, 4: {"fa", "fb"}, 8: {"g1", "g2"}}
	if got := l.byLine(t); !reflect.DeepEqual(got, want) {
		t.Errorf("tests run per mutated line:\n got %v\nwant %v\n%s", got, want, log.String())
	}
	for _, mu := range m.Mutants {
		if strings.HasPrefix(mu.Site.Desc, "line 12:") != mu.Uncovered {
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
