// +build !appengine,!gopherjs,windows

package logrus

import (
	"os"
	"syscall"

	sequences "github.com/konsorten/go-windows-terminal-sequences"
)

func (f *TextFormatter) initTerminal(entry *Entry) {
	switch v := entry.Logger.Out.(type) {
	case *os.File:
		handle := syscall.Handle(v.Fd())

		sequences.EnableVirtualTerminalProcessing(handle, true)
	}
}
