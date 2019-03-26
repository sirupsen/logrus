// +build !appengine,!js,windows

package terminal

import (
	"io"
	"os"
	"syscall"

	sequences "github.com/konsorten/go-windows-terminal-sequences"
)

func InitTerminal(w io.Writer) {
	switch v := w.(type) {
	case *os.File:
		sequences.EnableVirtualTerminalProcessing(syscall.Handle(v.Fd()), true)
	}
}
