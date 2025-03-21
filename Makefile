include .env

LOCAL_BIN:=$(CURDIR)/bin
REPO_MOCK_DIR:=internal/service/mocks
PKG:=github.com/sanchey92/jwt-example

build:
	@go build -o $(LOCAL_BIN)/app cmd/server/main.go
	@echo "Application built in $(LOCAL_BIN)/app"

run: build
	@$(LOCAL_BIN)/app

init-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@latest

clean:
	@rm -rf $(LOCAL_BIN)/*
	@echo "Binaries cleaned"

mock:
	@mkdir -p $(REPO_MOCK_DIR)
	@$(LOCAL_BIN)/mockgen -source=internal/service/auth.go -destination=$(REPO_MOCK_DIR)/user_repository.go -package=mocks
	@echo "Mocks generated in $(MOCK_DIR)"

clean-mocks:
	@rm -rf $(REPO_MOCK_DIR)
	@echo "Mocks cleaned"

test: mock
	@go test -v ./... -cover

test-race: mock
	@go test -v ./... -race -cover

test-short: mock
	@go test -v ./... -short

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v