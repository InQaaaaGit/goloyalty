.PHONY: build run test clean deps lint

# Переменные
BINARY_NAME=gophermart
BUILD_DIR=build
MAIN_FILE=main.go

# Сборка приложения
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Запуск приложения
run:
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_FILE)

# Установка зависимостей
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Запуск тестов
test:
	@echo "Running tests..."
	go test -v ./...

# Запуск тестов с покрытием
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Проверка кода линтером
lint:
	@echo "Running linter..."
	golangci-lint run

# Очистка
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

# Создание Docker образа
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

# Запуск в Docker
docker-run:
	@echo "Running in Docker..."
	docker run -p 8080:8080 $(BINARY_NAME)

# Помощь
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  deps          - Install dependencies"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  lint          - Run linter"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run in Docker"
	@echo "  help          - Show this help" 