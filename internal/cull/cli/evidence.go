package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
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

// notifySignals is signal.Notify (tests replace it).
var notifySignals = signal.Notify

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
	// ones are killed, nothing is written) and exits 128 + its number.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigs := make(chan os.Signal, 1)
	notifySignals(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(sigs)
	var got atomic.Value // os.Signal
	go func() {
		select {
		case s := <-sigs:
			got.Store(s)
			cancel()
		case <-ctx.Done():
		}
	}()

	res, err := evidence.Run(ctx, root, evidence.Options{
		Workers: workers, MaxMutants: maxMutants, RunTimeout: timeout, Stderr: errw,
		NoCache: boolFlag(fs, "no-cache"),
	})
	if s, ok := got.Load().(os.Signal); ok && err != nil {
		code := 130
		if n, ok := s.(syscall.Signal); ok {
			code = 128 + int(n)
		}
		return tools.Exitf(code, "cull evidence: interrupted (%v); nothing written", s)
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
