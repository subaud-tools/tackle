package serve

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/subaud-tools/tackle/internal/cull/store"
)

type fixture struct {
	t    *testing.T
	st   *store.Store
	srv  *Server
	h    http.Handler
	proj store.Project
	run  int64
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	t.Setenv("CULL_HOME", t.TempDir())
	st, err := store.Open(context.Background(), filepath.Join(t.TempDir(), "cull.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	f := &fixture{t: t, st: st}
	f.srv = New(st)
	// cull's routes with the event routes pageserve.Run adds in front of them.
	mux := http.NewServeMux()
	mux.Handle("/", f.srv.Handler())
	mux.HandleFunc("GET /api/poll", f.srv.events.Poll)
	mux.HandleFunc("GET /api/events", f.srv.events.Stream)
	f.h = mux
	f.proj, err = st.Project(context.Background(), "/repo")
	if err != nil {
		t.Fatal(err)
	}
	f.run = f.record()
	return f
}

func jev(p float64, key string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{"verdict":{"probabilities":{%q:%v}}}`, key, p))
}

func (f *fixture) record() int64 {
	f.t.Helper()
	return f.recordItems(defaultItems())
}

func defaultItems() []store.Item {
	return []store.Item{
		{ID: "g1", Kind: "group", Hash: "gh1", Verdict: "keep_separate", Rule: "review_band", Jev: jev(0.4, "consolidate"), Members: []string{"a", "b"}, Rows: [][]string{{"x", "y"}}},
		{ID: "t1", Kind: "test", Hash: "h1", File: "a_test.go", Name: "A", Verdict: "keep", Rule: "review_band", Jev: jev(0.3, "cut"), State: json.RawMessage(`{"s":1}`), Model: "m"},
		{ID: "t2", Kind: "test", Hash: "h2", File: "b_test.go", Name: "B", Verdict: "keep", Rule: "truncated", Jev: jev(0.6, "cut")},
		{ID: "t0", Kind: "test", Hash: "h0", File: "c_test.go", Name: "C", Verdict: "keep", Rule: "review_band"},
	}
}

func (f *fixture) recordItems(items []store.Item) int64 {
	f.t.Helper()
	id, err := f.st.RecordRun(context.Background(), store.Run{ProjectID: f.proj.ID, Mode: "all", Total: 9, Summary: map[string]int{"cut": 2}, QuestionsHash: "qh"}, items)
	if err != nil {
		f.t.Fatal(err)
	}
	return id
}

func (f *fixture) do(method, path, body string) *httptest.ResponseRecorder {
	f.t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	f.h.ServeHTTP(w, req)
	return w
}

func (f *fixture) review() map[string]any {
	f.t.Helper()
	w := f.do("GET", fmt.Sprintf("/api/review?project=%d", f.proj.ID), "")
	if w.Code != 200 {
		f.t.Fatalf("review %d %s", w.Code, w.Body)
	}
	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		f.t.Fatal(err)
	}
	return m
}

func (f *fixture) put(run int64, hash string) *httptest.ResponseRecorder {
	body := fmt.Sprintf(`{"project":%d,"run":%d,"answers":[{"id":"t1","hash":%q,"kind":"test","value":"cut","note":"n","via":"item","blind":true}]}`, f.proj.ID, run, hash)
	return f.do("PUT", "/api/answers", body)
}

type wireEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type pollResp struct {
	Cursor string      `json:"cursor"`
	Reset  bool        `json:"reset"`
	Events []wireEvent `json:"events"`
}

func (f *fixture) poll(since string) pollResp {
	f.t.Helper()
	w := f.do("GET", "/api/poll?since="+since, "")
	var r pollResp
	if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil {
		f.t.Fatal(err)
	}
	return r
}

func (f *fixture) events(since string) (string, []wireEvent) {
	f.t.Helper()
	r := f.poll(since)
	return r.Cursor, r.Events
}

func TestReviewShapeAndOrder(t *testing.T) {
	f := newFixture(t)
	m := f.review()
	run := m["run"].(map[string]any)
	if run["id"].(float64) != float64(f.run) || run["total"].(float64) != 9 || run["mode"] != "all" {
		t.Fatalf("run %v", run)
	}
	if m["project"].(map[string]any)["root"] != "/repo" || m["blind"] != false {
		t.Fatalf("project %v", m)
	}
	var ids []string
	for _, it := range m["items"].([]any) {
		ids = append(ids, it.(map[string]any)["id"].(string))
	}
	if strings.Join(ids, ",") != "t2,t1,t0,g1" {
		t.Fatalf("order %v", ids)
	}
	g := m["items"].([]any)[3].(map[string]any)
	if len(g["members"].([]any)) != 2 || g["rows"] == nil {
		t.Fatalf("group %v", g)
	}
	t1 := m["items"].([]any)[1].(map[string]any)
	if t1["model"] != "m" || t1["state"].(map[string]any)["s"].(float64) != 1 || t1["answer"] != nil {
		t.Fatalf("t1 %v", t1)
	}
}

func TestReviewUnknownAndEmptyProject(t *testing.T) {
	f := newFixture(t)
	if w := f.do("GET", "/api/review?project=999", ""); w.Code != 404 || !strings.Contains(w.Body.String(), `"no project"`) {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	p, _ := f.st.Project(context.Background(), "/empty")
	w := f.do("GET", fmt.Sprintf("/api/review?project=%d", p.ID), "")
	var m map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &m)
	if w.Code != 200 || m["run"] != nil || len(m["items"].([]any)) != 0 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
}

func TestPutStoresAndEmits(t *testing.T) {
	f := newFixture(t)
	if w := f.put(f.run, "h1"); w.Code != 204 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	it := f.review()["items"].([]any)[1].(map[string]any)
	a := it["answer"].(map[string]any)
	if a["value"] != "cut" || a["note"] != "n" || a["via"] != "item" || a["blind"] != true || a["answered_at"] == nil || a["sent_at"] != nil {
		t.Fatalf("answer %v", a)
	}
	got, _ := f.st.Answers(context.Background(), f.proj.ID)
	if got[store.Key{ItemID: "t1", Hash: "h1"}].Model != "m" || got[store.Key{ItemID: "t1", Hash: "h1"}].QuestionsHash != "qh" {
		t.Fatalf("stored %+v", got)
	}
	_, evs := f.events("")
	if len(evs) != 1 || evs[0].Type != "answers" || !strings.Contains(string(evs[0].Data), fmt.Sprintf(`"project":%d`, f.proj.ID)) {
		t.Fatalf("events %v", evs)
	}
}

func TestPutStaleAndInvalid(t *testing.T) {
	f := newFixture(t)
	for name, w := range map[string]*httptest.ResponseRecorder{
		"hash": f.put(f.run, "nope"),
	} {
		if w.Code != 409 || !strings.Contains(w.Body.String(), `"stale"`) {
			t.Fatalf("%s: %d %s", name, w.Code, w.Body)
		}
	}
	// one stale answer sinks the batch
	body := fmt.Sprintf(`{"project":%d,"run":%d,"answers":[{"id":"t1","hash":"h1","kind":"test","value":"cut","via":"item"},{"id":"t2","hash":"bad","kind":"test","value":"keep","via":"item"}]}`, f.proj.ID, f.run)
	if w := f.do("PUT", "/api/answers", body); w.Code != 409 {
		t.Fatalf("%d", w.Code)
	}
	if got, _ := f.st.Answers(context.Background(), f.proj.ID); len(got) != 0 {
		t.Fatalf("stored %v", got)
	}
	if _, evs := f.events(""); len(evs) != 0 {
		t.Fatalf("events %v", evs)
	}
	for _, b := range []string{
		`{"project":1,"run":1,"answers":[{"id":"t1","hash":"h1","kind":"test","value":"merge","via":"item"}]}`,
		`{"project":1,"run":1,"answers":[{"id":"t1","hash":"h1","kind":"test","value":"cut","via":"x"}]}`,
		`{"project":1,"run":1,"extra":1,"answers":[]}`,
		`{"project":1,"run":1,"answers":[]}`,
		`{"project":1,"run":1}`,
		fmt.Sprintf(`{"project":%d,"answers":[{"id":"t1","hash":"h1","kind":"test","value":"cut","via":"item"}]}`, f.proj.ID),
		`{"project":0,"run":0,"answers":[]}`,
		`not json`,
	} {
		if w := f.do("PUT", "/api/answers", b); w.Code != 400 {
			t.Fatalf("%s: %d %s", b, w.Code, w.Body)
		}
	}
}

func TestDeleteAndSend(t *testing.T) {
	f := newFixture(t)
	f.put(f.run, "h1")
	w := f.do("POST", "/api/send", fmt.Sprintf(`{"project":%d}`, f.proj.ID))
	if w.Code != 200 || strings.TrimSpace(w.Body.String()) != `{"sent":1,"to":""}` {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	if w := f.do("POST", "/api/send", fmt.Sprintf(`{"project":%d}`, f.proj.ID)); strings.TrimSpace(w.Body.String()) != `{"sent":0,"to":""}` {
		t.Fatalf("%s", w.Body)
	}
	if _, evs := f.events(""); len(evs) != 2 { // put + first send only
		t.Fatalf("events %v", evs)
	}
	a := f.review()["items"].([]any)[1].(map[string]any)["answer"].(map[string]any)
	if a["sent_at"] == nil {
		t.Fatalf("not sent: %v", a)
	}
	if w := f.do("DELETE", fmt.Sprintf("/api/answers?project=%d&id=t1&hash=h1", f.proj.ID), ""); w.Code != 204 {
		t.Fatalf("%d", w.Code)
	}
	if f.review()["items"].([]any)[1].(map[string]any)["answer"] != nil {
		t.Fatal("answer remains")
	}
	if w := f.do("DELETE", fmt.Sprintf("/api/answers?project=%d&id=t1&hash=h1", f.proj.ID), ""); w.Code != 204 {
		t.Fatalf("idempotent %d", w.Code)
	}
}

func TestRunEventWithinThreeSeconds(t *testing.T) {
	f := newFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go f.srv.Watch(ctx, 50*time.Millisecond)
	time.Sleep(150 * time.Millisecond) // first read recorded
	cur, evs := f.events("")
	if len(evs) != 0 {
		t.Fatalf("events at start %v", evs)
	}
	id := f.record()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if _, evs = f.events(cur); len(evs) > 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	want := fmt.Sprintf(`{"project":%d,"run":%d}`, f.proj.ID, id)
	if len(evs) != 1 || evs[0].Type != "run" || string(evs[0].Data) != want {
		t.Fatalf("events %v want %s", evs, want)
	}
}

// sse opens /api/events and returns a channel of "id|data" strings.
func sse(t *testing.T, url, lastID string) (<-chan string, func()) {
	t.Helper()
	return sseReq(t, url+"/api/events", lastID)
}

func sseReq(t *testing.T, u, lastID string) (<-chan string, func()) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	if lastID != "" {
		req.Header.Set("Last-Event-ID", lastID)
	}
	resp, err := http.DefaultClient.Do(req) //nolint:bodyclose // the reader goroutine below closes it
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	out := make(chan string, 100)
	go func() {
		defer close(out)
		defer func() { _ = resp.Body.Close() }()
		sc := bufio.NewScanner(resp.Body)
		var id string
		for sc.Scan() {
			l := sc.Text()
			if v, ok := strings.CutPrefix(l, "id: "); ok {
				id = v
			}
			if v, ok := strings.CutPrefix(l, "data: "); ok {
				out <- id + "|" + v
			}
		}
	}()
	return out, cancel
}

func next(t *testing.T, c <-chan string) string {
	t.Helper()
	select {
	case v, ok := <-c:
		if !ok {
			t.Fatal("stream closed")
		}
		return v
	case <-time.After(3 * time.Second):
		t.Fatal("no event")
	}
	return ""
}

func TestPutFromOldTabValidatesAgainstLatestRun(t *testing.T) {
	f := newFixture(t)
	old := f.run
	items := defaultItems()
	for i := range items {
		if items[i].ID == "t2" {
			items[i].Hash = "h2-new"
		}
	}
	latest := f.recordItems(items)
	if latest == old {
		t.Fatal("no new run")
	}
	// t1 unchanged in the latest run: accepted, stored against the latest run.
	if w := f.put(old, "h1"); w.Code != 204 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	got, _ := f.st.Answers(context.Background(), f.proj.ID)
	if a, ok := got[store.Key{ItemID: "t1", Hash: "h1"}]; !ok || a.RunID != latest {
		t.Fatalf("stored %+v want run %d", got, latest)
	}
	// t2's hash changed: stale, nothing stored.
	body := fmt.Sprintf(`{"project":%d,"run":%d,"answers":[{"id":"t2","hash":"h2","kind":"test","value":"keep","via":"item"}]}`, f.proj.ID, old)
	if w := f.do("PUT", "/api/answers", body); w.Code != 409 || !strings.Contains(w.Body.String(), `"stale"`) {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	if got, _ = f.st.Answers(context.Background(), f.proj.ID); len(got) != 1 {
		t.Fatalf("stored %v", got)
	}
}

func TestPutProjectWithoutRunIsStale(t *testing.T) {
	f := newFixture(t)
	p, _ := f.st.Project(context.Background(), "/empty")
	body := fmt.Sprintf(`{"project":%d,"run":1,"answers":[{"id":"t1","hash":"h1","kind":"test","value":"cut","via":"item"}]}`, p.ID)
	if w := f.do("PUT", "/api/answers", body); w.Code != 409 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
}

func itemIDs(m map[string]any) string {
	var ids []string
	for _, it := range m["items"].([]any) {
		ids = append(ids, it.(map[string]any)["id"].(string))
	}
	return strings.Join(ids, ",")
}

func TestReviewOrderTiesAndGroups(t *testing.T) {
	f := newFixture(t)
	f.recordItems([]store.Item{
		{ID: "tb", Kind: "test", Hash: "1", Jev: jev(0.5, "cut")},
		{ID: "ta", Kind: "test", Hash: "2", Jev: jev(0.5, "cut")},
		{ID: "g1", Kind: "group", Hash: "3", Jev: jev(0.4, "consolidate")},
		{ID: "g2", Kind: "group", Hash: "4", Jev: jev(0.7, "consolidate")},
		{ID: "tc", Kind: "test", Hash: "5", Jev: jev(0.9, "cut")},
	})
	if got := itemIDs(f.review()); got != "tc,ta,tb,g2,g1" {
		t.Fatalf("order %s", got)
	}
}

func TestSameIDDifferentHashGetsNoAnswer(t *testing.T) {
	f := newFixture(t)
	if w := f.put(f.run, "h1"); w.Code != 204 {
		t.Fatal(w.Code)
	}
	items := defaultItems()
	for i := range items {
		if items[i].ID == "t1" {
			items[i].Hash = "h1-changed"
		}
	}
	f.recordItems(items)
	m := f.review()
	for _, it := range m["items"].([]any) {
		if it.(map[string]any)["id"] == "t1" && it.(map[string]any)["answer"] != nil {
			t.Fatalf("answer attached to changed hash: %v", it)
		}
	}
	// the answer is still counted in the project's totals
	if ans := m["answered"].(map[string]any); ans["total"] != float64(1) || ans["sent"] != float64(0) {
		t.Fatalf("answered %v", ans)
	}
}

func TestAnsweredShape(t *testing.T) {
	f := newFixture(t)
	f.put(f.run, "h1")
	body := fmt.Sprintf(`{"project":%d,"run":%d,"answers":[{"id":"t2","hash":"h2","kind":"test","value":"keep","via":"item"}]}`, f.proj.ID, f.run)
	f.do("PUT", "/api/answers", body)
	f.do("POST", "/api/send", fmt.Sprintf(`{"project":%d}`, f.proj.ID))
	f.put(f.run, "h1") // re-answer: unsent again
	b, _ := json.Marshal(f.review()["answered"])
	if string(b) != `{"sent":1,"total":2}` {
		t.Fatalf("answered %s", b)
	}
}

func TestWatchQuietWhenUnchanged(t *testing.T) {
	f := newFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { f.srv.Watch(ctx, 10*time.Millisecond); close(done) }()
	time.Sleep(300 * time.Millisecond) // many ticks
	cancel()
	<-done
	if _, evs := f.events(""); len(evs) != 0 {
		t.Fatalf("events %v", evs)
	}
}

func TestReviewCursorCoversRaceWindow(t *testing.T) {
	f := newFixture(t)
	cur, _ := f.review()["cursor"].(string)
	if cur == "" {
		t.Fatal("review has no cursor")
	}
	f.put(f.run, "h1") // emitted after the cursor was read
	_, evs := f.events(cur)
	if len(evs) != 1 || evs[0].Type != "answers" {
		t.Fatalf("poll since review cursor: %+v", evs)
	}
	ts := httptest.NewServer(f.h)
	defer ts.Close()
	c, cancel := sse(t, ts.URL, cur)
	defer cancel()
	if v := next(t, c); !strings.Contains(v, `"answers"`) || strings.Contains(v, "reset") {
		t.Fatalf("sse from review cursor: %s", v)
	}
}

func TestReviewCursorOnNoRunPath(t *testing.T) {
	f := newFixture(t)
	p, _ := f.st.Project(context.Background(), "/empty")
	w := f.do("GET", fmt.Sprintf("/api/review?project=%d", p.ID), "")
	var m map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &m)
	if c, _ := m["cursor"].(string); c == "" || m["run"] != nil {
		t.Fatalf("%s", w.Body)
	}
}

func TestPutGroupDoesNotReplaceItemAnswer(t *testing.T) {
	f := newFixture(t)
	if w := f.put(f.run, "h1"); w.Code != 204 {
		t.Fatal(w.Code)
	}
	body := fmt.Sprintf(`{"project":%d,"run":%d,"answers":[{"id":"t1","hash":"h1","kind":"test","value":"keep","via":"group"}]}`, f.proj.ID, f.run)
	if w := f.do("PUT", "/api/answers", body); w.Code != 204 {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	got, _ := f.st.Answers(context.Background(), f.proj.ID)
	if a := got[store.Key{ItemID: "t1", Hash: "h1"}]; a.Value != "cut" || a.Via != "item" {
		t.Fatalf("%+v", a)
	}
}
