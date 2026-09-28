// Package mutate defines language-neutral mutation sites (a byte range in a
// file plus its replacement) and Apply/Sample helpers, plus a per-language
// Mutator interface for finding them.
package mutate

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sort"
)

// Site is one candidate mutation: replace src[Start:End] with Replacement.
// Start and End are UTF-8 byte offsets into the file. Desc is a short,
// human-readable description like "line 12: flip ==".
type Site struct {
	Start       int    `json:"start"`
	End         int    `json:"end"`
	Replacement string `json:"replacement"`
	Desc        string `json:"desc"`
}

// Apply returns src with s applied: everything outside [s.Start, s.End) is
// byte-for-byte unchanged. It errors if the site's range is not a valid
// span within src.
func Apply(src []byte, s Site) ([]byte, error) {
	if s.Start < 0 || s.End < s.Start || s.End > len(src) {
		return nil, fmt.Errorf("mutate: site [%d,%d) out of range for %d-byte source", s.Start, s.End, len(src))
	}
	out := make([]byte, 0, len(src)-(s.End-s.Start)+len(s.Replacement))
	out = append(out, src[:s.Start]...)
	out = append(out, []byte(s.Replacement)...)
	out = append(out, src[s.End:]...)
	return out, nil
}

// Sample returns at most max sites from sites, sorted by Start. The
// selection is deterministic given seed: the same sites and seed always
// pick the same subset. When len(sites) <= max, all sites are returned.
func Sample(sites []Site, max int, seed []byte) []Site {
	sorted := append([]Site(nil), sites...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Start < sorted[j].Start })
	if max < 0 {
		max = 0
	}
	if len(sorted) <= max {
		return sorted
	}

	h := sha256.Sum256(seed)
	rng := rand.New(rand.NewSource(int64(binary.BigEndian.Uint64(h[:8]))))
	idx := rng.Perm(len(sorted))[:max]
	sort.Ints(idx)

	out := make([]Site, 0, max)
	for _, i := range idx {
		out = append(out, sorted[i])
	}
	return out
}

// Mutator finds mutation Sites in one language's source files.
type Mutator interface {
	// Lang is the language this Mutator handles: "go", "python", or
	// "typescript".
	Lang() string
	// Sites returns the candidate mutation sites in src, a file at
	// relpath (project-root-relative, '/'-separated) under root. It
	// returns no sites (and no error) for files this Mutator must never
	// mutate, such as tests or generated code.
	Sites(root, relpath string, src []byte) ([]Site, error)
}

// For returns the Mutator for lang ("go", "python", "typescript"), or nil
// if lang is unknown or not yet implemented.
func For(lang string) Mutator {
	switch lang {
	case "go":
		return goMutator{}
	default:
		return nil
	}
}
