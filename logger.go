package logrus

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"
)

type Logger struct {
	// The logs are `io.Copy`'d to this in a mutex. It's common to set this to a
	// file, or leave it default which is `os.Stderr`. You can also set this to
	// something more adventorous, such as logging to Kafka.
	Out io.Writer
	// Hooks for the logger instance. These allow firing events based on logging
	// levels and log entries. For example, to send errors to an error tracking
	// service, log to StatsD or dump the core on fatal errors.
	Hooks LevelHooks
	// All log entries pass through the formatter before logged to Out. The
	// included formatters are `TextFormatter` and `JSONFormatter` for which
	// TextFormatter is the default. In development (when a TTY is attached) it
	// logs with colors, but to a file it wouldn't. You can easily implement your
	// own that implements the `Formatter` interface, see the `README` or included
	// formatters for examples.
	Formatter Formatter
	// The logging level the logger should log at. This is typically (and defaults
	// to) `logrus.Info`, which allows Info(), Warn(), Error() and Fatal() to be
	// logged. `logrus.Debug` is useful in
	Level Level
	// Used to sync writing to the log.
	mu sync.Mutex
}

// Creates a new logger. Configuration should be set by changing `Formatter`,
// `Out` and `Hooks` directly on the default logger instance. You can also just
// instantiate your own:
//
//    var log = &Logger{
//      Out: os.Stderr,
//      Formatter: new(JSONFormatter),
//      Hooks: make(LevelHooks),
//      Level: logrus.DebugLevel,
//    }
//
// It's recommended to make this a global instance called `log`.
func New() *Logger {
	return &Logger{
		Out:       os.Stderr,
		Formatter: new(TextFormatter),
		Hooks:     make(LevelHooks),
		Level:     InfoLevel,
	}
}

// Adds a field to the log entry, note that you it doesn't log until you call
// Debug, Print, Info, Warn, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (logger *Logger) WithField(key string, value interface{}) *Entry {
	return NewEntry(logger).WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *Logger) WithFields(fields Fields) *Entry {
	return NewEntry(logger).WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (logger *Logger) WithError(err error) *Entry {
	return NewEntry(logger).WithError(err)
}

type LevelParser interface {
	Parse(*string) (Level, error)
}

type RegexpParser struct {
	r *regexp.Regexp
}

func (pr *RegexpParser) prefixRegex() {
	//pr.r = regexp.MustCompile(`^\\[(?P<Level?\\w+)\\]`)
	pr.r = regexp.MustCompile(`^\[\w+\]`)
}

func (pr *RegexpParser) Parse(s *string) (Level, error) {
	//b := string(pr.r.Find([]byte(*s)))
	//_ = pr.r.FindStringSubmatch(*s)
	b := pr.r.Find([]byte(*s))
	return ParseLevel(string(b)[1 : len(b)-1])
}

type PrefixStrCmp struct{}

func (p *PrefixStrCmp) Parse(s *string) (Level, error) {
	str := *s
	prefix := str[:7]

	switch prefix {
	case "[INFO] ":
		return InfoLevel, nil
	case "[WARN] ":
		return WarnLevel, nil
	case "[ERROR]":
		return ErrorLevel, nil
	case "[FATAL]":
		return FatalLevel, nil
	case "[DEBUG]":
		return DebugLevel, nil
	case "[PANIC]":
		return PanicLevel, nil
	default:
		return DebugLevel, nil
	}
	return WarnLevel, fmt.Errorf("prefixstrcmp switch failed?")
}

// Output pulls information on the function caller before creating a log Entry
func (logger *Logger) Output(calldepth int, s string) error {
	var (
		file string
		line int
		ok   bool
	)
	ps := PrefixStrCmp{}
	_, file, line, ok = runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}
	entry := NewEntry(logger).WithFields(map[string]interface{}{"file": file, "line": strconv.Itoa(line)})
	lvl, err := ps.Parse(&s)
	if err != nil {
		NewEntry(logger).log(ErrorLevel, fmt.Sprintf("error parsing log line level: '%s'", s[:10]))
		entry.log(ErrorLevel, s)
		return err
	}
	entry.log(lvl, s)

	return nil
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	if logger.Level >= DebugLevel {
		NewEntry(logger).Debugf(format, args...)
	}
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	if logger.Level >= InfoLevel {
		NewEntry(logger).Infof(format, args...)
	}
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	NewEntry(logger).Printf(format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	if logger.Level >= WarnLevel {
		NewEntry(logger).Warnf(format, args...)
	}
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	if logger.Level >= WarnLevel {
		NewEntry(logger).Warnf(format, args...)
	}
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	if logger.Level >= ErrorLevel {
		NewEntry(logger).Errorf(format, args...)
	}
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	if logger.Level >= FatalLevel {
		NewEntry(logger).Fatalf(format, args...)
	}
	os.Exit(1)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	if logger.Level >= PanicLevel {
		NewEntry(logger).Panicf(format, args...)
	}
}

func (logger *Logger) Debug(args ...interface{}) {
	if logger.Level >= DebugLevel {
		NewEntry(logger).Debug(args...)
	}
}

func (logger *Logger) Info(args ...interface{}) {
	if logger.Level >= InfoLevel {
		NewEntry(logger).Info(args...)
	}
}

func (logger *Logger) Print(args ...interface{}) {
	NewEntry(logger).Info(args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	if logger.Level >= WarnLevel {
		NewEntry(logger).Warn(args...)
	}
}

func (logger *Logger) Warning(args ...interface{}) {
	if logger.Level >= WarnLevel {
		NewEntry(logger).Warn(args...)
	}
}

func (logger *Logger) Error(args ...interface{}) {
	if logger.Level >= ErrorLevel {
		NewEntry(logger).Error(args...)
	}
}

func (logger *Logger) Fatal(args ...interface{}) {
	if logger.Level >= FatalLevel {
		NewEntry(logger).Fatal(args...)
	}
	os.Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	if logger.Level >= PanicLevel {
		NewEntry(logger).Panic(args...)
	}
}

func (logger *Logger) Debugln(args ...interface{}) {
	if logger.Level >= DebugLevel {
		NewEntry(logger).Debugln(args...)
	}
}

func (logger *Logger) Infoln(args ...interface{}) {
	if logger.Level >= InfoLevel {
		NewEntry(logger).Infoln(args...)
	}
}

func (logger *Logger) Println(args ...interface{}) {
	NewEntry(logger).Println(args...)
}

func (logger *Logger) Warnln(args ...interface{}) {
	if logger.Level >= WarnLevel {
		NewEntry(logger).Warnln(args...)
	}
}

func (logger *Logger) Warningln(args ...interface{}) {
	if logger.Level >= WarnLevel {
		NewEntry(logger).Warnln(args...)
	}
}

func (logger *Logger) Errorln(args ...interface{}) {
	if logger.Level >= ErrorLevel {
		NewEntry(logger).Errorln(args...)
	}
}

func (logger *Logger) Fatalln(args ...interface{}) {
	if logger.Level >= FatalLevel {
		NewEntry(logger).Fatalln(args...)
	}
	os.Exit(1)
}

func (logger *Logger) Panicln(args ...interface{}) {
	if logger.Level >= PanicLevel {
		NewEntry(logger).Panicln(args...)
	}
}
