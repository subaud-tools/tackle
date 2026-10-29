package proj

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/sessiontheme"
)

// PiThemeName follows the remembered identity, never its display name or hue.
func PiThemeName(sessionID string) string { return "proj-session-" + sessionID }

// piAgentDir mirrors pi's PI_CODING_AGENT_DIR override and home shorthand.
func piAgentDir() (string, error) {
	dir := os.Getenv("PI_CODING_AGENT_DIR")
	if dir != "" && dir != "~" && !strings.HasPrefix(dir, "~/") {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if dir == "~" {
		return home, nil
	}
	if strings.HasPrefix(dir, "~/") {
		return filepath.Join(home, dir[2:]), nil
	}
	return filepath.Join(home, ".pi", "agent"), nil
}

// PreparePiTheme writes only when enabled. Detached work gets blue/deep until
// its first real attachment assigns a group and an unused shade in the record.
func PreparePiTheme(entry Entry) (string, error) {
	if !SessionThemesEnabled() {
		return "", nil
	}
	r, err := readRecord()
	if err != nil {
		return "", err
	}
	return preparePiTheme(entry, r)
}

func preparePiTheme(entry Entry, r Record) (string, error) {
	if entry.ID == "" || strings.ContainsAny(entry.ID, `/\\`) || entry.ID == "." || entry.ID == ".." {
		return "", fmt.Errorf("proj: invalid pi theme session identity %q", entry.ID)
	}
	a := sessiontheme.Assignment{Hue: "blue", Shade: "deep"}
	if entry.GroupID != "" {
		group, ok := r.Groups[entry.GroupID]
		if !ok {
			return "", fmt.Errorf("proj: missing theme group %s", entry.GroupID)
		}
		a = sessiontheme.Assignment{Hue: group.Hue, Shade: entry.Shade}
	}
	palette, err := sessiontheme.Resolve(a)
	if err != nil {
		return "", err
	}
	name := PiThemeName(entry.ID)
	b, err := sessiontheme.PiTheme(name, palette)
	if err != nil {
		return "", err
	}
	dir, err := piAgentDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "themes", name+".json")
	old, err := os.ReadFile(path)
	if err == nil && bytes.Equal(old, b) {
		return name, nil
	}
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if err := tools.WriteFileAtomic(path, b, 0600); err != nil {
		return "", err
	}
	return name, nil
}

// Remember before preparing a first launch, letting UpdateRecord allocate the
// same stable identity that snapshots and reconciliation will subsequently use.
func prepareLaunchPiTheme(socket, name, dir string) (Entry, string, error) {
	err := UpdateRecord(func(r *Record) bool {
		if _, ok := r.Sessions[name]; ok {
			return false
		}
		r.Sessions[name] = Entry{Socket: socket, Project: ProjectFromSocket(socket), Dir: dir, Agent: "pi"}
		return true
	})
	if err != nil {
		return Entry{}, "", err
	}
	r, err := readRecord()
	if err != nil {
		return Entry{}, "", err
	}
	e := r.Sessions[name]
	theme, err := PreparePiTheme(e)
	return e, theme, err
}
