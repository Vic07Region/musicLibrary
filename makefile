PROJECT_NAME := $(shell go list -m)
BASE_DIR:=$(CURDIR)
LOCAL_BIN:=$(CURDIR)/bin
GO_VERSION := 1.20.1

check-go-version:
	@if ! go version | awk '{print $$3}' | awk -F. '{if ($$1 < 1 || ($$1 == 1 && $$2 < 20) || ($$1 == 1 && $$2 == 20 && $$3 < 1)) exit 1; exit 0}'; then \
		echo "Требуется Go версии $(GO_VERSION) или выше. Установите его с https://go.dev/dl"; \
		exit 1; \
	fi

# установка зависимостей
install-deps: check-go-version
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest

# получение зависимостей
get-deps:
	go mod tidy
	go get -u github.com/swaggo/swag/cmd/swag

# Сборка проекта
build: check-go-version
	go build ./cmd/main.go

# Запуск проекта
run: build
	./main

# Генерация документации swag
swag-docs:
	swag init -g ./cmd/main.go -o docs

# Удаление артефактов сборки
clean:
	rm -rf $(LOCAL_BIN)

# Помощь: выводит список доступных команд
help:
	@echo "Доступные команды:"
	@echo "  install-deps   Устанавливает зависимости проекта"
	@echo "  get-deps       Загружает зависимости проекта"
	@echo "  build          Собирает проект"
	@echo "  run           Запускает проект"
	@echo "  clean          Очищает сгенерированные файлы"
	@echo "  swag-docs      Генерирует документацию swag"

.default: help