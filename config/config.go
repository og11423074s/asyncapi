package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Env string

const (
	EnvDev  Env = "dev"
	EnvTest Env = "test"
)

type Config struct {
	ApiServerPort    string `env:"API_SERVER_PORT"`
	ApiServerHost    string `env:"API_SERVER_HOST"`
	DatabaseName     string `env:"DB_NAME"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePort     string `env:"DB_PORT"`
	DatabasePortTest string `env:"DB_PORT_TEST"`
	DatabasePassword string `env:"DB_PASSWORD"`
	Env              Env    `env:"ENV" envDefault:"dev"`
	ProjectRoot      string `env:"PROJECT_ROOT"`
}

func (c *Config) DatabaseUrl() string {

	port := c.DatabasePort

	if c.Env == EnvTest {
		port = c.DatabasePortTest
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		port,
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
