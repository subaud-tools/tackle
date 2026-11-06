package proj

import (
	"fmt"

	"github.com/subaud-tools/tackle/internal/proj/pages"
)

// SplitResult is what `proj split --json` prints: Split is the session's split
// page address.
type SplitResult struct {
	Split string `json:"split"`
}

// Split opens right as a tab of session's split page on proj's page server,
// starting it when it is not running, with no Chrome, and returns the split
// page's address: proj's public address with PAGE_PUBLIC_DOMAIN set, else
// the loopback one. It refuses a session not on socket and a right the split
// page's frame may not show.
func Split(socket, session, right string) (string, error) {
	if !localPage(right) {
		return "", notOurPage()
	}
	if _, err := Run(socket, "has-session", "-t", "="+session); err != nil {
		return "", fmt.Errorf("no session %q on %s", session, socket)
	}
	addr, token, inst, err := ensurePages()
	if err != nil {
		return "", stepError("start the page server", err)
	}
	if err := openTab(session, right); err != nil {
		return "", err
	}
	return pages.ViewURL(addr, token, inst, socket, session), nil
}
