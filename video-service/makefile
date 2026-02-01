ifneq (,$(wildcard ./.env))
    include .env
    export
endif

MIGRATIONS_DIR := $(shell yq e '.storage.migrations_dir' config/app.yml)

.PHONY: tests test mocks run sqlc lint lint-fix tools all swag db-up

# Optional arguments
PKG ?= app
TEST ?= .
MIG_FILE_NAME ?= init

# Tools
MOCKGEN := mockgen
SQLC := sqlc
GOLANGCI_LINT := golangci-lint
SWAG := swag.exe

# Paths
MAIN := cmd/app/main.go

define go_test
	go test ./internal/$(1) -run $(2)
endef
define create_db_migration
	goose -dir "$(MIGRATIONS_DIR)" -s create $(1) sql
endef

#########################################
############### Tools ###################
#########################################
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

win-tools:
	scoop install eza
	scoop install fd
	scoop install golangci-lint
	scoop install grep	
	scoop install ripgrep
	scoop install yq

lint:
	$(GOLANGCI_LINT) fmt
	$(GOLANGCI_LINT) run

fix:
	$(GOLANGCI_LINT) fmt
	$(GOLANGCI_LINT) run --fix

fmt:
	$(GOLANGCI_LINT) fmt

swag: 
	${SWAG} init -g cmd/app/main.go

sqlc:
	$(SQLC) generate

#########################################
######### Docker commands ###############
#########################################
docker-db-up:
	docker run -d --rm --name video-postgres18 -p 5432:5432 -e POSTGRES_USER=gato -e POSTGRES_PASSWORD=root postgres:18-alpine

#########################################
######### Database commands #############
#########################################
db-up: docker-db-up 
# 	$(MAKE) db-create
# 	$(MAKE) db-init

db-create:
	docker exec -i video-postgres18 createdb -U gato -h 127.0.0.1 -p 5432 videodb 

db-drop:
	docker exec -i video-postgres18 dropdb -U gato videodb

db-init:
	@echo "Migrations Dir: $(MIGRATIONS_DIR)"
	goose -dir "$(MIGRATIONS_DIR)" postgres "$(DATABASE_URL)" up

db-con:
	docker exec -i video-postgres18 psql -d videodb -U gato

### Migration commands
db-mig-create:
	$(call create_db_migration,$(MIG_FILE_NAME))

## Run service (sqlc -> go run)
run: sqlc swag
	go run $(MAIN)

## Default
all: tests

## Generate mocks
mocks:
	$(MOCKGEN) -source="./internal/app/video_service.go" -destination="./internal/app/mock/video_service_mock.go" -mock_names=VideoService=MockVideoService
	$(MOCKGEN) -source="./internal/ports/video_repo.go" -destination="./internal/ports/mock/video_repo_mock.go" -mock_names=VideoRepository=MockVideoRepository
	$(MOCKGEN) -source="./internal/ports/id_gen.go" -destination="./internal/ports/mock/id_gen_mock.go" -mock_names=IDGen=MockIDGen

## Update mocks + run tests
tests: mocks
	go test ./...

test:
	$(call go_test,$(PKG),$(TEST))
