package app

import (
	"log"
	"testing"
	"time"
	"video-provider/internal/user-service/domain"
	mock_ports "video-provider/internal/user-service/ports/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockRepo := mock_ports.NewMockUserRepository(ctrl)
	mockHasher := mock_ports.NewMockPasswordHasher(ctrl)

	// Initialize UserService with the mocks
	userService := &UserService{Repo: mockRepo, log: *log.Default()}

	// Define test cases
	testCases := []struct {
		name       string
		cmd        RegisterUserCommand
		expectErr  bool
		expectedID uuid.UUID
	}{
		{
			name: "Valid registration",
			cmd: RegisterUserCommand{
				Email:    "test@example.com",
				Name:     "TestName",
				Lastname: "TestLastname",
				Password: "validPass",
			},
			expectErr:  false,
			expectedID: uuid.MustParse("d9fa522f-0006-464f-8d68-356ba1d6ad7d"),
		},
		{
			name: "Invalid password",
			cmd: RegisterUserCommand{
				Email:    "test@example.com",
				Name:     "TestName",
				Lastname: "TestLastname",
				Password: "invalidPass",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mocking the hasher to always return a valid hash for simplicity
			mockHasher.EXPECT().Hash(tc.cmd.Password).Return([]byte("hashedPass"), nil)

			// Mocking the repository create method
			expectedUser := &domain.User{Email: tc.cmd.Email, Name: tc.cmd.Name, LastName: tc.cmd.Lastname}
			mockRepo.EXPECT().Create(expectedUser, "hashedPass", "hashedPass").Return(tc.expectedID, nil)

			// Call the Register method and check for errors
			id, err := userService.Create(tc.cmd)
			if tc.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if id != tc.expectedID {
					t.Errorf("Expected ID %v, but got %v", tc.expectedID, id)
				}
			}
		})
	}
}

func TestUserService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock instances
	mockRepo := mock_ports.NewMockUserRepository(ctrl)
	userService := &UserService{Repo: mockRepo, log: *log.Default()}

	// Define test cases
	testCases := []struct {
		name         string
		id           uuid.UUID
		expectErr    bool
		expectedUser GetUserResult
	}{
		{
			name:      "Valid user retrieval",
			id:        uuid.MustParse("00000000-0000-0000-0000-000000000000"), // Replace with actual UUID if known
			expectErr: false,
			expectedUser: GetUserResult{
				Name:     "TestName",
				Email:    "test@example.com",
				Lastname: "TestLastname",
				CreateAt: time.Time{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mocking the repository find by ID method
			expectedUser := &domain.User{
				ID:        tc.id,
				Email:     "test@example.com",
				Name:      "TestName",
				LastName:  "TestLastname",
				CreatedAt: time.Time{},
			}
			mockRepo.EXPECT().FindByID(tc.id).Return(expectedUser, nil)

			// Call the Get method and check for errors
			user, err := userService.Get(tc.id)
			if tc.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if user.Name != tc.expectedUser.Name ||
					user.Email != tc.expectedUser.Email ||
					user.Lastname != tc.expectedUser.Lastname ||
					user.CreateAt != tc.expectedUser.CreateAt {
					t.Errorf("Expected user %+v, but got %+v", tc.expectedUser, user)
				}
			}
		})
	}
}
