package projcli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/subaud-tools/tools-common/setupbundle"
)

func TestSetupBundleCLIAndWrapper(t *testing.T) {
	var out, errw bytes.Buffer
	if code := Dispatch([]string{"commands", "--json"}, &out, &errw); code != 0 || !strings.Contains(out.String(), "setup-bundle") {
		t.Fatalf("build-only command missing: exit %d: %s %s", code, &out, &errw)
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	dest := filepath.Join(t.TempDir(), "proj_setup.tar.gz")
	out.Reset()
	if code := Dispatch([]string{"setup-bundle", "--version", "1.2.3", "--out", dest}, &out, &errw); code != 0 {
		t.Fatalf("exit %d: %s", code, &errw)
	}
	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	b, err := setupbundle.Read(data)
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := os.ReadFile("../../shell/proj.zsh")
	if err != nil {
		t.Fatal(err)
	}
	if b.Tool != "proj" || b.Version != "1.2.3" || len(b.Files) != 1 || !bytes.Equal(b.Files["proj.zsh"], canonical) {
		t.Fatalf("wrong bundle: %+v", b)
	}
	var manifest struct {
		Kempt    struct{ Spec int }
		Packages map[string]struct {
			Managed []struct{ From, To, Mode string } `toml:"managed-file"`
			Lines   []struct{ File, Line string }     `toml:"line-in-file"`
		}
	}
	md, err := toml.Decode(string(b.Manifest), &manifest)
	if err != nil {
		t.Fatal(err)
	}
	p := manifest.Packages["integration"]
	if len(md.Undecoded()) != 0 || manifest.Kempt.Spec != 2 || len(manifest.Packages) != 1 || len(p.Managed) != 1 || len(p.Lines) != 1 {
		t.Fatalf("unexpected declaration: %s", b.Manifest)
	}
	if p.Managed[0].From != "proj.zsh" || p.Managed[0].To != "~/.local/share/proj/proj.zsh" || p.Managed[0].Mode != "0644" || p.Lines[0].File != "~/.zshrc" || p.Lines[0].Line != `source "$HOME/.local/share/proj/proj.zsh"` {
		t.Fatalf("wrong targets: %s", b.Manifest)
	}
	entries, err := os.ReadDir(home)
	if err != nil || len(entries) != 0 {
		t.Fatalf("build mutated HOME: %v %v", entries, err)
	}
	wrapper := filepath.Join(home, ".local/share/proj/proj.zsh")
	if err := os.MkdirAll(filepath.Dir(wrapper), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wrapper, b.Files["proj.zsh"], 0o644); err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	stub := "#!/bin/sh\nprintf '%s\\n' 'printf x >> \"$HOME/executed\"' > \"$PROJ_EXEC\"\n"
	if err := os.WriteFile(filepath.Join(bin, "proj"), []byte(stub), 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("zsh", "-f", "-c", p.Lines[0].Line+"\nproj")
	cmd.Env = append(os.Environ(), "PATH="+bin+":"+os.Getenv("PATH"), "NO_AUTO_TMUX=1", "TMPDIR="+t.TempDir())
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("isolated wrapper: %v: %s", err, output)
	}
	executed, err := os.ReadFile(filepath.Join(home, "executed"))
	if err != nil || string(executed) != "x" {
		t.Fatalf("PROJ_EXEC did not execute exactly once: %q %v", executed, err)
	}
}
