package logrus_test

import (
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func ExampleSkipper() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Out = os.Stdout
	l.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	}

	l.AddSkipper(logrus.NewPackageSkipper("github.com/sirupsen/logrus_test"))
	l.Info("example of custom skipper")
	// Output:
	// {"file":"run_example.go","func":"runExample","level":"info","msg":"example of custom skipper"}
}
