package evidence

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/schuettc/tackle/internal/cull/cases"
	"github.com/schuettc/tackle/internal/cull/mutate"
	"github.com/schuettc/tackle/internal/cull/runner"
)

// Module is one non-test source file to mutate, the test files that
// exercise it, and the cull ids of the tests in those files that exercise
// it. Rel is root-relative and '/'-separated; it must never be a test file
// (the Python and TypeScript mutators don't skip test files themselves).
// Callees are the definitions in Rel that the targets call: when set,
// mutants are planted only inside them (see targetSites).
type Module struct {
	Rel, Lang string
	TestFiles []string
	Tests     []string // cull ids of tests exercising it
	Callees   []cases.Callee
}

// perCallee is how many mutants targeted sampling plants per callee.
const perCallee = 15

// Mutant is one planted bug and the measured tests that failed with it in
// place. KilledBy is nil when the mutant is excluded (its run timed out,
// failed to build/collect, or couldn't start): it tells nothing about
// individual tests. A non-nil empty KilledBy means no measured test caught
// it.
type Mutant struct {
	Module   string      `json:"module"`
	Site     mutate.Site `json:"site"`
	KilledBy []string    `json:"killed_by"`
}

// Matrix is the kill matrix over every module run. Baseline holds the
// measured tests: those that passed both baseline runs of their module's
// test files, in a module where at least one mutant was planted, and
// failed no baseline run anywhere. Flaky lists (sorted) the tests that
// failed at least one baseline run: never measured. Unmapped sums the
// first baseline run's unmapped native results over all modules.
// MeasuredByModule holds, per module, the Baseline tests that passed both
// of THAT module's baseline runs: only those are credited (or charged)
// for that module's mutants. Shadowed lists (sorted) the modules whose
// canary no measured test noticed: the tests load another copy of them,
// so none of their mutants run and none of their tests are measured.
// Dropped maps a module to the test files (sorted) left out of it because
// they did not run cleanly on their own (see modPlan.isolate).
type Matrix struct {
	Mutants  []Mutant        `json:"mutants"`
	Baseline map[string]bool `json:"baseline"`
	Flaky    []string        `json:"flaky"`
	Unmapped int             `json:"unmapped"`

	MeasuredByModule map[string]map[string]bool `json:"measured_by_module"`
	Shadowed         []string                   `json:"shadowed"`
	Dropped          map[string][]string        `json:"dropped"`
}

// Options tunes RunMatrix and Run. Zero values take defaults: Workers = runtime.NumCPU(),
// MaxMutants = 30 (for a module without Callees), RunTimeout = 60s,
// BaselineTimeout = 10m, Stderr = discard. Baseline runs get
// BaselineTimeout; canary and mutant runs get the larger of RunTimeout and
// 3× the module's slower baseline run.
// Evidence always runs the runner package's detected per-test command,
// never a project's test_command override.
type Options struct {
	Workers, MaxMutants         int
	RunTimeout, BaselineTimeout time.Duration
	Stderr                      io.Writer
}

// group is a set of test files run together by one Runner, with an index
// whose native keys are unambiguous within the group.
type group struct {
	r     runner.Runner
	files []string
	index map[string]string
}

// modPlan is one module ready to measure.
type modPlan struct {
	mod    Module
	src    []byte
	sites  []mutate.Site
	groups []*group

	ok       bool            // baseline valid
	shadowed bool            // canary unnoticed: tests load another copy
	measured map[string]bool // passed both baseline runs
	failed   map[string]bool // failed at least one baseline run
	dropped  []string        // test files that did not run cleanly alone
	unmapped int
	slowest  time.Duration // the slower baseline run (all groups)
}

// runTimeout is the timeout for p's canary and mutant runs: the larger of
// floor and 3× its slower baseline run.
func (p *modPlan) runTimeout(floor time.Duration) time.Duration {
	return max(floor, 3*p.slowest)
}

// plantHook, if set (tests only), is called with the mutated file's path
// right after each mutant is planted.
var plantHook func(file string)

var dupIDRe = regexp.MustCompile(` #\d+$`)

var idPrefix = map[string]string{"go": "go:", "python": "py:", "typescript": "ts:"}

// RunMatrix measures mods in sandbox copies of root (see Prepare) and returns
// their kill matrix. Per module: the test files run twice unmutated (the
// baseline); then a canary (see canarySource) is planted and the tests run
// once, and a module whose canary no measured test notices is shadowed
// (see Matrix.Shadowed); then each sampled mutant (see targetSites; for a module
// without Callees, at most MaxMutants seeded by the module's bytes) is written into a free sandbox, the tests are run, and
// the module's original bytes are restored before the sandbox is reused.
// The working tree is never written. Run holds an exclusive lock on the
// project's sandbox root (<base>/.lock) throughout; a concurrent Run or
// Prepare on the same project fails fast with an error naming it.
func RunMatrix(ctx context.Context, root string, mods []Module, opt Options) (Matrix, error) {
	if opt.Workers <= 0 {
		opt.Workers = runtime.NumCPU()
	}
	if opt.MaxMutants <= 0 {
		opt.MaxMutants = 30
	}
	if opt.RunTimeout <= 0 {
		opt.RunTimeout = 60 * time.Second
	}
	if opt.BaselineTimeout <= 0 {
		opt.BaselineTimeout = 10 * time.Minute
	}
	if opt.Stderr == nil {
		opt.Stderr = io.Discard
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return Matrix{}, err
	}
	var logMu sync.Mutex
	logf := func(format string, args ...any) {
		logMu.Lock()
		defer logMu.Unlock()
		fmt.Fprintf(opt.Stderr, "cull evidence: "+format+"\n", args...)
	}
	base, unlock, err := lockBase(abs)
	if err != nil {
		return Matrix{}, err
	}
	defer unlock()
	sbs, err := prepare(abs, base, opt.Workers, logf)
	if err != nil {
		return Matrix{}, err
	}

	var plans []*modPlan
	for _, mod := range mods {
		if p := planModule(abs, sbs[0].Dir, mod, opt.MaxMutants, logf); p != nil {
			plans = append(plans, p)
		}
	}

	pool := make(chan Sandbox, len(sbs))
	for _, sb := range sbs {
		pool <- sb
	}

	// Phase 1: baselines, one job per module.
	jobs := make([]func(Sandbox) error, 0, len(plans))
	for _, p := range plans {
		jobs = append(jobs, func(sb Sandbox) error {
			p.baseline(ctx, sb, opt.BaselineTimeout)
			if len(p.dropped) > 0 {
				logf("%s: dropped test files that did not run cleanly alone: %s", p.mod.Rel, strings.Join(p.dropped, ", "))
			}
			logf("baseline %s: %d measured, %d failing", p.mod.Rel, len(p.measured), len(p.failed))
			if !p.ok || len(p.measured) == 0 {
				return nil
			}
			loaded, why, err := p.canary(ctx, sb, p.runTimeout(opt.RunTimeout))
			if err != nil {
				return err
			}
			p.shadowed = !loaded
			if loaded {
				logf("canary %s: loaded (%s)", p.mod.Rel, why)
			} else {
				logf("canary %s: not loaded; %s", p.mod.Rel, reasonShadowed)
			}
			return nil
		})
	}
	if err := runJobs(ctx, pool, opt.Workers, jobs); err != nil {
		return Matrix{}, err
	}

	// Phase 2: mutants, one job per (module, site), results by index.
	var mutants []Mutant
	jobs = jobs[:0]
	for _, p := range plans {
		if !p.ok {
			logf("%s: baseline invalid; its tests are unmeasurable", p.mod.Rel)
			continue
		}
		if len(p.measured) == 0 || p.shadowed {
			continue
		}
		for i, s := range p.sites {
			idx := len(mutants)
			mutants = append(mutants, Mutant{Module: p.mod.Rel, Site: s})
			n := i + 1
			jobs = append(jobs, func(sb Sandbox) error {
				killed, why, err := p.mutant(ctx, sb, s, p.runTimeout(opt.RunTimeout))
				if err != nil {
					return err
				}
				mutants[idx].KilledBy = killed
				if killed == nil {
					logf("%s mutant %d/%d (%s): excluded (%s)", p.mod.Rel, n, len(p.sites), s.Desc, why)
				} else {
					logf("%s mutant %d/%d (%s): killed by %d", p.mod.Rel, n, len(p.sites), s.Desc, len(killed))
				}
				return nil
			})
		}
	}
	if err := runJobs(ctx, pool, opt.Workers, jobs); err != nil {
		return Matrix{}, err
	}

	m := Matrix{Mutants: mutants, Baseline: map[string]bool{}, MeasuredByModule: map[string]map[string]bool{}, Dropped: map[string][]string{}}
	flaky := map[string]bool{}
	for _, p := range plans {
		m.Unmapped += p.unmapped
		for id := range p.failed {
			flaky[id] = true
		}
	}
	for _, p := range plans {
		if p.shadowed {
			m.Shadowed = append(m.Shadowed, p.mod.Rel)
		}
		if len(p.dropped) > 0 {
			m.Dropped[p.mod.Rel] = p.dropped
		}
		if !p.ok || p.shadowed || len(p.sites) == 0 {
			continue
		}
		for id := range p.measured {
			if !flaky[id] {
				m.Baseline[id] = true
				if m.MeasuredByModule[p.mod.Rel] == nil {
					m.MeasuredByModule[p.mod.Rel] = map[string]bool{}
				}
				m.MeasuredByModule[p.mod.Rel][id] = true
			}
		}
	}
	for id := range flaky {
		m.Flaky = append(m.Flaky, id)
	}
	sort.Strings(m.Flaky)
	sort.Strings(m.Shadowed)
	// measured is per module, flakiness global: a test stable in one
	// module's baselines may have failed in another's. Only tests in the
	// final Baseline are credited.
	for i := range m.Mutants {
		if m.Mutants[i].KilledBy == nil {
			continue
		}
		kept := []string{}
		for _, id := range m.Mutants[i].KilledBy {
			if m.Baseline[id] {
				kept = append(kept, id)
			}
		}
		m.Mutants[i].KilledBy = kept
	}
	return m, nil
}

// runJobs runs jobs on at most workers goroutines, each job holding one
// sandbox from pool for its duration. The first job error (or ctx
// cancellation) stops further jobs and is returned.
func runJobs(ctx context.Context, pool chan Sandbox, workers int, jobs []func(Sandbox) error) error {
	next := make(chan func(Sandbox) error)
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr error
	)
	fail := func(err error) {
		mu.Lock()
		if firstErr == nil {
			firstErr = err
		}
		mu.Unlock()
	}
	failed := func() bool {
		mu.Lock()
		defer mu.Unlock()
		return firstErr != nil
	}
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range next {
				if failed() {
					continue
				}
				sb := <-pool
				err := job(sb)
				pool <- sb
				if err != nil {
					fail(err)
				}
			}
		}()
	}
	for _, job := range jobs {
		if ctx.Err() != nil || failed() {
			break
		}
		next <- job
	}
	close(next)
	wg.Wait()
	if firstErr != nil {
		return firstErr
	}
	return ctx.Err()
}

// planModule reads mod's source, samples its sites, and groups its test
// files by Runner. It returns nil (with a progress line) for a module that
// can't be measured at all.
func planModule(root, sbDir string, mod Module, max int, logf func(string, ...any)) *modPlan {
	if !copiedFile(sbDir, mod.Rel) {
		logf("%s: not a copied project file; skipped", mod.Rel)
		return nil
	}
	mut := mutate.For(mod.Lang)
	if mut == nil {
		logf("%s: no mutator for %q; skipped", mod.Rel, mod.Lang)
		return nil
	}
	src, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(mod.Rel)))
	if err != nil {
		logf("%s: %v; skipped", mod.Rel, err)
		return nil
	}
	sites, err := mut.Sites(root, mod.Rel, src)
	if err != nil {
		logf("%s: %v; skipped", mod.Rel, err)
		return nil
	}
	if len(mod.Callees) > 0 {
		sites = targetSites(src, sites, mod.Callees)
	} else {
		seed := sha256.Sum256(src)
		sites = mutate.Sample(sites, max, seed[:])
	}
	if len(sites) == 0 {
		logf("%s: no mutation sites; skipped", mod.Rel)
		return nil
	}
	p := &modPlan{mod: mod, src: src, sites: sites}

	byKey := map[string][]*group{}
	var keys []string
	for _, f := range mod.TestFiles {
		idx := nativeIndex(mod.Lang, f, mod.Tests)
		if len(idx) == 0 {
			continue
		}
		r, err := runner.For(root, mod.Lang, f)
		if err != nil {
			logf("%s: %v; its tests are unmeasurable", f, err)
			continue
		}
		key := fmt.Sprintf("%#v", r)
		if _, seen := byKey[key]; !seen {
			keys = append(keys, key)
		}
		var g *group
		for _, cand := range byKey[key] {
			if !overlaps(cand.index, idx) {
				g = cand
				break
			}
		}
		if g == nil {
			g = &group{r: r, index: map[string]string{}}
			byKey[key] = append(byKey[key], g)
		}
		g.files = append(g.files, f)
		for k, v := range idx {
			g.index[k] = v
		}
	}
	for _, k := range keys {
		p.groups = append(p.groups, byKey[k]...)
	}
	if len(p.groups) == 0 {
		logf("%s: no runnable tests; skipped", mod.Rel)
		return nil
	}
	return p
}

// targetSites keeps the sites fully inside each callee's definition (its
// Source text's first occurrence in src; the whole module when it isn't
// found), samples up to perCallee per callee (seeded by src and the
// callee's symbol), and returns their union sorted by Start.
func targetSites(src []byte, sites []mutate.Site, callees []cases.Callee) []mutate.Site {
	type key struct {
		start, end int
		repl       string
	}
	seen := map[key]bool{}
	var out []mutate.Site
	for _, c := range callees {
		lo, hi := 0, len(src)
		if c.Source != "" {
			if i := bytes.Index(src, []byte(c.Source)); i >= 0 {
				lo, hi = i, i+len(c.Source)
			}
		}
		var in []mutate.Site
		for _, s := range sites {
			if s.Start >= lo && s.End <= hi {
				in = append(in, s)
			}
		}
		seed := sha256.Sum256(append(append([]byte(nil), src...), c.Symbol...))
		for _, s := range mutate.Sample(in, perCallee, seed[:]) {
			k := key{s.Start, s.End, s.Replacement}
			if !seen[k] {
				seen[k] = true
				out = append(out, s)
			}
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Start < out[j].Start })
	return out
}

// nativeIndex maps runner-native names to the ids in tests that belong to
// testFile: the id with its "<lang>:<file>:" prefix removed (Go: the go
// test -json Test path; Python: Class::func or func; TypeScript: the
// " > "-joined title chain). Deduplicated ids ("... #2") have no native
// name of their own and are left out.
func nativeIndex(lang, testFile string, tests []string) map[string]string {
	prefix := idPrefix[lang] + testFile + ":"
	idx := map[string]string{}
	for _, id := range tests {
		if !strings.HasPrefix(id, prefix) || dupIDRe.MatchString(id) {
			continue
		}
		idx[strings.TrimPrefix(id, prefix)] = id
	}
	return idx
}

func overlaps(a, b map[string]string) bool {
	for k := range b {
		if _, ok := a[k]; ok {
			return true
		}
	}
	return false
}

// badRun explains why an Outcome can't be used, or returns "".
func badRun(o runner.Outcome) string {
	switch {
	case o.TimedOut:
		return "timed out"
	case o.Invalid:
		return "invalid: " + o.InvalidReason
	case o.Err != "":
		return "error: " + o.Err
	}
	return ""
}

// baseline runs every group twice unmutated in sb. If a run is unusable
// (Invalid, timed out, failed to start), it isolates the module's test
// files (see isolate).
func (p *modPlan) baseline(ctx context.Context, sb Sandbox, timeout time.Duration) {
	p.combinedBaseline(ctx, sb, timeout)
	if !p.ok {
		p.isolate(ctx, sb, timeout)
	}
}

// combinedBaseline runs every group twice unmutated in sb.
func (p *modPlan) combinedBaseline(ctx context.Context, sb Sandbox, timeout time.Duration) {
	p.ok = true
	p.failed = map[string]bool{}
	passes := map[string]int{}
	for run := 0; run < 2; run++ {
		start := time.Now()
		for _, g := range p.groups {
			p.clearBytecode(sb)
			o := g.r.Run(ctx, sb.Dir, g.files, g.index, timeout)
			if badRun(o) != "" {
				p.ok = false
			}
			if run == 0 {
				p.unmapped += o.Unmapped
			}
			for id := range o.Passed {
				passes[id]++
			}
			for id := range o.Failed {
				p.failed[id] = true
			}
		}
		p.slowest = max(p.slowest, time.Since(start))
	}
	p.measured = map[string]bool{}
	if !p.ok {
		return
	}
	for id, n := range passes {
		if n == 2 && !p.failed[id] {
			p.measured[id] = true
		}
	}
}

// isolate runs each test file alone, twice, unmutated in sb, and drops
// (into p.dropped) the files with an unusable run or a test that passed
// one run and failed the other. The kept files form the module's groups
// and their two runs its baseline; p.ok is false when none is kept.
// Failures seen in the combined baseline stay in p.failed.
func (p *modPlan) isolate(ctx context.Context, sb Sandbox, timeout time.Duration) {
	passes := map[string]int{}
	var dur [2]time.Duration
	p.unmapped = 0
	var groups []*group
	for _, g := range p.groups {
		kept := &group{r: g.r, index: map[string]string{}}
		for _, f := range g.files {
			idx := nativeIndex(p.mod.Lang, f, p.mod.Tests)
			fp, ff := map[string]int{}, map[string]bool{}
			bad, unmapped := false, 0
			for run := 0; run < 2 && !bad; run++ {
				p.clearBytecode(sb)
				start := time.Now()
				o := g.r.Run(ctx, sb.Dir, []string{f}, idx, timeout)
				dur[run] += time.Since(start)
				if badRun(o) != "" {
					bad = true
				}
				if run == 0 {
					unmapped = o.Unmapped
				}
				for id := range o.Passed {
					fp[id]++
				}
				for id := range o.Failed {
					ff[id] = true
				}
			}
			flaky := false
			for id := range ff {
				if fp[id] > 0 {
					flaky = true
				}
			}
			if bad || flaky {
				p.dropped = append(p.dropped, f)
				if flaky {
					for id := range ff {
						p.failed[id] = true
					}
				}
				continue
			}
			kept.files = append(kept.files, f)
			maps.Copy(kept.index, idx)
			p.unmapped += unmapped
			for id, n := range fp {
				passes[id] += n
			}
			for id := range ff {
				p.failed[id] = true
			}
		}
		if len(kept.files) > 0 {
			groups = append(groups, kept)
		}
	}
	sort.Strings(p.dropped)
	p.groups = groups
	p.slowest = max(dur[0], dur[1])
	p.ok = len(groups) > 0
	p.measured = map[string]bool{}
	if !p.ok {
		return
	}
	for id, n := range passes {
		if n == 2 && !p.failed[id] {
			p.measured[id] = true
		}
	}
}

// mutant plants s in sb, runs the groups holding measured tests, and
// restores the module. It returns the sorted measured tests that failed
// (non-nil), or nil and a reason when the run is excluded. An error means
// the sandbox couldn't be written or restored.
func (p *modPlan) mutant(ctx context.Context, sb Sandbox, s mutate.Site, timeout time.Duration) ([]string, string, error) {
	mutated, err := mutate.Apply(p.src, s)
	if err != nil {
		return nil, err.Error(), nil
	}
	file := filepath.Join(sb.Dir, filepath.FromSlash(p.mod.Rel))
	if err := os.WriteFile(file, mutated, 0o644); err != nil {
		return nil, "", fmt.Errorf("evidence: plant mutant in %s: %w", file, err)
	}
	if plantHook != nil {
		plantHook(file)
	}
	killed := []string{}
	why := ""
	for _, g := range p.groups {
		if !holdsMeasured(g, p.measured) {
			continue
		}
		p.clearBytecode(sb)
		o := g.r.Run(ctx, sb.Dir, g.files, g.index, timeout)
		if why = badRun(o); why != "" {
			break
		}
		for id := range o.Failed {
			if p.measured[id] {
				killed = append(killed, id)
			}
		}
	}
	if err := os.WriteFile(file, p.src, 0o644); err != nil {
		return nil, "", fmt.Errorf("evidence: restore %s: %w", file, err)
	}
	if why != "" {
		return nil, why, nil
	}
	sort.Strings(killed)
	return killed, "", nil
}

// canary plants a module-breaking change (see canarySource) in sb, runs
// the groups holding measured tests once, and restores the module. loaded
// is false only when every run was valid (not Invalid, timed out or
// failed to start) and no measured test failed: the tests then load some
// other copy of the module. An Invalid run (a Go build failure, a pytest
// collection error) counts as loaded. why says what showed it loaded. An
// error means the sandbox couldn't be written or restored.
func (p *modPlan) canary(ctx context.Context, sb Sandbox, timeout time.Duration) (loaded bool, why string, err error) {
	broken, ok := canarySource(p.mod.Lang, p.src)
	if !ok {
		return true, "no canary for " + p.mod.Lang, nil
	}
	file := filepath.Join(sb.Dir, filepath.FromSlash(p.mod.Rel))
	if err := os.WriteFile(file, broken, 0o644); err != nil {
		return false, "", fmt.Errorf("evidence: plant canary in %s: %w", file, err)
	}
	for _, g := range p.groups {
		if !holdsMeasured(g, p.measured) {
			continue
		}
		p.clearBytecode(sb)
		o := g.r.Run(ctx, sb.Dir, g.files, g.index, timeout)
		if why = badRun(o); why != "" {
			loaded = true
			break
		}
		for id := range o.Failed {
			if p.measured[id] {
				loaded, why = true, id+" failed"
				break
			}
		}
		if loaded {
			break
		}
	}
	if err := os.WriteFile(file, p.src, 0o644); err != nil {
		return false, "", fmt.Errorf("evidence: restore %s: %w", file, err)
	}
	return loaded, why, nil
}

func holdsMeasured(g *group, measured map[string]bool) bool {
	for _, id := range g.index {
		if measured[id] {
			return true
		}
	}
	return false
}

// clearBytecode removes a Python module's cached bytecode from sb before a
// run. The runner already runs with PYTHONDONTWRITEBYTECODE=1; this also
// drops any .pyc that got there anyway (e.g. copied from the working
// tree), since CPython's whole-second mtime + size check can't tell a
// same-size mutant from the version the .pyc was compiled from.
func (p *modPlan) clearBytecode(sb Sandbox) {
	if p.mod.Lang != "python" {
		return
	}
	rel := filepath.FromSlash(p.mod.Rel)
	stem := strings.TrimSuffix(filepath.Base(rel), filepath.Ext(rel))
	pyc, _ := filepath.Glob(filepath.Join(sb.Dir, filepath.Dir(rel), "__pycache__", stem+".*.pyc"))
	for _, f := range pyc {
		_ = os.Remove(f)
	}
}
