package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/schuettc/tackle/internal/cull/discover"
	"github.com/schuettc/tackle/internal/cull/evidence"
	tools "github.com/schuettc/tools-common"
)

var evidenceFlags = flags("evidence", "cull evidence [--workers 8] [--max-mutants 30] [--timeout 60s] [--no-cache] [--json]",
	"Settle the tests the last `cull check` left in the review band (.cull/last.json) by\n"+
		"planting small bugs in the project code they call, one at a time, in sandbox copies\n"+
		"(the working tree is never modified), and running the tests against each. A test that\n"+
		"is the only one catching some bug is kept; one that catches nothing only it catches is\n"+
		"cut; one that can't be measured stays in review with the reason. No other verdict\n"+
		"changes. Writes .cull/evidence.json and updates .cull/last.json. Makes no Jev calls.\n"+
		"Each bug runs only the tests whose coverage includes it (all of them when coverage\n"+
		"can't be measured). Module results are cached; a rerun redoes only modules whose\n"+
		"inputs changed (--no-cache: redo all).\n"+
		"Refuses if a test file changed since check.\n"+
		"Exit: 0 done, 2 refused or error, 130/143 interrupted (SIGINT/SIGTERM; nothing written).",
	func(fs *flag.FlagSet) {
		fs.Int("workers", 8, "parallel sandboxes")
		fs.Int("max-mutants", 30, "planted bugs per module whose targets' callees aren't known (at most)")
		fs.Duration("timeout", 60*time.Second, "minimum per mutant run timeout (3× the slower baseline run when longer)")
		fs.Bool("no-cache", false, "ignore cached module results and run everything again")
		fs.Bool("json", false, "print the result as JSON")
	})

// notifySignals is signal.Notify and exitNow os.Exit (tests replace
// them).
var (
	notifySignals = signal.Notify
	exitNow       = os.Exit
)

// syncWriter serializes writes to w.
type syncWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func (s *syncWriter) Write(b []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.w.Write(b)
}

// sigCode is the exit code for signal s (128 + its number).
func sigCode(s os.Signal) int {
	if n, ok := s.(syscall.Signal); ok {
		return 128 + int(n)
	}
	return 130
}

func runEvidence(args []string, out, errw io.Writer) error {
	fs := evidenceFlags()
	pos, err := parse(fs, args, out)
	if err != nil {
		return err
	}
	if len(pos) > 0 {
		return tools.UsageError{Msg: "unexpected arguments"}
	}
	get := func(name string) any { return fs.Lookup(name).Value.(flag.Getter).Get() }
	workers, maxMutants := get("workers").(int), get("max-mutants").(int)
	timeout := get("timeout").(time.Duration)
	switch {
	case workers <= 0:
		return tools.UsageError{Msg: "--workers must be positive"}
	case maxMutants <= 0:
		return tools.UsageError{Msg: "--max-mutants must be positive"}
	case timeout <= 0:
		return tools.UsageError{Msg: "--timeout must be positive"}
	}

	root, err := discover.Root(".")
	if err != nil {
		return tools.Exitf(2, "%v", err)
	}
	// Progress and warnings always go to stderr; --json changes stdout only.
	jsonMode := boolFlag(fs, "json")
	// A signal cancels the run at once (no new test run starts, in-flight
	// ones are killed, nothing is written) and exits 128 + its number. A
	// second one exits at once: every test's process group is killed and
	// every planted module restored first.
	errw = &syncWriter{w: errw}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// done ends the signal watcher; runEvidence returns only once it has
	// (a forced exit in progress finishes first).
	var watching sync.WaitGroup
	defer watching.Wait()
	done := make(chan struct{})
	defer close(done)
	sigs := make(chan os.Signal, 2)
	notifySignals(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(sigs)
	var got atomic.Value // os.Signal
	watching.Add(1)
	go func() {
		defer watching.Done()
		select {
		case s := <-sigs:
			got.Store(s)
			cancel()
		case <-done:
			return
		}
		select {
		case s := <-sigs:
			evidence.Abort()
			_, _ = fmt.Fprintf(errw, "cull evidence: interrupted again (%v); exiting now\n", s)
			exitNow(sigCode(s))
		case <-done:
		}
	}()

	res, err := evidence.Run(ctx, root, evidence.Options{
		Workers: workers, MaxMutants: maxMutants, RunTimeout: timeout, Stderr: errw,
		NoCache: boolFlag(fs, "no-cache"),
	})
	if s, ok := got.Load().(os.Signal); ok && err != nil {
		return tools.Exitf(sigCode(s), "cull evidence: interrupted (%v); nothing written", s)
	}
	if err != nil {
		return tools.Exitf(2, "%v", err)
	}
	if jsonMode {
		return tools.PrintJSON(out, res)
	}
	writeEvidence(out, res)
	return nil
}

// writeEvidence prints the settled counts, the unmeasured tests with their
// reasons, and what ran.
func writeEvidence(w io.Writer, r evidence.Result) {
	if len(r.Rules) == 0 {
		fmt.Fprintln(w, "nothing to measure: no review-band or unmeasured tests in .cull/last.json")
		return
	}
	n := map[string]int{}
	for _, rule := range r.Rules {
		n[rule]++
	}
	fmt.Fprintf(w, "kept as sole guards: %d\n", n[evidence.RuleSoleGuard])
	fmt.Fprintf(w, "cut as redundant: %d\n", n[evidence.RuleRedundant])
	fmt.Fprintf(w, "cut as catching nothing: %d\n", n[evidence.RuleCatchesNothing])
	fmt.Fprintf(w, "unmeasured: %d\n", len(r.Unmeasured))
	for _, id := range r.Unmeasured {
		fmt.Fprintf(w, "  %s: %s\n", id, r.Reasons[id])
	}
	fmt.Fprintf(w, "modules: %d, mutants: %d, duration: %s\n", r.Modules, r.Mutants, r.Duration.Round(time.Millisecond))
}
