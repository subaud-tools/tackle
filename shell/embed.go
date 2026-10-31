// Package shell exposes the canonical proj wrapper for release packaging.
package shell

import "embed"

//go:embed proj.zsh
var wrapper embed.FS

// Proj returns only the canonical wrapper, rooted at proj.zsh.
func Proj() embed.FS { return wrapper }
