package logger

import (
	"log/slog"
	"os"
)

func SetupLogger(debugMode bool) *slog.Logger {
	var handler slog.Handler

	level := slog.LevelInfo
	if debugMode {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	if debugMode {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
