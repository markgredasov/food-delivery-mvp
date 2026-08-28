package webhook

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PerAttemptTimeout time.Duration `envconfig:"PER_ATTEMPT_TIMEOUT" required:"true"`
	MaxAttempts       int           `envconfig:"MAX_ATTEMPTS" required:"true"`
	Backoff           time.Duration `envconfig:"BACKOFF" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("WEBHOOK", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get webhook config: %w", err)
		panic(err)
	}

	return config
}
