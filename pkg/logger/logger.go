// pkg/logger/logger.go 수정
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

func NewLogger(level string) *zap.Logger {
	cfg := zap.NewProductionConfig()
	lvl := zapcore.InfoLevel
	if err := lvl.UnmarshalText([]byte(level)); err == nil {
		cfg.Level = zap.NewAtomicLevelAt(lvl)
	}

	// 사용자 홈 디렉토리에 로그 저장
	homeDir, _ := os.UserHomeDir()
	logFilePath := filepath.Join(homeDir, "auth-service.log")

	// 출력 경로 설정 - 파일과 표준 출력 모두 사용
	cfg.OutputPaths = []string{"stdout", logFilePath}

	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build(
		zap.AddCaller(),
		zap.Fields(
			zap.String("service", "auth-service"),
			zap.String("version", "1.0.0"),
		),
	)

	if err != nil {
		// 로거 초기화 실패 시 기본 로거 반환
		defaultLogger, _ := zap.NewProduction()
		defaultLogger.Error("로거 초기화 실패", zap.Error(err))
		return defaultLogger
	}

	return logger
}
