package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Config represents logging configuration.
type Config struct {
	Level    string
	Channel  string // e.g. "single", "daily", "stack"
	Path     string // path to log directory
	Channels []string // used if Channel == "stack"
	MaxAge   int      // max days to keep log files (0 = unlimited)
	MaxSize  int64    // max file size in bytes (0 = unlimited)
}

// dailyWriter handles writing to a file that rotates daily.
type dailyWriter struct {
	mu       sync.Mutex
	dir      string
	prefix   string
	currDate string
	file     *os.File
	maxAge   int
	maxSize  int64
}

func newDailyWriter(dir, prefix string, maxAge int, maxSize int64) *dailyWriter {
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("failed to create log dir: %v\n", err)
	}
	w := &dailyWriter{
		dir:     dir,
		prefix:  prefix,
		maxAge:  maxAge,
		maxSize: maxSize,
	}
	w.cleanup()
	return w
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

func (w *dailyWriter) cleanup() {
	if w.maxAge <= 0 {
		return
	}
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return
	}
	cutoff := time.Now().AddDate(0, 0, -w.maxAge)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".log") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(w.dir, entry.Name()))
		}
	}
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
		writer = newDailyWriter(cfg.Path, "gow", cfg.MaxAge, cfg.MaxSize)
	case "stack":
		var writers []io.Writer
		for _, ch := range cfg.Channels {
			if ch == "daily" {
				writers = append(writers, newDailyWriter(cfg.Path, "gow", cfg.MaxAge, cfg.MaxSize))
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

const (
	RequestIDKey    ContextKey = "request_id"
	UserIDKey       ContextKey = "user_id"
	TenantIDKey     ContextKey = "tenant_id"
	ChannelKey      ContextKey = "channel"
	ExtraKey        ContextKey = "extra"
)

// WithContext enriches the logger with values extracted from context.
func WithContext(ctx context.Context, logger *slog.Logger) *slog.Logger {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		logger = logger.With("request_id", reqID)
	}
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		logger = logger.With("user_id", userID)
	}
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
		logger = logger.With("tenant_id", tenantID)
	}
	if channel, ok := ctx.Value(ChannelKey).(string); ok {
		logger = logger.With("channel", channel)
	}
	return logger
}

// With returns a logger enriched with the given key-value pairs.
func With(logger *slog.Logger, args ...any) *slog.Logger {
	return logger.With(args...)
}

// LogRequest logs an HTTP request with standard fields.
func LogRequest(r *http.Request, statusCode int, duration time.Duration, logger *slog.Logger) {
	logger.Info("HTTP Request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", statusCode,
		"duration_ms", duration.Milliseconds(),
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent(),
	)
}

// LogError logs an error with context.
func LogError(err error, logger *slog.Logger, args ...any) {
	logger.Error(err.Error(), args...)
}

// LogWithContext creates a logger enriched from an HTTP request context.
func LogWithContext(r *http.Request, logger *slog.Logger) *slog.Logger {
	return WithContext(r.Context(), logger)
}

// RequestMiddleware returns an HTTP middleware that enriches the context with request metadata.
func RequestMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract or generate request ID
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = generateRequestID()
			}

			// Add request ID to context
			ctx := context.WithValue(r.Context(), RequestIDKey, reqID)

			// Add user ID if present in header (for API auth)
			if userID := r.Header.Get("X-User-ID"); userID != "" {
				ctx = context.WithValue(ctx, UserIDKey, userID)
			}

			// Set response header
			w.Header().Set("X-Request-ID", reqID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func generateRequestID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), os.Getpid())
}

