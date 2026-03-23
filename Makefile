.PHONY: build run clean install-deps test

# Сборка бинарника
build:
    go build -o bin/t-mess ./cmd/t-mess

# Запуск
run:
    go run ./cmd/t-mess

# Очистка
clean:
    rm -rf bin/
    rm -rf ~/.config/t-mess/

# Установка зависимостей
install-deps:
    go mod download
    go mod tidy

# Тесты
test:
    go test -v ./...