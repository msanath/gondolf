// package slog is a wrapper around the log/slog package.
// It provides a way to create a logger and add it to the context,
// and a way to retrieve the logger from the context.
package ctxslog

import (
	"context"
	"log/slog"
)

// ctxSlogKey is the key used to store the logger in the context.
type ctxSlogKey struct{}

// NewContext returns a new context with the given logger by copying the
// provided context and adding the logger to it.
func NewContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxSlogKey{}, logger)
}

// FromContext returns the logger from the given context.
// If no logger is found, it returns a no-op logger.
func FromContext(ctx context.Context) *slog.Logger {
	v := ctx.Value(ctxSlogKey{})
	if v == nil {
		return newNoOpSlogLogger()
	}
	if l, ok := v.(*slog.Logger); ok {
		return l
	}
	return newNoOpSlogLogger()
}

// newNoOpSlogLogger returns a no-op logger.
func newNoOpSlogLogger() *slog.Logger {
	return slog.New(&slogNoOpHandler{})
}

type slogNoOpHandler struct{}

func (s *slogNoOpHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (s *slogNoOpHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (s *slogNoOpHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return s
}

func (s *slogNoOpHandler) WithGroup(_ string) slog.Handler {
	return s
}
