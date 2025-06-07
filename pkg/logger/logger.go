package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

type Logger struct {
	zap *zap.Logger
}

func InitLogger(env string) (*Logger, error) {
	var cfg zap.Config

	if env == "dev" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.Encoding = "console"
	} else {
		cfg = zap.NewProductionConfig()
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		cfg.Encoding = "json"
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(logger)
	globalLogger = logger

	l := &Logger{zap: logger}
	return l, nil
}

func (l *Logger) Sync() func() error {
	return l.zap.Sync
}
