package logging

import (
	"log/slog"
	"os"
)

// Setup creates and configures the default logger.
func Setup(level string) *slog.Logger {
	var l slog.Level
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: l,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger) // Set as global default for convenience

	return logger
}
