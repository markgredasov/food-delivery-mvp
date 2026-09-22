// Package logger provides a structured logging wrapper around uber-go/zap.
//
// It supports logging to stdout with configurable log levels.
// The package also includes context-based logger retrieval and graceful shutdown.
package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type key string

var loggerKey key = "log"

// Logger wraps zap.Logger and manages the underlying log file.
//
// It provides structured logging capabilities with caller information
// and supports both console and file output simultaneously.
type Logger struct {
	*zap.Logger
}

// FromContext retrieves a Logger instance from the given context.
//
// The logger must have been previously stored in the context with the key "logger".
// If no logger is found, it panics to indicate a programming error.
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(loggerKey).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

// NewLogger creates a new Logger instance with the provided configuration.
//
// It sets up both console and file logging with the specified log level.
// The log file is created in the configured folder with a timestamp-based name.
// The encoder uses a custom time format for consistent human-readable output.
func NewLogger(cfg Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
	}, nil
}

// With creates a new Logger with additional fields attached.
//
// The returned logger inherits all fields from the parent logger and adds
// the provided fields. It shares the same underlying log file.
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
	}
}

// InContext appends logger in context and returns context with logger as value.
func (l *Logger) InContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}
