package app

import (
	"context"
	"errors"
	"testing"
	"time"
	"video-provider/internal/pkg/shared"
	"video-provider/internal/user-service/domain"
	mock_ports "video-provider/internal/user-service/ports/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func TestUserService_Create(t *testing.T) {
	testCases := []struct {
		testName string
		id       uuid.UUID
		user     domain.User
		password string
		err      error
	}{
		{
			testName: "Valid registration",
			user: domain.User{
				Email:     "test@example.com",
				Name:      "TestName",
				LastName:  "TestLastname",
				CreatedAt: time.Now(),
			},
			password: "saAsd1231!!",
			err:      nil,
			id:       uuid.MustParse("d9fa522f-0006-464f-8d68-356ba1d6ad7d"),
		},
		{
			testName: "DB error encountered",
			user: domain.User{
				Email:     "test@example.com",
				Name:      "TestName",
				LastName:  "TestLastname",
				CreatedAt: time.Now(),
			},
			password: "saAsd1231!!",
			err:      errors.New("db error"),
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mock_ports.NewMockUserRepository(ctrl)
			mockHasher := mock_ports.NewMockPasswordHasher(ctrl)
			userService := NewUserService(mockRepo, mockHasher)

			mockHasher.EXPECT().Hash(tc.password).Return([]byte(tc.password), nil).Times(1)
			mockRepo.EXPECT().Create(gomock.Any(), tc.user, []byte(tc.password)).Return(tc.id, tc.err).Times(1)

			id, err := userService.Create(context.Background(), tc.user, tc.password)
			if tc.err != nil && err == nil {
				t.Error("Expected error but got none")
			} else if tc.err == nil && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if id != tc.id {
				t.Errorf("Expected ID %v, but got %v", tc.id, id)
			}
		})
	}
}

func TestUserService_Get(t *testing.T) {
	testCases := []struct {
		testName string
		id       uuid.UUID
		err      error
		errCode  shared.ErrorCode
		user     domain.User
	}{
		{
			testName: "Valid user retrieval",
			id:       uuid.MustParse("d9fa522f-0006-464f-8d68-356ba1d6ad7d"),
			err:      nil,
			user: domain.User{
				Name:      "TestName",
				Email:     "test@example.com",
				LastName:  "TestLastname",
				CreatedAt: time.Time{},
			},
		},
		{
			testName: "Invalid user retrieval",
			id:       uuid.MustParse("d9fa522f-0006-464f-8d68-356ba1d6ad7d"),
			err:      errors.New("db error"),
			errCode:  shared.ErrInternal,
			user:     domain.User{},
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mock_ports.NewMockUserRepository(ctrl)
			mockHasher := mock_ports.NewMockPasswordHasher(ctrl)
			userService := NewUserService(mockRepo, mockHasher)

			mockRepo.EXPECT().FindByID(gomock.Any(), tc.id).Return(tc.user, tc.err).Times(1)

			user, err := userService.Get(context.Background(), tc.id)
			if tc.err != nil && err == nil {
				t.Error("Expected error but got none")
			} else if tc.err == nil && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if user.Name != tc.user.Name ||
				user.Email != tc.user.Email ||
				user.LastName != tc.user.LastName ||
				user.CreatedAt != tc.user.CreatedAt {
				t.Errorf("Expected user %+v, but got %+v", tc.user, user)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	expectedToken := "success"
	expectedHash := []byte("P@ssword123")
	expectedNoMatchErr := errors.New("failed to compare password")

	testCases := []struct {
		testName   string
		inputEmail string
		inputPas   []byte
		repoCalls  int
		repoErr    error
		hashCalls  int
		hashErr    error
		resErr     error
		resToken   string
	}{
		{
			testName:   "Valid login",
			inputEmail: "test@example.com",
			inputPas:   expectedHash,
			repoCalls:  1,
			hashCalls:  1,
			resToken:   expectedToken,
		},
		{
			testName: "Empty email",
			resErr:   shared.NewError(shared.ErrInvalidInput, "email is required", nil),
		},
		{
			testName:   "Empty password",
			inputEmail: "test@example.com",
			resErr:     shared.NewError(shared.ErrInvalidInput, "password is required", nil),
		},
		{
			testName:   "Db error (email not found, etc)",
			inputEmail: "unexpectedEmail@domain.com",
			inputPas:   []byte("unexpectedPassword"),
			repoCalls:  1,
			repoErr:    pgx.ErrNoRows,
		},
		{
			testName:   "Wrong password",
			inputEmail: "test@example.com",
			inputPas:   []byte("unexpectedPassword"),
			repoCalls:  1,
			hashCalls:  1,
			hashErr:    expectedNoMatchErr,
			resErr:     shared.NewError(shared.ErrUnauthorized, "failed to compare password", nil),
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mock_ports.NewMockUserRepository(ctrl)
			mockHasher := mock_ports.NewMockPasswordHasher(ctrl)
			userService := NewUserService(mockRepo, mockHasher)

			mockRepo.EXPECT().GetPasswordHash(gomock.Any(), tc.inputEmail).Return(uuid.Nil, expectedHash, tc.repoErr).Times(tc.repoCalls)
			mockHasher.EXPECT().CompareHashAndPassword(expectedHash, tc.inputPas).Return(tc.hashErr).Times(tc.hashCalls)

			token, err := userService.Login(context.Background(), tc.inputEmail, tc.inputPas)
			if err == nil {
				if tc.repoErr == nil && tc.hashErr == nil && tc.resErr == nil && token != tc.resToken {
					t.Errorf("Expected token '%s', but got '%s'", tc.resToken, token)
				}
			} else {
				if tc.resErr != nil && err.Error() != tc.resErr.Error() {
					t.Errorf("Expected result rror '%s', but got '%s'", tc.resErr.Error(), err.Error())
				}
				if tc.hashErr != nil && err.Error() != tc.hashErr.Error() {
					t.Errorf("Expected hasher error'%s', but got '%s'", tc.hashErr.Error(), err.Error())
				}
				if tc.repoErr != nil && err.Error() != tc.repoErr.Error() {
					t.Errorf("Expected repo error '%s', but got '%s'", tc.repoErr.Error(), err.Error())
				}
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	testCases := []struct {
		testName string
		id       uuid.UUID
		user     domain.User
		err      error
	}{
		{
			testName: "Successful update",
			id:       uuid.MustParse("d9fa522f-0006-464f-8d68-356ba1d6ad7d"),
			user: domain.User{
				Name:     "UpdatedName",
				Email:    "updated@example.com",
				LastName: "UpdatedLastname",
			},
			err: nil,
		},
		{
			testName: "User not found",
			id:       uuid.MustParse("d9fa522f-0006-464f-8d68-356ba1d6ad7d"),
			user: domain.User{
				Name:     "UpdatedName",
				Email:    "updated@example.com",
				LastName: "UpdatedLastname",
			},
			err: shared.NewError(shared.ErrNotFound, "user not found", nil),
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mock_ports.NewMockUserRepository(ctrl)
			mockHasher := mock_ports.NewMockPasswordHasher(ctrl)
			userService := NewUserService(mockRepo, mockHasher)

			// Setup expectations based on test case
			if tc.err == nil {
				mockRepo.EXPECT().FindByID(gomock.Any(), tc.id).Return(domain.User{
					Name:     "OriginalName",
					Email:    "original@example.com",
					LastName: "OriginalLastname",
				}, nil)
				mockRepo.EXPECT().Update(gomock.Any(), tc.id, gomock.Any()).Return(nil)
			} else {
				mockRepo.EXPECT().FindByID(gomock.Any(), tc.id).Return(domain.User{}, shared.NewError(shared.ErrNotFound, "user not found", nil))
			}

			err := userService.Update(context.Background(), tc.id, tc.user)

			if tc.err != nil {
				if err == nil {
					t.Errorf("Expected error %v, but got nil", tc.err)
				} else if err.Error() != tc.err.Error() {
					t.Errorf("Expected error %v, but got %v", tc.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
