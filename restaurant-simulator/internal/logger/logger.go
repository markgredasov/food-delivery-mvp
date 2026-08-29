// Package logger provides a structured logging wrapper around uber-go/zap.
//
// It supports logging to both stdout and a file with configurable log levels.
// The package also includes context-based logger retrieval and graceful shutdown.
package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger and manages the underlying log file.
//
// It provides structured logging capabilities with caller information
// and supports both console and file output simultaneously.
type Logger struct {
	*zap.Logger

	file *os.File // underlying log file handle for cleanup
}

// FromContext retrieves a Logger instance from the given context.
//
// The logger must have been previously stored in the context with the key "logger".
// If no logger is found, it panics to indicate a programming error.
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("logger").(*Logger)
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

	if err := os.MkdirAll(cfg.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		cfg.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

// With creates a new Logger with additional fields attached.
//
// The returned logger inherits all fields from the parent logger and adds
// the provided fields. It shares the same underlying log file.
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		file:   l.file,
	}
}

// Close closes the underlying log file.
//
// This method should be called when the logger is no longer needed,
// typically in a defer statement after creating the logger.
// If closing the file fails, an error message is printed to stdout
// since the logger may not be available for logging errors.
func (l *Logger) Close() error {
	if err := l.file.Close(); err != nil {
		//nolint:forbidigo // logger was not closed properly so logging error to stdout
		fmt.Printf("failed to close application logger: %v\n", err)
	}

	return nil
}
