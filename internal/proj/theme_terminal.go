package proj

import (
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/subaud-tools/tools-common/sessiontheme"
)

// TerminalSequences changes only terminal colors; it contains no shell input,
// title mutation, global config reload, or cursor-motion commands. OSC 12 sets
// the cursor fill. Ghostty has no separate OSC for the text beneath the cursor.
func TerminalSequences(p sessiontheme.Palette) []byte {
	var out strings.Builder
	fmt.Fprintf(&out, "\x1b]10;%s\x1b\\\x1b]11;%s\x1b\\\x1b]12;%s\x1b\\", p.Foreground, p.Background, p.Cursor)
	for i, color := range p.ANSI {
		fmt.Fprintf(&out, "\x1b]4;%d;%s\x1b\\", i, color)
	}
	return []byte(out.String())
}

var writeTerminal = writeTerminalDevice

// Open only the observed macOS PTY device, never a symlink or a regular file.
// O_NONBLOCK prevents a disappeared/stalled client from blocking the sync lock.
func writeTerminalDevice(tty string, b []byte) error {
	suffix, ok := strings.CutPrefix(tty, "/dev/ttys")
	if !ok || suffix == "" || strings.Trim(suffix, "0123456789") != "" {
		return fmt.Errorf("proj: unsafe terminal path %q", tty)
	}
	f, err := os.OpenFile(tty, os.O_WRONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	if info.Mode()&os.ModeCharDevice == 0 || !ok || int64(st.Uid) != int64(os.Getuid()) {
		return fmt.Errorf("proj: not an owned terminal device: %s", tty)
	}
	n, err := f.Write(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return io.ErrShortWrite
	}
	return nil
}
