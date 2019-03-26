// +build !appengine,!js,!windows

package logrus

import (
	"io"
	"os"

	"golang.org/x/sys/unix"
)

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		_, err := unix.IoctlGetTermios(int(v.Fd()), ioctlReadTermios)

		return err == nil
	default:
		return false
	}
}
