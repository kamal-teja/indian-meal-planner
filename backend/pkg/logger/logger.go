package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Logger wraps slog.Logger with additional convenience methods
type Logger struct {
	*slog.Logger
}

// New creates a new logger instance
func New(level, format string) *Logger {
	var handler slog.Handler

	// Parse log level
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// Create handler based on format
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// Fatal logs a fatal error and exits
func (l *Logger) Fatal(msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}
