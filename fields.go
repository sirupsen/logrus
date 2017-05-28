package logrus

import (
	"fmt"
)

var (
	FieldsLogger *Logger
)

func (logger *Logger) SetFieldsLogger() {
	FieldsLogger = logger
}

func (f Fields) withFields(fields Fields) *Entry {
	entry := FieldsLogger.newEntry()
	defer FieldsLogger.releaseEntry(entry)
	return entry.WithFields(f)
}

func (f Fields) Debugf(format string, args ...interface{}) {
	if FieldsLogger.level() >= DebugLevel {
		f.withFields(f).debugf(format, args...)
	}
}

func (f Fields) Infof(format string, args ...interface{}) {
	if FieldsLogger.level() >= InfoLevel {
		f.withFields(f).infof(format, args...)
	}
}

func (f Fields) Printf(format string, args ...interface{}) {
	if FieldsLogger.level() >= InfoLevel {
		f.withFields(f).printf(format, args...)
	}
}

func (f Fields) Warnf(format string, args ...interface{}) {
	if FieldsLogger.level() >= WarnLevel {
		f.withFields(f).warnf(format, args...)
	}
}

func (f Fields) Warningf(format string, args ...interface{}) {
	if FieldsLogger.level() >= WarnLevel {
		f.withFields(f).warnf(format, args...)
	}
}

func (f Fields) Errorf(format string, args ...interface{}) {
	if FieldsLogger.level() >= ErrorLevel {
		f.withFields(f).errorf(format, args...)
	}
}

func (f Fields) Fatalf(format string, args ...interface{}) {
	if FieldsLogger.level() >= FatalLevel {
		f.withFields(f).fatalf(format, args...)
	}
	Exit(1)
}

func (f Fields) Panicf(format string, args ...interface{}) {
	if FieldsLogger.level() >= PanicLevel {
		f.withFields(f).panicf(format, args...)
	}
}

func (f Fields) Debug(args ...interface{}) {
	if FieldsLogger.level() >= DebugLevel {
		f.withFields(f).debug(args...)
	}
}

func (f Fields) Info(args ...interface{}) {
	if FieldsLogger.level() >= InfoLevel {
		f.withFields(f).info(args...)
	}
}

func (f Fields) Print(args ...interface{}) {
	if FieldsLogger.level() >= InfoLevel {
		f.withFields(f).info(args...)
	}
}

func (f Fields) Warn(args ...interface{}) {
	if FieldsLogger.level() >= WarnLevel {
		f.withFields(f).warn(args...)
	}
}

func (f Fields) Warning(args ...interface{}) {
	if FieldsLogger.level() >= WarnLevel {
		f.withFields(f).warn(args...)
	}
}

func (f Fields) Error(args ...interface{}) {
	if FieldsLogger.level() >= ErrorLevel {
		f.withFields(f).error(args...)
	}
}

func (f Fields) Fatal(args ...interface{}) {
	if FieldsLogger.level() >= FatalLevel {
		f.withFields(f).fatal(args...)
	}
	Exit(1)
}

func (f Fields) Panic(args ...interface{}) {
	if FieldsLogger.level() >= PanicLevel {
		f.withFields(f).panic(args...)
	}
	Exit(1)
}

func (f Fields) Debugln(args ...interface{}) {
	if FieldsLogger.level() >= DebugLevel {
		f.withFields(f).debugln(args...)
	}
}

func (f Fields) Infoln(args ...interface{}) {
	if FieldsLogger.level() >= InfoLevel {
		f.withFields(f).infoln(args...)
	}
}

func (f Fields) Println(args ...interface{}) {
	if FieldsLogger.level() >= InfoLevel {
		f.withFields(f).println(args...)
	}
}

func (f Fields) Warnln(args ...interface{}) {
	if FieldsLogger.level() >= WarnLevel {
		f.withFields(f).warnln(args...)
	}
}

func (f Fields) Warningln(args ...interface{}) {
	if FieldsLogger.level() >= WarnLevel {
		f.withFields(f).warnln(args...)
	}
}

func (f Fields) Errorln(args ...interface{}) {
	if FieldsLogger.level() >= ErrorLevel {
		f.withFields(f).errorln(args...)
	}
}

func (f Fields) Fatalln(args ...interface{}) {
	if FieldsLogger.level() >= FatalLevel {
		f.withFields(f).fatalln(args...)
	}
	Exit(1)
}

func (f Fields) Panicln(args ...interface{}) {
	if FieldsLogger.level() >= PanicLevel {
		f.withFields(f).panicln(args...)
	}
	Exit(1)
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
	panic(fmt.Sprint(args...))
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
