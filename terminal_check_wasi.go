//go:build wasi

package logrus

func isTerminal(fd int) bool {
	return false
}
