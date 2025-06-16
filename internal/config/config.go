package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	JWTSecret            string
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Флаги командной строки
	flag.StringVar(&cfg.RunAddress, "a", ":8080", "Server address")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "Database URI")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "Accrual system address")
	flag.Parse()

	// Переменные окружения (имеют приоритет над флагами)
	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		cfg.RunAddress = envRunAddress
	}
	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		cfg.DatabaseURI = envDatabaseURI
	}
	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		cfg.AccrualSystemAddress = envAccrualSystemAddress
	}

	// JWT секрет из переменной окружения
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "default-secret-key-change-in-production"
	}

	return cfg, nil
}
