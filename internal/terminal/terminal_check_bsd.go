// +build darwin dragonfly freebsd netbsd openbsd

package terminal

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TIOCGETA
