package logrus

type PreparedLogger interface {
	SetLevel(level Level)
	AsLevel(level Level) *LogEntry
	AsDebug() *LogEntry
	AsError() *LogEntry
	AsWarning() *LogEntry
	AsFatal() *LogEntry
	AsPanic() *LogEntry
}

