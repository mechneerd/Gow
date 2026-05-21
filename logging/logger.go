package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Config represents logging configuration.
type Config struct {
	Level    string
	Channel  string // e.g. "single", "daily", "stack"
	Path     string // path to log directory
	Channels []string // used if Channel == "stack"
}

// dailyWriter handles writing to a file that rotates daily.
type dailyWriter struct {
	mu       sync.Mutex
	dir      string
	prefix   string
	currDate string
	file     *os.File
}

func newDailyWriter(dir, prefix string) *dailyWriter {
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("failed to create log dir: %v\n", err)
	}
	return &dailyWriter{
		dir:    dir,
		prefix: prefix,
	}
}

func (w *dailyWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().Format("2006-01-02")
	if w.currDate != now || w.file == nil {
		if w.file != nil {
			w.file.Close()
		}
		path := filepath.Join(w.dir, fmt.Sprintf("%s-%s.log", w.prefix, now))
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return 0, err
		}
		w.file = f
		w.currDate = now
	}
	return w.file.Write(p)
}

// Setup configures and returns the application logger.
func Setup(cfg Config) *slog.Logger {
	var l slog.Level
	switch cfg.Level {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	var writer io.Writer

	switch cfg.Channel {
	case "daily":
		writer = newDailyWriter(cfg.Path, "gow")
	case "stack":
		// Example simplistic stack channel mapping
		var writers []io.Writer
		for _, ch := range cfg.Channels {
			if ch == "daily" {
				writers = append(writers, newDailyWriter(cfg.Path, "gow"))
			} else if ch == "stdout" {
				writers = append(writers, os.Stdout)
			}
		}
		writer = io.MultiWriter(writers...)
	default: // single / stdout
		writer = os.Stdout
	}

	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: l,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

// ContextKey is a type for context keys to avoid collisions.
type ContextKey string

const RequestIDKey ContextKey = "request_id"

// WithContext enriches the logger with values extracted from context.
func WithContext(ctx context.Context, logger *slog.Logger) *slog.Logger {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return logger.With("request_id", reqID)
	}
	return logger
}
