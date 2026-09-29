package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
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
// {"error", "contexts": {context: [lines]}, "stmts", "tests", "spawned",
// "ignored"} to $CULL_COV_OUT.
//
// Each test (setup, call, teardown of every fixture it uses) runs in its
// own forked child, pytest-forked style, which sends back (pickled, as
// pytest's serializable reports keep tuples) its reports and coverage
// and exits: no test changes the parent's state, so a
// memoized or lazily imported line counts for every test that runs it,
// not only the first. Collection runs once, in the parent (global).
// With CULL_COV_FORK=0 the tests run in order in one process instead
// (still one context per test phase): a test's lines then include what
// it ran because of the tests before it.
//
// An audit hook (and a wrap of multiprocessing's Process.start, whose
// spawn and forkserver methods raise no audit event) sees every process
// the tests start (the plugin's own forks exempt): one started in a
// test's call phase lists that test in "spawned" (its child's work can't
// be measured), one started anywhere else is an error, unless it runs a
// native program that can't run the project's Python (see _benign): that
// one is listed in "ignored".
//
// It starts coverage when pytest loads it (before collection imports
// anything) and refuses (error) if the target was already imported, if
// coverage's tracer (sys.gettrace() and threading's, as set by start)
// was replaced or removed after collection or any test phase, if a test
// phase ran outside a forked child, or if a fork failed or a child died
// without its data. The ctrace core is forced: sys.monitoring stops
// reporting a line after its first hit, which breaks per-test contexts.
const pyCovPlugin = `import json, os, pickle, re, sys, threading, warnings

_target = os.path.realpath(os.environ["CULL_COV_TARGET"])
_root = os.path.realpath(os.environ["CULL_COV_ROOT"])
_out = os.environ["CULL_COV_OUT"]
_fork = os.environ.get("CULL_COV_FORK") != "0"
_err = None
_cov = None
_tests = set()
_tracer = None
_ttracer = None
_phase = None
_spawned = set()
_ignored = set()  # native programs started outside a call phase
_own = False  # the plugin itself is forking
_child = False  # this process is a test's forked child
_ctxs = {}  # context -> lines, from the children
_n = 0
_SPAWN = frozenset(["subprocess.Popen", "os.fork", "os.forkpty", "os.posix_spawn", "os.posix_spawnp", "os.system", "os.startfile"])
# Interpreters, shells, launchers and wrappers: native, but they run
# anything.
_LAUNCHERS = frozenset(
    "python pypy node deno bun go sh bash zsh dash fish env make uv uvx poetry pipx pip tox nox npm npx pnpm yarn pytest git"
    " nice nohup timeout gtimeout time xargs sudo doas stdbuf caffeinate pixi micromamba mamba conda rye hatch pdm"
    " docker podman nix just task bazel cargo".split()
)
# Bytes of a native program that embeds Python.
_EMBEDS = (b"Py_Initialize", b"Py_Main", b"libpython")


def _find(name, path, cwd):
    """The executable file a spawn of name runs: name itself if it has a
    directory part (or path is None), else the first executable
    candidate in the path list (subprocess's search); relative ones
    against cwd. None if there is none."""
    if path is None or os.path.dirname(name):
        cands = [name]
    else:
        cands = [os.path.join(d, name) for d in path]
    for c in cands:
        c = os.path.join(cwd, c)
        if os.path.isfile(c) and os.access(c, os.X_OK):
            return c
    return None


def _native(p):
    """Whether file p is a native executable by its magic bytes: Mach-O
    (thin or fat), ELF or PE. A #! script is not."""
    with open(p, "rb") as fh:
        h = fh.read(8)
    if h[:4] in (b"\xfe\xed\xfa\xce", b"\xfe\xed\xfa\xcf", b"\xce\xfa\xed\xfe", b"\xcf\xfa\xed\xfe", b"\x7fELF"):
        return True
    if len(h) == 8 and h[:4] in (b"\xca\xfe\xba\xbe", b"\xca\xfe\xba\xbf"):
        return int.from_bytes(h[4:], "big") < 45  # fat Mach-O arch count, not a Java class version
    if len(h) == 8 and h[:4] in (b"\xbe\xba\xfe\xca", b"\xbf\xba\xfe\xca"):
        return int.from_bytes(h[4:], "little") < 45
    return h[:2] == b"MZ"


def _embeds(p):
    """Whether file p's bytes name Python's embedding API (read, never
    run)."""
    keep = max(len(e) for e in _EMBEDS) - 1
    tail = b""
    with open(p, "rb") as fh:
        while True:
            chunk = fh.read(1 << 20)
            if not chunk:
                return False
            buf = tail + chunk
            if any(e in buf for e in _EMBEDS):
                return True
            tail = buf[-keep:]


def _argv_runs(argv, cwd):
    """Whether any argv element may run code: a launcher's name, a .py
    file, -m, or an existing path inside the project root."""
    for a in argv or ():
        try:
            s = os.fsdecode(a)
        except TypeError:
            s = str(a)
        if s == "-m" or s.lower().endswith(".py") or _launcher(os.path.basename(s)):
            return True
        p = os.path.join(cwd, s)
        if os.path.exists(p):
            r = os.path.realpath(p)
            if r == _root or r.startswith(os.path.join(_root, "")):
                return True
    return False


def _launcher(name):
    b = name.lower()
    if b.endswith(".exe"):
        b = b[:-4]
    return re.sub(r"[-_]?[0-9][0-9.]*t?$", "", b) in _LAUNCHERS


def _benign(event, args):
    """The real path of the program a spawn outside a call phase runs, if
    it can't run the project's Python: a native executable (by magic
    bytes) outside the project root (symlinks resolved), not an
    interpreter, shell, launcher or wrapper, not embedding Python, and
    given no argument that may run code (see _argv_runs). None if it
    may, or can't be told."""
    try:
        cwd = os.getcwd()
        if event == "subprocess.Popen":
            exe, argv, pcwd, env = args
            name = os.fsdecode(exe)
            if pcwd is not None:
                cwd = os.path.join(cwd, os.fsdecode(pcwd))
            found = _find(name, os.get_exec_path(env), cwd)
        elif event == "os.posix_spawn":
            name, argv = os.fsdecode(args[0]), args[1]
            found = _find(name, None, cwd)
        elif event == "os.posix_spawnp":
            # The C library searches this process's PATH; require the
            # passed env's PATH (if any) to find the same file.
            name, argv, env = os.fsdecode(args[0]), args[1], args[2]
            found = _find(name, os.get_exec_path(), cwd)
            if env is not None and _find(name, os.get_exec_path(env), cwd) != found:
                return None
        else:
            return None
        if found is None:
            return None
        real = os.path.realpath(found)
        if real == _root or real.startswith(os.path.join(_root, "")):
            return None
        if any(_launcher(os.path.basename(p)) for p in (name, found, real)):
            return None
        if isinstance(argv, (str, bytes)) or _argv_runs(argv, cwd):
            return None
        return real if _native(real) and not _embeds(real) else None
    except Exception:
        return None


def _on_spawn(event, args=None):
    global _err
    if _own:
        return
    if _phase is not None and _phase[1] == "call":
        _spawned.add(_phase[0])
    elif _err is None:
        exe = _benign(event, args)
        if exe is not None:
            _ignored.add(exe)
            return
        where = "collection" if _phase is None else _phase[0] + " " + _phase[1]
        _err = "a process was spawned outside a test's call phase (%s during %s)" % (event, where)


def _audit(event, args):
    if event.startswith(("subprocess.", "os.")) and (event in _SPAWN or event.startswith(("os.exec", "os.spawn"))):
        _on_spawn(event, args)


sys.addaudithook(_audit)
try:
    import multiprocessing.process as _mp

    _mp_start = _mp.BaseProcess.start

    def _start(self, *a, **k):
        _on_spawn("multiprocessing.Process.start")
        return _mp_start(self, *a, **k)

    _mp.BaseProcess.start = _start
except Exception:
    pass


def _thread_trace():
    g = getattr(threading, "gettrace", None)
    return g() if g else None


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
        _tracer = sys.gettrace()
        _ttracer = _thread_trace()
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


def _check(where):
    global _err
    if _cov is None or _err is not None:
        return
    if sys.gettrace() is not _tracer or _thread_trace() is not _ttracer:
        _err = "coverage's tracer was replaced or removed during " + where


def pytest_collection_finish(session):
    _check("collection")


def _enter(item, phase):
    global _phase, _err
    if _cov is not None and _fork and not _child and _err is None:
        _err = "%s ran outside a forked child (another plugin ran the test protocol)" % _sel(item)
    _phase = (_sel(item), phase)
    _switch(_sel(item) + "|" + phase)


def _leave(item, phase):
    global _phase
    _phase = None
    _check(_sel(item) + " " + phase)


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_setup(item):
    _enter(item, "setup")
    yield
    _leave(item, "setup")


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_call(item):
    _tests.add(_sel(item))
    _enter(item, "call")
    yield
    _leave(item, "call")


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_teardown(item, nextitem):
    _enter(item, "teardown")
    yield
    _leave(item, "teardown")
    _switch("")


def _contexts():
    """Stops coverage and returns its target lines by context."""
    _cov.stop()
    data = _cov.get_data()
    ctxs = {}
    for f in data.measured_files():
        if os.path.realpath(f) != _target:
            continue
        for line, cs in (data.contexts_by_lineno(f) or {}).items():
            for c in cs:
                ctxs.setdefault(c, set()).add(line)
    return ctxs


def _run_child(item, res):
    from _pytest.runner import runtestprotocol

    reports = runtestprotocol(item, log=False, nextitem=None)
    ctxs = _contexts()
    hook = item.config.hook
    d = {
        "error": _err,
        "reports": [hook.pytest_report_to_serializable(config=item.config, report=r) for r in reports],
        "contexts": {c: sorted(ls) for c, ls in ctxs.items()},
        "tests": sorted(_tests),
        "spawned": sorted(_spawned),
        "ignored": sorted(_ignored),
    }
    with open(res + ".tmp", "wb") as fh:
        try:
            pickle.dump(d, fh)
        except Exception:
            fh.seek(0)
            fh.truncate()
            for r in d["reports"]:
                r["user_properties"] = [(str(k), str(v)) for k, v in r.get("user_properties", [])]
            pickle.dump(d, fh)
    os.rename(res + ".tmp", res)


@pytest.hookimpl(tryfirst=True)
def pytest_runtest_protocol(item, nextitem):
    global _err, _own, _child, _n
    if _cov is None or _err is not None or not _fork:
        return None
    if not hasattr(os, "fork"):
        _err = "per-test isolation needs os.fork"
        return None
    _n += 1
    res = "%s.%d.pickle" % (_out, _n)
    for f in (sys.stdout, sys.stderr):
        try:
            f.flush()
        except Exception:
            pass
    ihook = item.ihook
    ihook.pytest_runtest_logstart(nodeid=item.nodeid, location=item.location)
    _own = True
    try:
        with warnings.catch_warnings():
            warnings.simplefilter("ignore")
            pid = os.fork()
    except Exception as e:
        _own = False
        _err = "per-test isolation: fork failed (%s)" % e
        ihook.pytest_runtest_logfinish(nodeid=item.nodeid, location=item.location)
        return True
    if pid == 0:
        _own = False
        _child = True
        code = 0
        try:
            _run_child(item, res)
        except BaseException:
            code = 70
        os._exit(code)
    _own = False
    _, status = os.waitpid(pid, 0)
    try:
        with open(res, "rb") as fh:
            d = pickle.load(fh)
        os.remove(res)
    except Exception:
        d = None
    if d is None:
        if _err is None:
            _err = "the forked run of %s died without its data (wait status %d)" % (_sel(item), status)
    else:
        if d["error"] and _err is None:
            _err = d["error"]
        for c, ls in d["contexts"].items():
            _ctxs.setdefault(c, set()).update(ls)
        _tests.update(d["tests"])
        _spawned.update(d["spawned"])
        _ignored.update(d["ignored"])
        for r in d["reports"]:
            rep = item.config.hook.pytest_report_from_serializable(config=item.config, data=r)
            ihook.pytest_runtest_logreport(report=rep)
    ihook.pytest_runtest_logfinish(nodeid=item.nodeid, location=item.location)
    return True


def pytest_unconfigure(config):
    if _child:
        return
    out = {"error": _err, "contexts": {}, "stmts": [], "tests": sorted(_tests), "spawned": sorted(_spawned), "ignored": sorted(_ignored)}
    if _cov is not None:
        try:
            ctxs = _contexts()
            for c, ls in _ctxs.items():
                ctxs.setdefault(c, set()).update(ls)
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
	Spawned  []string         `json:"spawned"`
	Ignored  []string         `json:"ignored"`
}

// Cover measures per-test coverage of module with pyCovPlugin, in two
// passes over the files: each test isolated in a forked child, and all
// tests in order in one process. A test covers a line if either pass
// says so, and passed only if it passed in both. A unit is a test
// function (all its parametrizations; a mutant it covers runs its whole
// file, see EachRunner); lines are compared by the statement holding them
// (coverage.py's statement starts), so a continuation line counts with
// its statement. Lines run outside any test's call phase (imports,
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
	defer func() { _ = os.RemoveAll(tmp) }()
	if err := os.WriteFile(filepath.Join(tmp, "cull_cov_plugin.py"), []byte(pyCovPlugin), 0o600); err != nil {
		return nil, err
	}
	var env []string
	for _, kv := range pythonEnv(dir) {
		if old, ok := strings.CutPrefix(kv, "PYTHONPATH="); ok {
			kv = "PYTHONPATH=" + tmp + string(os.PathListSeparator) + old
		}
		env = append(env, kv)
	}
	env = append(env, "CULL_COV_TARGET="+target, "CULL_COV_ROOT="+absDir)
	index := map[string]string{}
	for _, idx := range files {
		maps.Copy(index, idx)
	}
	n := 0
	// pass runs the plugin over testFiles, forking per test or in order
	// (each pass with its own timeout).
	pass := func(testFiles []string, fork bool) (pyCovData, map[string]bool, error) {
		n++
		outPath := filepath.Join(tmp, fmt.Sprintf("out%d.json", n))
		junit := filepath.Join(tmp, fmt.Sprintf("junit%d.xml", n))
		mode := "CULL_COV_FORK=1"
		if !fork {
			mode = "CULL_COV_FORK=0"
		}
		argv := append([]string{verify.PythonInterpreter(dir), "-m", "pytest", "-q", "-p", "no:cacheprovider", "-p", "no:randomly", "-p", "cull_cov_plugin", "--junitxml=" + junit}, testFiles...)
		res := runProcessEnv(ctx, dir, argv, append(slices.Clone(env), "CULL_COV_OUT="+outPath, mode), timeout)
		var d pyCovData
		switch {
		case res.StartErr != nil:
			return d, nil, res.StartErr
		case ctx.Err() != nil:
			return d, nil, ctx.Err()
		case res.TimedOut:
			return d, nil, fmt.Errorf("the coverage run timed out")
		case res.ExitCode == 2 || res.ExitCode == 3 || res.ExitCode == 4:
			return d, nil, fmt.Errorf("pytest exited %d in the coverage run", res.ExitCode)
		}
		b, err := os.ReadFile(outPath)
		if err != nil {
			return d, nil, fmt.Errorf("the coverage plugin wrote no data (%s)", lastLineOf(res.Stderr, res.Stdout))
		}
		if err := json.Unmarshal(b, &d); err != nil {
			return d, nil, fmt.Errorf("unreadable coverage data: %v", err)
		}
		if d.Error != nil {
			return d, nil, fmt.Errorf("%s", *d.Error)
		}
		if len(d.Contexts) == 0 {
			return d, nil, fmt.Errorf("no line of %s ran in the coverage run", module)
		}
		o := newOutcome()
		readPyJUnit(dir, junit, index, &o)
		return d, o.Passed, nil
	}
	names := sortedFiles(files)
	iso, isoPassed, err := pass(names, true)
	if err != nil {
		return nil, err
	}
	ord, ordPassed, err := pass(names, false)
	if err != nil {
		return nil, err
	}
	// Cross-file check: each file alone, in order, must run the same
	// lines per test as after the files before it; its tests must pass.
	alonePassed := ordPassed
	if len(names) > 1 {
		alonePassed = map[string]bool{}
		for _, f := range names {
			alone, passed, err := pass([]string{f}, false)
			if err != nil {
				return nil, err
			}
			if t := pyCallDiff(ord, alone, f); t != "" {
				return nil, fmt.Errorf("%s runs other lines alone than after the test files before it", t)
			}
			maps.Copy(alonePassed, passed)
		}
	}
	d := mergePyCov(iso, ord)
	cov := pyCoverage(src, files, d)
	cov.Passed = map[string]bool{}
	for id := range isoPassed {
		if ordPassed[id] && alonePassed[id] {
			cov.Passed[id] = true
		}
	}
	cov.Ignored = d.Ignored
	return cov, nil
}

// pyCallDiff is the first test of file whose call phase ran other lines
// in a than in b ("" if none).
func pyCallDiff(a, b pyCovData, file string) string {
	calls := func(d pyCovData) map[string][]int {
		out := map[string][]int{}
		for c, ls := range d.Contexts {
			sel, phase, _ := cutLast(c, "|")
			if phase == "call" && strings.HasPrefix(sel, file+"::") {
				ls = slices.Clone(ls)
				sort.Ints(ls)
				out[sel] = slices.Compact(ls)
			}
		}
		return out
	}
	ca, cb := calls(a), calls(b)
	for _, m := range []map[string][]int{ca, cb} {
		for _, sel := range sortedKeys(m) {
			if !slices.Equal(ca[sel], cb[sel]) {
				return sel
			}
		}
	}
	return ""
}

// mergePyCov is the union of two passes' data.
func mergePyCov(a, b pyCovData) pyCovData {
	union := func(x, y []string) []string {
		set := map[string]bool{}
		for _, s := range append(slices.Clone(x), y...) {
			set[s] = true
		}
		out := slices.Collect(maps.Keys(set))
		sort.Strings(out)
		return out
	}
	d := pyCovData{Contexts: map[string][]int{}, Stmts: a.Stmts, Tests: union(a.Tests, b.Tests), Spawned: union(a.Spawned, b.Spawned), Ignored: union(a.Ignored, b.Ignored)}
	for _, src := range []map[string][]int{a.Contexts, b.Contexts} {
		for c, ls := range src {
			d.Contexts[c] = append(d.Contexts[c], ls...)
		}
	}
	return d
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
	spawned := map[string]bool{}
	for _, sel := range d.Spawned {
		spawned[unitSel(sel)] = true
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
		hit := hits(byUnit[sel], false)
		if spawned[sel] {
			hit = func(int, int) bool { return true }
		}
		cov.Units = append(cov.Units, Unit{File: file, IDs: []string{id}, hit: hit})
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
