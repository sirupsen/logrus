// +build !appengine,!js,!windows,!aix

package logrus

import (
	"io"
	"os"

	"golang.org/x/sys/unix"
)

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		// Returns whether the given file descriptor is a terminal.
		// Taken from golang.org/x/crypto/ssh/terminal.IsTerminal without importing the whole package
		_, err := unix.IoctlGetTermios(int(v.Fd()), unix.TIOCGETA)
		return err == nil
	default:
		return false
	}
}
