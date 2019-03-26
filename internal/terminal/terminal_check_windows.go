// +build !appengine,!js,windows

package terminal

import (
	"io"
	"os"
	"syscall"
)

func IsTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		var mode uint32
		err := syscall.GetConsoleMode(syscall.Handle(v.Fd()), &mode)
		return err == nil
	default:
		return false
	}
}
