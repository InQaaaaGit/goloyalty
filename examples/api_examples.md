# Примеры использования API GopherMart

## Регистрация и аутентификация

### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "login": "user@example.com",
    "password": "password123"
  }'
```

**Ответ:**
```json
{
  "status": "success"
}
```

### Вход в систему
```bash
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "user@example.com",
    "password": "password123"
  }'
```

**Ответ:**
```json
{
  "status": "success"
}
```

## Работа с заказами

### Загрузка номера заказа
```bash
curl -X POST http://localhost:8080/api/user/orders \
  -H "Content-Type: text/plain" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d "12345678903"
```

**Возможные ответы:**

Успешная загрузка:
```json
{
  "id": 1,
  "user_id": 1,
  "number": "12345678903",
  "status": "NEW",
  "accrual": 0,
  "uploaded_at": "2024-01-15T10:30:00Z"
}
```

Неправильный формат номера:
```json
{
  "error": "invalid order number format"
}
```

Заказ уже загружен другим пользователем:
```json
{
  "error": "order already uploaded by another user"
}
```

### Получение списка заказов
```bash
curl -X GET http://localhost:8080/api/user/orders \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Ответ:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "number": "12345678903",
    "status": "PROCESSED",
    "accrual": 500,
    "uploaded_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "user_id": 1,
    "number": "45678912301",
    "status": "NEW",
    "accrual": 0,
    "uploaded_at": "2024-01-15T11:00:00Z"
  }
]
```

## Работа с балансом

### Получение баланса
```bash
curl -X GET http://localhost:8080/api/user/balance \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Ответ:**
```json
{
  "current": 1500,
  "withdrawn": 500
}
```

### Списание средств
```bash
curl -X POST http://localhost:8080/api/user/balance/withdraw \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "order": "2377225624",
    "sum": 751
  }'
```

**Возможные ответы:**

Успешное списание:
```json
{
  "id": 1,
  "user_id": 1,
  "order": "2377225624",
  "sum": 751,
  "processed_at": "2024-01-15T12:00:00Z"
}
```

Недостаточно средств:
```json
{
  "error": "insufficient funds"
}
```

Неправильный формат номера заказа:
```json
{
  "error": "invalid order number format"
}
```

## Коды ошибок

| HTTP Status | Описание |
|-------------|----------|
| 200 | Успешный запрос |
| 201 | Ресурс создан |
| 202 | Запрос принят к обработке |
| 400 | Неверный формат запроса |
| 401 | Не авторизован |
| 402 | Недостаточно средств |
| 409 | Конфликт (заказ уже загружен) |
| 422 | Неверный формат номера заказа |
| 500 | Внутренняя ошибка сервера |

## Валидация номеров заказов

Система использует алгоритм Луна для проверки корректности номеров заказов:

- Номер должен содержать только цифры
- Длина номера должна быть не менее 2 символов
- Контрольная сумма должна быть корректной

**Примеры валидных номеров:**
- `12345678903`
- `45678912301`
- `2377225624`

**Примеры невалидных номеров:**
- `1234567890` (неверная контрольная сумма)
- `abc123` (содержит буквы)
- `1` (слишком короткий) 