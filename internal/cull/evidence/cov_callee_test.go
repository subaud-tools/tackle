package evidence

import (
	"fmt"
	"strings"
	"testing"

	"github.com/schuettc/tackle/internal/cull/cases"
	"github.com/schuettc/tackle/internal/cull/mutate"
)

// A callee whose span is located but holds no definition of its own
// (Defs reports only an enclosing class, or nothing) samples its span
// directly: perCallee sites, never zero.
func TestCalleeWithoutDefSampled(t *testing.T) {
	var b strings.Builder
	b.WriteString("class K:\n")
	body := "    def f(self, x):\n"
	for i := 0; i < 20; i++ {
		body += fmt.Sprintf("        x = x + %d\n", i)
	}
	body += "        return x\n"
	b.WriteString(body)
	src := []byte(b.String())
	var sites []mutate.Site
	for off := strings.Index(b.String(), "x + "); off >= 0; {
		sites = append(sites, mutate.Site{Start: off + 2, End: off + 3, Desc: fmt.Sprintf("s%d", off)})
		n := strings.Index(b.String()[off+1:], "x + ")
		if n < 0 {
			break
		}
		off += 1 + n
	}
	callee := []cases.Callee{{Symbol: "K.f", Source: strings.TrimPrefix(body, "    ")}}
	for name, defs := range map[string][]mutate.Def{
		"none":      nil,
		"enclosing": {{Name: "K", Kind: "class", Start: 0, End: len(src)}},
		// A def inside the callee that holds none of its sites.
		"elsewhere": {{Name: "K", Kind: "class", Start: 0, End: len(src)}, {Name: "g", Kind: "function", Start: len(src) - 1, End: len(src)}},
	} {
		got := targetSites("m.py", src, sites, callee, defs, nolog)
		if len(got) != perCallee {
			t.Errorf("%s: %d sites, want %d", name, len(got), perCallee)
		}
	}
}
