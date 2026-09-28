package mutate

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

// requireTS skips the test unless node is on PATH and $CULL_TS names a
// typescript package: cull's evidence-gathering tools run inside plain Go
// checkouts with no node_modules of their own, so, unlike the Go tests,
// these are opt-in (run with CULL_TS set to exercise them for real).
func requireTS(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not found in PATH, skipping")
	}
	if os.Getenv("CULL_TS") == "" {
		t.Skip("CULL_TS not set, skipping")
	}
}

const tsOperatorsFixture = `import { something } from "./util";

export function helper(a: number, b: number, ok: boolean, done: boolean): number {
  if (a == b) {
  }
  if (a != b) {
  }
  if (a === b) {
  }
  if (a !== b) {
  }
  if (a < b) {
  }
  if (a > b) {
  }
  if (a <= b) {
  }
  if (a >= b) {
  }
  if (ok && !done) {
  }
  if (ok || !done) {
  }
  if (!ok) {
  }
  const s = a + b;
  const d = a - b;
  const p = a * b;
  const n = 41;
  const f = 3.5;
  const greeting = "hi";
  const flagT = true;
  const flagF = false;
  return s;
}

export function other(): number {
  return 1 + 2;
}
`

var tsDescSuffix = regexp.MustCompile(`^line \d+: (.+)$`)

func TestTSSitesOperators(t *testing.T) {
	requireTS(t)
	m := For("typescript")
	if m == nil {
		t.Fatal(`For("typescript") = nil`)
	}
	sites, err := m.Sites(".", "fixture.ts", []byte(tsOperatorsFixture))
	if err != nil {
		t.Fatalf("Sites: %v", err)
	}

	want := map[string]int{
		"flip ==":          1,
		"flip !=":          1,
		"flip ===":         1,
		"flip !==":         1,
		"flip <":           1,
		"flip >":           1,
		"flip <=":          1,
		"flip >=":          1,
		"flip &&":          1,
		"flip ||":          1,
		"drop !":           3,
		"flip +":           2,
		"flip -":           1,
		"flip *":           1,
		"+1 int":           3,
		"+1 float":         1,
		"mutate string":    1,
		"flip true":        1,
		"flip false":       1,
		"return undefined": 2,
	}
	got := map[string]int{}
	total := 0
	for _, s := range sites {
		sub := tsDescSuffix.FindStringSubmatch(s.Desc)
		if sub == nil {
			t.Errorf("Desc %q does not match expected format", s.Desc)
			continue
		}
		got[sub[1]]++
		total++
	}
	if total != 26 {
		t.Errorf("len(sites) = %d, want 26; got desc counts %+v", total, got)
	}
	for k, wantN := range want {
		if got[k] != wantN {
			t.Errorf("desc %q: got %d, want %d", k, got[k], wantN)
		}
	}

	src := []byte(tsOperatorsFixture)
	for _, s := range sites {
		text := string(src[s.Start:s.End])
		if text == `"./util"` {
			t.Errorf("mutated the import specifier: %+v", s)
		}
		mutated, err := Apply(src, s)
		if err != nil {
			t.Fatalf("Apply(%+v): %v", s, err)
		}
		parseTSFixture(t, mutated, s.Desc)
	}
}

func TestTSSitesUnicodeOffsets(t *testing.T) {
	requireTS(t)
	const src = `const greeting = "héllo wörld";

export function f(a: number, b: number): boolean {
  return a == b;
}
`
	m := For("typescript")
	sites, err := m.Sites(".", "fixture.ts", []byte(src))
	if err != nil {
		t.Fatalf("Sites: %v", err)
	}
	var eq *Site
	for i := range sites {
		if sites[i].Desc == "line 4: flip ==" {
			eq = &sites[i]
		}
	}
	if eq == nil {
		t.Fatalf("no comparison site found among %+v", sites)
	}
	out, err := Apply([]byte(src), *eq)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	want := `const greeting = "héllo wörld";

export function f(a: number, b: number): boolean {
  return a != b;
}
`
	if string(out) != want {
		t.Fatalf("Apply with unicode-preceded offsets:\ngot:  %q\nwant: %q", out, want)
	}
}

func TestTSSitesMissingToolchain(t *testing.T) {
	t.Setenv("PATH", "")
	m := For("typescript")
	sites, err := m.Sites(".", "fixture.ts", []byte(tsOperatorsFixture))
	if err != nil {
		t.Fatalf("Sites: %v", err)
	}
	if sites != nil {
		t.Fatalf("Sites with no node on PATH = %+v, want nil", sites)
	}
}

// parseTSFixture requires the typescript package found via $CULL_TS (or
// node_modules/typescript) and fails t if src has any parse diagnostics.
// desc identifies the mutation under test in failure messages.
func parseTSFixture(t *testing.T, src []byte, desc string) {
	t.Helper()
	tsDir := os.Getenv("CULL_TS")
	script := `
const ts = require(process.argv[1]);
const chunks = [];
process.stdin.on("data", (c) => chunks.push(c));
process.stdin.on("end", () => {
  const src = Buffer.concat(chunks).toString("utf8");
  const sf = ts.createSourceFile("fixture.ts", src, ts.ScriptTarget.Latest, true, ts.ScriptKind.TS);
  const n = (sf.parseDiagnostics && sf.parseDiagnostics.length) || 0;
  if (n > 0) {
    process.stderr.write(String(sf.parseDiagnostics[0].messageText));
    process.exit(1);
  }
});
`
	cmd := exec.Command("node", "-e", script, tsDir)
	cmd.Stdin = bytes.NewReader(src)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("parse failed for mutation %q: %s\n%s\nsrc:\n%s", desc, out, err, src)
	}
}
