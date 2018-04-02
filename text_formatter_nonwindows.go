// +build !appengine,!gopherjs,!windows

package logrus

func (f *TextFormatter) initTerminal(entry *Entry) {
}
