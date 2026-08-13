package slog_test

import (
	"log/slog"
	"os"

	"github.com/sirupsen/logrus"
	lslog "github.com/sirupsen/logrus/hooks/slog"
)

// LoggerLevel exposes a Logrus logger's configured level as a slog.Leveler.
//
// This lets slog handlers follow changes to the Logrus logger's level without
// keeping a separate slog.LevelVar in sync.
type LoggerLevel struct {
	Logger *logrus.Logger
}

func (l LoggerLevel) Level() slog.Level {
	return lslog.Level(l.Logger.GetLevel()).Level()
}

// ExampleLevel_dynamic demonstrates using Level to share dynamic level
// configuration between Logrus and slog.
func ExampleLevel_dynamic() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: LoggerLevel{Logger: logger},
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				return slog.Attr{} // Suppress timestamps for the example.
			}
			return attr
		},
	}))

	slogger.Info("before")

	// The slog handler follows changes to the Logrus logger's level without
	// requiring a separate slog.LevelVar to be updated.
	logger.SetLevel(logrus.WarnLevel)

	slogger.Info("ignored")
	slogger.Warn("after")

	// Output:
	// level=INFO msg=before
	// level=WARN msg=after
}
