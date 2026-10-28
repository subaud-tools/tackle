package serve

import (
	"flag"
	"testing"

	"github.com/subaud-tools/tackle/internal/cull/judge"
	"github.com/subaud-tools/tools-common/wiregen"
)

// updateWire rewrites web/wire.d.ts instead of asserting it. `just cull-wire`.
var updateWire = flag.Bool("update", false, "rewrite web/wire.d.ts from the Go structs")

// wireRoots is the contract: what the page reads off GET /api/review, and
// the two shapes an item's state is stored in (the page narrows state by
// kind, web/kinds.d.ts). Everything they reach follows.
var wireRoots = []any{
	review{},
	judge.State{},
	judge.GroupState{},
}

const wireHeader = `// wire.d.ts — the Go↔page contract, generated. Do not edit.
//
// Written by internal/cull/serve/wire_test.go from the json tags of review
// (GET /api/review), judge.State and judge.GroupState, and every type they
// reach. An item's state and jev are unknown here: Go stores them as bytes
// whose shape depends on the item's kind, which web/kinds.d.ts narrows.
// Regenerate with ` + "`just cull-wire`" + `.
`

func TestTheWireTypeIsCommitted(t *testing.T) {
	// AllowUnknown: itemBase.State and Jev are json.RawMessage, sent as stored.
	opts := wiregen.Options{Header: wireHeader, Regenerate: "just cull-wire", AllowUnknown: true}
	if err := wiregen.CheckCommitted("../web/wire.d.ts", wireRoots, opts, *updateWire); err != nil {
		t.Fatal(err)
	}
}
