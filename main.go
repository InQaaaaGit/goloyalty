package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gophermart/internal/config"
	"gophermart/internal/handlers"
	"gophermart/internal/middleware"
	"gophermart/internal/repository"
	"gophermart/internal/service"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к базе данных
	db, err := sql.Open("postgres", cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверка подключения к БД
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Инициализация репозитория
	repo := repository.New(db)

	// Инициализация базы данных
	ctx := context.Background()
	if err := repo.InitDB(ctx); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

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
			log.Fatalf("Failed to start server: %v", err)
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
}
