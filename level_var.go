package logrus

// Can be used with the `flag' package to parse Logrus level from
// command-line args.
//
// Example:
//		log := logrus.New()
//		lvl := logrus.NewLevelVar(log)
//		flag.Var(lvl, "log-level", "Sets the log level")
//		flag.Parse()
//
type LevelVar struct {
	l *Logger
}

// The logger should have it's default Level already set
func NewLevelVar(log *Logger) *LevelVar {
	return &LevelVar{log}
}

func (v *LevelVar) Set(str string) (err error) {
	var l Level
	if l, err = ParseLevel(str); err != nil {
		return err
	}
	v.l.Level = l
	return
}

func (v *LevelVar) String() string {
	return v.l.Level.String()
}
