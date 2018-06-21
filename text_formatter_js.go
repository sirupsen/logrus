// +build js

package logrus

import (
	"io"
)

func (f *TextFormatter) checkIfTerminal(w io.Writer) bool {
	return false
}
