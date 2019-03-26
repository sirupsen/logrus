// +build !appengine,!js,!windows

package terminal

import (
	"io"
	"os"

	"golang.org/x/sys/unix"
)

func IsTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		_, err := unix.IoctlGetTermios(int(v.Fd()), ioctlReadTermios)

		return err == nil
	default:
		return false
	}
}
