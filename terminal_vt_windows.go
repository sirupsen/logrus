//go:build windows

package logrus

import "golang.org/x/sys/windows"

func enableVirtualTerminalProcessing(fd uintptr) bool {
	h := windows.Handle(fd)

	var mode uint32
	if err := windows.GetConsoleMode(h, &mode); err != nil {
		return false
	}
	if mode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != 0 {
		return true
	}
	if err := windows.SetConsoleMode(h, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err != nil {
		return false
	}
	return true
}
