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
	"slices"
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
// mutants are planted only inside them and the definitions they reach in
// Rel (see targetSites).
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
// it. Uncovered marks a mutant never run because coverage showed no
// measured test executes its code: KilledBy is empty (non-nil), and it
// counts as planted for no test (see measure).
type Mutant struct {
	Module    string      `json:"module"`
	Site      mutate.Site `json:"site"`
	KilledBy  []string    `json:"killed_by"`
	Uncovered bool        `json:"uncovered,omitempty"`
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
// Unconfirmed lists (sorted) the modules whose canary run timed out or
// failed to start: whether the tests load them is unknown, so they are
// treated like Shadowed ones.
// Dropped maps a module to the test files (sorted) left out of it because
// they did not run cleanly on their own (see modPlan.isolate).
// Unloaded maps a module to the tests (sorted) that passed a valid canary
// run of it: they don't load this copy, so they are not measured there.
// NoSites lists (sorted) the modules skipped for having no mutation site
// (inside the targets' callees, when given).
type Matrix struct {
	Mutants  []Mutant        `json:"mutants"`
	Baseline map[string]bool `json:"baseline"`
	Flaky    []string        `json:"flaky"`
	Unmapped int             `json:"unmapped"`

	MeasuredByModule map[string]map[string]bool `json:"measured_by_module"`
	Shadowed         []string                   `json:"shadowed"`
	Unconfirmed      []string                   `json:"unconfirmed"`
	Dropped          map[string][]string        `json:"dropped"`
	Unloaded         map[string][]string        `json:"unloaded"`
	NoSites          []string                   `json:"no_sites"`
}

// Options tunes RunMatrix and Run. Zero values take defaults: Workers = runtime.NumCPU(),
// MaxMutants = 30 (for a module without Callees), RunTimeout = 60s,
// BaselineTimeout = 10m, Stderr = discard. Baseline runs get
// BaselineTimeout; canary and mutant runs get the larger of RunTimeout and
// 3× the module's slower baseline run.
// Evidence always runs the runner package's detected per-test command,
// never a project's test_command override. NoCoverage turns off
// coverage-guided test selection: every mutant runs every test file.
// NoCache ignores (and doesn't write) the result cache (see resultCache).
type Options struct {
	Workers, MaxMutants         int
	RunTimeout, BaselineTimeout time.Duration
	Stderr                      io.Writer
	NoCoverage, NoCache         bool
}

// group is a set of test files run together by one Runner, with an index
// whose native keys are unambiguous within the group.
type group struct {
	r     runner.Runner
	files []string
	index map[string]string
	cov   *runner.Coverage // per-test coverage of the module (nil: none)
}

// modPlan is one module ready to measure.
type modPlan struct {
	mod    Module
	src    []byte
	sites  []mutate.Site
	groups []*group

	ok       bool            // baseline valid
	shadowed bool            // canary unnoticed: tests load another copy
	unsure   bool            // canary run timed out or failed to start
	measured map[string]bool // passed both baseline runs
	failed   map[string]bool // failed at least one baseline run
	dropped  []string        // test files that did not run cleanly alone
	unloaded []string        // measured tests that passed a valid canary run
	unmapped int
	slowest  time.Duration // the slower baseline run (all groups)

	selective bool     // every group has coverage: mutants run only covering tests
	funcs     [][2]int // function and method spans (sites outside run every test)

	key           string   // result cache key ("": no cache)
	cached        bool     // results restored from the cache
	cachedMutants []Mutant // the cached mutants
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
// (see Matrix.Shadowed); then per-test coverage of the module is measured
// (see modPlan.coverage; the full test set, unmutated); then each sampled
// mutant (see targetSites; for a module without Callees, at most
// MaxMutants seeded by the module's bytes) is written into a free
// sandbox, the tests whose coverage includes it are run (every test for
// module-level code, or when coverage is unavailable; none: the mutant is
// Uncovered and not run), and the module's original bytes are restored
// before the sandbox is reused.
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
	sbs, copied, err := prepare(abs, base, opt.Workers, logf)
	if err != nil {
		return Matrix{}, err
	}
	var cache *resultCache
	if !opt.NoCache {
		cache = newResultCache(abs, base, copied, opt)
	}

	var plans []*modPlan
	var noSites []string
	for _, mod := range mods {
		p, empty := planModule(abs, sbs[0].Dir, mod, opt.MaxMutants, logf)
		if p != nil {
			if cache != nil {
				p.key = cache.keyFor(p)
				if e, ok := cache.get(p.key, p.mod.Rel); ok {
					p.restore(e)
					logf("%s: unchanged since a previous run; reusing its %d mutant results", p.mod.Rel, len(e.Mutants))
				}
			}
			plans = append(plans, p)
		}
		if empty {
			noSites = append(noSites, mod.Rel)
		}
	}

	pool := make(chan Sandbox, len(sbs))
	for _, sb := range sbs {
		pool <- sb
	}

	// Phase 1: baselines, one job per module.
	jobs := make([]func(Sandbox) error, 0, len(plans))
	for _, p := range plans {
		if p.cached {
			continue
		}
		jobs = append(jobs, func(sb Sandbox) error {
			p.baseline(ctx, sb, opt.BaselineTimeout)
			if err := ctx.Err(); err != nil {
				return err // cancelled: its results mean nothing
			}
			if len(p.dropped) > 0 {
				logf("%s: dropped test files that did not run cleanly alone: %s", p.mod.Rel, strings.Join(p.dropped, ", "))
			}
			logf("baseline %s: %d measured, %d failing", p.mod.Rel, len(p.measured), len(p.failed))
			if !p.ok || len(p.measured) == 0 {
				return nil
			}
			st, why, err := p.canary(ctx, sb, p.runTimeout(opt.RunTimeout))
			if err != nil {
				return err
			}
			if err := ctx.Err(); err != nil {
				return err
			}
			if len(p.unloaded) > 0 {
				logf("canary %s: passed by %s; they don't load this copy", p.mod.Rel, strings.Join(p.unloaded, ", "))
			}
			switch st {
			case canaryLoaded:
				logf("canary %s: loaded (%s)", p.mod.Rel, why)
				if !opt.NoCoverage {
					p.coverage(ctx, sb, opt.BaselineTimeout, logf)
					if err := ctx.Err(); err != nil {
						return err
					}
				}
			case canaryShadowed:
				p.shadowed = true
				logf("canary %s: not loaded; %s", p.mod.Rel, reasonShadowed)
			default:
				p.unsure = true
				logf("canary %s: %s; %s", p.mod.Rel, why, reasonUnconfirmed)
			}
			return nil
		})
	}
	if err := runJobs(ctx, pool, opt.Workers, jobs); err != nil {
		return Matrix{}, err
	}

	// Phase 2: mutants, one job per (module, site), results by index.
	var mutants []Mutant
	span := map[*modPlan][2]int{} // p's mutants: mutants[lo:hi]
	jobs = jobs[:0]
	for _, p := range plans {
		lo := len(mutants)
		switch {
		case p.cached:
			for _, mu := range p.cachedMutants {
				mu.KilledBy = slices.Clone(mu.KilledBy)
				mutants = append(mutants, mu)
			}
		case !p.ok:
			logf("%s: baseline invalid; its tests are unmeasurable", p.mod.Rel)
		case len(p.measured) == 0 || p.shadowed || p.unsure:
		default:
			for i, s := range p.sites {
				idx := len(mutants)
				mutants = append(mutants, Mutant{Module: p.mod.Rel, Site: s})
				n := i + 1
				if p.selective && p.uncovered(s) {
					mutants[idx].KilledBy, mutants[idx].Uncovered = []string{}, true
					logf("%s mutant %d/%d (%s): uncovered (no measured test runs it); not run", p.mod.Rel, n, len(p.sites), s.Desc)
					continue
				}
				jobs = append(jobs, func(sb Sandbox) error {
					killed, why, err := p.mutant(ctx, sb, s, p.runTimeout(opt.RunTimeout))
					if err != nil {
						return err
					}
					if err := ctx.Err(); err != nil {
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
		span[p] = [2]int{lo, len(mutants)}
	}
	if err := runJobs(ctx, pool, opt.Workers, jobs); err != nil {
		return Matrix{}, err
	}
	if cache != nil {
		for _, p := range plans {
			sp := span[p]
			if muts := mutants[sp[0]:sp[1]]; !p.cached && cacheable(p, muts) {
				if err := cache.put(entry(p.key, p, muts)); err != nil {
					logf("%s: result not cached: %v", p.mod.Rel, err)
				}
			}
		}
	}

	m := Matrix{Mutants: mutants, Baseline: map[string]bool{}, MeasuredByModule: map[string]map[string]bool{}, Dropped: map[string][]string{}, Unloaded: map[string][]string{}, NoSites: noSites}
	sort.Strings(m.NoSites)
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
		if p.unsure {
			m.Unconfirmed = append(m.Unconfirmed, p.mod.Rel)
		}
		if len(p.dropped) > 0 {
			m.Dropped[p.mod.Rel] = p.dropped
		}
		if len(p.unloaded) > 0 {
			m.Unloaded[p.mod.Rel] = p.unloaded
		}
		if !p.ok || p.shadowed || p.unsure || len(p.sites) == 0 {
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
	sort.Strings(m.Unconfirmed)
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
// cancellation) stops further jobs at once, even ones already handed to
// a worker, and is returned.
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
				if failed() || ctx.Err() != nil {
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
// can't be measured at all; noSites is true when that's because it has no
// mutation site.
func planModule(root, sbDir string, mod Module, max int, logf func(string, ...any)) (p *modPlan, noSites bool) {
	if !copiedFile(sbDir, mod.Rel) {
		logf("%s: not a copied project file; skipped", mod.Rel)
		return nil, false
	}
	mut := mutate.For(mod.Lang)
	if mut == nil {
		logf("%s: no mutator for %q; skipped", mod.Rel, mod.Lang)
		return nil, false
	}
	src, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(mod.Rel)))
	if err != nil {
		logf("%s: %v; skipped", mod.Rel, err)
		return nil, false
	}
	sites, err := mut.Sites(root, mod.Rel, src)
	if err != nil {
		logf("%s: %v; skipped", mod.Rel, err)
		return nil, false
	}
	if len(mod.Callees) > 0 {
		defs, err := mutate.Defs(mod.Lang, root, mod.Rel, src)
		if err != nil {
			logf("%s: definitions unavailable (%v); mutating the callees only", mod.Rel, err)
		}
		sites = targetSites(mod.Rel, src, sites, mod.Callees, defs, logf)
	} else {
		seed := sha256.Sum256(src)
		sites = mutate.Sample(sites, max, seed[:])
	}
	if len(sites) == 0 {
		logf("%s: no mutation sites; skipped", mod.Rel)
		return nil, true
	}
	p = &modPlan{mod: mod, src: src, sites: sites}

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
		return nil, false
	}
	return p, false
}

// targetSites plants up to perCallee sites (seeded by src and a name) in
// each callee and in each definition reached (see reach) from the defs
// inside a located callee's definition (every occurrence of its Source
// text in src, matched with CRLF line endings read as LF and mapped back
// to src's bytes; the whole module, logged, when it isn't found). Each
// def samples only its own sites, those not inside a def nested in it
// (a class: its field defaults and body statements; a function: not its
// nested helpers), and a callee only its sites not inside any def within
// it (a callee holding no def of its own: every site in its span). It
// returns their union, deduplicated and sorted by Start.
func targetSites(rel string, src []byte, sites []mutate.Site, callees []cases.Callee, defs []mutate.Def, logf func(string, ...any)) []mutate.Site {
	type key struct {
		start, end int
		repl       string
	}
	type span [2]int
	within := func(a, b span) bool { return a[0] >= b[0] && a[1] <= b[1] }
	// inner[i]: the smallest def span holding site i (ok false: none).
	inner := make([]span, len(sites))
	hasInner := make([]bool, len(sites))
	for i, s := range sites {
		for _, d := range defs {
			sp := span{d.Start, d.End}
			if within(span{s.Start, s.End}, sp) && (!hasInner[i] || sp[1]-sp[0] < inner[i][1]-inner[i][0]) {
				inner[i], hasInner[i] = sp, true
			}
		}
	}
	norm, orig := lfOffsets(src)
	seen := map[key]bool{}
	var out []mutate.Site
	sample := func(keep func(i int) bool, name string) {
		var in []mutate.Site
		for i, s := range sites {
			if keep(i) {
				in = append(in, s)
			}
		}
		seed := sha256.Sum256(append(append([]byte(nil), src...), name...))
		for _, s := range mutate.Sample(in, perCallee, seed[:]) {
			k := key{s.Start, s.End, s.Replacement}
			if !seen[k] {
				seen[k] = true
				out = append(out, s)
			}
		}
	}
	var seeds []int
	for _, c := range callees {
		var spans []span
		for _, sp := range sourceSpans(norm, orig, c.Source) {
			spans = append(spans, span(sp))
		}
		if len(spans) == 0 {
			logf("%s: definition of %s not found; sampling the whole module", rel, c.Symbol)
			sample(func(int) bool { return true }, c.Symbol)
			continue
		}
		in := func(x span) bool {
			for _, sp := range spans {
				if within(x, sp) {
					return true
				}
			}
			return false
		}
		var own []int // the defs inside this callee's spans
		for i, d := range defs {
			if in(span{d.Start, d.End}) {
				own = append(own, i)
			}
		}
		if len(own) == 0 {
			// No definition of its own: sample its span directly.
			sample(func(i int) bool { return in(span{sites[i].Start, sites[i].End}) }, c.Symbol)
			continue
		}
		sample(func(i int) bool {
			return in(span{sites[i].Start, sites[i].End}) && (!hasInner[i] || !in(inner[i]))
		}, c.Symbol)
		seeds = append(seeds, own...)
	}
	done := map[span]bool{} // defs sharing one span share one sample
	for _, j := range reach(defs, seeds) {
		sp := span{defs[j].Start, defs[j].End}
		if !done[sp] {
			done[sp] = true
			sample(func(i int) bool { return hasInner[i] && inner[i] == sp }, qualName(defs[j]))
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Start < out[j].Start })
	return out
}

// lfOffsets returns src with each CRLF read as LF, and for each offset in
// it (and its end) the matching offset in src.
func lfOffsets(src []byte) (norm []byte, orig []int) {
	norm = make([]byte, 0, len(src))
	orig = make([]int, 0, len(src)+1)
	for i := 0; i < len(src); i++ {
		if src[i] == '\r' && i+1 < len(src) && src[i+1] == '\n' {
			continue
		}
		norm = append(norm, src[i])
		orig = append(orig, i)
	}
	orig = append(orig, len(src))
	return norm, orig
}

// sourceSpans lists the [start, end) offsets in the original bytes of
// every non-overlapping occurrence of source (CRLF read as LF) in norm.
func sourceSpans(norm []byte, orig []int, source string) [][2]int {
	pat := []byte(strings.ReplaceAll(source, "\r\n", "\n"))
	if len(pat) == 0 {
		return nil
	}
	var spans [][2]int
	for off := 0; ; {
		i := bytes.Index(norm[off:], pat)
		if i < 0 {
			return spans
		}
		lo := off + i
		hi := lo + len(pat)
		// The span ends after the last source byte in src (a matched
		// trailing LF covers its CR).
		spans = append(spans, [2]int{orig[lo], orig[hi-1] + 1})
		off = hi
	}
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
	if !p.ok && ctx.Err() == nil {
		p.isolate(ctx, sb, timeout)
	}
}

// combinedBaseline runs every group twice unmutated in sb.
func (p *modPlan) combinedBaseline(ctx context.Context, sb Sandbox, timeout time.Duration) {
	p.failed = map[string]bool{}
	p.measured = map[string]bool{}
	var passes map[string]int
	p.ok, passes, p.slowest, p.unmapped = p.runGroups(ctx, sb, timeout)
	if !p.ok {
		return
	}
	for id, n := range passes {
		if n == 2 && !p.failed[id] {
			p.measured[id] = true
		}
	}
}

// runGroups runs every group twice unmutated in sb, adding every failure
// to p.failed. It stops at the first unusable run (ok false: its results
// aren't used, so a second pass would only cost time). passes counts each
// test's passing runs; slowest is the slower run (summed over groups);
// unmapped is the first run's unmapped results.
func (p *modPlan) runGroups(ctx context.Context, sb Sandbox, timeout time.Duration) (ok bool, passes map[string]int, slowest time.Duration, unmapped int) {
	ok = true
	passes = map[string]int{}
	for run := 0; run < 2 && ok; run++ {
		start := time.Now()
		for _, g := range p.groups {
			if ctx.Err() != nil {
				return false, passes, slowest, unmapped
			}
			p.clearBytecode(sb)
			o := g.r.Run(ctx, sb.Dir, g.files, g.index, timeout)
			for id := range o.Failed {
				p.failed[id] = true
			}
			if badRun(o) != "" {
				ok = false
				break
			}
			if run == 0 {
				unmapped += o.Unmapped
			}
			for id := range o.Passed {
				passes[id]++
			}
		}
		slowest = max(slowest, time.Since(start))
	}
	return ok, passes, slowest, unmapped
}

// isolate runs each test file alone, twice, unmutated in sb, and drops
// (into p.dropped) the files with an unusable run or a test that passed
// one run and failed the other. The kept files form the module's groups
// and are then run together twice more (see runGroups): a test is
// measured only if it passed all four runs. p.ok is false when no file is
// kept or a together run is unusable. Failures seen in any baseline run
// stay in p.failed.
func (p *modPlan) isolate(ctx context.Context, sb Sandbox, timeout time.Duration) {
	passes := map[string]int{}
	var dur [2]time.Duration
	var groups []*group
	for _, g := range p.groups {
		kept := &group{r: g.r, index: map[string]string{}}
		for _, f := range g.files {
			idx := nativeIndex(p.mod.Lang, f, p.mod.Tests)
			fp, ff := map[string]int{}, map[string]bool{}
			bad := false
			var fdur [2]time.Duration
			for run := 0; run < 2 && !bad; run++ {
				if ctx.Err() != nil {
					p.ok = false
					return
				}
				p.clearBytecode(sb)
				start := time.Now()
				o := g.r.Run(ctx, sb.Dir, []string{f}, idx, timeout)
				fdur[run] = time.Since(start)
				if badRun(o) != "" {
					bad = true
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
			dur[0] += fdur[0]
			dur[1] += fdur[1]
			maps.Copy(kept.index, idx)
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
	// The kept files run together in every canary and mutant run: run
	// them together twice unmutated too, so a failure seen only in that
	// combination marks its test flaky instead of killing every mutant.
	ok, together, slowest, unmapped := p.runGroups(ctx, sb, timeout)
	p.slowest = max(p.slowest, slowest)
	p.unmapped = unmapped
	if !ok {
		p.ok = false
		return
	}
	for id, n := range passes {
		if n == 2 && together[id] == 2 && !p.failed[id] {
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
		units, all := p.selection(g, s)
		if !all && len(units) == 0 {
			continue // no test of this group runs the mutated code
		}
		if ctx.Err() != nil {
			why = "cancelled"
			break
		}
		p.clearBytecode(sb)
		var o runner.Outcome
		if all {
			o = g.r.Run(ctx, sb.Dir, g.files, g.index, timeout)
		} else {
			o = g.r.(runner.Coverer).RunUnits(ctx, sb.Dir, units, g.index, timeout)
		}
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

// coverage measures, per group, what each test executes of the module
// (see runner.Coverer), after the canary, with the final measured set.
// If any group holding measured tests can't be measured, or a measured
// test is missing from its data, it logs why once and leaves the module
// running every test file on every mutant (p.selective false).
func (p *modPlan) coverage(ctx context.Context, sb Sandbox, timeout time.Duration, logf func(string, ...any)) {
	fallback := func(why string) {
		for _, g := range p.groups {
			g.cov = nil
		}
		logf("%s: no per-test coverage (%s); running every test file per mutant", p.mod.Rel, why)
	}
	for _, g := range p.groups {
		if !holdsMeasured(g, p.measured) {
			continue
		}
		c, ok := g.r.(runner.Coverer)
		if !ok {
			fallback("the test runner can't measure it")
			return
		}
		files := map[string]map[string]string{}
		for _, f := range g.files {
			files[f] = nativeIndex(p.mod.Lang, f, p.mod.Tests)
		}
		p.clearBytecode(sb)
		cov, err := c.Cover(ctx, sb.Dir, files, p.mod.Rel, timeout)
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			fallback(err.Error())
			return
		}
		have := map[string]bool{}
		for _, u := range cov.Units {
			for _, id := range u.IDs {
				have[id] = true
			}
		}
		for _, id := range sortedKeys(p.measured) {
			if slices.Contains(sortedValues(g.index), id) && !have[id] {
				fallback("no coverage data for " + id)
				return
			}
		}
		g.cov = cov
	}
	defs, err := mutate.Defs(p.mod.Lang, sb.Dir, p.mod.Rel, p.src)
	if err != nil {
		fallback("definitions unavailable: " + err.Error())
		return
	}
	for _, d := range defs {
		if d.Kind == "function" || d.Kind == "method" {
			p.funcs = append(p.funcs, [2]int{d.Start, d.End})
		}
	}
	p.selective = true
}

// selection is what a mutant at s runs of g: all (every test file: the
// module isn't selective, s is outside every function or method, or its
// code runs for every test), else the units whose coverage includes it
// (none: skip g).
func (p *modPlan) selection(g *group, s mutate.Site) (units []runner.Unit, all bool) {
	if !p.selective || g.cov == nil || !p.inFunc(s) || g.cov.Global(s.Start, s.End) {
		return nil, true
	}
	for _, u := range g.cov.Units {
		if u.Covers(s.Start, s.End) && slices.ContainsFunc(u.IDs, func(id string) bool { return p.measured[id] }) {
			units = append(units, u)
		}
	}
	return units, false
}

// uncovered reports whether no group holding measured tests would run
// anything for a mutant at s.
func (p *modPlan) uncovered(s mutate.Site) bool {
	for _, g := range p.groups {
		if !holdsMeasured(g, p.measured) {
			continue
		}
		if units, all := p.selection(g, s); all || len(units) > 0 {
			return false
		}
	}
	return true
}

// inFunc reports whether s lies inside a function or method.
func (p *modPlan) inFunc(s mutate.Site) bool {
	for _, f := range p.funcs {
		if s.Start >= f[0] && s.End <= f[1] {
			return true
		}
	}
	return false
}

func sortedValues(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

// canaryStatus is what a canary run showed.
type canaryStatus int

const (
	canaryLoaded      canaryStatus = iota // a measured test failed, or a run was Invalid
	canaryShadowed                        // every run valid, no measured test failed
	canaryUnconfirmed                     // otherwise: a run timed out or failed to start
)

// canary plants a module-breaking change (see canarySource) in sb, runs
// the groups holding measured tests once, and restores the module. A
// measured test that passes a valid run doesn't load this copy: it moves
// from p.measured to p.unloaded. The
// status is canaryLoaded when a measured test failed (even in a run that
// then timed out) or a run was Invalid (a Go build failure, a pytest
// collection error); else canaryUnconfirmed when a run timed out or
// failed to start; else canaryShadowed: the tests load some other copy
// of the module. why says what showed it. An error means the sandbox
// couldn't be written or restored.
func (p *modPlan) canary(ctx context.Context, sb Sandbox, timeout time.Duration) (st canaryStatus, why string, err error) {
	broken, ok := canarySource(p.mod.Lang, p.src)
	if !ok {
		return canaryLoaded, "no canary for " + p.mod.Lang, nil
	}
	file := filepath.Join(sb.Dir, filepath.FromSlash(p.mod.Rel))
	if err := os.WriteFile(file, broken, 0o644); err != nil {
		return 0, "", fmt.Errorf("evidence: plant canary in %s: %w", file, err)
	}
	st = canaryShadowed
	var passed []string
	for _, g := range p.groups {
		if !holdsMeasured(g, p.measured) {
			continue
		}
		if ctx.Err() != nil {
			st, why = canaryUnconfirmed, "cancelled"
			break
		}
		p.clearBytecode(sb)
		o := g.r.Run(ctx, sb.Dir, g.files, g.index, timeout)
		if badRun(o) == "" {
			for id := range o.Passed {
				if p.measured[id] {
					passed = append(passed, id)
				}
			}
		}
		if st == canaryLoaded {
			continue
		}
		for id := range o.Failed {
			if p.measured[id] {
				st, why = canaryLoaded, id+" failed"
				break
			}
		}
		if st == canaryLoaded {
			continue
		}
		if o.Invalid && !o.TimedOut && o.Err == "" {
			st, why = canaryLoaded, badRun(o)
			continue
		}
		if r := badRun(o); r != "" {
			st, why = canaryUnconfirmed, r
		}
	}
	if err := os.WriteFile(file, p.src, 0o644); err != nil {
		return 0, "", fmt.Errorf("evidence: restore %s: %w", file, err)
	}
	sort.Strings(passed)
	for _, id := range passed {
		delete(p.measured, id)
	}
	p.unloaded = passed
	return st, why, nil
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
