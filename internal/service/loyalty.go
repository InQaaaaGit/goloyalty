package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gophermart/internal/models"
	"gophermart/internal/repository"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
)

type loyaltyService struct {
	repo                 repository.Repository
	accrualSystemAddress string
	jwtSecret            string
}

func New(repo repository.Repository, accrualSystemAddress string) Service {
	return &loyaltyService{
		repo:                 repo,
		accrualSystemAddress: accrualSystemAddress,
		jwtSecret:            "default-secret-key-change-in-production", // В реальном приложении загружается из конфига
	}
}

func (s *loyaltyService) Register(ctx context.Context, login, password string) (*models.User, string, error) {
	// Проверяем, что пользователь не существует
	existingUser, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, "", fmt.Errorf("user already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаем пользователя
	user, err := s.repo.CreateUser(ctx, login, string(hashedPassword))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Генерируем JWT токен
	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *loyaltyService) Login(ctx context.Context, login, password string) (*models.User, string, error) {
	// Получаем пользователя
	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Генерируем JWT токен
	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *loyaltyService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *loyaltyService) UploadOrder(ctx context.Context, userID int64, orderNumber string) (*models.Order, error) {
	// Проверяем формат номера заказа (алгоритм Луна)
	if !s.isValidOrderNumber(orderNumber) {
		return nil, fmt.Errorf("invalid order number format")
	}

	// Проверяем, не загружен ли уже этот заказ
	existingOrder, err := s.repo.GetOrderByNumber(ctx, orderNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing order: %w", err)
	}

	if existingOrder != nil {
		if existingOrder.UserID == userID {
			return existingOrder, nil // Заказ уже загружен этим пользователем
		} else {
			return nil, fmt.Errorf("order already uploaded by another user")
		}
	}

	// Создаем новый заказ
	order, err := s.repo.CreateOrder(ctx, userID, orderNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (s *loyaltyService) GetOrders(ctx context.Context, userID int64) ([]models.Order, error) {
	return s.repo.GetOrdersByUserID(ctx, userID)
}

func (s *loyaltyService) GetBalance(ctx context.Context, userID int64) (*models.Balance, error) {
	return s.repo.GetUserBalance(ctx, userID)
}

func (s *loyaltyService) Withdraw(ctx context.Context, userID int64, order string, sum float64) (*models.Withdrawal, error) {
	// Проверяем формат номера заказа
	if !s.isValidOrderNumber(order) {
		return nil, fmt.Errorf("invalid order number format")
	}

	// Получаем текущий баланс
	balance, err := s.repo.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	// Проверяем достаточность средств
	if balance.Current < sum {
		return nil, fmt.Errorf("insufficient funds")
	}

	// Создаем списание
	withdrawal, err := s.repo.CreateWithdrawal(ctx, userID, order, sum)
	if err != nil {
		return nil, fmt.Errorf("failed to create withdrawal: %w", err)
	}

	// Обновляем баланс
	newCurrent := balance.Current - sum
	newWithdrawn := balance.Withdrawn + sum
	err = s.repo.UpdateUserBalance(ctx, userID, newCurrent, newWithdrawn)
	if err != nil {
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	return withdrawal, nil
}

func (s *loyaltyService) GetWithdrawals(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	return s.repo.GetWithdrawalsByUserID(ctx, userID)
}

func (s *loyaltyService) ProcessOrders(ctx context.Context) error {
	// Получаем все заказы со статусом NEW
	// В реальной реализации здесь была бы логика получения заказов из БД
	// Для демонстрации просто возвращаем nil
	return nil
}

// Вспомогательные методы

func (s *loyaltyService) generateJWT(userID int64) (string, error) {
	// Генерируем случайный ID для токена
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	tokenID := hex.EncodeToString(randomBytes)

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"jti":     tokenID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *loyaltyService) isValidOrderNumber(number string) bool {
	// Проверяем, что номер не пустой и имеет минимальную длину
	if len(number) < 2 {
		return false
	}

	// Проверяем, что номер состоит только из цифр
	for _, char := range number {
		if char < '0' || char > '9' {
			return false
		}
	}

	// Алгоритм Луна
	sum := 0
	alternate := false

	// Проходим по цифрам справа налево
	for i := len(number) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(number[i]))

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

func (s *loyaltyService) checkAccrualSystem(orderNumber string) (*models.AccrualResponse, error) {
	if s.accrualSystemAddress == "" {
		return nil, fmt.Errorf("accrual system address not configured")
	}

	url := fmt.Sprintf("%s/api/orders/%s", s.accrualSystemAddress, orderNumber)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request accrual system: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var accrualResp models.AccrualResponse
		// В реальной реализации здесь был бы парсинг JSON ответа
		return &accrualResp, nil
	case http.StatusNoContent:
		return nil, nil
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("rate limit exceeded")
	default:
		return nil, fmt.Errorf("accrual system error: %d", resp.StatusCode)
	}
}
