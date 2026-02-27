//go:build !appengine

package logrus

import (
	"io"
	"os"

	"golang.org/x/term"
)

func checkIfTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		fd := f.Fd()
		maxInt := uintptr(^uint(0) >> 1)
		if fd > maxInt || !term.IsTerminal(int(fd)) {
			return false
		}

		// On Windows consoles, ANSI escape sequences are not processed
		// unless ENABLE_VIRTUAL_TERMINAL_PROCESSING is set.
		return enableVirtualTerminalProcessing(fd)
	}
	return false
}