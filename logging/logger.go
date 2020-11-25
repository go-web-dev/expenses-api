package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/steevehook/expenses-rest-api/config"
)

// Logger represents application logger
var Logger *zap.Logger

// Init initializes application logger
func Init(cfg *config.Manager) error {
	var logLevel zapcore.Level
	err := logLevel.Set(cfg.LoggingLevel())
	if err != nil {
		return err
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(logLevel)
	zapConfig.OutputPaths = cfg.LoggingOutput()
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := zapConfig.Build()
	defer func() {
		_ = logger.Sync()
	}()
	Logger = logger
	return err
}
