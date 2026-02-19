ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# default config if not passed (example: make run CONFIG=dev)
CONFIG ?= local
CONFIG_PATH := config/$(CONFIG)-service.yml

# ifeq (, $(shell which yq))
# 	$(error "Tool not found: 'yq'.")
# endif

ifeq ("$(wildcard $(CONFIG_PATH))","")
	$(error "Config not found: $(CONFIG_PATH)!")
endif

# func to call yq
get-cfg = $(shell yq e $(1) $(CONFIG_PATH))

# load everything
DB_NAME      	:= $(call get-cfg, '.db.name')
DB_VENDOR    	:= $(call get-cfg, '.db.vendor')
DB_VERSION   	:= $(call get-cfg, '.db.version')
DB_HOST      	:= $(call get-cfg, '.db.host')
DB_PORT      	:= $(call get-cfg, '.db.port')
DB_MAX_CONN 	:= $(call get-cfg, '.db.maxconns')
MIGRATIONS_DIR 	:= $(call get-cfg, '.db.migrationdir')

DB_CONTAINER_NAME := $(DB_NAME)-$(DB_VENDOR)-$(DB_VERSION)

MAIN = cmd/app.go
PKG ?= app
TEST ?= .
MIG_FILE_NAME ?= init

MOCKGEN := mockgen
SQLC := sqlc
GOLANGCI_LINT := golangci-lint
SWAG := swag.exe

#   --- Common Commands ---
.PHONY: run
run:
	@echo "Checking config: $(CONFIG_PATH)..."
	@echo "Starting service with DB: $(DB_NAME)"
	@echo "Container name will be: $(DB_CONTAINER_NAME)"
	go run cmd/$(CONFIG)-service/app.go -config=$(CONFIG_PATH)

.PHONY: generate
generate: sqlc swag

.PHONY: lint
lint:
	$(GOLANGCI_LINT) fmt
	$(GOLANGCI_LINT) run
	@echo "Formatted"

.PHONY: swag
swag: 
	${SWAG} init -g $(MAIN)
	@echo "Swagger docs generated"

.PHONY: sqlc
sqlc:
	$(SQLC) generate
	@echo "SQLC generated"


.PHONY: test
test:
	$(call go_test,$(PKG),$(TEST))

define go_test
	go test ./internal/$(1) -run $(2)
endef

.PHONY: tests
tests: mocks
	go test ./...

#   ---  Usage scanario --- 
# make db-up CONFIG=video
# make db-init CONFIG=video
# make run CONFIG=video

#   --- Docker ---
.PHONY: db-up 
db-up:
	docker run -d --rm --name $(DB_CONTAINER_NAME) -p $(DB_PORT):$(DB_PORT) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) postgres:$(DB_VERSION)
	@echo "Docker contained started with name $(DB_CONTAINER_NAME) on $(DB_PORT)"

.PHONY: db-init
db-init:
	docker exec -i $(DB_CONTAINER_NAME) createdb -U $(POSTGRES_USER) -h $(DB_HOST) -p $(DB_PORT) $(DB_NAME)
	@echo "CreateDB for $(DB_NAME)"
	$(MAKE) migrate-up

.PHONY: migrate-up
migrate-up:
	goose -dir "$(MIGRATIONS_DIR)" postgres "$(DB_VENDOR)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" up
	@echo "Migrate-up finished"

.PHONY: migrate-down
migrate-down:
	goose -dir "$(MIGRATIONS_DIR)" postgres "$(DB_VENDOR)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" down
	@echo "Migrate-down finished"

.PHONY: db-status
db-status:
	@echo "--- Configuration: $(CONFIG)-service ---"
	@echo "Target DB: $(DB_NAME) on $(DB_HOST):$(DB_PORT)"
	@echo "Container: $(DB_CONTAINER_NAME)"
	@echo "Migrations: $(MIGRATIONS_DIR)"

# 	--- Tools ---
.PHONY: tools  req-win-tools opt-win-tools
tools:
	irm get.scoop.sh | iex
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

req-win-tools:
	scoop install pwsh
	scoop install yq

opt-win-tools:
	scoop install fd
	scoop install ripgrep

