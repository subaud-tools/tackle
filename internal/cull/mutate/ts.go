package mutate

import "github.com/schuettc/tackle/internal/cull/extract/ts"

// tsMutator finds mutation sites in TypeScript source by running the
// extract/ts package's embedded helper in --sites mode. It never imports
// the extract package's own Extractor machinery; ts.Sites is a standalone
// entry point that talks to the same embedded tsext.cjs.
type tsMutator struct{}

func (tsMutator) Lang() string { return "typescript" }

// Sites returns no sites (and no error) if node is not on PATH, or no
// typescript package can be found: a missing toolchain makes the file
// unmeasurable, not a bug.
func (tsMutator) Sites(root, relpath string, src []byte) ([]Site, error) {
	found, err := ts.Sites(root, relpath, src)
	if err != nil {
		return nil, err
	}
	if len(found) == 0 {
		return nil, nil
	}
	out := make([]Site, len(found))
	for i, s := range found {
		out[i] = Site{Start: s.Start, End: s.End, Replacement: s.Replacement, Desc: s.Desc}
	}
	return out, nil
}
