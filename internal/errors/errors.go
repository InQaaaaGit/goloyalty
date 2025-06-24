package errors

import (
	"errors"
	"fmt"
)

// Типизированные ошибки приложения
var (
	// Ошибки аутентификации
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")

	// Ошибки заказов
	ErrInvalidOrderNumberFormat          = errors.New("invalid order number format")
	ErrOrderAlreadyUploadedByAnotherUser = errors.New("order already uploaded by another user")
	ErrOrderAlreadyUploaded              = errors.New("order already uploaded by this user")

	// Ошибки баланса
	ErrInsufficientFunds = errors.New("insufficient funds")

	// Ошибки валидации
	ErrValidationFailed = errors.New("validation failed")

	// Ошибки базы данных
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrDatabaseQuery      = errors.New("database query failed")

	// Ошибки конфигурации
	ErrConfigLoad = errors.New("configuration load failed")

	// Ошибки миграций
	ErrMigrationFailed = errors.New("migration failed")

	// Ошибки HTTP
	ErrInvalidRequestFormat = errors.New("invalid request format")
	ErrMethodNotAllowed     = errors.New("method not allowed")
	ErrInternalServer       = errors.New("internal server error")
)

// AppError представляет ошибку приложения с дополнительным контекстом
type AppError struct {
	Code    string
	Message string
	Err     error
}

// Error возвращает сообщение об ошибке
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap возвращает вложенную ошибку
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError создает новую ошибку приложения
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Is проверяет, является ли ошибка определенным типом
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As проверяет, можно ли привести ошибку к определенному типу
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Wrap оборачивает ошибку с дополнительным сообщением
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf оборачивает ошибку с форматированным сообщением
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(format+": %w", append(args, err)...)
}
