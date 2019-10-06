// +build js nacl plan9 gopherjs

package logrus

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return false
}
