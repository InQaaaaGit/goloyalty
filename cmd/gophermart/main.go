package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gophermart/internal/config"
	"gophermart/internal/handlers"
	"gophermart/internal/middleware"
	"gophermart/internal/migrations"
	"gophermart/internal/repository"
	"gophermart/internal/service"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
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

	// Подключение к базе данных
	db, err := sql.Open("postgres", cfg.DatabaseURI)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Проверка подключения к БД
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Запуск миграций
	migrator := migrations.New(db)
	if err := migrator.Run(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Инициализация репозитория
	repo := repository.New(db)

	// Инициализация сервиса
	svc := service.New(repo, cfg.AccrualSystemAddress)

	// Инициализация обработчиков
	h := handlers.New(svc)

	// Создание роутера
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Compress)

	// API маршруты
	r.Route("/api", func(r chi.Router) {
		// Публичные маршруты
		r.Post("/user/register", h.Register)
		r.Post("/user/login", h.Login)

		// Защищенные маршруты
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Post("/user/orders", h.UploadOrder)
			r.Get("/user/orders", h.GetOrders)
			r.Get("/user/balance", h.GetBalance)
			r.Post("/user/balance/withdraw", h.Withdraw)
			r.Get("/user/withdrawals", h.GetWithdrawals)
		})
	})

	// Создание HTTP сервера
	srv := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: r,
	}

	// Канал для graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		log.Printf("Server starting on %s", cfg.RunAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-done
	log.Println("Server shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
	return nil
}
