# Система лояльности "Гофермарт"

HTTP API для системы накопительной лояльности интернет-магазина "Гофермарт".

## Описание

Система позволяет пользователям:
- Регистрироваться и аутентифицироваться
- Загружать номера заказов для получения баллов лояльности
- Просматривать историю заказов и начислений
- Проверять текущий баланс баллов
- Списывать баллы для оплаты новых заказов

## Архитектура

Проект построен по принципам Clean Architecture с разделением на слои:

- **Handlers** - HTTP обработчики
- **Service** - бизнес-логика
- **Repository** - работа с базой данных
- **Models** - модели данных
- **Middleware** - промежуточное ПО (аутентификация, сжатие)

## Технологии

- **Go 1.21+** - основной язык
- **PostgreSQL** - база данных
- **Chi** - HTTP роутер
- **JWT** - аутентификация
- **bcrypt** - хеширование паролей

## Установка и запуск

### Требования

- Go 1.21 или выше
- PostgreSQL 12 или выше

### Конфигурация

Система поддерживает конфигурацию через переменные окружения или флаги командной строки:

```bash
# Адрес запуска сервера
RUN_ADDRESS=:8080

# Подключение к базе данных
DATABASE_URI=postgres://user:password@localhost:5432/gophermart?sslmode=disable

# Адрес системы расчёта начислений
ACCRUAL_SYSTEM_ADDRESS=http://localhost:8081

# Секретный ключ для JWT (опционально)
JWT_SECRET=your-secret-key
```

### Запуск

```bash
# Клонирование репозитория
git clone https://github.com/InQaaaaGit/go_loyalty.git
cd go_loyalty

# Установка зависимостей
go mod download

# Запуск с переменными окружения
export DATABASE_URI="postgres://user:password@localhost:5432/gophermart?sslmode=disable"
export ACCRUAL_SYSTEM_ADDRESS="http://localhost:8081"
go run main.go

# Или с флагами командной строки
go run main.go -a :8080 -d "postgres://user:password@localhost:5432/gophermart?sslmode=disable" -r "http://localhost:8081"
```

## API Endpoints

### Регистрация пользователя
```
POST /api/user/register
Content-Type: application/json

{
    "login": "user@example.com",
    "password": "password123"
}
```

### Аутентификация
```
POST /api/user/login
Content-Type: application/json

{
    "login": "user@example.com",
    "password": "password123"
}
```

### Загрузка номера заказа
```
POST /api/user/orders
Content-Type: text/plain
Authorization: Bearer <token>

12345678903
```

### Получение списка заказов
```
GET /api/user/orders
Authorization: Bearer <token>
```

### Получение баланса
```
GET /api/user/balance
Authorization: Bearer <token>
```

### Списание средств
```
POST /api/user/balance/withdraw
Content-Type: application/json
Authorization: Bearer <token>

{
    "order": "2377225624",
    "sum": 751
}
```

### История списаний
```
GET /api/user/withdrawals
Authorization: Bearer <token>
```

## Структура базы данных

### Таблица users
- `id` - уникальный идентификатор пользователя
- `login` - логин пользователя (уникальный)
- `password_hash` - хеш пароля
- `created_at` - дата создания

### Таблица orders
- `id` - уникальный идентификатор заказа
- `user_id` - ID пользователя
- `number` - номер заказа (уникальный)
- `status` - статус обработки (NEW, PROCESSING, INVALID, PROCESSED)
- `accrual` - сумма начисления
- `uploaded_at` - дата загрузки

### Таблица user_balances
- `user_id` - ID пользователя (первичный ключ)
- `current_balance` - текущий баланс
- `withdrawn_balance` - сумма списанных средств

### Таблица withdrawals
- `id` - уникальный идентификатор списания
- `user_id` - ID пользователя
- `order_number` - номер заказа
- `sum` - сумма списания
- `processed_at` - дата обработки

## Особенности реализации

### Алгоритм Луна
Номера заказов проверяются с помощью алгоритма Луна для валидации контрольной суммы.

### JWT Аутентификация
Используется JWT токены для аутентификации пользователей с временем жизни 24 часа.

### Сжатие данных
Поддерживается сжатие HTTP ответов с помощью gzip.

### Фоновая обработка
Заказы обрабатываются в фоновом режиме с периодической проверкой системы начислений.

## Тестирование

```bash
# Запуск тестов
go test ./...

# Запуск тестов с покрытием
go test -cover ./...
```

## Лицензия

MIT License
