package logger

import (
	"context"
	"log/slog"

	"github.com/asyauqi15/payslip-system/internal/constant"
)

// WithRequestID returns a logger with the request ID from context included
func WithRequestID(ctx context.Context) *slog.Logger {
	requestID := GetRequestID(ctx)
	if requestID == "" {
		return slog.Default()
	}
	return slog.With("request_id", requestID)
}

// GetRequestID extracts the request ID from context
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestID, ok := ctx.Value(constant.ContextKeyRequestID).(string); ok {
		return requestID
	}
	return ""
}

// WithFields returns a logger with request ID and additional fields
func WithFields(ctx context.Context, fields ...any) *slog.Logger {
	logger := WithRequestID(ctx)
	if len(fields) > 0 {
		logger = logger.With(fields...)
	}
	return logger
}

// Info logs an info message with request ID
func Info(ctx context.Context, msg string, args ...any) {
	WithRequestID(ctx).Info(msg, args...)
}

// Error logs an error message with request ID
func Error(ctx context.Context, msg string, args ...any) {
	WithRequestID(ctx).Error(msg, args...)
}

// Warn logs a warning message with request ID
func Warn(ctx context.Context, msg string, args ...any) {
	WithRequestID(ctx).Warn(msg, args...)
}

// Debug logs a debug message with request ID
func Debug(ctx context.Context, msg string, args ...any) {
	WithRequestID(ctx).Debug(msg, args...)
}
