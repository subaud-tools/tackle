package proj

import "syscall"

// RaiseFileLimit sets the process's soft RLIMIT_NOFILE to want, or to the hard
// limit if that is lower, keeping the hard limit. It always calls Setrlimit,
// even when the soft limit already reads higher: Go raises its own soft limit
// before main and hands children (the tmux servers proj starts) the original
// one back unless the program sets the limit itself.
func RaiseFileLimit(want uint64) error {
	var lim syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &lim); err != nil {
		return err
	}
	lim.Cur = min(want, lim.Max)
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &lim)
}
