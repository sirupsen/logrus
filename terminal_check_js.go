//go:build js || tinygo
// +build js tinygo

package logrus

func isTerminal(fd int) bool {
	return false
}
