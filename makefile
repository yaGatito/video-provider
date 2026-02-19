ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# API_NAME ?= $(shell yq e '.api.name' config/app.yml)
# API_HOST ?= $(shell yq e '.api.host' config/app.yml)	
# API_PORT ?= $(shell yq e '.api.port' config/app.yml)

DB_NAME := $(shell yq e '.db.name' config/app.yml)
DB_HOST := $(shell yq e '.db.host' config/app.yml)
DB_PORT := $(shell yq e '.db.port' config/app.yml)
DB_MAX_CONNS := $(shell yq e '.db.maxconns' config/app.yml)
DB_VERSION := $(shell yq e '.db.version' config/app.yml)
DB_CONTAINER_NAME := $(DB_NAME)-$(DB_VENDOR)-$(DB_VERSION)
DB_VENDOR := $(shell yq e '.db.vendor' config/app.yml)
MIGRATIONS_DIR := $(shell yq e '.db.migrationdir' config/app.yml)

DATABASE_URL ?= "$(DB_VENDOR)://$(POSTGRES_USER):$(POSTGRES_PWD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable&pool_max_conns=$(DB_MAX_CONNS)&pool_max_conn_lifetime=1h30m"

MAIN = cmd/app.go
PKG ?= app
TEST ?= .
MIG_FILE_NAME ?= init

MOCKGEN := mockgen
SQLC := sqlc
GOLANGCI_LINT := golangci-lint
SWAG := swag.exe

.PHONY: tools run test tests swag sqlc db-up db-drop db-init lint # db-mig-create

#   --- Common Commands ---
run: lint
	go run $(MAIN)

generate: sqlc swag

lint:
	$(GOLANGCI_LINT) fmt
	$(GOLANGCI_LINT) run

swag: 
	${SWAG} init -g $(MAIN)

sqlc:
	$(SQLC) generate


test:
	$(call go_test,$(PKG),$(TEST))

define go_test
	go test ./internal/$(1) -run $(2)
endef

tests: mocks
	go test ./...

#   --- Database ---
db-up:
	@echo calling  db_up
	$(call db_up)
	@echo called
# 	docker run -d --rm --name $(DB_CONTAINER_NAME) -p $(DB_PORT):$(DB_PORT) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PWD) postgres:$(DB_VERSION)
# 	docker exec -i $(DB_CONTAINER_NAME) createdb -U $(POSTGRES_USER) -h $(DB_HOST) -p $(DB_PORT) $(DB_NAME)
# 	$(MAKE) db-init

define db_up
	docker run -d --rm --name $(DB_CONTAINER_NAME) -p $(DB_PORT):$(DB_PORT) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PWD) postgres:$(DB_VERSION)
	docker exec -i $(DB_CONTAINER_NAME) createdb -U $(POSTGRES_USER) -h $(DB_HOST) -p $(DB_PORT) $(DB_NAME)
	goose -dir "$(MIGRATIONS_DIR)" postgres "$(DATABASE_URL)" up
endef


db-init:
	goose -dir "$(MIGRATIONS_DIR)" postgres "$(DATABASE_URL)" up

# db-mig-create:
# 	$(call create_db_migration,$(MIG_FILE_NAME))

# define create_db_migration
# 	goose -dir "$(MIGRATIONS_DIR)" -s create $(1) sql
# endef

db-drop:
	docker exec -i $(DB_CONTAINER_NAME) dropdb -U $(POSTGRES_USER) $(DB_NAME)


# 	--- Tools ---
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

