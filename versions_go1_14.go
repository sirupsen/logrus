// +build go1.14

package logrus

import "runtime"

// funcName returns the function name that logrus calls
func funcName(pcs []uintptr) string {
	return runtime.FuncForPC(pcs[0]).Name()
}
