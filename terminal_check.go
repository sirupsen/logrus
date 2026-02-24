//go:build !windows && !appengine

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
		if fd > maxInt {
			return false
		}
		return term.IsTerminal(int(fd))
	}
	return false
}
