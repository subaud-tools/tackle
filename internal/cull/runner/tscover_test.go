package runner

import "testing"

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
