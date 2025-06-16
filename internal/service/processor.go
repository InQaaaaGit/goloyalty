package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gophermart/internal/models"
	"gophermart/internal/repository"
)

type OrderProcessor struct {
	repo                 repository.Repository
	accrualSystemAddress string
}

func NewOrderProcessor(repo repository.Repository, accrualSystemAddress string) *OrderProcessor {
	return &OrderProcessor{
		repo:                 repo,
		accrualSystemAddress: accrualSystemAddress,
	}
}

// ProcessOrders запускает фоновую обработку заказов
func (p *OrderProcessor) ProcessOrders(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second) // Проверяем каждые 10 секунд
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := p.processNewOrders(ctx); err != nil {
				log.Printf("Error processing orders: %v", err)
			}
		}
	}
}

// processNewOrders обрабатывает новые заказы
func (p *OrderProcessor) processNewOrders(ctx context.Context) error {
	// В реальной реализации здесь была бы логика получения заказов со статусом NEW
	// Для демонстрации просто возвращаем nil
	return nil
}

// checkAccrualSystem проверяет заказ в системе начислений
func (p *OrderProcessor) checkAccrualSystem(orderNumber string) (*models.AccrualResponse, error) {
	if p.accrualSystemAddress == "" {
		return nil, fmt.Errorf("accrual system address not configured")
	}

	url := fmt.Sprintf("%s/api/orders/%s", p.accrualSystemAddress, orderNumber)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request accrual system: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var accrualResp models.AccrualResponse
		if err := json.NewDecoder(resp.Body).Decode(&accrualResp); err != nil {
			return nil, fmt.Errorf("failed to decode accrual response: %w", err)
		}
		return &accrualResp, nil
	case http.StatusNoContent:
		return nil, nil
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("rate limit exceeded")
	default:
		return nil, fmt.Errorf("accrual system error: %d", resp.StatusCode)
	}
}

// updateOrderStatus обновляет статус заказа
func (p *OrderProcessor) updateOrderStatus(ctx context.Context, orderID int64, status string, accrual *float64) error {
	return p.repo.UpdateOrderStatus(ctx, orderID, status, accrual)
}

// updateUserBalance обновляет баланс пользователя
func (p *OrderProcessor) updateUserBalance(ctx context.Context, userID int64, accrual float64) error {
	// Получаем текущий баланс
	balance, err := p.repo.GetUserBalance(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user balance: %w", err)
	}

	// Обновляем баланс
	newCurrent := balance.Current + accrual
	err = p.repo.UpdateUserBalance(ctx, userID, newCurrent, balance.Withdrawn)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %w", err)
	}

	return nil
}
