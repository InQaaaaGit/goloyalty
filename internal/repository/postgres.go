package repository

import (
	"context"
	"database/sql"
	"fmt"
	"gophermart/internal/models"
	"time"
)

type postgresRepo struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateUser(ctx context.Context, login, hashedPassword string) (*models.User, error) {
	var user models.User
	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id, login`
	err := r.db.QueryRowContext(ctx, query, login, hashedPassword).Scan(&user.ID, &user.Login)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &user, nil
}

func (r *postgresRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	query := `SELECT id, login, password_hash FROM users WHERE login = $1`
	err := r.db.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}
	return &user, nil
}

func (r *postgresRepo) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	query := `SELECT id, login, password_hash FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

func (r *postgresRepo) CreateOrder(ctx context.Context, userID int64, number string) (*models.Order, error) {
	var order models.Order
	query := `INSERT INTO orders (user_id, number, status, uploaded_at) VALUES ($1, $2, $3, $4) RETURNING id, user_id, number, status, uploaded_at`
	now := time.Now()
	err := r.db.QueryRowContext(ctx, query, userID, number, models.OrderStatusNew, now).Scan(
		&order.ID, &order.UserID, &order.Number, &order.Status, &order.UploadedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return &order, nil
}

func (r *postgresRepo) GetOrderByNumber(ctx context.Context, number string) (*models.Order, error) {
	var order models.Order
	query := `SELECT id, user_id, number, status, accrual, uploaded_at FROM orders WHERE number = $1`
	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&order.ID, &order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get order by number: %w", err)
	}
	return &order, nil
}

func (r *postgresRepo) GetOrdersByUserID(ctx context.Context, userID int64) ([]models.Order, error) {
	query := `SELECT id, user_id, number, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by user id: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *postgresRepo) UpdateOrderStatus(ctx context.Context, orderID int64, status string, accrual *float64) error {
	query := `UPDATE orders SET status = $1, accrual = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, accrual, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}

func (r *postgresRepo) GetUserBalance(ctx context.Context, userID int64) (*models.Balance, error) {
	var balance models.Balance
	query := `SELECT current_balance, withdrawn_balance FROM user_balances WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		if err == sql.ErrNoRows {
			// Создаем запись с нулевым балансом
			balance = models.Balance{Current: 0, Withdrawn: 0}
			return &balance, nil
		}
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}
	return &balance, nil
}

func (r *postgresRepo) UpdateUserBalance(ctx context.Context, userID int64, current, withdrawn float64) error {
	query := `INSERT INTO user_balances (user_id, current_balance, withdrawn_balance) VALUES ($1, $2, $3) 
			  ON CONFLICT (user_id) DO UPDATE SET current_balance = $2, withdrawn_balance = $3`
	_, err := r.db.ExecContext(ctx, query, userID, current, withdrawn)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %w", err)
	}
	return nil
}

func (r *postgresRepo) CreateWithdrawal(ctx context.Context, userID int64, order string, sum float64) (*models.Withdrawal, error) {
	var withdrawal models.Withdrawal
	query := `INSERT INTO withdrawals (user_id, order_number, sum, processed_at) VALUES ($1, $2, $3, $4) RETURNING id, user_id, order_number, sum, processed_at`
	now := time.Now()
	err := r.db.QueryRowContext(ctx, query, userID, order, sum, now).Scan(
		&withdrawal.ID, &withdrawal.UserID, &withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create withdrawal: %w", err)
	}
	return &withdrawal, nil
}

func (r *postgresRepo) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	query := `SELECT id, user_id, order_number, sum, processed_at FROM withdrawals WHERE user_id = $1 ORDER BY processed_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawals by user id: %w", err)
	}
	defer rows.Close()

	var withdrawals []models.Withdrawal
	for rows.Next() {
		var withdrawal models.Withdrawal
		err := rows.Scan(&withdrawal.ID, &withdrawal.UserID, &withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan withdrawal: %w", err)
		}
		withdrawals = append(withdrawals, withdrawal)
	}
	return withdrawals, nil
}

func (r *postgresRepo) InitDB(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			number VARCHAR(255) UNIQUE NOT NULL,
			status VARCHAR(50) NOT NULL,
			accrual DECIMAL(10,2),
			uploaded_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS user_balances (
			user_id INTEGER PRIMARY KEY REFERENCES users(id),
			current_balance DECIMAL(10,2) DEFAULT 0,
			withdrawn_balance DECIMAL(10,2) DEFAULT 0,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS withdrawals (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			order_number VARCHAR(255) NOT NULL,
			sum DECIMAL(10,2) NOT NULL,
			processed_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		_, err := r.db.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}
