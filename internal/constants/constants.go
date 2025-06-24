package constants

import "time"

// HTTP статусы
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusAccepted            = 202
	StatusNoContent           = 204
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusPaymentRequired     = 402
	StatusConflict            = 409
	StatusUnprocessableEntity = 422
	StatusInternalServerError = 500
	StatusMethodNotAllowed    = 405
)

// Временные константы
const (
	DefaultJWTExpiration = 24 * time.Hour
	DefaultCookieMaxAge  = 86400 // 24 часа в секундах
)

// Сообщения об ошибках
const (
	ErrMsgMethodNotAllowed     = "Method not allowed"
	ErrMsgUnauthorized         = "Unauthorized"
	ErrMsgInvalidRequestFormat = "Invalid request format"
	ErrMsgInternalServerError  = "Internal server error"
	ErrMsgUserAlreadyExists    = "User already exists"
	ErrMsgInvalidCredentials   = "Invalid credentials"
	ErrMsgInvalidOrderNumber   = "Invalid order number format"
	ErrMsgOrderAlreadyUploaded = "Order already uploaded by another user"
	ErrMsgInsufficientFunds    = "Insufficient funds"
	ErrMsgOrderRequired        = "Order number is required"
	ErrMsgFailedToReadBody     = "Failed to read request body"
)

// HTTP заголовки
const (
	HeaderContentType     = "Content-Type"
	HeaderAuthorization   = "Authorization"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderContentEncoding = "Content-Encoding"
)

// MIME типы
const (
	MimeTypeJSON = "application/json"
	MimeTypeText = "text/plain"
)

// Префиксы
const (
	BearerPrefix = "Bearer "
)

// Настройки безопасности
const (
	MinPasswordLength    = 8
	MaxPasswordLength    = 72
	MinLoginLength       = 3
	MaxLoginLength       = 72
	MinOrderNumberLength = 2
)
