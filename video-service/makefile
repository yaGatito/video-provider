.PHONY: tests test mocks run sqlc all

# Optional arguments
PKG ?= app
TEST ?= .

# Tools
MOCKGEN := mockgen
SQLC := sqlc

# Paths
MOCKS_DIR := internal/mocks
MAIN := cmd/app/main.go

define go_test
	go test ./internal/$(1) -run $(2)
endef

## Update mocks + run tests
tests: mocks
	go test ./...

test:
	$(call go_test,$(PKG),$(TEST))

tools:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0

## Generate mocks
mocks:
	$(MOCKGEN) -source="./internal/app/video_service.go" -destination="./internal/app/mock/video_service_mock.go" -mock_names=VideoService=MockVideoService
	$(MOCKGEN) -source="./internal/ports/video_repo.go" -destination="./internal/ports/mock/video_repo_mock.go" -mock_names=VideoRepository=MockVideoRepository
	$(MOCKGEN) -source="./internal/ports/id_gen.go" -destination="./internal/ports/mock/id_gen_mock.go" -mock_names=IDGen=MockIDGen

## Generate sqlc code
sqlc:
	$(SQLC) generate

## Run service (sqlc -> go run)
run: sqlc
	go run $(MAIN)

## Default
all: tests