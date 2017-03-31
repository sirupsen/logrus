package logrus

import (
	"bufio"
	"io"
	"runtime"
	"bytes"
)

func (logger *Logger) Writer() *io.PipeWriter {
	return logger.WriterLevel(InfoLevel)
}

func (logger *Logger) WriterLevel(level Level) *io.PipeWriter {
	return NewEntry(logger).WriterLevel(level)
}

func (entry *Entry) Writer() *io.PipeWriter {
	return entry.WriterLevel(InfoLevel)
}

func (entry *Entry) WriterLevel(level Level) *io.PipeWriter {
	reader, writer := io.Pipe()

	go entry.writerScanner(reader)
	runtime.SetFinalizer(writer, writerFinalizer)

	return writer
}

func (entry *Entry) writerScanner(reader *io.PipeReader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		txt := scanner.Text()
		printFunc := entry.get_level_function(get_level(txt))
		printFunc(txt)
	}
	if err := scanner.Err(); err != nil {
		entry.Errorf("Error while reading from Writer: %s", err)
	}
	reader.Close()
}

func writerFinalizer(writer *io.PipeWriter) {
	writer.Close()
}

func (entry *Entry) get_level_function(level Level) func(args ...interface{}) {
	var printFunc func(args ...interface{})

	switch level {
	case DebugLevel:
		printFunc = entry.Debug
	case InfoLevel:
		printFunc = entry.Info
	case WarnLevel:
		printFunc = entry.Warn
	case ErrorLevel:
		printFunc = entry.Error
	case FatalLevel:
		printFunc = entry.Fatal
	case PanicLevel:
		printFunc = entry.Panic
	default:
		printFunc = entry.Print
	}

	return printFunc
}

func get_level(line string) Level {
	var lvl string
	line_b := []byte(line)
	x := bytes.IndexByte(line_b, '[')
	if x >= 0 {
		y := bytes.IndexByte(line_b[x:], ']')
		if y >= 0 {
			lvl = string(line_b[x+1 : x+y])
		}
	}
	level, err := ParseLevel(lvl)
	if err != nil {
		level = InfoLevel
	}
	return level
}
