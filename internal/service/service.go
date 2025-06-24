package service

import (
	"context"
	"gophermart/internal/models"
)

type Service interface {
	// Пользователи
	Register(ctx context.Context, login, password string) (*models.User, string, error)
	Login(ctx context.Context, login, password string) (*models.User, string, error)
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)

	// Заказы
	UploadOrder(ctx context.Context, userID int64, orderNumber string) (*models.Order, error)
	GetOrders(ctx context.Context, userID int64) ([]models.Order, error)

	// Баланс
	GetBalance(ctx context.Context, userID int64) (*models.Balance, error)
	Withdraw(ctx context.Context, userID int64, order string, sum float64) (*models.Withdrawal, error)
	GetWithdrawals(ctx context.Context, userID int64) ([]models.Withdrawal, error)

	// Фоновая обработка заказов
	ProcessOrders(ctx context.Context) error
}
