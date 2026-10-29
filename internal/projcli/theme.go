package projcli

import (
	"flag"
	"io"
	"slices"

	"github.com/subaud-tools/tackle/internal/proj"
	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/sessiontheme"
)

type themeFlags struct {
	session, window, hue *string
	json                 *bool
}

func newThemeFlags() (*flag.FlagSet, *themeFlags) {
	fs := flag.NewFlagSet("theme", flag.ContinueOnError)
	f := &themeFlags{session: fs.String("session", "", "agent session ID to query"), window: fs.String("window", "", "native Ghostty window ID"), hue: fs.String("hue", "", "window hue family"), json: fs.Bool("json", false, "emit JSON")}
	tools.SetUsage(fs, "proj theme --session <agent-id> --json | --window <native-window-id> --hue <hue> --json", "Queries a remembered session theme or sets a known window's hue.")
	return fs, f
}

func cmdTheme(args []string, out, errw io.Writer) error {
	fs, f := newThemeFlags()
	if err := tools.ParseFlags(fs, args, out); err != nil {
		return err
	}
	query := *f.session != "" && *f.window == "" && *f.hue == ""
	set := *f.session == "" && *f.window != "" && slices.Contains(sessiontheme.Hue("").Values(), *f.hue)
	if !*f.json || fs.NArg() != 0 || (!query && !set) {
		return tools.UsageError{Msg: "use --session <agent-id> --json or --window <native-window-id> --hue <hue> --json"}
	}
	if query {
		q, err := proj.ThemeForAgent(*f.session)
		if err != nil {
			return err
		}
		return writeJSON(out, q)
	}
	if err := proj.SetWindowHue(*f.window, sessiontheme.Hue(*f.hue)); err != nil {
		return err
	}
	return writeJSON(out, struct {
		Window string           `json:"window"`
		Hue    sessiontheme.Hue `json:"hue"`
	}{*f.window, sessiontheme.Hue(*f.hue)})
}
