package logger

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger() *slog.Logger {
	var handler slog.Handler

	levelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))

	var level slog.Level

	switch strings.ToLower(levelStr) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError

	}

	handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)
}
