include .env

LOCAL_BIN:=$(CURDIR)/bin
REPO_MOCK_DIR:=internal/service/mocks
PKG:=github.com/sanchey92/jwt-example

.PHONY: build
build:
	@go build -o $(LOCAL_BIN)/app cmd/server/main.go
	@echo "Application built in $(LOCAL_BIN)/app"

.PHONY: run
run: build
	@$(LOCAL_BIN)/app

.PHONY: init-deps
init-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@latest

.PHONY: clean
clean:
	@rm -rf $(LOCAL_BIN)/*
	@echo "Binaries cleaned"

.PHONY: mock
mock:
	@mkdir -p $(REPO_MOCK_DIR)
	@$(LOCAL_BIN)/mockgen -source=internal/service/auth.go -destination=$(REPO_MOCK_DIR)/repository_mock.go -package=mocks
	@echo "Mocks generated in $(MOCK_DIR)"

.PHONY: clean-mocks
clean-mocks:
	@rm -rf $(REPO_MOCK_DIR)
	@echo "Mocks cleaned"

.PHONY: test
test: mock
	@go test -v ./... -cover

.PHONY: test-race
test-race: mock
	@go test -v ./... -race -cover

.PHONY: test-short
test-short: mock
	@go test -v ./... -short

.PHONY: local-migration-status
local-migration-status:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

.PHONY: local-migration-up
local-migration-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

.PHONY: local-migration-down
local-migration-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v

.PHONY: help
help:
	@echo "Доступные команды:"
	@echo "  make build                - Собрать приложение"
	@echo "  make run                 - Собрать и запустить приложение"
	@echo "  make init-deps           - Установить зависимости (goose, mockgen)"
	@echo "  make clean               - Очистить бинарные файлы"
	@echo "  make mock                - Сгенерировать моки для репозиториев"
	@echo "  make clean-mocks         - Удалить сгенерированные моки"
	@echo "  make test                - Запустить тесты с покрытием"
	@echo "  make test-race           - Запустить тесты с проверкой гонки данных"
	@echo "  make test-short          - Запустить короткие тесты"
	@echo "  make local-migration-*   - Управление миграциями (status/up/down)"
	@echo "  make help                - Показать эту справку"