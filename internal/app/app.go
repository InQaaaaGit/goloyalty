package app

import (
	"database/sql"
	"log"
	"net/http"

	"gophermart/internal/config"
	"gophermart/internal/handlers"
	"gophermart/internal/middleware"
	"gophermart/internal/migrations"
	"gophermart/internal/repository"
	"gophermart/internal/service"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

// App представляет основное приложение
type App struct {
	Config   *config.Config
	DB       *sql.DB
	Repo     repository.Repository
	Service  service.Service
	Handlers *handlers.Handler
	Router   *chi.Mux
	Server   *http.Server
}

// New создает новый экземпляр приложения
func New(cfg *config.Config) (*App, error) {
	app := &App{
		Config: cfg,
	}

	// Инициализация базы данных
	if err := app.initDatabase(); err != nil {
		return nil, err
	}

	// Инициализация репозитория
	app.Repo = repository.New(app.DB)

	// Инициализация сервиса
	app.Service = service.New(app.Repo, cfg.AccrualSystemAddress)

	// Инициализация обработчиков
	app.Handlers = handlers.New(app.Service)

	// Инициализация роутера
	app.Router = app.initRouter()

	// Инициализация HTTP сервера
	app.Server = &http.Server{
		Addr:    cfg.RunAddress,
		Handler: app.Router,
	}

	return app, nil
}

// initDatabase инициализирует подключение к базе данных и запускает миграции
func (app *App) initDatabase() error {
	// Подключение к базе данных
	db, err := sql.Open("postgres", app.Config.DatabaseURI)
	if err != nil {
		return err
	}

	// Проверка подключения к БД
	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}

	app.DB = db

	// Запуск миграций
	migrator := migrations.New(db)
	if err := migrator.Run(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// initRouter инициализирует роутер с маршрутами
func (app *App) initRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Compress)

	// API маршруты
	r.Route("/api", func(r chi.Router) {
		// Публичные маршруты
		r.Post("/user/register", app.Handlers.Register)
		r.Post("/user/login", app.Handlers.Login)

		// Защищенные маршруты
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Post("/user/orders", app.Handlers.UploadOrder)
			r.Get("/user/orders", app.Handlers.GetOrders)
			r.Get("/user/balance", app.Handlers.GetBalance)
			r.Post("/user/balance/withdraw", app.Handlers.Withdraw)
			r.Get("/user/withdrawals", app.Handlers.GetWithdrawals)
		})
	})

	return r
}

// Shutdown корректно завершает работу приложения
func (app *App) Shutdown() error {
	if app.DB != nil {
		return app.DB.Close()
	}
	return nil
}
