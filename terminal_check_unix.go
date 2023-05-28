//go:build (linux || aix || zos) && (!js || !tinygo)
// +build linux aix zos
// +build !js !tinygo

package logrus

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TCGETS

func isTerminal(fd int) bool {
	_, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	return err == nil
}
