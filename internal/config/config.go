package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:":8080"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecret            string `env:"JWT_SECRET" envDefault:"default-secret-key-change-in-production"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Парсим флаги командной строки
	flag.StringVar(&cfg.RunAddress, "a", ":8080", "Server address")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "Database URI")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "Accrual system address")
	flag.Parse()

	// Парсим переменные окружения (имеют приоритет над флагами)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
