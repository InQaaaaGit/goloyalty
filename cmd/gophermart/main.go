package main

import (
	"fmt"
	"log"
	"os"

	"gophermart/internal/app"
	"gophermart/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Application error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Создание и инициализация приложения
	application, err := app.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	// Запуск сервера с graceful shutdown
	return application.RunWithGracefulShutdown()
}
