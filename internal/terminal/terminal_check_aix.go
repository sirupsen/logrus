// +build aix

package terminal

import (
	"golang.org/x/sys/unix"
)

const ioctlReadTermios = unix.TCGETS
