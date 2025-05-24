package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseName     string `env:"DB_NAME"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePort     string `env:"DB_PORT"`
	DatabasePassword string `env:"DB_PASSWORD"`
}

func (c *Config) DatabaseUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseName)
}

func New() (*Config, error) {
	var cfg Config
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}
	return &cfg, nil
}
