PROJECT_NAME:= $(shell go list -m)
BASE_DIR:=$(CURDIR)
LOCAL_BIN:=$(CURDIR)/bin

# установка зависимостей
install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest

# получение зависимостей
get-deps:
	go mod tidy
	go get -u github.com/swaggo/swag/cmd/swag

# Сборка проекта
build:
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
	@echo "  swag-docs      Генерирует документацию swagger"

.default: help