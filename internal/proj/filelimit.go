package proj

import "syscall"

// RaiseFileLimit raises the process's soft RLIMIT_NOFILE to want, or to the
// hard limit if that is lower, keeping the hard limit. A soft limit already at
// least want is left alone. Setting it through syscall.Setrlimit is what makes
// child processes (the tmux servers proj starts) inherit it: Go raises its own
// limit at startup but hands children the original one otherwise.
func RaiseFileLimit(want uint64) error {
	var lim syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &lim); err != nil {
		return err
	}
	if lim.Cur >= want {
		return nil
	}
	lim.Cur = min(want, lim.Max)
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &lim)
}
