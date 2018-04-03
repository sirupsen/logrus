// +build darwin freebsd openbsd netbsd dragonfly
// +build !appengine,!gopherjs

package logrus

import (
	"io"

	"golang.org/x/sys/unix"
)

const ioctlReadTermios = unix.TIOCGETA

type Termios unix.Termios

func initTerminal(w io.Writer) {
}
