// +build !appengine,!js,!windows

package logrus

import (
	"io"
	"os"

	"github.com/sirupsen/logrus/terminal"
)

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
