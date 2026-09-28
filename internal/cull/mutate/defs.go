package mutate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"

	"github.com/schuettc/tackle/internal/cull/extract/python"
	"github.com/schuettc/tackle/internal/cull/extract/ts"
)

// Def is one definition in a source file: a module-level function, type,
// class or variable (Owner ""), or a method (Owner: its class, or its Go
// receiver's type). Kind is "function", "class", "type", "var" or
// "method". Start and End are UTF-8 byte offsets of its span; Refs lists
// (sorted, deduplicated) the names used anywhere inside it.
type Def struct {
	Name, Owner, Kind string
	Start, End        int
	Refs              []string
}

// Defs returns the definitions in src, a file at relpath under root, as
// lang's parser reports them: go/ast for Go, the Python and TypeScript
// helpers' --defs mode otherwise. A missing toolchain or an unknown lang
// yields none and no error.
func Defs(lang, root, relpath string, src []byte) ([]Def, error) {
	switch lang {
	case "go":
		return goDefs(relpath, src)
	case "python":
		ds, err := python.Defs(root, relpath, src)
		out := make([]Def, len(ds))
		for i, d := range ds {
			out[i] = Def{d.Name, d.Owner, d.Kind, d.Start, d.End, d.Refs}
		}
		return out, err
	case "typescript":
		ds, err := ts.Defs(root, relpath, src)
		out := make([]Def, len(ds))
		for i, d := range ds {
			out[i] = Def{d.Name, d.Owner, d.Kind, d.Start, d.End, d.Refs}
		}
		return out, err
	}
	return nil, nil
}

// goDefs lists src's top-level funcs, methods (Owner: the receiver's type
// name; the receiver itself is not a ref), types, vars and consts. Spans
// start at the declaration itself, doc comments excluded.
func goDefs(relpath string, src []byte) ([]Def, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, relpath, src, 0)
	if err != nil {
		return nil, fmt.Errorf("mutate: parse %s: %w", relpath, err)
	}
	off := func(p token.Pos) int { return fset.Position(p).Offset }
	refs := func(nodes ...ast.Node) []string {
		seen := map[string]bool{}
		for _, n := range nodes {
			if n == nil {
				continue
			}
			ast.Inspect(n, func(x ast.Node) bool {
				if id, ok := x.(*ast.Ident); ok {
					seen[id.Name] = true
				}
				return true
			})
		}
		out := make([]string, 0, len(seen))
		for k := range seen {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}
	var out []Def
	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			def := Def{Name: d.Name.Name, Kind: "function", Start: off(d.Pos()), End: off(d.End())}
			if d.Recv != nil && len(d.Recv.List) > 0 {
				def.Owner, def.Kind = recvType(d.Recv.List[0].Type), "method"
			}
			var body ast.Node
			if d.Body != nil {
				body = d.Body
			}
			def.Refs = refs(d.Type, body)
			out = append(out, def)
		case *ast.GenDecl:
			for _, s := range d.Specs {
				// A lone unparenthesized spec spans its keyword too.
				start, end := off(s.Pos()), off(s.End())
				if !d.Lparen.IsValid() {
					start, end = off(d.Pos()), off(d.End())
				}
				switch s := s.(type) {
				case *ast.TypeSpec:
					out = append(out, Def{Name: s.Name.Name, Kind: "type", Start: start, End: end, Refs: refs(s)})
				case *ast.ValueSpec:
					for _, n := range s.Names {
						out = append(out, Def{Name: n.Name, Kind: "var", Start: start, End: end, Refs: refs(s)})
					}
				}
			}
		}
	}
	return out, nil
}

// recvType is a receiver type's name: T for T, *T, T[K] or *T[K].
func recvType(x ast.Expr) string {
	for {
		switch t := x.(type) {
		case *ast.StarExpr:
			x = t.X
		case *ast.IndexExpr:
			x = t.X
		case *ast.IndexListExpr:
			x = t.X
		case *ast.ParenExpr:
			x = t.X
		case *ast.Ident:
			return t.Name
		default:
			return "?"
		}
	}
}
