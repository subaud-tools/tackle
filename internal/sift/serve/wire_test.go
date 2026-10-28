package serve

import (
	"flag"
	"testing"

	"github.com/subaud-tools/tools-common/wiregen"
)

// updateWire rewrites web/wire.d.ts instead of asserting it. `just sift-wire`.
var updateWire = flag.Bool("update", false, "rewrite web/wire.d.ts from the Go structs")

// wireRoots is the contract: what the page reads off GET /api/review (its
// rows and files also come back from a decision), and what the agent reads
// off GET /api/agent/status. Everything they reach follows.
var wireRoots = []any{
	review{},
	status{},
}

const wireHeader = `// wire.d.ts — the Go↔page contract, generated. Do not edit.
//
// Written by internal/sift/serve/wire_test.go from the json tags of review
// (GET /api/review) and status (GET /api/agent/status), and every type they
// reach. Regenerate with ` + "`just sift-wire`" + `.
`

func TestTheWireTypeIsCommitted(t *testing.T) {
	opts := wiregen.Options{Header: wireHeader, Regenerate: "just sift-wire"}
	if err := wiregen.CheckCommitted("../web/wire.d.ts", wireRoots, opts, *updateWire); err != nil {
		t.Fatal(err)
	}
}
