package proj

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Transcripts locates agent conversation files by session id. Which store an
// id lives in is how proj knows the agent a session is actually running:
// @proj_agent only records what proj launched, and a pane can be switched to
// another agent afterwards.
type Transcripts struct {
	PiRoot     string // pi's session dir: <PiRoot>/<cwd-key>/<timestamp>_<id>.jsonl
	ClaudeRoot string // Claude Code's config dir: <ClaudeRoot>/projects/<cwd-key>/<id>.jsonl

	// idx caches the stores' ids after the first Locate. A pointer, so copies
	// of one value share a single read; nil (a literal value) reads per call.
	idx *transcriptIndex
}

// transcriptIndex maps each conversation id in the stores to "pi" or "claude".
type transcriptIndex struct {
	once sync.Once
	ids  map[string]string
}

// safeID is the id shape both announcers write. Anything else is refused
// before it reaches a glob pattern.
var safeID = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// DefaultTranscripts honours each agent's own override
// ($PI_CODING_AGENT_SESSION_DIR, $CLAUDE_CONFIG_DIR) and falls back to its
// default location under $HOME.
func DefaultTranscripts() Transcripts {
	home := os.Getenv("HOME")
	t := Transcripts{
		PiRoot:     os.Getenv("PI_CODING_AGENT_SESSION_DIR"),
		ClaudeRoot: os.Getenv("CLAUDE_CONFIG_DIR"),
	}
	if t.PiRoot == "" {
		t.PiRoot = filepath.Join(home, ".pi", "agent", "sessions")
	}
	if t.ClaudeRoot == "" {
		t.ClaudeRoot = filepath.Join(home, ".claude")
	}
	t.idx = &transcriptIndex{}
	return t
}

// Locate reports which agent's store holds the conversation id. The stores
// are read once per Transcripts value from DefaultTranscripts; files added
// after the first call are not seen.
func (t Transcripts) Locate(id string) (agent string, ok bool) {
	if !safeID.MatchString(id) {
		return "", false
	}
	var ids map[string]string
	if t.idx == nil {
		ids = t.readIndex()
	} else {
		t.idx.once.Do(func() { t.idx.ids = t.readIndex() })
		ids = t.idx.ids
	}
	agent, ok = ids[id]
	return agent, ok
}

// readIndex lists both stores once. pi is written last so it wins over Claude,
// as the per-id search it replaces checked pi first.
func (t Transcripts) readIndex() map[string]string {
	ids := map[string]string{}
	for _, dir := range subdirs(filepath.Join(t.ClaudeRoot, "projects")) {
		for _, name := range jsonlNames(dir) {
			ids[strings.TrimSuffix(name, ".jsonl")] = "claude"
		}
	}
	for _, dir := range subdirs(t.PiRoot) {
		for _, name := range jsonlNames(dir) {
			base := strings.TrimSuffix(name, ".jsonl")
			if i := strings.LastIndex(base, "_"); i >= 0 {
				ids[base[i+1:]] = "pi"
			}
		}
	}
	return ids
}

// subdirs returns root's immediate subdirectories, symlinks followed as the
// glob it replaces followed them, or nil when unreadable.
func subdirs(root string) []string {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		p := filepath.Join(root, e.Name())
		if e.Type()&os.ModeSymlink != 0 {
			if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
				continue
			}
		} else if !e.IsDir() {
			continue
		}
		out = append(out, p)
	}
	return out
}

// jsonlNames returns the names of dir's *.jsonl files.
func jsonlNames(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".jsonl") {
			out = append(out, e.Name())
		}
	}
	return out
}
