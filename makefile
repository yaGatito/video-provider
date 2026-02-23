ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# default config if not passed (example: make run CONFIG=dev)
CONFIG ?= local
SERVICE_NAME = $(CONFIG)-service
CONFIG_PATH := config/$(SERVICE_NAME).yml


# ifeq ("$(wildcard $(CONFIG_PATH))","")
# 	$(error "Config not found: $(CONFIG_PATH)!")
# endif

# funcs
get-cfg = $(shell yq e $(1) $(CONFIG_PATH))
log = @echo [::MAKEFILE::] $(1)

### BE AWARE BEFORE ANY CHANGES HERE
DB_NAME      	:= $(call get-cfg, '.db.name')
DB_VENDOR    	:= $(call get-cfg, '.db.vendor')
DB_VERSION   	:= $(call get-cfg, '.db.version')
DB_HOST      	:= $(call get-cfg, '.db.host')
DB_PORT      	:= $(call get-cfg, '.db.port')
# DB_MAX_CONN 	:= $(call get-cfg, '.db.maxconns')
MIGRATIONS_DIR 	:= $(call get-cfg, '.db.migrationdir')
DB_URL := $(DB_VENDOR)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

DB_CONTAINER_NAME := $(DB_NAME)-$(DB_VENDOR)-$(DB_VERSION)
MAIN = cmd/$(SERVICE_NAME)/app.go

PKG ?= app
TEST ?= .

MOCKGEN := mockgen
SQLC := sqlc
GOLANGCI_LINT := golangci-lint
SWAG := swag.exe

#   --- Common Commands ---
.PHONY: setup
setup:
	$(call log, "Starting user database...")
	$(MAKE) db-up CONFIG=user
	timeout /t 5 /nobreak >nul
	$(call log, "Initializing user database...")
	$(MAKE) db-init CONFIG=user
	
	$(call log, "Starting video database...")
	$(MAKE) db-up CONFIG=video
	timeout /t 5 /nobreak >nul
	$(call log, "Initializing video database...")
	$(MAKE) db-init CONFIG=video
	
	$(call log, "Running migrations for user database...")
	$(MAKE) migrate-up CONFIG=user
	
	$(call log, "Running migrations for video database...")
	$(MAKE) migrate-up CONFIG=video
	
# 	$(call log, "Running user-service...")
# 	$(MAKE) run CONFIG=user

# 	$(call log, "Running video-service...")
# 	$(MAKE) run CONFIG=video

.PHONY: front
front:
	$(call log, "Starting frontend application...")
	cd ./web/watch-ua
	npm start

.PHONY: run
run:
	$(call log, "Checking config: $(CONFIG_PATH)...")
	go run cmd/$(SERVICE_NAME)/app.go -config=$(CONFIG_PATH)

.PHONY: generate
generate: sqlc swag mocks

.PHONY: lint
lint:
	$(GOLANGCI_LINT) fmt
	$(GOLANGCI_LINT) run
	$(call log, "Formatted")

.PHONY: swag
swag: 
	$(call log, "Swagger generate: $(MAIN)")
	$(call log, "Swagger output: docs")
	${SWAG} init -g $(MAIN) -o docs

.PHONY: sqlc
sqlc:
	$(call log, "SQLC generate by file: internal/$(SERVICE_NAME)/adapters/postgres/sqlc.yml")
	$(SQLC) generate -f "internal/$(SERVICE_NAME)/adapters/postgres/sqlc.yml"

.PHONY: mocks
mocks:
ifeq ("$(CONFIG)","video")
	$(MOCKGEN) -source="./internal/$(SERVICE_NAME)/app/video_service.go" -destination="./internal/$(SERVICE_NAME)/app/mock/video_service_mock.go" -mock_names=VideoService=MockVideoService
	$(MOCKGEN) -source="./internal/$(SERVICE_NAME)/ports/video_repo.go" -destination="./internal/$(SERVICE_NAME)/ports/mock/video_repo_mock.go" -mock_names=VideoRepository=MockVideoRepository
	$(MOCKGEN) -source="./internal/$(SERVICE_NAME)/ports/id_gen.go" -destination="./internal/$(SERVICE_NAME)/ports/mock/id_gen_mock.go" -mock_names=IDGen=MockIDGen
	$(call log, "$(SERVICE_NAME) mocks generated")
endif

ifeq ("$(CONFIG)","user")
	$(MOCKGEN) -source="./internal/$(SERVICE_NAME)/ports/user_repo.go" -destination="./internal/$(SERVICE_NAME)/ports/mock/user_repo_mock.go" -mock_names=UserRepository=MockUserRepository
	$(call log, "$(SERVICE_NAME) mocks generated")
endif

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
	docker run --rm -d --name $(DB_CONTAINER_NAME) -p $(DB_PORT):$(DB_PORT) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) postgres:$(DB_VERSION) -p $(DB_PORT)
	$(call log, "Docker contained started with name $(DB_CONTAINER_NAME) on $(DB_PORT)")

.PHONY: db-init
db-init:
	docker exec -it $(DB_CONTAINER_NAME) createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(DB_NAME) -p $(DB_PORT)
	$(call log, "CreateDB for $(DB_NAME)")
# 	$(MAKE) migrate-up

#   --- Database migrations ---
.PHONY: migrate-up
migrate-up:
	$(call log, "goose -dir "$(MIGRATIONS_DIR)" postgres "$(DB_URL)" up")
	goose -dir "$(MIGRATIONS_DIR)" postgres "$(DB_URL)" up
	$(call log, "Migrate-up finished")

.PHONY: migrate-down
migrate-down:
	goose -dir "$(MIGRATIONS_DIR)" postgres "$(DB_URL)" down
	$(call log, "Migrate-down finished")

.PHONY: migrate-init
migrate-init:	
	goose -dir "$(MIGRATIONS_DIR)" -s create init sql
	$(call log, "Migrate-init finished")

.PHONY: db-status
db-status:
	$(call log, "Configuration: $(SERVICE_NAME) ---")
	$(call log, "Target DB: $(DB_NAME) on $(DB_HOST):$(DB_PORT)")
	$(call log, "Container: $(DB_CONTAINER_NAME)")
	$(call log, "Migrations: $(MIGRATIONS_DIR)")

# 	--- Tools --- 
.PHONY: go-win-tools req-win-tools opt-win-tools
go-win-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

req-win-tools:
	irm get.scoop.sh | iex
	scoop install pwsh
	scoop install yq

opt-win-tools:
	scoop install fd
	scoop install ripgrep
