package proj

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"

	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/sessiontheme"
)

// ThemeGroup is a remembered window identity, independent of its native ID.
type ThemeGroup struct {
	Hue sessiontheme.Hue `json:"hue"`
}

func initThemeMaps(r *Record) {
	if r.Groups == nil {
		r.Groups = map[string]ThemeGroup{}
	}
	if r.NativeGroups == nil {
		r.NativeGroups = map[string]string{}
	}
}

func newThemeGroup(r *Record) string {
	initThemeMaps(r)
	id := rand.Text()
	hues := sessiontheme.Hue("").Values()
	r.Groups[id] = ThemeGroup{Hue: sessiontheme.Hue(hues[len(r.Groups)%len(hues)])}
	return id
}

// placeThemeSessions assigns only previously unassigned shades. Moves retain
// their slot, even when it collides with a session in the destination group.
func placeThemeSessions(r *Record, group string, names []string) {
	slots := sessiontheme.Shade("").Values()
	for _, name := range names {
		e, ok := r.Sessions[name]
		if !ok {
			continue
		}
		if e.ID == "" {
			e.ID = rand.Text()
		}
		if e.Shade == "" {
			used := map[sessiontheme.Shade]bool{}
			count := 0
			for _, other := range r.Sessions {
				if other.GroupID == group && other.Shade != "" {
					used[other.Shade] = true
					count++
				}
			}
			e.Shade = sessiontheme.Shade(slots[count%len(slots)])
			for _, slot := range slots {
				if !used[sessiontheme.Shade(slot)] {
					e.Shade = sessiontheme.Shade(slot)
					break
				}
			}
		}
		e.GroupID = group
		r.Sessions[name] = e
	}
}

// assignWindowGroup is the record-side seam for observed native windows. The
// observer must hold withThemeSyncLock across observation and this update.
func assignWindowGroup(r *Record, native string, names []string) {
	initThemeMaps(r)
	group := r.NativeGroups[native]
	if group == "" {
		group = newThemeGroup(r)
		r.NativeGroups[native] = group
	}
	placeThemeSessions(r, group, names)
}

// completeThemeRecord upgrades v1 and retains layout identity when the existing
// layout saver rebuilds its name lists. It never assigns a group to unplaced work.
func completeThemeRecord(r *Record) bool {
	changed := r.Version != recordVersion
	r.Version = recordVersion
	initThemeMaps(r)
	names := make([]string, 0, len(r.Sessions))
	for name := range r.Sessions {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		e := r.Sessions[name]
		if e.ID == "" {
			e.ID = rand.Text()
			r.Sessions[name] = e
			changed = true
		}
	}
	groups := make([]string, len(r.Layout.Windows))
	for i, win := range r.Layout.Windows {
		if i < len(r.Layout.GroupIDs) {
			groups[i] = r.Layout.GroupIDs[i]
		}
		if groups[i] == "" {
			for _, name := range win {
				if id := r.Sessions[name].GroupID; id != "" {
					groups[i] = id
					break
				}
			}
		}
		if groups[i] == "" {
			groups[i] = newThemeGroup(r)
			changed = true
		}
		var unassigned []string
		for _, name := range win {
			e, ok := r.Sessions[name]
			if ok && e.GroupID == "" {
				unassigned = append(unassigned, name)
			}
		}
		if len(unassigned) > 0 {
			changed = true
			placeThemeSessions(r, groups[i], unassigned)
		}
	}
	if len(groups) != len(r.Layout.GroupIDs) {
		changed = true
	}
	r.Layout.GroupIDs = groups
	return changed
}

// withThemeSyncLock serializes restore window creation with theme reconciliation.
// Lock order is theme-sync then UpdateRecord; never acquire it inside UpdateRecord.
func withThemeSyncLock(fn func() error) error {
	p := filepath.Join(tools.StateDir("proj"), "theme-sync.lock")
	if err := tools.EnsureDir(filepath.Dir(p)); err != nil {
		return err
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN) //nolint:errcheck
	return fn()
}

// ThemeForAgent is read-only: even a v1 record is not upgraded by a query.
func ThemeForAgent(agentSessionID string) (sessiontheme.Query, error) {
	r, err := readRecord()
	if err != nil {
		return sessiontheme.Query{}, err
	}
	if agentSessionID == "" {
		return sessiontheme.Query{}, nil
	}
	for _, e := range r.Sessions {
		if e.Conversation != agentSessionID || e.ID == "" || e.GroupID == "" || e.Shade == "" {
			continue
		}
		group, ok := r.Groups[e.GroupID]
		if !ok {
			return sessiontheme.Query{}, fmt.Errorf("proj: missing theme group %s", e.GroupID)
		}
		a := sessiontheme.Assignment{Hue: group.Hue, Shade: e.Shade}
		p, err := sessiontheme.Resolve(a)
		if err != nil {
			return sessiontheme.Query{}, err
		}
		return sessiontheme.Query{Assigned: true, Theme: &sessiontheme.Resolved{SessionID: e.ID, GroupID: e.GroupID, Assignment: a, Palette: p, Header: sessiontheme.HeaderFor(p)}}, nil
	}
	return sessiontheme.Query{}, nil
}

// SetWindowHue edits a known native binding; it does not observe or allocate.
func SetWindowHue(nativeWindowID string, hue sessiontheme.Hue) error {
	if _, err := sessiontheme.Resolve(sessiontheme.Assignment{Hue: hue, Shade: "midnight"}); err != nil {
		return err
	}
	return withThemeSyncLock(func() error {
		var found bool
		err := UpdateRecord(func(r *Record) bool {
			group, ok := r.NativeGroups[nativeWindowID]
			if !ok {
				return false
			}
			old, ok := r.Groups[group]
			if !ok {
				return false
			}
			found = true
			if old.Hue == hue {
				return false
			}
			r.Groups[group] = ThemeGroup{Hue: hue}
			return true
		})
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("proj: unknown window %q", nativeWindowID)
		}
		return nil
	})
}
