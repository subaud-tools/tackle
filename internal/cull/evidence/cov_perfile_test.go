package evidence

import "testing"

// The reviewer's set_mode repro: test_1 sets mode 1 and runs f without
// checking it, test_2 resets mode 0, test_3 checks f only in mode 0. A
// mutant in f selects test_1 and test_3; run together, test_3 would see
// mode 1 and catch nothing. Each selected file runs in its own process,
// so test_3 stays a sole guard, as in the full run.
const modePySrc = "_mode = [0]\n\n\ndef set_mode(m):\n    _mode[0] = m\n\n\ndef f(x):\n    return x > 0\n"

func TestCoverageSelectedFilesSeparatePython(t *testing.T) {
	needCoverage(t)
	root := newRepo(t, map[string]string{
		"calc.py":   modePySrc,
		"test_1.py": "from calc import f, set_mode\n\n\ndef test_a():\n    set_mode(1)\n    f(1)\n",
		"test_2.py": "from calc import set_mode\n\n\ndef test_b():\n    set_mode(0)\n",
		"test_3.py": "import calc\nfrom calc import f\n\n\ndef test_c():\n    if calc._mode[0] == 0:\n        assert f(1)\n        assert not f(0)\n",
	})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"test_1.py", "test_2.py", "test_3.py"}, map[string]verdict{"test_a": band, "test_b": band, "test_c": band})
	settleExpect(t, root, "")
	wantSole(t, root, "test_c")
}

const modeGoSrc = "package calc\n\nvar mode int\n\nfunc SetMode(m int) { mode = m }\n\nfunc F(x int) bool {\n\treturn x > 0\n}\n"

// The reviewer's os.Exit repro: a per-file `go tool test2json` invocation
// that skips go test's own binary flags (-test.paniconexit0 in
// particular) lets a mutant's os.Exit(0) look like a pass instead of the
// kill the full run would record, wrongly cutting the guarding test.
const exitGoSrc = "package calc\n\nimport \"os\"\n\nfunc Stop(x int) {\n\tif x > 100 {\n\t\tos.Exit(0)\n\t}\n}\n\nfunc Add(a, b int) int {\n\treturn a + b\n}\n"

func TestCoverageSelectedFilesSeparateGo(t *testing.T) {
	needTools(t)
	root := newRepo(t, map[string]string{
		"go.mod":    "module fixture\n\ngo 1.22\n",
		"calc.go":   modeGoSrc,
		"a_test.go": "package calc\n\nimport \"testing\"\n\nfunc TestA(t *testing.T) {\n\tSetMode(1)\n\tF(1)\n}\n",
		"b_test.go": "package calc\n\nimport \"testing\"\n\nfunc TestB(t *testing.T) { SetMode(0) }\n",
		"c_test.go": "package calc\n\nimport \"testing\"\n\nfunc TestC(t *testing.T) {\n\tif mode == 0 && (!F(1) || F(0)) {\n\t\tt.Fatal(\"c\")\n\t}\n}\n",
	})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"a_test.go", "b_test.go", "c_test.go"}, map[string]verdict{"TestA": band, "TestB": band, "TestC": band})
	settleExpect(t, root, "")
	wantSole(t, root, "TestC")
}

// The reviewer's os.Exit repro: a mutant that makes Stop's guard trip on
// Stop(5) calls os.Exit(0) mid-test. go test always passes
// -test.paniconexit0 to the test binary, which turns that os.Exit(0)
// into a failure instead of a silent, successful-looking process exit.
// A per-file `go tool test2json` invocation that omits the flag misses
// that failure and wrongly cuts TestS.
func TestCoverageSelectedFilesSeparateGoPanicOnExit0(t *testing.T) {
	needTools(t)
	root := newRepo(t, map[string]string{
		"go.mod":    "module fixture\n\ngo 1.22\n",
		"calc.go":   exitGoSrc,
		"s_test.go": "package calc\n\nimport \"testing\"\n\nfunc TestS(t *testing.T) {\n\tStop(5)\n}\n",
		"a_test.go": "package calc\n\nimport \"testing\"\n\nfunc TestA(t *testing.T) {\n\tif Add(2, 3) != 5 {\n\t\tt.Fatal(\"a\")\n\t}\n}\n",
	})
	band := verdict{"review", "review_band"}
	writeLast(t, root, []string{"s_test.go", "a_test.go"}, map[string]verdict{"TestS": band, "TestA": band})
	settleExpect(t, root, "")
	wantSole(t, root, "TestS")
}
