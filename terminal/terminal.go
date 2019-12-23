// +build aix darwin dragonfly freebsd linux,!appengine netbsd openbsd

package terminal

import "golang.org/x/sys/unix"

// IsTerminal returns whether the given file descriptor is a terminal.
func IsTerminal(fd int) bool {
	_, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	return err == nil
}
