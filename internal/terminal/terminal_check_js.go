// +build js

package terminal

import (
	"io"
)

func IsTerminal(w io.Writer) bool {
	return false
}