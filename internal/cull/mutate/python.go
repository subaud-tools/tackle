package mutate

import "github.com/schuettc/tackle/internal/cull/extract/python"

// pythonMutator finds mutation sites in Python source by running the
// extract/python package's embedded helper in --sites mode. It never
// imports the extract package's own Extractor machinery; python.Sites is
// a standalone entry point that talks to the same embedded pyext.py.
type pythonMutator struct{}

func (pythonMutator) Lang() string { return "python" }

// Sites returns no sites (and no error) if python3 is not on PATH: a
// missing interpreter makes the file unmeasurable, not a bug.
func (pythonMutator) Sites(root, relpath string, src []byte) ([]Site, error) {
	found, err := python.Sites(root, relpath, src)
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
