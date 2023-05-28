//go:build js || nacl || plan9 || tinygo
// +build js nacl plan9 tinygo

package logrus

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return false
}
