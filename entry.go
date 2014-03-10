package logrus

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

type Entry struct {
	logger *Logger
	Data   Fields
}

var baseTimestamp time.Time

func init() {
	baseTimestamp = time.Now()
}

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		logger: logger,
		// Default is three fields, give a little extra room
		Data: make(Fields, 5),
	}
}

func (entry *Entry) Reader() (*bytes.Buffer, error) {
	serialized, err := entry.logger.Formatter.Format(entry)
	return bytes.NewBuffer(serialized), err
}

func (entry *Entry) String() (string, error) {
	reader, err := entry.Reader()
	if err != nil {
		return "", err
	}

	return reader.String(), err
}

func (entry *Entry) WithField(key string, value interface{}) *Entry {
	entry.Data[key] = value
	return entry
}

func (entry *Entry) WithFields(fields Fields) *Entry {
	for key, value := range fields {
		entry.WithField(key, value)
	}
	return entry
}

func (entry *Entry) log(level string, msg string) string {
	entry.Data["time"] = time.Now().String()
	entry.Data["level"] = level
	entry.Data["msg"] = msg

	reader, err := entry.Reader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v", err)
	}

	entry.logger.mu.Lock()
	defer entry.logger.mu.Unlock()

	_, err = io.Copy(entry.logger.Out, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v", err)
	}

	return reader.String()
}

func (entry *Entry) Debug(args ...interface{}) {
	if Level >= LevelDebug {
		entry.log("debug", fmt.Sprint(args...))
		entry.logger.Hooks.Fire(LevelDebug, entry)
	}
}

func (entry *Entry) Info(args ...interface{}) {
	if Level >= LevelInfo {
		entry.log("info", fmt.Sprint(args...))
		entry.logger.Hooks.Fire(LevelInfo, entry)
	}
}

func (entry *Entry) Print(args ...interface{}) {
	if Level >= LevelInfo {
		entry.log("info", fmt.Sprint(args...))
		entry.logger.Hooks.Fire(LevelInfo, entry)
	}
}

func (entry *Entry) Warn(args ...interface{}) {
	if Level >= LevelWarn {
		entry.log("warning", fmt.Sprint(args...))
		entry.logger.Hooks.Fire(LevelWarn, entry)
	}
}

func (entry *Entry) Error(args ...interface{}) {
	if Level >= LevelError {
		entry.log("error", fmt.Sprint(args...))
		entry.logger.Hooks.Fire(LevelError, entry)
	}
}

func (entry *Entry) Fatal(args ...interface{}) {
	if Level >= LevelFatal {
		entry.log("fatal", fmt.Sprint(args...))
		entry.logger.Hooks.Fire(LevelFatal, entry)
	}
	os.Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	if Level >= LevelPanic {
		msg := entry.log("panic", fmt.Sprint(args...))
		entry.logger.Hooks.Fire(LevelPanic, entry)
		panic(msg)
	}
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Debug(fmt.Sprintf(format, args...))
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Info(fmt.Sprintf(format, args...))
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Print(fmt.Sprintf(format, args...))
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.Warn(fmt.Sprintf(format, args...))
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.Warn(fmt.Sprintf(format, args...))
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.Print(fmt.Sprintf(format, args...))
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.Fatal(fmt.Sprintf(format, args...))
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	entry.Panic(fmt.Sprintf(format, args...))
}

// Entry Println family functions

func (entry *Entry) Debugln(args ...interface{}) {
	entry.Debug(fmt.Sprint(args...))
}

func (entry *Entry) Infoln(args ...interface{}) {
	entry.Info(fmt.Sprint(args...))
}

func (entry *Entry) Println(args ...interface{}) {
	entry.Print(fmt.Sprint(args...))
}

func (entry *Entry) Warnln(args ...interface{}) {
	entry.Warn(fmt.Sprint(args...))
}

func (entry *Entry) Warningln(args ...interface{}) {
	entry.Warn(fmt.Sprint(args...))
}

func (entry *Entry) Errorln(args ...interface{}) {
	entry.Error(fmt.Sprint(args...))
}

func (entry *Entry) Fatalln(args ...interface{}) {
	entry.Fatal(fmt.Sprint(args...))
}

func (entry *Entry) Panicln(args ...interface{}) {
	entry.Panic(fmt.Sprint(args...))
}
