package integrationsetup_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/subaud-tools/tools-common/setupbundle"
	commonskills "github.com/subaud-tools/tools-common/skills"
)

func TestReleaseWorkflowSetupSelection(t *testing.T) {
	data, err := os.ReadFile("../../.github/workflows/release.yml")
	if err != nil {
		t.Fatal(err)
	}
	workflow := string(data)
	const declaration = `SETUP_BUNDLE_TOOLS: '["proj", "sift", "imprint", "cull"]'`
	if strings.Count(workflow, declaration) != 1 {
		t.Fatal("declare the setup-producing tools once")
	}
	if strings.Count(workflow, "contains(fromJSON(env.SETUP_BUNDLE_TOOLS), steps.rel.outputs.tool)") != 3 {
		t.Fatal("setup-go, producer and extra-assets must use the same tool selection")
	}
}

// Exercise the actual cmd adapters: an internal-only registration would not
// publish the declaration embedded by the release's entry point.
func TestReleaseBundles(t *testing.T) {
	for _, tool := range []string{"sift", "imprint", "cull"} {
		t.Run(tool, func(t *testing.T) {
			bin := filepath.Join(t.TempDir(), tool)
			build := exec.Command("go", "build", "-o", bin, "../../cmd/"+tool)
			if output, err := build.CombinedOutput(); err != nil {
				t.Fatalf("build: %v: %s", err, output)
			}
			home := t.TempDir()
			env := append(os.Environ(), "HOME="+home, "XDG_CONFIG_HOME="+home+"/.config", "XDG_STATE_HOME="+home+"/.state", "XDG_DATA_HOME="+home+"/.data", "XDG_CACHE_HOME="+home+"/.cache", "PATH=/usr/bin:/bin", "CLAUDE_CONFIG_DIR=", "CODEX_HOME=", "PI_CODING_AGENT_DIR=")
			dest := filepath.Join(t.TempDir(), tool+"_setup.tar.gz")
			run := func(args ...string) {
				t.Helper()
				cmd := exec.Command(bin, args...)
				cmd.Env = env
				cmd.Dir = home
				if output, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("%v: %v: %s", args, err, output)
				}
			}
			run("setup-bundle", "--version", "1.2.3-rc.4", "--out", dest)
			data, err := os.ReadFile(dest)
			if err != nil {
				t.Fatal(err)
			}
			b, err := setupbundle.Read(data)
			if err != nil {
				t.Fatal(err)
			}
			if b.Tool != tool || b.Version != "1.2.3-rc.4" {
				t.Fatalf("wrong identity: %+v", b)
			}
			var m struct {
				Kempt    struct{ Spec int }
				Packages map[string]map[string]any
			}
			if _, err := toml.Decode(string(b.Manifest), &m); err != nil {
				t.Fatal(err)
			}
			p := m.Packages["integration"]
			if m.Kempt.Spec != 2 || len(m.Packages) != 1 || p == nil {
				t.Fatalf("wrong namespace: %s", b.Manifest)
			}
			var canonical struct{ Packages map[string]map[string]any }
			if _, err := toml.DecodeFile("../../cmd/"+tool+"/kempt.toml", &canonical); err != nil {
				t.Fatal(err)
			}
			original := canonical.Packages[tool]
			delete(original, "download")
			delete(original, "notes")
			if tool == "sift" {
				original["notes"] = []any{"Then, once: `sift init --root DIR` (the directories to audit), and `sift doctor` to check the setup."}
			}
			managed := p["managed-file"]
			delete(p, "managed-file")
			if !reflect.DeepEqual(p, original) {
				t.Fatalf("registrations/constraints drifted: got %#v want %#v", p, original)
			}
			if tool == "cull" {
				if len(b.Files) != 0 || managed != nil {
					t.Fatalf("cull has invented assets: %+v", b)
				}
			} else {
				body, err := os.ReadFile("../../skills/" + tool + "/SKILL.md")
				if err != nil {
					t.Fatal(err)
				}
				if len(b.Files) != 1 || !bytes.Equal(b.Files[tool+"/SKILL.md"], body) {
					t.Fatal("bundle leaks a sibling or changes skill bytes")
				}
				rows, ok := managed.([]map[string]any)
				if !ok {
					t.Fatalf("missing managed files: %#v", managed)
				}
				targets := commonskills.UserTargets()
				want := make(map[string]bool)
				for _, dir := range targets {
					want[dir+"/"+tool+"/SKILL.md"] = true
				}
				if len(rows) != len(want) {
					t.Fatalf("wrong targets: %#v", rows)
				}
				for _, row := range rows {
					to, _ := row["to"].(string)
					if !want[to] || row["from"] != tool+"/SKILL.md" || row["mode"] != "0644" {
						t.Fatalf("wrong target: %#v", row)
					}
					delete(want, to)
				}
			}
			run("setup-bundle", "--version", "1.2.3-rc.4", "--out", dest)
			repeated, err := os.ReadFile(dest)
			if err != nil || !bytes.Equal(data, repeated) {
				t.Fatal("bundle is not deterministic")
			}
			entries, err := os.ReadDir(home)
			if err != nil || len(entries) != 0 {
				t.Fatalf("producer touched HOME: %v %v", entries, err)
			}
			if tool != "cull" {
				run("skills", "install", "--agent", "claude-code")
				installed, err := os.ReadFile(filepath.Join(home, ".claude/skills", tool, "SKILL.md"))
				if err != nil || !bytes.Equal(installed, b.Files[tool+"/SKILL.md"]) {
					t.Fatalf("direct install differs: %v", err)
				}
			}
		})
	}
}
