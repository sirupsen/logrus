// Based on ssh/terminal/util_solaris.go:
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build solaris

package logrus

import (
	"golang.org/x/sys/unix"
	"syscall"
)

// IsTerminal returns true if the given file descriptor is a terminal.
// see: http://src.illumos.org/source/xref/illumos-gate/usr/src/lib/libbc/libc/gen/common/isatty.c
func IsTerminal() bool {
	fd := syscall.Stderr
	var termio unix.Termio
	err := unix.IoctlSetTermio(fd, unix.TCGETA, &termio)
	return err == nil
}
