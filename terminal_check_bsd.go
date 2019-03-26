// +build darwin dragonfly freebsd netbsd openbsd

package logrus

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TIOCGETA
