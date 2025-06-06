package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey string

const (
	requestIDKey ctxKey = "request_id"
	loggerKey    ctxKey = "logger"
)

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{zap: l.zap.With(fields...)}
}

func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.With(zap.String("request_id", requestID))
}

func (l *Logger) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.L()
	}

	// Prefer logger in context
	if lgr, ok := ctx.Value(loggerKey).(*zap.Logger); ok && lgr != nil {
		return lgr
	}

	// Fallback: global logger with optional request_id
	lgr := zap.L()
	if reqID, ok := ctx.Value(requestIDKey).(string); ok && reqID != "" {
		return lgr.With(zap.String("request_id", reqID))
	}
	return lgr
}
