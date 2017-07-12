package stack

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	vendor = "/vendor/"
)

var (
	goPath     = fmt.Sprintf("%s/src/", os.Getenv("GOPATH"))
	goPath_len = len(goPath)
)

// To hook (that is, switch mode) to support the print source code line of information,
// such as which one package, which file, which function, which line
type CodeLineHook struct {
	LogLevel logrus.Level
}

type filePos struct {
	Pkg  string `json:"pkg"`
	File string `json:"file"`
	Func string `json:"func"`
	Line int    `json:"line"`
}

func (hook *CodeLineHook) Fire(entry *logrus.Entry) (_ error) {
	if pc, fullPath, line, ok := runtime.Caller(entry.Skip); ok {
		funcName := runtime.FuncForPC(pc).Name()
		relativePath := fullPath
		if temp := vendorPath(fullPath); temp != "" {
			relativePath = temp
		}

		pkg, file := path.Split(relativePath)
		if pkg != "" {
			pkg = pkg[:len(pkg)-1]
		}
		entry.Data["pos"] = filePos{
			Pkg:  pkg,
			File: file,
			Func: path.Base(funcName),
			Line: line,
		}
	}
	return
}

func (hook *CodeLineHook) Levels() []logrus.Level {
	levels := make([]logrus.Level, hook.LogLevel+1)
	for i, _ := range levels {
		levels[i] = logrus.Level(i)
	}
	return levels
}

func vendorPath(fullPath string) string {
	if i := strings.Index(fullPath, vendor); i != -1 {
		return fullPath[i+len(vendor):]
	}
	return trimGoPath(fullPath)
}

func trimGoPath(fullPath string) string {
	if i := strings.Index(fullPath, goPath); i != -1 {
		return fullPath[i+goPath_len:]
	}
	return ""
}
