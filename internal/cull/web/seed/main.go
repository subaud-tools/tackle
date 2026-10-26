// Command seed is a dev-only helper for the review page's probe: it opens the
// cull store under a CULL_HOME, registers a project and records ONE run whose
// items are those of a GET /api/review fixture, so the real API reproduces
// the fixture. It is not part of the shipped binary.
//
//	go run ./internal/cull/web/seed [-mutate id] [-drop id] [-sends] <CULL_HOME> <review.json> <project root>
//
// -mutate changes the named item's hash and -drop leaves it out: use them on a
// second call to record a newer run the page has not seen. -sends then
// records one send in each state the page shows, through the store's own
// claims, reports and replies.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/subaud-tools/tackle/internal/cull/store"
	"github.com/subaud-tools/tools-common/channelmcp"
	"github.com/subaud-tools/tools-common/msgstate"
)

type fixtureItem struct {
	ID      string          `json:"id"`
	Kind    string          `json:"kind"`
	Hash    string          `json:"hash"`
	File    string          `json:"file"`
	Name    string          `json:"name"`
	Verdict string          `json:"verdict"`
	Rule    string          `json:"rule"`
	Model   string          `json:"model"`
	State   json.RawMessage `json:"state"`
	Jev     json.RawMessage `json:"jev"`
	Rows    [][]string      `json:"rows"`
	Members []string        `json:"members"`
}

type fixture struct {
	Run struct {
		Mode    string         `json:"mode"`
		Base    string         `json:"base"`
		Total   int            `json:"total"`
		Summary map[string]int `json:"summary"`
	} `json:"run"`
	Items []fixtureItem `json:"items"`
}

// Seed records the fixture's run for root in the store at path, applying the
// optional mutate/drop, and returns the project id and the run id.
func Seed(ctx context.Context, path string, data []byte, root, mutate, drop string) (int64, int64, error) {
	var fx fixture
	if err := json.Unmarshal(data, &fx); err != nil {
		return 0, 0, fmt.Errorf("fixture: %w", err)
	}
	st, err := store.Open(ctx, path)
	if err != nil {
		return 0, 0, err
	}
	defer func() { _ = st.Close() }()
	p, err := st.Project(ctx, root)
	if err != nil {
		return 0, 0, err
	}
	items := make([]store.Item, 0, len(fx.Items))
	for _, it := range fx.Items {
		if it.ID == drop {
			continue
		}
		if it.ID == mutate {
			it.Hash = "sha256:mutated" + it.Hash[len("sha256:"):]
		}
		items = append(items, store.Item{ID: it.ID, Kind: it.Kind, Hash: it.Hash, File: it.File, Name: it.Name,
			Verdict: it.Verdict, Rule: it.Rule, State: it.State, Jev: it.Jev, Model: it.Model, Rows: it.Rows, Members: it.Members})
	}
	run, err := st.RecordRun(ctx, store.Run{ProjectID: p.ID, Mode: fx.Run.Mode, Base: fx.Run.Base, QuestionsHash: "probe",
		At: time.Now(), Total: fx.Run.Total, Summary: fx.Run.Summary}, items)
	return p.ID, run, err
}

// SeedSends records, oldest first, a send with no reply, one held while the
// agent is mid-turn, one with the agent, one working, answered, declined
// and failed, and last one no session has claimed. Each is the fixture's
// first test answered keep again, sent, claimed by the session "seed" and
// moved as the harness and the agent would move it.
func SeedSends(ctx context.Context, path string, data []byte, projectID, runID int64) error {
	var fx fixture
	if err := json.Unmarshal(data, &fx); err != nil {
		return fmt.Errorf("fixture: %w", err)
	}
	var first *fixtureItem
	for i := range fx.Items {
		if fx.Items[i].Kind == "test" {
			first = &fx.Items[i]
			break
		}
	}
	if first == nil {
		return fmt.Errorf("-sends needs a test in the fixture")
	}
	st, err := store.Open(ctx, path)
	if err != nil {
		return err
	}
	defer func() { _ = st.Close() }()
	const who = "seed"
	landed := func(sd store.Send) error {
		_, err := st.ReportSends(ctx, who, channelmcp.ReportLanded, sd.Claim)
		return err
	}
	// reply lands the send and has the agent reply state.
	reply := func(state msgstate.State, text string) func(store.Send) error {
		return func(sd store.Send) error {
			if err := landed(sd); err != nil {
				return err
			}
			_, err := st.ReplySends(ctx, who, []int64{sd.ID}, state, text)
			return err
		}
	}
	steps := []func(store.Send) error{
		func(sd store.Send) error {
			if err := landed(sd); err != nil {
				return err
			}
			_, err := st.ReportSends(ctx, who, channelmcp.ReportTurnEnded, "")
			return err
		},
		func(sd store.Send) error {
			_, err := st.ReportSends(ctx, who, channelmcp.ReportHeld, sd.Claim)
			return err
		},
		landed,
		reply(msgstate.Working, "cutting the answered tests"),
		reply(msgstate.Answered, "3 tests cut, the suite passes"),
		reply(msgstate.Declined, "Court asked to wait for the next run"),
		reply(msgstate.Failed, "the suite failed after the cut; rolled back"),
		nil, // sent, and no session has it
	}
	for _, step := range steps {
		a := store.Answer{ItemID: first.ID, Hash: first.Hash, Kind: first.Kind, Value: "keep", Via: "item"}
		if err := st.SaveAnswers(ctx, projectID, runID, []store.Answer{a}); err != nil {
			return err
		}
		if _, err := st.Send(ctx, projectID, ""); err != nil {
			return err
		}
		if step == nil {
			continue
		}
		sd, ok, err := st.ClaimSend(ctx, projectID, who, []string{who}, time.Now())
		if err != nil || !ok {
			return fmt.Errorf("claim: %v %w", ok, err)
		}
		if err := step(sd); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	mutate := flag.String("mutate", "", "item id whose hash to change")
	drop := flag.String("drop", "", "item id to leave out")
	sends := flag.Bool("sends", false, "record a send in each state")
	flag.Parse()
	if flag.NArg() != 3 {
		fmt.Fprintln(os.Stderr, "usage: seed [-mutate id] [-drop id] [-sends] <CULL_HOME> <review.json> <project root>")
		os.Exit(2)
	}
	home, fixturePath, root := flag.Arg(0), flag.Arg(1), flag.Arg(2)
	if err := os.Setenv("CULL_HOME", home); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	pid, run, err := Seed(context.Background(), store.Path(), data, root, *mutate, *drop)
	if err == nil && *sends {
		err = SeedSends(context.Background(), store.Path(), data, pid, run)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "seed:", err)
		os.Exit(1)
	}
	fmt.Printf("{\"project\":%d,\"run\":%d}\n", pid, run)
}
