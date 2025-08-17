package logger

import "go.uber.org/zap"

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.zap.Sugar().Infof(format, args...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.zap.Sugar().Warnf(format, args...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.zap.Sugar().Errorf(format, args...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.zap.Sugar().Debugf(format, args...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.zap.Sugar().Fatalf(format, args...)
}

// InfoStruct logs a struct as a JSON string at Info level
func (l *Logger) InfoStruct(msg string, v any) {
	l.zap.Info(msg, zap.Any("data", toJSON(v)))
}

// ErrorStruct logs a struct as a JSON string at Error level
func (l *Logger) ErrorStruct(msg string, v any) {
	l.zap.Error(msg, zap.Any("data", toJSON(v)))
}

// DebugStruct logs a struct as a JSON string at Debug level
func (l *Logger) DebugStruct(msg string, v any) {
	l.zap.Debug(msg, zap.Any("data", toJSON(v)))
}

// FatalStruct logs a struct as a JSON string at Debug level
func (l *Logger) FatalStruct(msg string, v any) {
	l.zap.Fatal(msg, zap.Any("data", toJSON(v)))
}
