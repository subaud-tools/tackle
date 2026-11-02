// Package integrationsetup derives release integrations from the canonical
// sift, imprint and cull declarations embedded by their cmd entry points.
package integrationsetup

import (
	"bytes"
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing/fstest"

	"github.com/BurntSushi/toml"
	"github.com/subaud-tools/tackle/skills"
	"github.com/subaud-tools/tools-common/setupbundle"
	commonskills "github.com/subaud-tools/tools-common/skills"
)

// Source keeps the selected registrations and platform constraints, but leaves
// the owning binary download to the parent manifest. It never reads user state.
func Source(tool string, declaration []byte) (setupbundle.Source, error) {
	var manifest struct{ Packages map[string]map[string]any }
	if _, err := toml.Decode(string(declaration), &manifest); err != nil {
		return setupbundle.Source{}, err
	}
	p := manifest.Packages[tool]
	if p == nil {
		return setupbundle.Source{}, fmt.Errorf("setup declaration has no package %s", tool)
	}
	delete(p, "download")
	if notes, ok := p["notes"].([]any); ok {
		kept := make([]string, 0, len(notes))
		for _, note := range notes {
			text := note.(string)
			if !strings.Contains(text, "skills install") {
				kept = append(kept, text)
			}
		}
		delete(p, "notes")
		if len(kept) > 0 {
			p["notes"] = kept
		}
	}
	source := setupbundle.Source{}
	if tool != "cull" {
		path := tool + "/SKILL.md"
		body, err := skills.FS.ReadFile(path)
		if err != nil {
			return source, err
		}
		// The shared embed also holds the sibling skill. Export only this path,
		// retaining its prefix so managed-file sources match the archive inventory.
		source.Files = fstest.MapFS{path: &fstest.MapFile{Data: body, Mode: 0o644}}
		targets := commonskills.UserTargets()
		var files []map[string]string
		for _, agent := range slices.Sorted(maps.Keys(targets)) {
			files = append(files, map[string]string{"from": path, "to": targets[agent] + "/" + path, "mode": "0644"})
		}
		p["managed-file"] = files
	}
	var out bytes.Buffer
	err := toml.NewEncoder(&out).Encode(map[string]any{"kempt": map[string]int{"spec": 2}, "packages": map[string]any{"integration": p}})
	source.Manifest = out.Bytes()
	return source, err
}
