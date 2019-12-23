// +build windows

package terminal

import (
	"golang.org/x/sys/windows"
)

func IsTerminal(fd int) bool {
	var st uint32
	err := windows.GetConsoleMode(windows.Handle(fd), &st)
	return err == nil
}
