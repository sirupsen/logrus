package logrus

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"
	"time"

	errors "github.com/go-errors/errors"
)

var (
	// std is the name of the standard logger in stdlib `log`
	std = New()
)

func StandardLogger() *Logger {
	return std
}

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
	std.SetOutput(out)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter Formatter) {
	std.SetFormatter(formatter)
}

// SetReportCaller sets whether the standard logger will include the calling
// method as a field.
func SetReportCaller(include bool) {
	std.SetReportCaller(include)
}

// SetLevel sets the standard logger level.
func SetLevel(level Level) {
	std.SetLevel(level)
}

// GetLevel returns the standard logger level.
func GetLevel() Level {
	return std.GetLevel()
}

// IsLevelEnabled checks if the log level of the standard logger is greater than the level param
func IsLevelEnabled(level Level) bool {
	return std.IsLevelEnabled(level)
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook Hook) {
	std.AddHook(hook)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *Entry {
	return std.WithField(ErrorKey, err)
}

// WithContext creates an entry from the standard logger and adds a context to it.
func WithContext(ctx context.Context) *Entry {
	return std.WithContext(ctx)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *Entry {
	return std.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields Fields) *Entry {
	return std.WithFields(fields)
}

// WithTime creates an entry from the standard logger and overrides the time of
// logs generated with it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithTime(t time.Time) *Entry {
	return std.WithTime(t)
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	std.Trace(args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	std.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	std.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	std.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	std.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	std.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	std.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	std.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

// If error is non-nil, print error log via Error
func PrintOnError(err error, args ...interface{}) {
	if err != nil {
		if len(args) > 0 {
			args[0] = fmt.Sprintf("Error: %v (%v)", errorText(err), args[0])
			Error(args...)
		} else {
			Error("Error: " + errorText(err))
		}
	}
}

// If error is non-nil, panic via Panic
func PanicOnError(err error, args ...interface{}) {
	if err != nil {
		if len(args) > 0 {
			args[0] = fmt.Sprintf("Error: %v (%v)", errorText(err), args[0])
			Panic(args...)
		} else {
			Panic("Error: " + errorText(err))
		}
	}
}

// TraceFn logs a message from a func at level Trace on the standard logger.
func TraceFn(fn LogFunction) {
	std.TraceFn(fn)
}

// DebugFn logs a message from a func at level Debug on the standard logger.
func DebugFn(fn LogFunction) {
	std.DebugFn(fn)
}

// PrintFn logs a message from a func at level Info on the standard logger.
func PrintFn(fn LogFunction) {
	std.PrintFn(fn)
}

// InfoFn logs a message from a func at level Info on the standard logger.
func InfoFn(fn LogFunction) {
	std.InfoFn(fn)
}

// WarnFn logs a message from a func at level Warn on the standard logger.
func WarnFn(fn LogFunction) {
	std.WarnFn(fn)
}

// WarningFn logs a message from a func at level Warn on the standard logger.
func WarningFn(fn LogFunction) {
	std.WarningFn(fn)
}

// ErrorFn logs a message from a func at level Error on the standard logger.
func ErrorFn(fn LogFunction) {
	std.ErrorFn(fn)
}

// PanicFn logs a message from a func at level Panic on the standard logger.
func PanicFn(fn LogFunction) {
	std.PanicFn(fn)
}

// FatalFn logs a message from a func at level Fatal on the standard logger then the process will exit with status set to 1.
func FatalFn(fn LogFunction) {
	std.FatalFn(fn)
}

// If error is non-nil, print error log via ErrorFn
func PrintOnErrorFn(err error, fn LogFunction) {
	if err != nil {
		Errorf("Error: %v", errorText(err))
		ErrorFn(fn)
	}
}

// If error is non-nil, panic via PanicFn
func PanicOnErrorFn(err error, fn LogFunction) {
	if err != nil {
		Errorf("Error: %v", errorText(err))
		PanicFn(fn)
	}
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	std.Tracef(format, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	std.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

// If error is non-nil, print error log via Errorf
func PrintOnErrorf(err error, format string, args ...interface{}) {
	if err != nil {
		format = "Error: " + errorText(err) + " (" + format + ")"
		Errorf(format, args...)
	}
}

// If error is non-nil, panic via Panicf
func PanicOnErrorf(err error, format string, args ...interface{}) {
	if err != nil {
		format = "Error: " + errorText(err) + " (" + format + ")"
		Panicf(format, args...)
	}
}

// Traceln logs a message at level Trace on the standard logger.
func Traceln(args ...interface{}) {
	std.Traceln(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	std.Debugln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	std.Println(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	std.Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	std.Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	std.Warningln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	std.Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	std.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	std.Fatalln(args...)
}

// If error is non-nil, print error log via Errorln
func PrintOnErrorln(err error, args ...interface{}) {
	if err != nil {
		if len(args) > 0 {
			args[0] = fmt.Sprintf("Error: %v (%v)", errorText(err), args[0])
			Errorln(args...)
		} else {
			Errorln("Error: " + errorText(err))
		}
	}
}

// If error is non-nil, print error log via Errorln
func PanicOnErrorln(err error, args ...interface{}) {
	if err != nil {
		if len(args) > 0 {
			args[0] = fmt.Sprintf("Error: %v (%v)", errorText(err), args[0])
			Panicln(args...)
		} else {
			Panicln("Error: " + errorText(err))
		}
	}
}

func errorTextWithStackTrace(errorString string, errorWithStack interface{ StackFrames() []errors.StackFrame }) string {
	errorText := errorString + "\n"
	for _, frame := range errorWithStack.StackFrames() {
		errorText += frame.String()
	}
	return errorText + "\n"
}

func errorText(err error) string {
	if errorWithStack, ok := err.(interface{ StackFrames() []errors.StackFrame }); ok {
		return errorTextWithStackTrace(err.Error(), errorWithStack)
	} else {
		return errorTextWithStackTrace(err.Error(), errors.New(err))
	}
}

func CatchPanics() {
	// Panic handling - make sure panic gets reported in logs
	if err := recover(); err != nil {
		logEntry, ok := err.(*Entry)
		if ok { // Only report this to the user if it is of a type other than *BailInfo
			Fatalf("%s \n %v", logEntry.Message, string(debug.Stack()))
		} else {
			Fatalf("%s \n %v", err, string(debug.Stack()))
		}
	}
}
