package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gophermart/internal/config"
	"gophermart/internal/errors"
	"gophermart/internal/models"
	"gophermart/internal/repository"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gophermart/internal/constants"

	"github.com/golang-jwt/jwt/v5"
)

type loyaltyService struct {
	repo                 repository.Repository
	accrualSystemAddress string
	config               *config.Config
}

func New(repo repository.Repository, accrualSystemAddress string, cfg *config.Config) Service {
	return &loyaltyService{
		repo:                 repo,
		accrualSystemAddress: accrualSystemAddress,
		config:               cfg,
	}
}

func (s *loyaltyService) Register(ctx context.Context, login, password string) (*models.User, string, error) {
	// Валидация логина
	if err := s.validateLogin(login); err != nil {
		return nil, "", err
	}

	// Валидация пароля
	if err := s.validatePassword(password); err != nil {
		return nil, "", err
	}

	// Проверяем, что пользователь не существует
	existingUser, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to check existing user")
	}
	if existingUser != nil {
		return nil, "", errors.ErrUserAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to hash password")
	}

	// Создаем пользователя
	user, err := s.repo.CreateUser(ctx, login, string(hashedPassword))
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to create user")
	}

	// Генерируем JWT токен
	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to generate token")
	}

	return user, token, nil
}

func (s *loyaltyService) Login(ctx context.Context, login, password string) (*models.User, string, error) {
	// Валидация логина
	if err := s.validateLogin(login); err != nil {
		return nil, "", err
	}

	// Валидация пароля
	if err := s.validatePassword(password); err != nil {
		return nil, "", err
	}

	// Получаем пользователя
	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to get user")
	}
	if user == nil {
		return nil, "", errors.ErrInvalidCredentials
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.ErrInvalidCredentials
	}

	// Генерируем JWT токен
	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to generate token")
	}

	return user, token, nil
}

func (s *loyaltyService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *loyaltyService) UploadOrder(ctx context.Context, userID int64, orderNumber string) (*models.Order, error) {
	// Проверяем формат номера заказа (алгоритм Луна)
	if !s.isValidOrderNumber(orderNumber) {
		return nil, errors.ErrInvalidOrderNumberFormat
	}

	// Проверяем, не загружен ли уже этот заказ
	existingOrder, err := s.repo.GetOrderByNumber(ctx, orderNumber)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check existing order")
	}

	if existingOrder != nil {
		if existingOrder.UserID == userID {
			return existingOrder, nil // Заказ уже загружен этим пользователем
		} else {
			return nil, errors.ErrOrderAlreadyUploadedByAnotherUser
		}
	}

	// Создаем новый заказ
	order, err := s.repo.CreateOrder(ctx, userID, orderNumber)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create order")
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
		return nil, errors.ErrInvalidOrderNumberFormat
	}

	// Получаем текущий баланс
	balance, err := s.repo.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user balance")
	}

	// Проверяем достаточность средств
	if balance.Current < sum {
		return nil, errors.ErrInsufficientFunds
	}

	// Создаем списание
	withdrawal, err := s.repo.CreateWithdrawal(ctx, userID, order, sum)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create withdrawal")
	}

	// Обновляем баланс
	newCurrent := balance.Current - sum
	newWithdrawn := balance.Withdrawn + sum
	err = s.repo.UpdateUserBalance(ctx, userID, newCurrent, newWithdrawn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update user balance")
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

	// Создаем JWT токен с настраиваемым временем жизни
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"jti":     tokenID,
		"exp":     time.Now().Add(s.config.JWTExpiration).Unix(),
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *loyaltyService) isValidOrderNumber(number string) bool {
	// Проверяем, что номер не пустой и имеет минимальную длину
	if len(number) < constants.MinOrderNumberLength {
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

func (s *loyaltyService) validatePassword(password string) error {
	// Проверяем минимальную длину пароля
	if len(password) < constants.MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", constants.MinPasswordLength)
	}

	// Проверяем максимальную длину пароля (bcrypt ограничение - 72 байта)
	if len(password) > constants.MaxPasswordLength {
		return fmt.Errorf("password must not exceed %d characters", constants.MaxPasswordLength)
	}

	// Проверяем, что пароль не пустой
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	return nil
}

func (s *loyaltyService) validateLogin(login string) error {
	// Проверяем минимальную длину логина
	if len(login) < constants.MinLoginLength {
		return fmt.Errorf("login must be at least %d characters long", constants.MinLoginLength)
	}

	// Проверяем максимальную длину логина
	if len(login) > constants.MaxLoginLength {
		return fmt.Errorf("login must not exceed %d characters", constants.MaxLoginLength)
	}

	// Проверяем, что логин не пустой
	if login == "" {
		return fmt.Errorf("login cannot be empty")
	}

	return nil
}
