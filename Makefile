# EventFlowEngine — Makefile для разработки Go-пакета

.PHONY: all build test vet fmt clean coverage

# Переменные
GO ?= go
PACKAGE = ./engine/...

# Цель по умолчанию
all: fmt vet build test

# Сборка пакета
build:
	$(GO) build $(PACKAGE)

# Запуск тестов
test:
	$(GO) test $(PACKAGE) -v -count=1

# Запуск тестов с покрытием
coverage:
	$(GO) test $(PACKAGE) -v -count=1 -coverprofile=coverage.out
	$(GO) tool cover -func=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html

# Проверка форматирования
fmt:
	$(GO) fmt $(PACKAGE)

# Статический анализ
vet:
	$(GO) vet $(PACKAGE)

# Полная проверка
check: fmt vet build test

# Очистка
clean:
	$(GO) clean
	rm -f coverage.out coverage.html

# Зависимости
tidy:
	$(GO) mod tidy

# Просмотр зависимостей
graph:
	$(GO) mod graph
