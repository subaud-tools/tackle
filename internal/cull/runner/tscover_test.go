package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// V8 offsets are UTF-16 units; the innermost range decides; the script
// may carry node's sourceURL line after the source.
func TestV8Helpers(t *testing.T) {
	src := []byte("const s = \"é😀\";\nfunction f() { return 1; }\n")
	u := utf16Offsets(src)
	if u[len(src)] != len(src)-3 {
		// é: 2 bytes → 1 unit; 😀: 4 bytes → 2 units.
		t.Fatalf("utf16 len %d", u[len(src)])
	}
	n := u[len(src)]
	url := "file:///x/calc.ts"
	if !aligned([]v8Range{{0, n, 1}}, n, url) || !aligned([]v8Range{{0, n + len("//# sourceURL="+url) + 2, 1}}, n, url) ||
		aligned([]v8Range{{0, n + 7, 1}}, n, url) {
		t.Error("aligned")
	}
	ranges := []v8Range{{0, 40, 1}, {20, 30, 0}, {25, 28, 2}}
	for off, want := range map[int]bool{5: true, 21: false, 26: true, 50: true} {
		if v8Ran(ranges, off) != want {
			t.Errorf("v8Ran(%d) != %v", off, want)
		}
	}
}

// A site runs if any of its bytes ran, not only its ends: here both ends
// sit in unexecuted ranges and the middle ran.
func TestV8SiteRanMiddle(t *testing.T) {
	src := []byte("aaaaaaaaaaaaaaaaaaaa")
	u16 := utf16Offsets(src)
	ranges := []v8Range{{0, 20, 1}, {2, 6, 0}, {12, 18, 0}}
	if !v8SiteRan(ranges, u16, 4, 14) {
		t.Error("a site whose middle ran reads as not run")
	}
	if v8SiteRan(ranges, u16, 3, 5) {
		t.Error("a site inside an unexecuted range reads as run")
	}
}

// Two coverage files (the test's process and a node child) for one
// script with disjoint hits: each group decides alone, so a count of 0
// in one doesn't hide a count of 1 in the other.
func TestV8GroupsPerScript(t *testing.T) {
	dir := t.TempDir()
	src := "export function f(x: number): boolean {\n  return x > 3;\n}\n\nexport function g(x: number): number {\n  return x * 2;\n}\n"
	mod := filepath.Join(dir, "calc.ts")
	if err := os.WriteFile(mod, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	real, err := filepath.EvalSymlinks(mod)
	if err != nil {
		t.Fatal(err)
	}
	cov := filepath.Join(dir, "cov")
	if err := os.Mkdir(cov, 0o700); err != nil {
		t.Fatal(err)
	}
	fOff, gOff := strings.Index(src, "x > 3"), strings.Index(src, "x * 2")
	fs, gs := strings.Index(src, "function f"), strings.Index(src, "function g")
	write := func(name string, fc, gc int) {
		body := fmt.Sprintf(`{"result":[{"url":%q,"functions":[{"functionName":"","ranges":[{"startOffset":0,"endOffset":%d,"count":1}]},{"functionName":"f","ranges":[{"startOffset":%d,"endOffset":%d,"count":%d}]},{"functionName":"g","ranges":[{"startOffset":%d,"endOffset":%d,"count":%d}]}]}]}`,
			"file://"+real, len(src), fs, strings.Index(src, "}\n")+1, fc, gs, len(src)-1, gc)
		if err := os.WriteFile(filepath.Join(cov, name), []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write("coverage-1.json", 1, 0)
	write("coverage-2.json", 0, 1)
	groups, found, err := readV8Dir(cov, real)
	if err != nil || !found || len(groups) != 2 {
		t.Fatalf("groups %v found %v err %v", groups, found, err)
	}
	u16 := utf16Offsets([]byte(src))
	if !allAligned(groups, u16[len(src)], "file://"+real) {
		t.Error("groups not aligned")
	}
	if !v8GroupsSiteRan(groups, u16, fOff, fOff+5) || !v8GroupsSiteRan(groups, u16, gOff, gOff+5) {
		t.Error("a site that ran in one process reads as not run")
	}
	// One misaligned group (a transformed copy) makes the file cover all.
	if allAligned(append(groups, [][]v8Range{{{0, len(src) + 7, 1}}}...), u16[len(src)], "file://"+real) {
		t.Error("a misaligned group passes")
	}
}
