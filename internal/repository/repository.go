package repository

import (
	"context"
	"gophermart/internal/models"
)

type Repository interface {
	// Пользователи
	CreateUser(ctx context.Context, login, hashedPassword string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)

	// Заказы
	CreateOrder(ctx context.Context, userID int64, number string) (*models.Order, error)
	GetOrderByNumber(ctx context.Context, number string) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string, accrual *float64) error

	// Баланс
	GetUserBalance(ctx context.Context, userID int64) (*models.Balance, error)
	UpdateUserBalance(ctx context.Context, userID int64, current, withdrawn float64) error

	// Списания
	CreateWithdrawal(ctx context.Context, userID int64, order string, sum float64) (*models.Withdrawal, error)
	GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]models.Withdrawal, error)

	// Инициализация БД
	InitDB(ctx context.Context) error
}
