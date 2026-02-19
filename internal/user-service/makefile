CONFIG_NAME = userservice

include ../../makefile

## Generate mocks (Specific to user-service)
mocks:
	$(MOCKGEN) -source="./internal/ports/user_repo.go" -destination="./internal/ports/mock/user_repo_mock.go" -mock_names=UserRepository=MockUserRepository
