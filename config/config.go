package config

import (
	"log"
	"mediary/internal"
	"os"
	"strconv"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv    string `env:"APP_ENV" envDefault:"development"`
	Port      string `env:"PORT" envDefault:"8080"`
	DBHost    string `env:"DB_HOST,required"`
	DBPort    int    `env:"DB_PORT" envDefault:"5432"`
	DBUser    string `env:"DB_USER,required"`
	DBPass    string `env:"DB_PASS,required"`
	DBName    string `env:"DB_NAME,required"`
	DBSSLMode string `env:"DB_SSLMODE" envDefault:"disable"`
	JWTSecret string `env:"JWT_SECRET,required"`
}

func Load(envFile string) (*Config, error) {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(envFile); err != nil {
			log.Println("No .env file found, continuing with OS envirornment variables")
		}
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, internal.NewConfigError("Error loading environment variables: %v", err)
	}

	validateConfig(&cfg)
	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.PortInt() <= 0 || cfg.PortInt() > 65535 {
		return internal.NewConfigError("PORT must be a valid number between 1-65535")
	}
	if cfg.DBPort <= 0 || cfg.DBPort > 65535 {
		return internal.NewConfigError("DB_PORT must be a valid number between 1-65535")
	}
	return nil
}

func (c *Config) PortInt() int {
	p, err := strconv.Atoi(c.Port)
	if err != nil {
		return 0
	}
	return p
}
