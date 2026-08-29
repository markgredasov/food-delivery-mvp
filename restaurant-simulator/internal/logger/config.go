// Package logger provides structured logging wrapper around uber-go/zap.
package logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config holds logger configuration parameters.
//
// Values are loaded from environment variables with the prefix "LOGGER".
type Config struct {
	// Level is the minimum log level (debug, info, warn, error, etc.).
	// Required. Loaded from LOGGER_LEVEL.
	Level string `envconfig:"LEVEL" required:"true"`

	// Folder is the directory path for log files.
	// Required. Loaded from LOGGER_FOLDER.
	Folder string `envconfig:"FOLDER" required:"true"`
}

// NewConfig loads logger configuration from environment variables.
//
// Expects LOGGER_LEVEL and LOGGER_FOLDER to be set.
// Returns error if variables are missing or cannot be processed.
func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

// NewConfigMust loads logger configuration and panics on error.
//
// Useful for mandatory configuration during application initialization.
// Panics if environment variables are missing or invalid.
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get logger config: %w", err)
		panic(err)
	}

	return config
}
