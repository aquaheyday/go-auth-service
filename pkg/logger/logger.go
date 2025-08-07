package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level string) *zap.Logger {
	cfg := zap.NewProductionConfig()
	lvl := zapcore.InfoLevel
	if err := lvl.UnmarshalText([]byte(level)); err == nil {
		cfg.Level = zap.NewAtomicLevelAt(lvl)
	}
	logger, _ := cfg.Build()
	return logger
}
