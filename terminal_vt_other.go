//go:build !windows

package logrus

func enableVirtualTerminalProcessing(fd uintptr) bool {
	return true
}
