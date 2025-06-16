package service

import (
	"context"
	"fmt"
	"gophermart/internal/repository"
	"time"
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
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := p.processNewOrders(ctx); err != nil {
				// В реальной реализации здесь было бы логирование ошибки
				fmt.Printf("Error processing orders: %v\n", err)
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
