package projcli

import (
	"flag"
	"fmt"
	"io"
	"slices"

	"github.com/subaud-tools/tackle/internal/proj"
	tools "github.com/subaud-tools/tools-common"
	"github.com/subaud-tools/tools-common/sessiontheme"
)

type themeFlags struct {
	session, window, hue *string
	json, enabled        *bool
}

func newThemeFlags() (*flag.FlagSet, *themeFlags) {
	fs := flag.NewFlagSet("theme", flag.ContinueOnError)
	f := &themeFlags{session: fs.String("session", "", "agent session ID to query"), window: fs.String("window", "", "native Ghostty window ID"), hue: fs.String("hue", "", "window hue family"), json: fs.Bool("json", false, "emit JSON"), enabled: fs.Bool("enabled", false, "query whether session themes are enabled")}
	tools.SetUsage(fs, "proj theme --enabled --json | --session <agent-id> --json | --window <native-window-id> --hue <hue> --json", "Queries a remembered session theme or sets a known window's hue.")
	return fs, f
}

func cmdTheme(args []string, out, errw io.Writer) error {
	fs, f := newThemeFlags()
	if err := tools.ParseFlags(fs, args, out); err != nil {
		return err
	}
	enabled := *f.enabled && *f.session == "" && *f.window == "" && *f.hue == ""
	query := !*f.enabled && *f.session != "" && *f.window == "" && *f.hue == ""
	set := !*f.enabled && *f.session == "" && *f.window != "" && slices.Contains(sessiontheme.Hue("").Values(), *f.hue)
	if !*f.json || fs.NArg() != 0 || (!enabled && !query && !set) {
		return tools.UsageError{Msg: "use --enabled --json, --session <agent-id> --json or --window <native-window-id> --hue <hue> --json"}
	}
	if enabled {
		return writeJSON(out, struct {
			Enabled bool `json:"enabled"`
		}{proj.SessionThemesEnabled()})
	}
	if query {
		q, err := proj.ThemeForAgent(*f.session)
		if err != nil {
			return err
		}
		return writeJSON(out, q)
	}
	if !proj.SessionThemesEnabled() {
		return fmt.Errorf("proj: session themes are disabled")
	}
	if err := reconcileThemes(true); err != nil {
		return err
	}
	if err := proj.SetWindowHue(*f.window, sessiontheme.Hue(*f.hue)); err != nil {
		return err
	}
	if err := reconcileThemes(true); err != nil {
		return err
	}
	return writeJSON(out, struct {
		Window string           `json:"window"`
		Hue    sessiontheme.Hue `json:"hue"`
	}{*f.window, sessiontheme.Hue(*f.hue)})
}
