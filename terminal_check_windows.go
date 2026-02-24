//go:build windows && !appengine

package logrus

import (
	"io"
	"os"

	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

func checkIfTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		fd := f.Fd()
		maxInt := uintptr(^uint(0) >> 1)
		if fd > maxInt {
			return false
		}
		if !term.IsTerminal(int(f.Fd())) {
			return false
		}

		h := windows.Handle(f.Fd())
		var mode uint32
		if err := windows.GetConsoleMode(h, &mode); err != nil {
			return false
		}
		if mode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != 0 {
			return true
		}
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		if err := windows.SetConsoleMode(h, mode); err != nil {
			return false
		}
	}
	return false
}
