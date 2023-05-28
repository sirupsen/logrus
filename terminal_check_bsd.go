//go:build (darwin || dragonfly || freebsd || netbsd || openbsd) && (!js || !tinygo)
// +build darwin dragonfly freebsd netbsd openbsd
// +build !js !tinygo

package logrus

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TIOCGETA

func isTerminal(fd int) bool {
	_, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	return err == nil
}
