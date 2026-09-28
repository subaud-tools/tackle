package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/schuettc/tackle/internal/cull/verify"
)

// pyCovPlugin is a pytest plugin measuring per-test line coverage of one
// file ($CULL_COV_TARGET) with coverage.py dynamic contexts (the
// mechanism pytest-cov's --cov-context=test uses): each test's setup,
// call and teardown get their own context, named by the test's node id
// with its file path relative to the run's working directory. It writes
// {"error", "contexts": {context: [lines]}, "stmts", "tests"} to
// $CULL_COV_OUT. It starts coverage when pytest loads it (before
// collection imports anything) and refuses (error) if the target was
// already imported. The ctrace core is forced: sys.monitoring stops
// reporting a line after its first hit, which breaks per-test contexts.
const pyCovPlugin = `import json, os, sys

_target = os.path.realpath(os.environ["CULL_COV_TARGET"])
_out = os.environ["CULL_COV_OUT"]
_err = None
_cov = None
_tests = set()
try:
    import coverage
except Exception as e:
    _err = "coverage.py is not importable (%s)" % e
if _err is None:
    for _m in list(sys.modules.values()):
        _f = getattr(_m, "__file__", None)
        if _f and os.path.realpath(_f) == _target:
            _err = "the module was imported before coverage could start"
            break
if _err is None:
    try:
        os.environ["COVERAGE_CORE"] = "ctrace"
        _cov = coverage.Coverage(data_file=None, config_file=False, include=[_target])
        _cov.start()
    except Exception as e:
        _err = "coverage.py failed to start (%s)" % e
        _cov = None

import pytest


def _sel(item):
    path = getattr(item, "path", None) or item.fspath
    rel = os.path.relpath(str(path), os.getcwd()).replace(os.sep, "/")
    nodeid = item.nodeid
    i = nodeid.find("::")
    return rel + (nodeid[i:] if i >= 0 else "")


def _switch(ctx):
    if _cov is not None:
        _cov.switch_context(ctx)


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_setup(item):
    _switch(_sel(item) + "|setup")
    yield


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_call(item):
    _tests.add(_sel(item))
    _switch(_sel(item) + "|call")
    yield


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_teardown(item, nextitem):
    _switch(_sel(item) + "|teardown")
    yield
    _switch("")


def pytest_unconfigure(config):
    out = {"error": _err, "contexts": {}, "stmts": [], "tests": sorted(_tests)}
    if _cov is not None:
        try:
            _cov.stop()
            data = _cov.get_data()
            ctxs = {}
            for f in data.measured_files():
                if os.path.realpath(f) != _target:
                    continue
                for line, cs in (data.contexts_by_lineno(f) or {}).items():
                    for c in cs:
                        ctxs.setdefault(c, set()).add(line)
            out["contexts"] = {c: sorted(ls) for c, ls in ctxs.items()}
            _cov.clear_exclude()
            out["stmts"] = sorted(_cov.analysis2(_target)[1])
        except Exception as e:
            out["error"] = "coverage.py failed (%s)" % e
    with open(_out, "w") as fh:
        json.dump(out, fh)
`

type pyCovData struct {
	Error    *string          `json:"error"`
	Contexts map[string][]int `json:"contexts"`
	Stmts    []int            `json:"stmts"`
	Tests    []string         `json:"tests"`
}

// Cover measures per-test coverage of module with pyCovPlugin. A unit
// is a test function (all its parametrizations, selected by its node id
// without the [...] suffix); lines are compared by the statement holding
// them (coverage.py's statement starts), so a continuation line counts
// with its statement. Lines run outside any test's call phase (imports,
// collection, fixture setup and teardown) are global.
func (r pythonRunner) Cover(ctx context.Context, dir string, files map[string]map[string]string, module string, timeout time.Duration) (*Coverage, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	target := filepath.Join(absDir, filepath.FromSlash(module))
	src, err := os.ReadFile(target)
	if err != nil {
		return nil, err
	}
	tmp, err := os.MkdirTemp("", "cull-pycov-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)
	if err := os.WriteFile(filepath.Join(tmp, "cull_cov_plugin.py"), []byte(pyCovPlugin), 0o600); err != nil {
		return nil, err
	}
	outPath := filepath.Join(tmp, "out.json")
	var env []string
	for _, kv := range pythonEnv(dir) {
		if old, ok := strings.CutPrefix(kv, "PYTHONPATH="); ok {
			kv = "PYTHONPATH=" + tmp + string(os.PathListSeparator) + old
		}
		env = append(env, kv)
	}
	env = append(env, "CULL_COV_TARGET="+target, "CULL_COV_OUT="+outPath)
	argv := append([]string{verify.PythonInterpreter(dir), "-m", "pytest", "-q", "-p", "no:cacheprovider", "-p", "cull_cov_plugin"}, sortedFiles(files)...)
	res := runProcessEnv(ctx, dir, argv, env, timeout)
	switch {
	case res.StartErr != nil:
		return nil, res.StartErr
	case ctx.Err() != nil:
		return nil, ctx.Err()
	case res.TimedOut:
		return nil, fmt.Errorf("the coverage run timed out")
	case res.ExitCode == 2 || res.ExitCode == 3 || res.ExitCode == 4:
		return nil, fmt.Errorf("pytest exited %d in the coverage run", res.ExitCode)
	}
	b, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("the coverage plugin wrote no data (%s)", lastLineOf(res.Stderr, res.Stdout))
	}
	var d pyCovData
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, fmt.Errorf("unreadable coverage data: %v", err)
	}
	if d.Error != nil {
		return nil, fmt.Errorf("%s", *d.Error)
	}
	if len(d.Contexts) == 0 {
		return nil, fmt.Errorf("no line of %s ran in the coverage run", module)
	}
	return pyCoverage(src, files, d), nil
}

// pyCoverage builds a Coverage from the plugin's data.
func pyCoverage(src []byte, files map[string]map[string]string, d pyCovData) *Coverage {
	starts := lineStarts(src)
	stmts := append([]int(nil), d.Stmts...)
	sort.Ints(stmts)
	// stmtOf is the statement holding line l: the last statement start at
	// or before it (0: none).
	stmtOf := func(l int) int {
		i := sort.SearchInts(stmts, l+1)
		if i == 0 {
			return 0
		}
		return stmts[i-1]
	}
	global := map[int]bool{}
	byUnit := map[string]map[int]bool{}
	for _, sel := range d.Tests {
		if byUnit[unitSel(sel)] == nil {
			byUnit[unitSel(sel)] = map[int]bool{}
		}
	}
	for c, lines := range d.Contexts {
		sel, phase, _ := cutLast(c, "|")
		set := global
		if phase == "call" && sel != "" {
			u := unitSel(sel)
			if byUnit[u] == nil {
				byUnit[u] = map[int]bool{}
			}
			set = byUnit[u]
		}
		for _, l := range lines {
			set[stmtOf(l)] = true
		}
	}
	hits := func(set map[int]bool, zero bool) func(int, int) bool {
		return func(start, end int) bool {
			for l := lineOf(starts, start); l <= lastLine(starts, start, end); l++ {
				s := stmtOf(l)
				if (s == 0 && zero) || set[s] {
					return true
				}
			}
			return false
		}
	}
	cov := &Coverage{global: hits(global, true)}
	sels := make([]string, 0, len(byUnit))
	for u := range byUnit {
		sels = append(sels, u)
	}
	sort.Strings(sels)
	for _, sel := range sels {
		file, rest, ok := strings.Cut(sel, "::")
		if !ok {
			continue
		}
		id, ok := pyUnitID(files[file], rest)
		if !ok {
			continue
		}
		cov.Units = append(cov.Units, Unit{Sel: sel, File: file, IDs: []string{id}, hit: hits(byUnit[sel], false)})
	}
	return cov
}

// unitSel is a node id selector without its parametrize suffix.
func unitSel(sel string) string {
	file, rest, ok := strings.Cut(sel, "::")
	if !ok {
		return sel
	}
	return file + "::" + stripParam(rest)
}

// pyUnitID maps a node id's part after the file ("Class::Inner::func")
// to its cull id through a file's index (native keys "Class.Inner::func").
func pyUnitID(index map[string]string, rest string) (string, bool) {
	parts := strings.Split(rest, "::")
	native := parts[len(parts)-1]
	if len(parts) > 1 {
		native = strings.Join(parts[:len(parts)-1], ".") + "::" + native
	}
	id, ok := index[native]
	return id, ok
}

// RunUnits runs the units' node ids (deduplicated) like Run.
func (r pythonRunner) RunUnits(ctx context.Context, dir string, units []Unit, index map[string]string, timeout time.Duration) Outcome {
	seen := map[string]bool{}
	var sels []string
	for _, u := range units {
		if !seen[u.Sel] {
			seen[u.Sel] = true
			sels = append(sels, u.Sel)
		}
	}
	return r.Run(ctx, dir, sels, index, timeout)
}

func cutLast(s, sep string) (before, after string, found bool) {
	i := strings.LastIndex(s, sep)
	if i < 0 {
		return s, "", false
	}
	return s[:i], s[i+len(sep):], true
}

// lastLineOf is the last non-empty line of the first non-empty buffer.
func lastLineOf(bufs ...[]byte) string {
	for _, b := range bufs {
		if s := strings.TrimSpace(string(b)); s != "" {
			return s[strings.LastIndex(s, "\n")+1:]
		}
	}
	return "no output"
}
