package helpers

import (
	"context"
	"log/slog"
)

type loggerContextKey struct{}

// WithLogger returns a new context with the given logger attached.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

// LoggerFrom returns the request-scoped logger stored in ctx, or slog.Default().
func LoggerFrom(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerContextKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
