package logger

import (
	"context"
)

// NoOpLogger is a logger that discards all log messages.
type NoOpLogger struct{}

func (l *NoOpLogger) SetDebug(enabled bool)                                       {}
func (l *NoOpLogger) Info(ctx context.Context, msg string, keysAndValues ...any)  {}
func (l *NoOpLogger) Warn(ctx context.Context, msg string, keysAndValues ...any)  {}
func (l *NoOpLogger) Error(ctx context.Context, msg string, keysAndValues ...any) {}
func (l *NoOpLogger) Debug(ctx context.Context, msg string, keysAndValues ...any) {}
