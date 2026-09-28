package evidence

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/schuettc/tackle/internal/cull/cases"
	"github.com/schuettc/tackle/internal/cull/extract"
	"github.com/schuettc/tackle/internal/cull/mutate"
	"github.com/schuettc/tackle/internal/cull/verify"
	"github.com/schuettc/tackle/internal/version"
	tools "github.com/schuettc/tools-common"
)

// cacheFormat changes whenever a cache entry's meaning does.
const cacheFormat = "cull-evidence-cache-1"

// resultCache holds each module's matrix results under
// <sandbox base>/cache/<key>.json (0600, written atomically), keyed by
// everything that can change them (see keyFor). Only clean results are
// stored (see cacheable): anything a hit can't prove identical misses.
type resultCache struct {
	dir   string
	root  string
	files map[string]string // copied project file → sha256 (see fingerprint)
	stamp []string          // cull build and environment
	opt   Options
}

// cacheEntry is one module's stored results.
type cacheEntry struct {
	Key      string   `json:"key"`
	Module   string   `json:"module"`
	Shadowed bool     `json:"shadowed"`
	Measured []string `json:"measured"`
	Unloaded []string `json:"unloaded"`
	Unmapped int      `json:"unmapped"`
	Mutants  []Mutant `json:"mutants"`
}

// cacheKey is what keyFor hashes.
type cacheKey struct {
	Format                      string
	Stamp                       []string
	Module, Lang                string
	Tests, TestFiles            []string
	Callees                     []cases.Callee
	Sites                       []mutate.Site
	Runners                     []string
	Files                       map[string]string
	MaxMutants                  int
	RunTimeout, BaselineTimeout time.Duration
	NoCoverage                  bool
}

// newResultCache fingerprints the copies in w0 (the first sandbox, what
// the runs see) of the copied files of the project at root.
func newResultCache(root, w0, base string, copied []string, opt Options) *resultCache {
	return &resultCache{dir: filepath.Join(base, "cache"), root: root, files: fingerprint(w0, copied), stamp: buildStamp(root), opt: opt}
}

// fingerprint hashes every copied project file under root, except .cull/ (evidence's
// own output) and Python bytecode.
func fingerprint(root string, copied []string) map[string]string {
	out := map[string]string{}
	for _, rel := range copied {
		if strings.HasPrefix(rel, ".cull/") || strings.HasSuffix(rel, ".pyc") || strings.Contains("/"+rel, "/__pycache__/") {
			continue
		}
		p := filepath.Join(root, filepath.FromSlash(rel))
		var data []byte
		if fi, err := os.Lstat(p); err == nil && fi.Mode()&os.ModeSymlink != 0 {
			t, _ := os.Readlink(p)
			data = []byte("symlink:" + t)
		} else if b, err := os.ReadFile(p); err == nil {
			data = b
		} else {
			data = []byte("unreadable:" + err.Error())
		}
		sum := sha256.Sum256(data)
		out[rel] = hex.EncodeToString(sum[:])
	}
	return out
}

var exeHash = sync.OnceValue(func() string {
	exe, err := os.Executable()
	if err != nil {
		return "unknown:" + err.Error()
	}
	b, err := os.ReadFile(exe)
	if err != nil {
		return "unknown:" + err.Error()
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
})

// buildStamp is cull's version and build, the environment that picks the
// test toolchains, and their version output.
func buildStamp(root string) []string {
	py := verify.PythonInterpreter(root)
	s := []string{version.Number(), version.Commit(), exeHash(), runtime.GOOS + "/" + runtime.GOARCH,
		"python=" + py, "go-version=" + toolVersion("go", "version"), "node-version=" + toolVersion("node", "--version"),
		"python-version=" + toolVersion(py, "--version")}
	for _, k := range []string{"PATH", "VIRTUAL_ENV", "PYTHONPATH", "GOFLAGS", "GOTOOLCHAIN", "CGO_ENABLED", "NODE_OPTIONS", "CULL_TS"} {
		s = append(s, k+"="+os.Getenv(k))
	}
	return s
}

// toolVersion is argv's trimmed output, or why it failed.
func toolVersion(argv ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Env = verify.FilteredEnv()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "unavailable: " + err.Error()
	}
	return strings.TrimSpace(string(out))
}

// keyFor is p's cache key: the sha256 of its module's bytes, the module's
// test files and tests, the callees and the sampled sites, the runners,
// the relevant options, cull's build and environment, and every other
// project file except the test files of other modules (a Go test file in
// the same package as one of p's test files counts: go test builds it).
// Any source or config change anywhere misses: p's tests may import it.
func (c *resultCache) keyFor(p *modPlan) string {
	k := cacheKey{Format: cacheFormat, Stamp: c.stamp, Module: p.mod.Rel, Lang: p.mod.Lang,
		Tests: p.mod.Tests, TestFiles: p.mod.TestFiles, Callees: p.mod.Callees, Sites: p.sites,
		Files: map[string]string{}, MaxMutants: c.opt.MaxMutants, RunTimeout: c.opt.RunTimeout,
		BaselineTimeout: c.opt.BaselineTimeout, NoCoverage: c.opt.NoCoverage}
	srcSum := sha256.Sum256(p.src)
	k.Files["\x00module"] = hex.EncodeToString(srcSum[:])
	for _, g := range p.groups {
		k.Runners = append(k.Runners, fmt.Sprintf("%#v %v", g.r, g.files))
	}
	testDirs := map[string]bool{}
	for _, f := range p.mod.TestFiles {
		testDirs[path.Dir(f)] = true
	}
	for rel, sum := range c.files {
		if extract.ForFile(rel) != nil && !slices.Contains(p.mod.TestFiles, rel) &&
			(!strings.HasSuffix(rel, "_test.go") || !testDirs[path.Dir(rel)]) {
			continue
		}
		k.Files[rel] = sum
	}
	for _, f := range p.mod.TestFiles {
		if _, ok := k.Files[f]; !ok {
			k.Files[f] = "missing"
		}
	}
	data, _ := json.Marshal(k)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// get returns the entry stored under key, if any and intact.
func (c *resultCache) get(key, module string) (cacheEntry, bool) {
	b, err := os.ReadFile(filepath.Join(c.dir, key+".json"))
	if err != nil {
		return cacheEntry{}, false
	}
	var e cacheEntry
	if json.Unmarshal(b, &e) != nil || e.Key != key || e.Module != module {
		return cacheEntry{}, false
	}
	for _, mu := range e.Mutants {
		if mu.KilledBy == nil || mu.Module != module {
			return cacheEntry{}, false
		}
	}
	return e, true
}

// put stores e under its key (0600, atomically).
func (c *resultCache) put(e cacheEntry) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	if err := tools.EnsureDir(c.dir); err != nil {
		return err
	}
	return tools.WriteFileAtomic(filepath.Join(c.dir, e.Key+".json"), data, 0o600)
}

// cacheable reports whether p's results are clean enough to reuse: a
// valid baseline and canary, no dropped file, no test failing any
// baseline run (flaky ones might measure next time), and no excluded
// mutant (timeouts and build failures are retried).
func cacheable(p *modPlan, muts []Mutant) bool {
	if !p.ok || p.unsure || len(p.dropped) > 0 || len(p.failed) > 0 {
		return false
	}
	for _, mu := range muts {
		if mu.KilledBy == nil {
			return false
		}
	}
	return true
}

// entry is p's results as a cache entry.
func entry(key string, p *modPlan, muts []Mutant) cacheEntry {
	return cacheEntry{Key: key, Module: p.mod.Rel, Shadowed: p.shadowed, Measured: sortedKeys(p.measured),
		Unloaded: p.unloaded, Unmapped: p.unmapped, Mutants: muts}
}

// restore fills p from a cache entry: its baseline, canary and coverage
// are not run again.
func (p *modPlan) restore(e cacheEntry) {
	p.cached = true
	p.ok = true
	p.shadowed = e.Shadowed
	p.measured = map[string]bool{}
	for _, id := range e.Measured {
		p.measured[id] = true
	}
	p.failed = map[string]bool{}
	p.unloaded = e.Unloaded
	p.unmapped = e.Unmapped
	p.cachedMutants = e.Mutants
}
