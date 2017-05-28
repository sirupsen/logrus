package logrus

import (
	"fmt"
)

var (
	FieldsLogger *Logger
)

const (
	SKIP_4 = 4
	SKIP_5 = 5
	SKIP_6 = 6
	SKIP_7 = 7
)

func (logger *Logger) SetFieldsLogger() {
	FieldsLogger = logger
}

func (f Fields) withFields(fields Fields, skip int) *Entry {
	entry := FieldsLogger.newEntry()
	defer FieldsLogger.releaseEntry(entry)
	return entry.WithFields(f).WithSkip(skip)
}

func (f Fields) Debugf(format string, args ...interface{}) {
	if FieldsLogger.level() >= DebugLevel {
		f.withFields(f, SKIP_6).debugf(format, args...)
	}
}

func (f Fields) Infof(format string, args ...interface{}) {
	if FieldsLogger.level() >= InfoLevel {
		f.withFields(f, SKIP_6).infof(format, args...)
	}
}

func (f Fields) Printf(format string, args ...interface{}) {
	if FieldsLogger.level() >= InfoLevel {
		f.withFields(f, SKIP_7).printf(format, args...)
	}
}

func (f Fields) Warnf(format string, args ...interface{}) {
	if FieldsLogger.level() >= WarnLevel {
		f.withFields(f, SKIP_6).warnf(format, args...)
	}
}

func (f Fields) Warningf(format string, args ...interface{}) {
	if FieldsLogger.level() >= WarnLevel {
		f.withFields(f, SKIP_6).warnf(format, args...)
	}
}

func (f Fields) Errorf(format string, args ...interface{}) {
	if FieldsLogger.level() >= ErrorLevel {
		f.withFields(f, SKIP_6).errorf(format, args...)
	}
}

func (f Fields) Fatalf(format string, args ...interface{}) {
	if FieldsLogger.level() >= FatalLevel {
		f.withFields(f, SKIP_6).fatalf(format, args...)
	}
}

func (f Fields) Panicf(format string, args ...interface{}) {
	if FieldsLogger.level() >= PanicLevel {
		f.withFields(f, SKIP_6).panicf(format, args...)
	}
}

func (f Fields) Debug(args ...interface{}) {
	if FieldsLogger.level() >= DebugLevel {
		f.withFields(f, SKIP_5).debug(args...)
	}
}

func (f Fields) Info(args ...interface{}) {
	if FieldsLogger.level() >= InfoLevel {
		f.withFields(f, SKIP_5).info(args...)
	}
}

func (f Fields) Print(args ...interface{}) {
	if FieldsLogger.level() >= InfoLevel {
		f.withFields(f, SKIP_5).info(args...)
	}
}

func (f Fields) Warn(args ...interface{}) {
	if FieldsLogger.level() >= WarnLevel {
		f.withFields(f, SKIP_5).warn(args...)
	}
}

func (f Fields) Warning(args ...interface{}) {
	if FieldsLogger.level() >= WarnLevel {
		f.withFields(f, SKIP_5).warn(args...)
	}
}

func (f Fields) Error(args ...interface{}) {
	if FieldsLogger.level() >= ErrorLevel {
		f.withFields(f, SKIP_5).error(args...)
	}
}

func (f Fields) Fatal(args ...interface{}) {
	if FieldsLogger.level() >= FatalLevel {
		f.withFields(f, SKIP_5).fatal(args...)
	}
}

func (f Fields) Panic(args ...interface{}) {
	if FieldsLogger.level() >= PanicLevel {
		f.withFields(f, SKIP_5).panic(args...)
	}
}

func (f Fields) Debugln(args ...interface{}) {
	if FieldsLogger.level() >= DebugLevel {
		f.withFields(f, SKIP_6).debugln(args...)
	}
}

func (f Fields) Infoln(args ...interface{}) {
	if FieldsLogger.level() >= InfoLevel {
		f.withFields(f, SKIP_6).infoln(args...)
	}
}

func (f Fields) Println(args ...interface{}) {
	if FieldsLogger.level() >= InfoLevel {
		f.withFields(f, SKIP_7).println(args...)
	}
}

func (f Fields) Warnln(args ...interface{}) {
	if FieldsLogger.level() >= WarnLevel {
		f.withFields(f, SKIP_6).warnln(args...)
	}
}

func (f Fields) Warningln(args ...interface{}) {
	if FieldsLogger.level() >= WarnLevel {
		f.withFields(f, SKIP_6).warnln(args...)
	}
}

func (f Fields) Errorln(args ...interface{}) {
	if FieldsLogger.level() >= ErrorLevel {
		f.withFields(f, SKIP_6).errorln(args...)
	}
}

func (f Fields) Fatalln(args ...interface{}) {
	if FieldsLogger.level() >= FatalLevel {
		f.withFields(f, SKIP_6).fatalln(args...)
	}
}

func (f Fields) Panicln(args ...interface{}) {
	if FieldsLogger.level() >= PanicLevel {
		f.withFields(f, SKIP_6).panicln(args...)
	}
}

// The entry object should not be added to the log level to judge.
// Should be in the top function call to determine, then the best performance.
func (entry *Entry) print(args ...interface{}) {
	entry.info(args...)
}

func (entry *Entry) debug(args ...interface{}) {
	entry.log(DebugLevel, fmt.Sprint(args...))
}

func (entry *Entry) info(args ...interface{}) {
	entry.log(InfoLevel, fmt.Sprint(args...))
}

func (entry *Entry) warn(args ...interface{}) {
	entry.log(WarnLevel, fmt.Sprint(args...))
}

func (entry *Entry) error(args ...interface{}) {
	entry.log(ErrorLevel, fmt.Sprint(args...))
}

func (entry *Entry) fatal(args ...interface{}) {
	entry.log(FatalLevel, fmt.Sprint(args...))
}

func (entry *Entry) panic(args ...interface{}) {
	entry.log(PanicLevel, fmt.Sprint(args...))
}

// Entry Printf family functions
func (entry *Entry) debugf(format string, args ...interface{}) {
	entry.debug(fmt.Sprintf(format, args...))
}

func (entry *Entry) infof(format string, args ...interface{}) {
	entry.info(fmt.Sprintf(format, args...))
}

func (entry *Entry) warnf(format string, args ...interface{}) {
	entry.warn(fmt.Sprintf(format, args...))
}

func (entry *Entry) errorf(format string, args ...interface{}) {
	entry.error(fmt.Sprintf(format, args...))
}

func (entry *Entry) fatalf(format string, args ...interface{}) {
	entry.fatal(fmt.Sprintf(format, args...))
}

func (entry *Entry) panicf(format string, args ...interface{}) {
	entry.panic(fmt.Sprintf(format, args...))

}

// Entry Println family functions

func (entry *Entry) debugln(args ...interface{}) {
	entry.debug(entry.sprintlnn(args...))
}

func (entry *Entry) infoln(args ...interface{}) {
	entry.info(entry.sprintlnn(args...))
}

func (entry *Entry) warnln(args ...interface{}) {
	entry.warn(entry.sprintlnn(args...))

}

func (entry *Entry) errorln(args ...interface{}) {
	entry.error(entry.sprintlnn(args...))
}

func (entry *Entry) fatalln(args ...interface{}) {
	entry.fatal(entry.sprintlnn(args...))
}

func (entry *Entry) panicln(args ...interface{}) {
	entry.panic(entry.sprintlnn(args...))
}

func (entry *Entry) println(args ...interface{}) {
	entry.infoln(args...)
}

func (entry *Entry) printf(format string, args ...interface{}) {
	entry.infof(format, args...)
}
