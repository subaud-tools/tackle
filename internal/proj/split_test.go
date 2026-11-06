package proj

import (
	"context"
	"strings"
	"testing"

	"github.com/subaud-tools/tackle/internal/proj/pages"
	"github.com/subaud-tools/tools-common/localweb"
	"github.com/subaud-tools/tools-common/pageserve"
)

// On a public domain proj serve's entry advertises its public address; the
// calls proj makes to it go to its loopback address.
func TestEnsurePagesUsesTheLoopbackAddressOnAPublicDomain(t *testing.T) {
	usePagesClient(t, &pageserve.Client{
		Find: func() (pageserve.Entry, error) {
			return pageserve.Entry{Base: "https://proj.ws.example.com", Local: "http://127.0.0.1:1", Token: "tok"}, nil
		},
		Stop:  func(context.Context) error { return nil },
		Start: func() (pageserve.Entry, error) { return pageserve.Entry{}, nil },
	})
	if addr, _, _, err := ensurePages(); err != nil || addr != "127.0.0.1:1" {
		t.Fatalf("ensurePages = %q, %v; want the loopback address", addr, err)
	}
}

func TestLocalPageAcceptsAPublicPageOnlyWithTheSetting(t *testing.T) {
	const public = "https://diazo.ws.example.com/plan?x=1"
	if localPage(public) {
		t.Fatal("a public page was accepted with the setting unset")
	}
	t.Setenv(localweb.PublicDomainEnv, "ws.example.com")
	if !localPage(public) || !localPage("http://127.0.0.1:4321/plan") {
		t.Fatal("refused a page the split shows")
	}
	for _, bad := range []string{"https://diazo.other.example.com/", "http://diazo.ws.example.com/", "https://a.b.ws.example.com/", "https://evil.example/"} {
		if localPage(bad) {
			t.Errorf("accepted %q", bad)
		}
	}
}

// proj split opens the tab in the session's split page through proj serve,
// with no Chrome, and answers the page's address.
func TestSplitOpensTheTabAndAnswersThePage(t *testing.T) {
	socket, name, _ := terminalServer(t)
	calls := fakeChrome(t, "", nil)
	opened := fakePageServer(t)
	addr, token, inst, _ := ensurePages()
	got, err := Split(socket, name, "http://127.0.0.1:4321/plan?t=x")
	if err != nil {
		t.Fatal(err)
	}
	if len(*calls) != 0 {
		t.Fatalf("ran Chrome: %+v", *calls)
	}
	if len(*opened) != 1 || (*opened)[0].Session != name || (*opened)[0].Address != "http://127.0.0.1:4321/plan?t=x" {
		t.Fatalf("opened tabs %+v", *opened)
	}
	if want := pages.ViewURL(addr, token, inst, socket, name); got != want {
		t.Fatalf("Split = %q, want %q", got, want)
	}
	if !strings.HasPrefix(got, "http://"+addr+"/term?") {
		t.Fatalf("Split = %q, want the loopback page", got)
	}
}

func TestSplitOnAPublicDomainAnswersThePublicPage(t *testing.T) {
	socket, name, _ := terminalServer(t)
	fakeChrome(t, "", nil)
	opened := fakePageServer(t)
	t.Setenv(localweb.PublicDomainEnv, "ws.example.com")
	got, err := Split(socket, name, "https://diazo.ws.example.com/plan")
	if err != nil {
		t.Fatal(err)
	}
	if len(*opened) != 1 || (*opened)[0].Address != "https://diazo.ws.example.com/plan" {
		t.Fatalf("opened tabs %+v", *opened)
	}
	if !strings.HasPrefix(got, "https://proj.ws.example.com/term?") || !pages.TermMatches(got, socket, name) {
		t.Fatalf("Split = %q, want proj's public page", got)
	}
}

func TestSplitRefusesAMissingSessionAndAForeignPage(t *testing.T) {
	socket, name, _ := terminalServer(t)
	fakeChrome(t, "", nil)
	opened := fakePageServer(t)
	if _, err := Split(socket, "no/such", "http://127.0.0.1:4321/plan"); err == nil || !strings.Contains(err.Error(), "no/such") {
		t.Fatalf("missing session: %v", err)
	}
	if _, err := Split(socket, name, "https://evil.example/"); err == nil {
		t.Fatal("opened a foreign page")
	}
	if len(*opened) != 0 {
		t.Fatalf("opened tabs %+v", *opened)
	}
}
