package proj

import (
	"os"
	"path/filepath"
	"testing"
)

func touch(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// transcriptTree lays out one pi and one Claude Code transcript the way each
// agent stores them, and returns roots pointing at it.
func transcriptTree(t *testing.T) Transcripts {
	t.Helper()
	root := t.TempDir()
	tr := Transcripts{PiRoot: filepath.Join(root, "pi"), ClaudeRoot: filepath.Join(root, "claude")}
	touch(t, filepath.Join(tr.PiRoot, "--x--", "2026-09-29T14-03-27-902Z_01a0ed7a.jsonl"))
	touch(t, filepath.Join(tr.ClaudeRoot, "projects", "-x", "78834978.jsonl"))
	return tr
}

func TestLocate(t *testing.T) {
	tr := transcriptTree(t)
	cases := []struct {
		id, agent string
		ok        bool
	}{
		{"01a0ed7a", "pi", true},
		{"78834978", "claude", true},
		{"nope", "", false},
	}
	for _, c := range cases {
		agent, ok := tr.Locate(c.id)
		if agent != c.agent || ok != c.ok {
			t.Errorf("Locate(%q) = (%q, %v) want (%q, %v)", c.id, agent, ok, c.agent, c.ok)
		}
	}
}

func TestLocateRejectsUnsafeID(t *testing.T) {
	tr := transcriptTree(t)
	// A file literally named "*" must not make a glob id resolve.
	touch(t, filepath.Join(tr.ClaudeRoot, "projects", "-x", "*.jsonl"))
	for _, id := range []string{"*", "../x", "", "a b", "0?a"} {
		if agent, ok := tr.Locate(id); ok {
			t.Errorf("Locate(%q) = (%q, true); want rejected", id, agent)
		}
	}
}

func TestDefaultTranscriptsEnv(t *testing.T) {
	t.Setenv("PI_CODING_AGENT_SESSION_DIR", "/p")
	t.Setenv("CLAUDE_CONFIG_DIR", "/c")
	roots := func(t Transcripts) Transcripts { return Transcripts{PiRoot: t.PiRoot, ClaudeRoot: t.ClaudeRoot} }
	if got := roots(DefaultTranscripts()); got != (Transcripts{PiRoot: "/p", ClaudeRoot: "/c"}) {
		t.Fatalf("with env: %+v", got)
	}
	t.Setenv("PI_CODING_AGENT_SESSION_DIR", "")
	t.Setenv("CLAUDE_CONFIG_DIR", "")
	t.Setenv("HOME", "/h")
	want := Transcripts{PiRoot: "/h/.pi/agent/sessions", ClaudeRoot: "/h/.claude"}
	if got := roots(DefaultTranscripts()); got != want {
		t.Fatalf("defaults: %+v want %+v", got, want)
	}
}

func TestTranscriptIndexMatchesSearch(t *testing.T) {
	root := t.TempDir()
	tr := Transcripts{PiRoot: filepath.Join(root, "pi"), ClaudeRoot: filepath.Join(root, "claude")}
	touch(t, filepath.Join(tr.PiRoot, "a1", "20261025_x1.jsonl"))
	touch(t, filepath.Join(tr.ClaudeRoot, "projects", "k", "x2.jsonl"))
	touch(t, filepath.Join(tr.PiRoot, "a1", "20261025_x3.jsonl"))
	touch(t, filepath.Join(tr.ClaudeRoot, "projects", "k", "x3.jsonl"))
	cases := []struct {
		id, agent string
		ok        bool
	}{
		{"x1", "pi", true},
		{"x2", "claude", true},
		{"x3", "pi", true},
		{"missing", "", false},
		{"../etc", "", false},
	}
	for _, c := range cases {
		agent, ok := tr.Locate(c.id)
		if agent != c.agent || ok != c.ok {
			t.Errorf("Locate(%q) = (%q, %v) want (%q, %v)", c.id, agent, ok, c.agent, c.ok)
		}
	}
}

func TestTranscriptIndexReadsOnce(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PI_CODING_AGENT_SESSION_DIR", filepath.Join(root, "pi"))
	t.Setenv("CLAUDE_CONFIG_DIR", filepath.Join(root, "claude"))
	tr := DefaultTranscripts()
	touch(t, filepath.Join(tr.PiRoot, "a1", "20261025_x1.jsonl"))
	if _, ok := tr.Locate("x1"); !ok {
		t.Fatal("Locate(x1) not found")
	}
	touch(t, filepath.Join(tr.PiRoot, "a1", "20261025_late.jsonl"))
	if agent, ok := tr.Locate("late"); ok {
		t.Fatalf("Locate(late) = (%q, true) after the index was built; want the cached read", agent)
	}
}
