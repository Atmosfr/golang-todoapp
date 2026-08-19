package core_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr            string        `envconfig:"ADDR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	AllowedOrigins  []string      `envconfig:"ALLOWED_ORIGINS" required:"true"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("HTTP", &cfg); err != nil {
		return Config{}, fmt.Errorf("new server config: %w", err)
	}
	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get http server config: %w", err)
		panic(err)
	}
	return cfg
}
