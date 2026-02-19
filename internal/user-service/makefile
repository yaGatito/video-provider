SERVICE_NAME := user-service
DB_PORT := 5433
DB_NAME := userdb
CONTAINER_NAME := user-postgres18

include ../makefile

## Generate mocks (Specific to user-service)
mocks:
	$(MOCKGEN) -source="./internal/ports/user_repo.go" -destination="./internal/ports/mock/user_repo_mock.go" -mock_names=UserRepository=MockUserRepository
