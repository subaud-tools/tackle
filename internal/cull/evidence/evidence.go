package evidence

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/schuettc/tackle/internal/cull/cases"
	"github.com/schuettc/tackle/internal/cull/check"
	"github.com/schuettc/tackle/internal/cull/discover"
	"github.com/schuettc/tackle/internal/cull/extract"
	"github.com/schuettc/tackle/internal/cull/extract/golang"
	_ "github.com/schuettc/tackle/internal/cull/extract/python"
	_ "github.com/schuettc/tackle/internal/cull/extract/ts"
	"github.com/schuettc/tackle/internal/cull/mutate"
	"github.com/schuettc/tackle/internal/cull/runner"
	tools "github.com/schuettc/tools-common"
)

// Evidence is one settled test's measurement, as stored in last.json.
type Evidence = check.Evidence

// Result is what one Run settled. Settled maps each test whose verdict
// changed to its new verdict (keep or cut); Rules maps every measured-or-
// not target to its new rule; Unmeasured (sorted) lists the targets left
// in review, with Reasons saying why. Modules and Mutants count what the
// matrix ran (Mutants includes excluded ones).
type Result struct {
	Settled    map[string]string `json:"settled"`
	Rules      map[string]string `json:"rules"`
	Unmeasured []string          `json:"unmeasured"`
	Reasons    map[string]string `json:"reasons"`
	Modules    int               `json:"modules"`
	Mutants    int               `json:"mutants"`
	Duration   time.Duration     `json:"duration"`
}

// Rules and reasons written to last.json.
const (
	RuleSoleGuard      = "evidence_sole_guard"
	RuleRedundant      = "evidence_redundant"
	RuleCatchesNothing = "evidence_catches_nothing"
	RuleUnmeasured     = "evidence_unmeasured"

	reasonNoCode     = "no project code to plant bugs in"
	reasonFlaky      = "flaky"
	reasonNoBaseline = "baseline did not run"
	reasonExcluded   = "every planted bug was excluded (timed out or did not build)"
	reasonDupName    = "no runner-reportable name (duplicate test name)"
	reasonNotFound   = "not found when the suite was extracted again"
	reasonDropped    = "their test file did not run cleanly"
)

// extractContext is large enough that extraction never drops a callee for
// its context budget: evidence needs every callee file, not a prompt.
const extractContext = 1 << 24

// afterMatrixHook, if set (tests only), runs right after the matrix.
var afterMatrixHook func()

// Run settles the review-band tests of root's .cull/last.json by planted
// bugs. Targets are tests with verdict "review" and rule "review_band" or
// "evidence_unmeasured" (so a re-run retries unmeasured ones).
// It refuses (an error, nothing written) if a target's file changed since
// cull check, or if a target file, a mutated module or one of its test
// files changed while the matrix ran. It mutates only non-test source
// files the targets call (their extracted callees), only inside the
// definitions they call (at most 15 mutants each), runs every test in
// those modules' test files against each mutant, and settles each target:
// sole guard of any bug → keep; measured with no sole guard → cut;
// unmeasured → stays review (rule evidence_unmeasured). A target is
// unmeasured if any module it exercises is unmeasured for it (see
// gateModules), whatever the others show. No other test is
// touched. It refuses if .cull/last.json itself changed since it was
// loaded. It writes .cull/evidence.json (the matrix) and rewrites
// .cull/last.json, both atomically and 0600. No target: nothing is run or
// written.
func Run(ctx context.Context, root string, opt Options) (Result, error) {
	start := time.Now()
	res := Result{Settled: map[string]string{}, Rules: map[string]string{}, Reasons: map[string]string{}, Unmeasured: []string{}}
	if opt.Stderr == nil {
		opt.Stderr = io.Discard
	}
	lastRel := ".cull/last.json"
	lastHash, _ := hashFile(root, lastRel) // before loading: a change in between only causes a refusal
	rep, err := check.LoadReport(root)
	if err != nil {
		return res, err
	}
	var targets []int
	for i, t := range rep.Tests {
		if t.Err == "" && t.Verdict == "review" && (t.Rule == "review_band" || t.Rule == RuleUnmeasured) {
			targets = append(targets, i)
		}
	}
	if len(targets) == 0 {
		res.Duration = time.Since(start)
		return res, nil
	}

	// Preflight [RF 5]: every target file as check saw it.
	targetFiles := map[string]bool{}
	for _, i := range targets {
		targetFiles[rep.Tests[i].File] = true
	}
	var changed []string
	for _, f := range sortedKeys(targetFiles) {
		h, err := hashFile(root, f)
		if fi, ok := rep.Files[f]; err != nil || !ok || h != fi.SHA256 {
			changed = append(changed, f)
		}
	}
	if len(changed) > 0 {
		return res, fmt.Errorf("refused: %s changed since cull check; run cull check again", strings.Join(changed, ", "))
	}

	cfg, _, err := check.LoadConfig(root)
	if err != nil {
		return res, err
	}
	files, err := discover.Suite(root, "", cfg.Exclude)
	if err != nil {
		return res, err
	}
	all, err := extractSuite(root, files, opt.Stderr)
	if err != nil {
		return res, err
	}
	spanSameNamedMethods(root, all, opt.Stderr)

	byID := map[string]cases.TestCase{}
	testsByFile := map[string][]string{}
	modTestFiles := map[string]map[string]bool{}
	modLang := map[string]string{}
	testMods := map[string][]string{} // test id → its modules (callee files), sorted
	for _, tc := range all {
		byID[tc.ID] = tc
		testsByFile[tc.File] = append(testsByFile[tc.File], tc.ID)
		seen := map[string]bool{}
		for _, c := range tc.Callees {
			f := c.File
			if seen[f] || !projectSource(root, f) {
				continue
			}
			seen[f] = true
			if modTestFiles[f] == nil {
				modTestFiles[f] = map[string]bool{}
				modLang[f] = tc.Lang
			}
			modTestFiles[f][tc.File] = true
			testMods[tc.ID] = append(testMods[tc.ID], f)
		}
		sort.Strings(testMods[tc.ID])
	}

	wanted := map[string]bool{}
	for _, i := range targets {
		for _, f := range testMods[rep.Tests[i].ID] {
			wanted[f] = true
		}
	}
	var mods []Module
	inputs := map[string]bool{}
	for f := range targetFiles {
		inputs[f] = true
	}
	for _, rel := range sortedKeys(wanted) {
		mod := Module{Rel: rel, Lang: modLang[rel], TestFiles: sortedKeys(modTestFiles[rel])}
		for _, tf := range mod.TestFiles {
			mod.Tests = append(mod.Tests, testsByFile[tf]...)
			inputs[tf] = true
		}
		sort.Strings(mod.Tests)
		mod.Callees = targetCallees(rel, targets, rep.Tests, byID)
		inputs[rel] = true
		mods = append(mods, mod)
	}
	pre := map[string]string{}
	for f := range inputs {
		pre[f], _ = hashFile(root, f)
	}

	m, err := RunMatrix(ctx, root, mods, opt)
	if err != nil {
		return res, err
	}
	if afterMatrixHook != nil {
		afterMatrixHook()
	}
	changed = nil
	for _, f := range sortedKeys(inputs) {
		if h, err := hashFile(root, f); err != nil || h != pre[f] {
			changed = append(changed, f)
		}
	}
	if len(changed) > 0 {
		return res, fmt.Errorf("refused: %s changed while cull evidence ran; nothing written, run cull check again", strings.Join(changed, ", "))
	}

	res.Modules, res.Mutants = len(mods), len(m.Mutants)
	byModule := map[string][]Mutant{}
	for _, mu := range m.Mutants {
		byModule[mu.Module] = append(byModule[mu.Module], mu)
	}
	flaky := map[string]bool{}
	for _, id := range m.Flaky {
		flaky[id] = true
	}
	shadowed := map[string]string{} // module → why its canary kept it unmeasured
	for _, rel := range m.Shadowed {
		shadowed[rel] = reasonShadowed
	}
	for _, rel := range m.Unconfirmed {
		shadowed[rel] = reasonUnconfirmed
	}
	view := matrixView{byModule: byModule, measuredBy: m.MeasuredByModule, flaky: flaky, shadowed: shadowed,
		dropped: m.Dropped, unloaded: m.Unloaded, noSites: map[string]bool{}}
	for _, rel := range m.NoSites {
		view.noSites[rel] = true
	}
	sites := map[string]bool{} // module → has mutation sites (computed lazily)
	hasSites := func(rel string) bool {
		if v, ok := sites[rel]; ok {
			return v
		}
		v := false
		if mut := mutate.For(modLang[rel]); mut != nil {
			if src, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel))); err == nil {
				s, err := mut.Sites(root, rel, src)
				v = err == nil && len(s) > 0
			}
		}
		sites[rel] = v
		return v
	}

	for _, i := range targets {
		t := &rep.Tests[i]
		ev := measure(t.ID, testMods[t.ID], byModule, m.MeasuredByModule, m.Baseline)
		ev = gateModules(ev, t.ID, t.File, testMods[t.ID], view)
		if !ev.Measured && ev.Reason == "" {
			ev.Reason = unmeasuredReason(root, t.ID, byID, testMods[t.ID], byModule, m.Baseline, flaky, shadowed, m.Dropped, hasSites)
		}
		t.Verdict, t.Rule = settle(ev)
		t.Reasons = dropEvidenceReasons(t.Reasons)
		if ev.Measured {
			res.Settled[t.ID] = t.Verdict
			t.Reasons = append(t.Reasons, fmt.Sprintf("evidence: caught %d/%d, sole %d", ev.Caught, ev.Planted, ev.Sole))
		} else {
			res.Unmeasured = append(res.Unmeasured, t.ID)
			res.Reasons[t.ID] = ev.Reason
			t.Reasons = append(t.Reasons, fmt.Sprintf("evidence: not measured (%s)", ev.Reason))
		}
		res.Rules[t.ID] = t.Rule
		t.Evidence = &ev
	}
	sort.Strings(res.Unmeasured)
	rep.Recount()

	if h, err := hashFile(root, lastRel); err != nil || h != lastHash {
		return res, fmt.Errorf("refused: .cull/last.json changed while cull evidence ran; nothing written, run cull evidence again")
	}
	if err := writeEvidenceJSON(root, mods, m); err != nil {
		return res, err
	}
	if err := check.WriteReport(root, rep); err != nil {
		return res, err
	}
	res.Duration = time.Since(start)
	return res, nil
}

// targetCallees lists (deduplicated) the callees in rel of the targets as
// extracted again now: the definitions evidence plants mutants in.
func targetCallees(rel string, targets []int, tests []check.TestResult, byID map[string]cases.TestCase) []cases.Callee {
	seen := map[cases.Callee]bool{}
	var out []cases.Callee
	for _, i := range targets {
		for _, c := range byID[tests[i].ID].Callees {
			if c.File == rel && !seen[c] {
				seen[c] = true
				out = append(out, c)
			}
		}
	}
	return out
}

// spanSameNamedMethods replaces each Go method callee whose name is
// declared on several receivers in its package with every one of them
// (logged once per package and name): extraction matches a method call
// by selector name alone, so it may have picked the wrong receiver's.
// The methods are looked up once per package and name, but each callee is
// replaced only if its own Source is one of them: a same-package function
// shares the bare name as its Symbol and keeps its own span.
func spanSameNamedMethods(root string, all []cases.TestCase, stderr io.Writer) {
	type key struct{ dir, name string }
	cache := map[key][]cases.Callee{}
	logged := map[key]bool{}
	for i := range all {
		tc := &all[i]
		if tc.Lang != "go" {
			continue
		}
		var out []cases.Callee
		for _, c := range tc.Callees {
			k := key{path.Dir(c.File), c.Symbol}
			ms, ok := cache[k]
			if !ok {
				ms = golang.Methods(root, filepath.Join(root, filepath.FromSlash(k.dir)), c.Symbol)
				cache[k] = ms
			}
			if len(ms) < 2 || !slices.ContainsFunc(ms, func(m cases.Callee) bool { return m.Source == c.Source }) {
				out = append(out, c)
				continue
			}
			if !logged[k] {
				logged[k] = true
				var where []string
				for _, m := range ms {
					where = append(where, m.Symbol+" in "+m.File)
				}
				fmt.Fprintf(stderr, "cull evidence: method %s is declared on %d receivers in %s; mutating every one (%s)\n",
					c.Symbol, len(ms), k.dir, strings.Join(where, ", "))
			}
			out = append(out, ms...)
		}
		tc.Callees = out
	}
}

// measure computes id's evidence over the non-excluded mutants of those of
// mods where id was measured (measuredBy[module][id]): a module whose
// baseline didn't measure id says nothing about it. Measured needs id in
// the baseline and at least one such mutant.
func measure(id string, mods []string, byModule map[string][]Mutant, measuredBy map[string]map[string]bool, baseline map[string]bool) Evidence {
	var ev Evidence
	others := map[string]bool{}
	for _, rel := range mods {
		if !measuredBy[rel][id] {
			continue
		}
		for _, mu := range byModule[rel] {
			if mu.KilledBy == nil {
				continue
			}
			ev.Planted++
			desc := mu.Module + " " + mu.Site.Desc
			caught := false
			for _, k := range mu.KilledBy {
				if k == id {
					caught = true
				} else {
					others[k] = true
				}
			}
			if caught {
				ev.Caught++
				ev.CaughtBy = append(ev.CaughtBy, desc)
				if len(mu.KilledBy) == 1 {
					ev.Sole++
					ev.SoleFor = append(ev.SoleFor, desc)
				}
			}
		}
	}
	ev.OtherTests = sortedKeys(others)
	ev.Measured = baseline[id] && ev.Planted > 0
	return ev
}

// dropEvidenceReasons removes a previous run's "evidence: " reasons, so a
// retried target carries only the latest one.
func dropEvidenceReasons(rs []string) []string {
	out := rs[:0:0]
	for _, r := range rs {
		if !strings.HasPrefix(r, "evidence: ") {
			out = append(out, r)
		}
	}
	return out
}

// settle maps a measurement to a verdict and rule: sole guard of any bug
// → keep; measured, caught some but sole for none → cut (redundant);
// measured, caught none → cut (catches nothing); unmeasured → review.
func settle(ev Evidence) (verdict, rule string) {
	switch {
	case ev.Measured && ev.Sole > 0:
		return "keep", RuleSoleGuard
	case ev.Measured && ev.Caught > 0:
		return "cut", RuleRedundant
	case ev.Measured:
		return "cut", RuleCatchesNothing
	}
	return "review", RuleUnmeasured
}

// testSupport reports whether rel is test support code that must never be
// mutated: a conftest.py, or any file under a directory named tests, test,
// __tests__, testdata or spec.
func testSupport(rel string) bool {
	segs := strings.Split(rel, "/")
	if segs[len(segs)-1] == "conftest.py" {
		return true
	}
	for _, d := range segs[:len(segs)-1] {
		switch d {
		case "tests", "test", "__tests__", "testdata", "spec":
			return true
		}
	}
	return false
}

// unmeasuredReason says why id has no usable measurement.
func unmeasuredReason(root, id string, byID map[string]cases.TestCase, mods []string,
	byModule map[string][]Mutant, baseline, flaky map[string]bool, shadowed map[string]string, dropped map[string][]string,
	hasSites func(string) bool) string {
	tc, ok := byID[id]
	if !ok {
		return reasonNotFound
	}
	if len(mods) == 0 {
		return reasonNoCode
	}
	if dupIDRe.MatchString(id) {
		return reasonDupName
	}
	if _, err := runner.For(root, tc.Lang, tc.File); err != nil {
		return err.Error()
	}
	for _, rel := range mods {
		if slices.Contains(dropped[rel], tc.File) {
			return reasonDropped
		}
	}
	if flaky[id] {
		return reasonFlaky
	}
	total := 0
	for _, rel := range mods {
		total += len(byModule[rel])
	}
	if total == 0 {
		for _, rel := range mods {
			if r := shadowed[rel]; r != "" {
				return r
			}
		}
		for _, rel := range mods {
			if hasSites(rel) {
				return reasonNoBaseline
			}
		}
		return reasonNoCode
	}
	if !baseline[id] {
		return reasonNoBaseline
	}
	return reasonExcluded
}

// projectSource reports whether a callee file is non-test source inside
// the project: a clean relative path to a regular file that no extractor
// claims as a test file and that is not test support (see testSupport).
func projectSource(root, rel string) bool {
	if rel == "" || path.IsAbs(rel) || filepath.IsAbs(rel) || path.Clean(rel) != rel ||
		rel == ".." || strings.HasPrefix(rel, "../") {
		return false
	}
	if extract.ForFile(rel) != nil || testSupport(rel) {
		return false
	}
	fi, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel)))
	return err == nil && fi.Mode().IsRegular()
}

// extractSuite runs every registered extractor over the files it matches,
// with a context budget large enough to keep every callee.
func extractSuite(root string, files []string, stderr io.Writer) ([]cases.TestCase, error) {
	var all []cases.TestCase
	for _, e := range extract.Registry() {
		var fs []string
		for _, f := range files {
			if e.Match(f) {
				fs = append(fs, f)
			}
		}
		if len(fs) == 0 {
			continue
		}
		r, err := e.Extract(root, fs, extractContext)
		if err != nil {
			return nil, err
		}
		for _, sk := range r.Skipped {
			fmt.Fprintf(stderr, "cull evidence: skipped %s: %s\n", sk.File, sk.Reason)
		}
		all = append(all, r.Cases...)
	}
	return all, nil
}

// evidenceFile is .cull/evidence.json: the modules measured and their
// full kill matrix.
type evidenceFile struct {
	Modules     []moduleJSON        `json:"modules"`
	Mutants     []Mutant            `json:"mutants"`
	Baseline    map[string]bool     `json:"baseline"`
	Flaky       []string            `json:"flaky"`
	Unmapped    int                 `json:"unmapped"`
	Shadowed    []string            `json:"shadowed"`
	Unconfirmed []string            `json:"unconfirmed"`
	Dropped     map[string][]string `json:"dropped"`
	Unloaded    map[string][]string `json:"unloaded"`
	NoSites     []string            `json:"no_sites"`
}

type moduleJSON struct {
	Rel       string   `json:"rel"`
	Lang      string   `json:"lang"`
	TestFiles []string `json:"test_files"`
	Tests     []string `json:"tests"`
}

func writeEvidenceJSON(root string, mods []Module, m Matrix) error {
	f := evidenceFile{Mutants: m.Mutants, Baseline: m.Baseline, Flaky: m.Flaky, Unmapped: m.Unmapped, Shadowed: m.Shadowed, Unconfirmed: m.Unconfirmed, Dropped: m.Dropped, Unloaded: m.Unloaded, NoSites: m.NoSites}
	for _, mod := range mods {
		f.Modules = append(f.Modules, moduleJSON{mod.Rel, mod.Lang, mod.TestFiles, mod.Tests})
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Join(root, ".cull")
	if err := tools.EnsureDir(dir); err != nil {
		return err
	}
	return tools.WriteFileAtomic(filepath.Join(dir, "evidence.json"), data, 0o600)
}

func hashFile(root, rel string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// matrixView is what the matrix says about each module, for moduleGap.
type matrixView struct {
	byModule   map[string][]Mutant
	measuredBy map[string]map[string]bool
	flaky      map[string]bool     // test ids
	shadowed   map[string]string   // module → shadowed/unconfirmed reason
	dropped    map[string][]string // module → dropped test files
	unloaded   map[string][]string // module → tests that passed its canary
	noSites    map[string]bool     // modules with no site to plant a bug in
}

// moduleGap says why module rel is unmeasured for test id (in test file
// file), or "" when it is measured for id or has no site to plant a bug
// in (so it can hide no catch).
func moduleGap(id, file, rel string, v matrixView) string {
	switch {
	case v.shadowed[rel] != "":
		return v.shadowed[rel]
	case slices.Contains(v.dropped[rel], file):
		return reasonDropped
	case slices.Contains(v.unloaded[rel], id):
		return reasonShadowed
	case v.noSites[rel]:
		return ""
	case !v.measuredBy[rel][id]:
		if v.flaky[id] {
			return reasonFlaky
		}
		return reasonNoBaseline
	}
	for _, mu := range v.byModule[rel] {
		if mu.KilledBy != nil {
			return ""
		}
	}
	return reasonExcluded
}

// gateModules leaves a measured ev measured only if every module in mods
// is measured for id (see moduleGap): a catch the target alone makes in
// an unmeasured module would go unseen, so it is never settled on the
// other modules alone. Otherwise ev becomes unmeasured with the first
// such module's reason (its counts are kept).
func gateModules(ev Evidence, id, file string, mods []string, v matrixView) Evidence {
	if !ev.Measured {
		return ev
	}
	for _, rel := range mods {
		if r := moduleGap(id, file, rel, v); r != "" {
			ev.Measured, ev.Reason = false, r
			return ev
		}
	}
	return ev
}
