package testutils

import (
	"time"
)

// SlowTestWriter is a io.Writer data sync to help test
// for race conditions.
type SlowTestWriter struct {
	counter int
}

// Write pretend to write data (with irregular intermittent delay)
func (stw *SlowTestWriter) Write(p []byte) (int, error) {
	// Random-ish delay to highlight concurrency issues.
	if ((stw.counter % 7) % 4) > 2 {
		time.Sleep(2 * time.Millisecond)
	}
	stw.counter++
	return len(p), nil
}
