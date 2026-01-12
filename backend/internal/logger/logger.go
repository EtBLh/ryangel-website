package logger

import (
	"strings"

	"github.com/ryangel/ryangel-backend/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New builds a zap logger configured from application settings.
func New(cfg *config.Config) (*zap.Logger, error) {
	var logConfig zap.Config
	if strings.EqualFold(cfg.AppEnv, "production") {
		logConfig = zap.NewProductionConfig()
	} else {
		logConfig = zap.NewDevelopmentConfig()
	}

	logConfig.Level = zap.NewAtomicLevelAt(parseLevel(cfg.LogLevel))
	return logConfig.Build()
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
