package models

import (
	"time"
)

// User представляет пользователя системы
type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"-"` // Пароль не возвращается в JSON
}

// Order представляет заказ пользователя
type Order struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// Balance представляет баланс пользователя
type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// Withdrawal представляет списание средств
type Withdrawal struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

// RegisterRequest запрос на регистрацию
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginRequest запрос на вход
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// WithdrawRequest запрос на списание средств
type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

// AccrualResponse ответ от системы начислений
type AccrualResponse struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

// Статусы заказов
const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

// Статусы системы начислений
const (
	AccrualStatusRegistered = "REGISTERED"
	AccrualStatusInvalid    = "INVALID"
	AccrualStatusProcessing = "PROCESSING"
	AccrualStatusProcessed  = "PROCESSED"
)
