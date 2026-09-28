package mutate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// generatedPattern matches the standard Go "generated file" marker
// (https://go.dev/s/generatedcode): a comment line of the form
// "// Code generated ... DO NOT EDIT.".
var generatedPattern = regexp.MustCompile(`(?m)^// Code generated .* DO NOT EDIT\.$`)

// goMutator finds mutation sites in Go source using go/ast byte offsets.
type goMutator struct{}

func (goMutator) Lang() string { return "go" }

// Sites never mutates _test.go files or generated files, and never touches
// import paths or struct tags.
func (goMutator) Sites(root, relpath string, src []byte) ([]Site, error) {
	if strings.HasSuffix(relpath, "_test.go") {
		return nil, nil
	}
	if generatedPattern.Match(src) {
		return nil, nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, relpath, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("mutate: parse %s: %w", relpath, err)
	}

	skip := map[*ast.BasicLit]bool{}
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.IMPORT {
			continue
		}
		for _, spec := range gd.Specs {
			if is, ok := spec.(*ast.ImportSpec); ok {
				skip[is.Path] = true
			}
		}
	}
	ast.Inspect(f, func(n ast.Node) bool {
		st, ok := n.(*ast.StructType)
		if !ok || st.Fields == nil {
			return true
		}
		for _, field := range st.Fields.List {
			if field.Tag != nil {
				skip[field.Tag] = true
			}
		}
		return true
	})

	var sites []Site
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.BinaryExpr:
			if repl, name, ok := flipBinaryOp(node.Op); ok {
				start := fset.Position(node.OpPos).Offset
				end := start + len(node.Op.String())
				sites = append(sites, Site{
					Start:       start,
					End:         end,
					Replacement: repl,
					Desc:        fmt.Sprintf("line %d: flip %s", fset.Position(node.OpPos).Line, name),
				})
			}
		case *ast.UnaryExpr:
			if node.Op == token.NOT {
				start := fset.Position(node.OpPos).Offset
				end := start + len(token.NOT.String())
				sites = append(sites, Site{
					Start:       start,
					End:         end,
					Replacement: "",
					Desc:        fmt.Sprintf("line %d: drop !", fset.Position(node.OpPos).Line),
				})
			}
		case *ast.Ident:
			if node.Name == "true" || node.Name == "false" {
				repl := "false"
				if node.Name == "false" {
					repl = "true"
				}
				start := fset.Position(node.Pos()).Offset
				end := start + len(node.Name)
				sites = append(sites, Site{
					Start:       start,
					End:         end,
					Replacement: repl,
					Desc:        fmt.Sprintf("line %d: flip %s", fset.Position(node.Pos()).Line, node.Name),
				})
			}
		case *ast.BasicLit:
			if skip[node] {
				return true
			}
			start := fset.Position(node.Pos()).Offset
			end := start + len(node.Value)
			line := fset.Position(node.Pos()).Line
			switch node.Kind {
			case token.INT:
				if repl, ok := bumpIntLiteral(node.Value); ok {
					sites = append(sites, Site{Start: start, End: end, Replacement: repl, Desc: fmt.Sprintf("line %d: +1 int", line)})
				}
			case token.FLOAT:
				if repl, ok := bumpFloatLiteral(node.Value); ok {
					sites = append(sites, Site{Start: start, End: end, Replacement: repl, Desc: fmt.Sprintf("line %d: +1 float", line)})
				}
			case token.STRING:
				if repl, ok := mutateStringLiteral(node.Value); ok {
					sites = append(sites, Site{Start: start, End: end, Replacement: repl, Desc: fmt.Sprintf("line %d: mutate string", line)})
				}
			}
		}
		return true
	})

	return sites, nil
}

// flipBinaryOp returns the replacement token text and a display name for op,
// or ok=false if op is not one of the mutated comparison/logical/arithmetic
// operators.
func flipBinaryOp(op token.Token) (repl, name string, ok bool) {
	switch op {
	case token.EQL:
		return "!=", "==", true
	case token.NEQ:
		return "==", "!=", true
	case token.LSS:
		return ">=", "<", true
	case token.GEQ:
		return "<", ">=", true
	case token.GTR:
		return "<=", ">", true
	case token.LEQ:
		return ">", "<=", true
	case token.LAND:
		return "||", "&&", true
	case token.LOR:
		return "&&", "||", true
	case token.ADD:
		return "-", "+", true
	case token.SUB:
		return "+", "-", true
	case token.MUL:
		return "+", "*", true
	default:
		return "", "", false
	}
}

// bumpIntLiteral returns lit's value plus one, formatted with the same base
// prefix as lit, or ok=false if lit cannot be parsed as an integer.
func bumpIntLiteral(lit string) (string, bool) {
	clean := strings.ReplaceAll(lit, "_", "")
	n, ok := new(big.Int).SetString(clean, 0)
	if !ok {
		return "", false
	}
	n.Add(n, big.NewInt(1))
	switch {
	case strings.HasPrefix(clean, "0x") || strings.HasPrefix(clean, "0X"):
		return "0x" + n.Text(16), true
	case strings.HasPrefix(clean, "0o") || strings.HasPrefix(clean, "0O"):
		return "0o" + n.Text(8), true
	case strings.HasPrefix(clean, "0b") || strings.HasPrefix(clean, "0B"):
		return "0b" + n.Text(2), true
	case len(clean) > 1 && clean[0] == '0':
		return "0" + n.Text(8), true
	default:
		return n.Text(10), true
	}
}

// bumpFloatLiteral returns lit's value plus one, or ok=false if lit cannot
// be parsed as a float.
func bumpFloatLiteral(lit string) (string, bool) {
	v, err := strconv.ParseFloat(strings.ReplaceAll(lit, "_", ""), 64)
	if err != nil {
		return "", false
	}
	return strconv.FormatFloat(v+1, 'g', -1, 64), true
}

// mutateStringLiteral inserts "_mut" just before the closing quote/backtick
// of lit (a Go string literal including its delimiters).
func mutateStringLiteral(lit string) (string, bool) {
	if len(lit) < 2 {
		return "", false
	}
	quote := lit[len(lit)-1]
	if lit[0] != quote {
		return "", false
	}
	return lit[:len(lit)-1] + "_mut" + string(quote), true
}
