ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# default config if not passed (example: make run config=video)
config ?= video
SERVICE_NAME = $(config)-service

ifeq ($(OS),Windows_NT)
EXE := .exe
SLEEP_5 := powershell -NoProfile -Command "Start-Sleep -Seconds 5"
else
EXE :=
SLEEP_5 := sleep 5
endif

log = @echo [::MAKEFILE::] $(1)

# Use .env values directly
ifeq ($(config),user)
  DB_NAME	= $(USER_DB_NAME)
  DB_HOST   = $(USER_DB_HOST)
  DB_PORT   = $(USER_DB_PORT)
  API_PORT  = $(USER_API_PORT)
endif

ifeq ($(config),video)
  DB_NAME   = $(VIDEO_DB_NAME)
  DB_HOST   = $(VIDEO_DB_HOST)
  DB_PORT   = $(VIDEO_DB_PORT)
  API_PORT  = $(VIDEO_API_PORT)
endif

MIGRATIONS_DIR 		= internal/$(SERVICE_NAME)/adapters/postgres/sql/migrations
DEFAULT_DB_PORT 	= 5432
DB_VENDOR 			= postgres
DB_VERSION 			= 18-alpine
DB_URL 				= $(DB_VENDOR)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
DB_CONTAINER_NAME 	= $(config)-db-$(DB_VENDOR)-$(DB_VERSION)

MAIN 	= internal/$(SERVICE_NAME)/cmd/start.go
PKG 	?= app
TEST 	?= .

MOCKGEN 		:= mockgen$(EXE)
SQLC 			:= sqlc$(EXE)
GOLANGCI_LINT 	:= golangci-lint$(EXE)
SWAG 			:= swag$(EXE)

.DEFAULT_GOAL 	:= help

#   --- Common Commands ---
.PHONY: help
help:
	@echo "Usage: make <target> [CONFIG=user|video]"
	@echo ""
	@echo "Common targets:"
	@echo "  make bootstrap      	- check required local tools"
	@echo "  make setup          	- run DB containers + init + migrations for user and video"
	@echo "  make go-run config=... - run service (user or video)"
	@echo "  make do-run config=... - run service via docker (user or video)"
	@echo "  make image      	 	- build docker image (user or video)"
	@echo "  make web 			 	- run web application"
	@echo "  make gen            	- run generations scripts: sqlc, swagger, gqlgen, mocks"
	@echo "  make test           	- run selected package test"
	@echo "  make tests          	- run all tests"
	@echo "  make db-status      	- print DB settings for current config"

.PHONY: bootstrap
bootstrap:
	@go version
	@docker --version
	@$(MOCKGEN) -version
	@$(SQLC) version
	@$(GOLANGCI_LINT) version
	@goose -version

.PHONY: setup
setup:
	$(MAKE) gen config=video
	$(call log, "Generations for video service completed")

	$(MAKE) gen config=user
	$(call log, "Generations for user service completed")

	docker-compose up --build
	$(call log, "Docker compose started")

	$(SLEEP_5)
	$(call log, "Running migrations for user database...")
	$(MAKE) migrate-up config=user
	
	$(call log, "Running migrations for video database...")
	$(MAKE) migrate-up config=video

.PHONY: compose
compose:
	docker-compose up --build
	$(call log, "Docker compose started")

.PHONY: web
web:
	$(call log, "Starting frontend application...")
	cd ./web && npm start

.PHONY: go-run
go-run:
	go run $(MAIN)

.PHONY: gen
gen: sqlc swag mocks

.PHONY: lint
lint:
	cd ./internal/$(SERVICE_NAME) && $(GOLANGCI_LINT) run -c ../../.golangci.yml
	$(GOLANGCI_LINT) fmt
	$(call log, "Formatted")
	cd ../../

.PHONY: swag
swag: 
	$(call log, "Swagger generate: $(MAIN)")
	$(call log, "Swagger output: docs")
	${SWAG} init -g cmd/start.go -o internal/$(SERVICE_NAME)/docs  --dir internal/$(SERVICE_NAME)

.PHONY: sqlc
sqlc:
	$(call log, "SQLC generate by file: internal/$(SERVICE_NAME)/adapters/postgres/sqlc.yml")
	$(SQLC) generate -f "./internal/$(SERVICE_NAME)/adapters/postgres/sqlc.yml"

.PHONY: mocks
mocks:
ifeq ("$(config)","video")
	$(MOCKGEN) -source="./internal/$(SERVICE_NAME)/app/service.go" -destination="./internal/$(SERVICE_NAME)/app/mock/service_mock.go" -mock_names=VideoService=MockVideoService
	$(MOCKGEN) -source="./internal/$(SERVICE_NAME)/ports/video_repo.go" -destination="./internal/$(SERVICE_NAME)/ports/mock/video_repo_mock.go" -mock_names=VideoRepository=MockVideoRepository
	$(call log, "$(SERVICE_NAME) mocks generated")
endif

ifeq ("$(config)","user")
	$(MOCKGEN) -source="./internal/$(SERVICE_NAME)/app/service.go" -destination="./internal/$(SERVICE_NAME)/app/mock/service_mock.go" -mock_names=UserInteractor=MockUserInteractor
	$(MOCKGEN) -source="./internal/$(SERVICE_NAME)/ports/user_repo.go" -destination="./internal/$(SERVICE_NAME)/ports/mock/user_repo_mock.go" -mock_names=UserRepository=MockUserRepository
	$(MOCKGEN) -source="./internal/$(SERVICE_NAME)/ports/hash_gen.go" -destination="./internal/$(SERVICE_NAME)/ports/mock/hash_gen_mock.go" -mock_names=PasswordHasher=MockPasswordHasher
	$(call log, "$(SERVICE_NAME) mocks generated")
endif

.PHONY: coverage
coverage:
	cd internal/ && go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out && rm coverage.out && cd ../

.PHONY: tests
tests: mocks
	go test ./...

#   --- Docker ---

.PHONY: do-run
do-run:
	docker build -D -t $(SERVICE_NAME) -f internal/$(SERVICE_NAME)/Dockerfile .
	docker rm -f $(SERVICE_NAME)
# 	docker run --rm -p 8081:8081 --env-file .env $(SERVICE_NAME)
	docker run  --name $(SERVICE_NAME) --rm -p $(API_PORT):$(API_PORT) $(SERVICE_NAME)

# .PHONY: db-up
# db-up:
# 	docker rm -f $(DB_CONTAINER_NAME)
# 	docker run --rm -d --name $(DB_CONTAINER_NAME) -p $(DB_PORT):$(DB_PORT) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) postgres:$(DB_VERSION) -p $(DB_PORT)
# 	$(call log, "Docker container started with name $(DB_CONTAINER_NAME) on $(DB_PORT)")

# .PHONY: db-down
# db-down:
# 	docker rm -f $(DB_CONTAINER_NAME)
# 	$(call log, "Docker container removed: $(DB_CONTAINER_NAME)")

# .PHONY: db-init
# db-init:
# 	docker exec $(DB_CONTAINER_NAME) createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(DB_NAME) -p $(DEFAULT_DB_PORT)
# 	$(call log, "CreateDB for $(DB_NAME)")

#   --- Database migrations ---
.PHONY: migrate-up
migrate-up:
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
.PHONY: go-tools
go-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
