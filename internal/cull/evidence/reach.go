package evidence

import (
	"sort"

	"github.com/schuettc/tackle/internal/cull/mutate"
)

// reach returns (sorted) the indexes of the defs reached from seeds within
// one module, seeds included: transitively, every def named by a reached
// def's refs (module-level functions, types and variables, and methods of
// any owner by name alone), and for each class or type so reached, all of
// its methods (dunders, constructors, Go methods on it). A reached method
// reaches its owner class too (cls(...), type(self)(...), replace(self)
// construct it unnamed). Every "init" def (Go) is reached whenever
// anything is: defs are one file's. Resolution is by
// name only: over-inclusion is accepted, dynamic dispatch may be missed.
func reach(defs []mutate.Def, seeds []int) []int {
	byName := map[string][]int{}
	byOwner := map[string][]int{}
	for i, d := range defs {
		byName[d.Name] = append(byName[d.Name], i)
		if d.Owner != "" {
			byOwner[d.Owner] = append(byOwner[d.Owner], i)
		}
	}
	seen := map[int]bool{}
	var queue []int
	visit := func(i int) {
		if !seen[i] {
			seen[i] = true
			queue = append(queue, i)
		}
	}
	for _, i := range seeds {
		visit(i)
	}
	if len(queue) > 0 {
		for i, d := range defs {
			if d.Kind == "init" {
				visit(i)
			}
		}
	}
	for len(queue) > 0 {
		d := defs[queue[0]]
		queue = queue[1:]
		if d.Owner == "" && (d.Kind == "class" || d.Kind == "type") {
			for _, j := range byOwner[d.Name] {
				visit(j)
			}
		}
		if d.Owner != "" {
			for _, j := range byName[d.Owner] {
				if defs[j].Owner == "" {
					visit(j)
				}
			}
		}
		for _, r := range d.Refs {
			for _, j := range byName[r] {
				visit(j)
			}
		}
	}
	out := make([]int, 0, len(seen))
	for i := range seen {
		out = append(out, i)
	}
	sort.Ints(out)
	return out
}

// qualName is d's name for seeding its sample: Owner.Name for a method.
func qualName(d mutate.Def) string {
	if d.Owner != "" {
		return d.Owner + "." + d.Name
	}
	return d.Name
}
