package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/subaud-tools/tools-common/localweb"
	"github.com/subaud-tools/tools-common/pageserve"
)

// cull report, which the Claude Code plugin's mod runs, posts the report to
// the running cull serve for the session it runs in.
func TestReportPostsToServe(t *testing.T) {
	t.Setenv("CULL_HOME", t.TempDir())
	t.Setenv("CLAUDE_CODE_SESSION_ID", "s1")
	t.Setenv("AGENT_SESSION_ID", "")
	t.Setenv("AGENT_SESSION_CHILD", "")
	got := make(chan map[string]string, 3)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/agent/report" || r.Header.Get(localweb.TokenHeader) != "tok" {
			http.NotFound(w, r)
			return
		}
		var b map[string]string
		_ = json.NewDecoder(r.Body).Decode(&b)
		got <- b
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)
	writeEntry(t, pageserve.Entry{Base: srv.URL, Token: "tok"})
	for _, c := range []struct {
		args  []string
		event string
		token string
	}{
		{[]string{"held", "0123456789abcdef"}, "held", "0123456789abcdef"},
		{[]string{"landed", "0123456789abcdef"}, "landed", "0123456789abcdef"},
		{[]string{"turn-ended"}, "turn_ended", ""},
	} {
		if code, out, errw := run(t, "", append([]string{"report"}, c.args...)...); code != 0 {
			t.Fatalf("report %v: code %d: %s %s", c.args, code, out, errw)
		}
		select {
		case b := <-got:
			if b["session"] != "s1" || b["event"] != c.event || b["token"] != c.token {
				t.Fatalf("report %v posted %v", c.args, b)
			}
		default:
			t.Fatalf("report %v reached no serve", c.args)
		}
	}
}
