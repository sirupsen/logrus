// +build darwin freebsd openbsd netbsd dragonfly
// +build !appengine

package logrus

import "github.com/golang/sys/unix"

const ioctlReadTermios = unix.TIOCGETA

type Termios unix.Termios
