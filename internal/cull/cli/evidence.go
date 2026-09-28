package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/schuettc/tackle/internal/cull/discover"
	"github.com/schuettc/tackle/internal/cull/evidence"
	tools "github.com/schuettc/tools-common"
)

var evidenceFlags = flags("evidence", "cull evidence [--workers 8] [--max-mutants 30] [--timeout 60s] [--json]",
	"Settle the tests the last `cull check` left in the review band (.cull/last.json) by\n"+
		"planting small bugs in the project code they call, one at a time, in sandbox copies\n"+
		"(the working tree is never modified), and running the tests against each. A test that\n"+
		"is the only one catching some bug is kept; one that catches nothing only it catches is\n"+
		"cut; one that can't be measured stays in review with the reason. No other verdict\n"+
		"changes. Writes .cull/evidence.json and updates .cull/last.json. Makes no Jev calls.\n"+
		"Refuses if a test file changed since check.\n"+
		"Exit: 0 done, 2 refused or error.",
	func(fs *flag.FlagSet) {
		fs.Int("workers", 8, "parallel sandboxes")
		fs.Int("max-mutants", 30, "planted bugs per module (at most)")
		fs.Duration("timeout", 60*time.Second, "per test run timeout")
		fs.Bool("json", false, "print the result as JSON")
	})

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
	jsonMode := boolFlag(fs, "json")
	var progress io.Writer = errw
	if jsonMode {
		progress = io.Discard
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	res, err := evidence.Run(ctx, root, evidence.Options{
		Workers: workers, MaxMutants: maxMutants, RunTimeout: timeout, Stderr: progress,
	})
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
		fmt.Fprintln(w, "nothing to measure: no review-band tests in .cull/last.json")
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
