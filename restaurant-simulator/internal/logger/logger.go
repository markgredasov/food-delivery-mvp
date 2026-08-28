// Package logger builds the simulator's zap.Logger.
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New builds a JSON zap.Logger at the given level, defaulting to "info".
func New(level string) (*zap.Logger, error) {
	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = zapcore.InfoLevel
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(lvl)
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return cfg.Build()
}
