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
	"sort"
	"strings"
	"time"

	"github.com/schuettc/tackle/internal/cull/cases"
	"github.com/schuettc/tackle/internal/cull/check"
	"github.com/schuettc/tackle/internal/cull/discover"
	"github.com/schuettc/tackle/internal/cull/extract"
	_ "github.com/schuettc/tackle/internal/cull/extract/golang"
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
)

// extractContext is large enough that extraction never drops a callee for
// its context budget: evidence needs every callee file, not a prompt.
const extractContext = 1 << 24

// afterMatrixHook, if set (tests only), runs right after the matrix.
var afterMatrixHook func()

// Run settles the review-band tests of root's .cull/last.json by planted
// bugs. Targets are tests with verdict "review" and rule "review_band".
// It refuses (an error, nothing written) if a target's file changed since
// cull check, or if a target file, a mutated module or one of its test
// files changed while the matrix ran. It mutates only non-test source
// files the targets call (their extracted callees), runs every test in
// those modules' test files against each mutant, and settles each target:
// sole guard of any bug → keep; measured with no sole guard → cut;
// unmeasured → stays review (rule evidence_unmeasured). No other test is
// touched. It writes .cull/evidence.json (the matrix) and rewrites
// .cull/last.json, both atomically and 0600. No target: nothing is run or
// written.
func Run(ctx context.Context, root string, opt Options) (Result, error) {
	start := time.Now()
	res := Result{Settled: map[string]string{}, Rules: map[string]string{}, Reasons: map[string]string{}, Unmeasured: []string{}}
	if opt.Stderr == nil {
		opt.Stderr = io.Discard
	}
	rep, err := check.LoadReport(root)
	if err != nil {
		return res, err
	}
	var targets []int
	for i, t := range rep.Tests {
		if t.Err == "" && t.Verdict == "review" && t.Rule == "review_band" {
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
		ev := measure(t.ID, testMods[t.ID], byModule, m.Baseline)
		if !ev.Measured {
			ev.Reason = unmeasuredReason(root, t.ID, byID, testMods[t.ID], byModule, m.Baseline, flaky, hasSites)
		}
		switch {
		case ev.Measured && ev.Sole > 0:
			t.Verdict, t.Rule = "keep", RuleSoleGuard
		case ev.Measured && ev.Caught > 0:
			t.Verdict, t.Rule = "cut", RuleRedundant
		case ev.Measured:
			t.Verdict, t.Rule = "cut", RuleCatchesNothing
		default:
			t.Rule = RuleUnmeasured
		}
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

	if err := writeEvidenceJSON(root, mods, m); err != nil {
		return res, err
	}
	if err := check.WriteReport(root, rep); err != nil {
		return res, err
	}
	res.Duration = time.Since(start)
	return res, nil
}

// measure computes id's evidence over the non-excluded mutants of mods.
// Measured needs id in the baseline and at least one such mutant.
func measure(id string, mods []string, byModule map[string][]Mutant, baseline map[string]bool) Evidence {
	var ev Evidence
	others := map[string]bool{}
	for _, rel := range mods {
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

// unmeasuredReason says why id has no usable measurement.
func unmeasuredReason(root, id string, byID map[string]cases.TestCase, mods []string,
	byModule map[string][]Mutant, baseline, flaky map[string]bool, hasSites func(string) bool) string {
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
	if flaky[id] {
		return reasonFlaky
	}
	total := 0
	for _, rel := range mods {
		total += len(byModule[rel])
	}
	if total == 0 {
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
// claims as a test file.
func projectSource(root, rel string) bool {
	if rel == "" || path.IsAbs(rel) || filepath.IsAbs(rel) || path.Clean(rel) != rel ||
		rel == ".." || strings.HasPrefix(rel, "../") {
		return false
	}
	if extract.ForFile(rel) != nil {
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
	Modules  []moduleJSON    `json:"modules"`
	Mutants  []Mutant        `json:"mutants"`
	Baseline map[string]bool `json:"baseline"`
	Flaky    []string        `json:"flaky"`
	Unmapped int             `json:"unmapped"`
}

type moduleJSON struct {
	Rel       string   `json:"rel"`
	Lang      string   `json:"lang"`
	TestFiles []string `json:"test_files"`
	Tests     []string `json:"tests"`
}

func writeEvidenceJSON(root string, mods []Module, m Matrix) error {
	f := evidenceFile{Mutants: m.Mutants, Baseline: m.Baseline, Flaky: m.Flaky, Unmapped: m.Unmapped}
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
