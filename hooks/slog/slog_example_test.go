package slog_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/sirupsen/logrus"
	lslog "github.com/sirupsen/logrus/hooks/slog"
)

// ExampleNewHandler demonstrates using slog alongside an existing Logrus logger.
func ExampleNewHandler() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})

	slog.SetDefault(slog.New(lslog.NewHandler(logger, nil)))

	// Both slog and Logrus write through the same Logrus backend.
	slog.Info("hello from slog", "source", "slog")
	logger.WithField("source", "logrus").Info("hello from logrus")

	// Output:
	// level=info msg="hello from slog" source=slog
	// level=info msg="hello from logrus" source=logrus
}

// ExampleNewHandler_options demonstrates configuring source reporting and
// custom slog-to-Logrus level mapping.
//
// Logrus writes to stderr by default, so the example has no stdout output.
// The log output will look similar to:
//
//	time="2026-01-02T03:04:05Z" level=info msg="regular info" func=github.com/sirupsen/logrus/hooks/slog_test.ExampleNewHandler_options file="/src/hooks/slog/slog_example_test.go:62" animal=walrus
//	time="2026-01-02T03:04:06Z" level=warning msg="custom level" func=github.com/sirupsen/logrus/hooks/slog_test.ExampleNewHandler_options file="/src/hooks/slog/slog_example_test.go:63" animal=walrus
func ExampleNewHandler_options() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})

	slogger := slog.New(lslog.NewHandler(logger, &lslog.HandlerOptions{
		// Preserve slog's source location as Logrus caller information.
		AddSource: true,

		// Map this custom slog level to Logrus WarnLevel.
		LevelMapper: func(level slog.Level) logrus.Level {
			if level == slog.LevelInfo+1 {
				return logrus.WarnLevel
			}
			return logrus.InfoLevel
		},
	}))

	slogger.Info("regular info", "animal", "walrus")
	slogger.Log(context.Background(), slog.LevelInfo+1, "custom level", "animal", "walrus")

	// Output:
}

// ExampleNewHook demonstrates forwarding existing Logrus logging to an slog
// logger, allowing applications to migrate their logging backend independently
// from code that still uses Logrus.
func ExampleNewHook() {
	var slogOutput bytes.Buffer
	slogger := slog.New(slog.NewTextHandler(&slogOutput, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return attr
		},
	}))

	logger := logrus.New()
	// Discard Logrus's own output; the hook forwards the entry to slog.
	logger.SetOutput(io.Discard)
	logger.AddHook(lslog.NewHook(slogger, nil))

	// Log through Logrus; the hook forwards this to slog.
	logger.WithField("animal", "walrus").Info("hello from logrus")

	// Show what was emitted by the slog handler.
	fmt.Print(slogOutput.String())

	// Output:
	// level=INFO msg="hello from logrus" animal=walrus
}

// ExampleNewHook_migration demonstrates using both adapters while migrating
// from Logrus to slog. Records can cross the bridge without being forwarded
// back into a Logrus logger they have already passed through.
func ExampleNewHook_migration() {
	var logrusOutput bytes.Buffer
	var slogOutput bytes.Buffer

	legacy := logrus.New()
	legacy.SetOutput(&logrusOutput)
	legacy.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})

	slogger := slog.New(slog.NewTextHandler(&slogOutput, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				return slog.Attr{} // Suppress timestamps for the example.
			}
			return attr
		},
	}))

	// Existing Logrus code is forwarded to the slog backend.
	legacy.AddHook(lslog.NewHook(slogger, nil))

	// New slog code can continue using the existing Logrus backend.
	bridged := slog.New(lslog.NewHandler(legacy, nil))

	legacy.Info("hello from logrus")
	bridged.Info("hello from slog")

	fmt.Println("Logrus output:")
	fmt.Print(logrusOutput.String())

	fmt.Println("slog output:")
	fmt.Print(slogOutput.String())

	// Output:
	// Logrus output:
	// level=info msg="hello from logrus"
	// level=info msg="hello from slog"
	// slog output:
	// level=INFO msg="hello from logrus"
	// level=INFO msg="hello from slog"
}
