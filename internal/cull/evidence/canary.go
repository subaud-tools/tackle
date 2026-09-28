package evidence

import (
	"strings"
)

// reasonShadowed is the unmeasured reason for a module whose canary no
// measured test noticed: the tests import some other copy of it.
const reasonShadowed = "tests don't load this copy of the module (an installed package or path setting shadows it)"

// canarySource returns src with a module-breaking change planted: any test
// that really loads this copy of the module fails (or fails to build or
// collect). Python: raise RuntimeError as the first statement after the
// docstring and any `from __future__` imports. Go: a reference to an
// undefined identifier, so the package fails to compile. TypeScript: a
// throw as the first statement after the leading imports. ok is false for
// any other language.
func canarySource(lang string, src []byte) ([]byte, bool) {
	switch lang {
	case "python":
		return insertAtLine(src, pyCanaryLine(splitLines(string(src))), "raise RuntimeError(\"cull canary\")\n"), true
	case "go":
		s := string(src)
		if s != "" && !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		return []byte(s + "var _ = cullCanaryUndefined\n"), true
	case "typescript":
		return insertAtLine(src, tsCanaryLine(splitLines(string(src))), "throw new Error(\"cull canary\");\n"), true
	}
	return nil, false
}

// splitLines splits s into lines, each keeping its trailing "\n".
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.SplitAfter(strings.TrimSuffix(s, "\n"), "\n")
}

// insertAtLine inserts text before line i of src (at the end when i is
// past the last line, adding the missing newline first).
func insertAtLine(src []byte, i int, text string) []byte {
	lines := splitLines(string(src))
	if i >= len(lines) {
		s := string(src)
		if s != "" && !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		return []byte(s + text)
	}
	off := 0
	for _, l := range lines[:i] {
		off += len(l)
	}
	return []byte(string(src[:off]) + text + string(src[off:]))
}

// pyCanaryLine is the index of the first line after the module docstring
// and any `from __future__` imports (blank and comment lines skipped).
func pyCanaryLine(lines []string) int {
	i := 0
	trivia := func() {
		for i < len(lines) {
			t := strings.TrimSpace(lines[i])
			if t != "" && !strings.HasPrefix(t, "#") {
				return
			}
			i++
		}
	}
	trivia()
	if i < len(lines) {
		if q, rest := pyStringStart(lines[i]); q != "" {
			i = pyStringEnd(lines, i, q, rest) + 1
		}
	}
	for {
		trivia()
		if i >= len(lines) || !strings.HasPrefix(lines[i], "from __future__ ") {
			return i
		}
		i = pyStatementEnd(lines, i) + 1
	}
}

// pyStringStart reports the quote a top-level string literal on line opens
// with ("" if the line doesn't start with one) and the text after it.
func pyStringStart(line string) (quote, rest string) {
	s := strings.TrimLeft(line, "rRuUbBfF")
	if len(line)-len(s) > 2 {
		return "", ""
	}
	for _, q := range []string{`"""`, `'''`, `"`, `'`} {
		if strings.HasPrefix(s, q) {
			return q, s[len(q):]
		}
	}
	return "", ""
}

// pyStringEnd is the line where the string opened on line i (with rest
// following the opening quote) closes.
func pyStringEnd(lines []string, i int, q, rest string) int {
	for {
		if closes(rest, q) || i >= len(lines)-1 {
			return i
		}
		i++
		rest = lines[i]
	}
}

// closes reports whether s holds an unescaped q.
func closes(s, q string) bool {
	for j := 0; j+len(q) <= len(s); j++ {
		if s[j] == '\\' {
			j++
			continue
		}
		if s[j:j+len(q)] == q {
			return true
		}
	}
	return false
}

// pyStatementEnd is the last line of the simple statement starting on line
// i: through a closing parenthesis, or backslash continuations.
func pyStatementEnd(lines []string, i int) int {
	if strings.Contains(lines[i], "(") && !strings.Contains(lines[i], ")") {
		for i < len(lines)-1 && !strings.Contains(lines[i], ")") {
			i++
		}
		return i
	}
	for i < len(lines)-1 && strings.HasSuffix(strings.TrimRight(lines[i], " \t\r\n"), "\\") {
		i++
	}
	return i
}

// tsCanaryLine is the index of the first line after the leading comments,
// directives and import declarations. An import ends on the first line
// holding a quote: the module specifier always comes last.
func tsCanaryLine(lines []string) int {
	inBlock := false
	for i := 0; i < len(lines); i++ {
		t := strings.TrimSpace(lines[i])
		switch {
		case inBlock:
			inBlock = !strings.Contains(t, "*/")
		case t == "" || strings.HasPrefix(t, "//") || strings.HasPrefix(t, "#!"):
		case strings.HasPrefix(t, "/*"):
			inBlock = !strings.Contains(t[2:], "*/")
		case tsDirective(t):
		case tsImport(t):
			for i < len(lines)-1 && !strings.ContainsAny(lines[i], `'"`) {
				i++
			}
		default:
			return i
		}
	}
	return len(lines)
}

func tsDirective(t string) bool {
	t = strings.TrimSuffix(t, ";")
	return t == `"use strict"` || t == `'use strict'`
}

func tsImport(t string) bool {
	rest, ok := strings.CutPrefix(t, "import")
	return ok && rest != "" && strings.ContainsAny(rest[:1], " \t{*'\"")
}
