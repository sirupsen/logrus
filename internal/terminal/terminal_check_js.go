// +build js

package logrus

import (
	"io"
)

func IsTerminal(w io.Writer) bool {
	return false
}
