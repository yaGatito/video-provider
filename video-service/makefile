.PHONY: tests test mocks run sqlc lint lint-fix tools all  

# Optional arguments
PKG ?= app
TEST ?= .

# Tools
MOCKGEN := mockgen
SQLC := sqlc
GOLANGCI_LINT := golangci-lint

# Paths
MOCKS_DIR := internal/mocks
MAIN := cmd/app/main.go

define go_test
	go test ./internal/$(1) -run $(2)
endef

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0

lint:
	$(GOLANGCI_LINT) run
lint-fix:
	$(GOLANGCI_LINT) run --fix
fmt:
	$(GOLANGCI_LINT) help

## Generate sqlc code
sqlc:
	$(SQLC) generate

## Run service (sqlc -> go run)
run: sqlc
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

### endpoints
# curl.exe -X POST "http://localhost:8081/v1/videos/pub/d9fa522f-0006-464f-8d68-356ba1d6ad7d" -H "Content-Type: application/json" -d '{"topic":"huy sosal","description":"sadasd"}'
