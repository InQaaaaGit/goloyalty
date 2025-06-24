package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run запускает HTTP сервер с graceful shutdown
func (app *App) Run() error {
	// Канал для graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		log.Printf("Server starting on %s", app.Config.RunAddress)
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-done
	log.Println("Server shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
	return nil
}

// RunWithGracefulShutdown запускает сервер с дополнительной обработкой graceful shutdown
func (app *App) RunWithGracefulShutdown() error {
	// Канал для graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		log.Printf("Server starting on %s", app.Config.RunAddress)
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-done
	log.Println("Server shutting down...")

	// Graceful shutdown сервера
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Graceful shutdown приложения
	if err := app.Shutdown(); err != nil {
		log.Printf("Application shutdown error: %v", err)
	}

	log.Println("Server and application exited")
	return nil
}
