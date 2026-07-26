package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string `envconfig:"LEVEL" required:"true"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

func NewLoggerConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("LOGGER", &cfg); err != nil {
		return Config{}, fmt.Errorf("process logger envconfig: %w", err)
	}
	return cfg, nil
}

func NewLoggerConfigMust() Config {
	cfg, err := NewLoggerConfig()
	if err != nil {
		panic(fmt.Errorf("new logger config must: %w", err))
	}
	return cfg
}
