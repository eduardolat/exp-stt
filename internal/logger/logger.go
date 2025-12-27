package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger is a simple logging interface that supports different log levels.
type Logger interface {
	// SetDebug enables or disables debug level logging.
	SetDebug(enabled bool)
	// Info logs an informational message.
	Info(ctx context.Context, msg string, keysAndValues ...any)
	// Warn logs a warning message.
	Warn(ctx context.Context, msg string, keysAndValues ...any)
	// Error logs an error message.
	Error(ctx context.Context, msg string, keysAndValues ...any)
	// Debug logs a debug message. Only logged if debug level is enabled.
	Debug(ctx context.Context, msg string, keysAndValues ...any)
}

type slogLogger struct {
	slogger     *slog.Logger
	enableDebug bool
}

// NewStdLogger creates a new text Logger that writes to the standard output.
func NewSlogLogger(enableDebug bool) Logger {
	handlerOptions := &slog.HandlerOptions{}
	if enableDebug {
		handlerOptions.Level = slog.LevelDebug
	}

	slogger := slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))

	return &slogLogger{
		slogger:     slogger,
		enableDebug: enableDebug,
	}
}

func (l *slogLogger) SetDebug(enabled bool) {
	l.enableDebug = enabled
}

func (l *slogLogger) Info(ctx context.Context, msg string, keysAndValues ...any) {
	l.slogger.InfoContext(ctx, msg, keysAndValues...)
}

func (l *slogLogger) Warn(ctx context.Context, msg string, keysAndValues ...any) {
	l.slogger.WarnContext(ctx, msg, keysAndValues...)
}

func (l *slogLogger) Error(ctx context.Context, msg string, keysAndValues ...any) {
	l.slogger.ErrorContext(ctx, msg, keysAndValues...)
}

func (l *slogLogger) Debug(ctx context.Context, msg string, keysAndValues ...any) {
	if !l.enableDebug {
		return
	}
	l.slogger.DebugContext(ctx, msg, keysAndValues...)
}
