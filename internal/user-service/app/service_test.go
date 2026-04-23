// internal/user-service/app/service_test.go
package app_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
	"video-provider/common/auth"
	"video-provider/common/shared"
	"video-provider/user-service/app"
	"video-provider/user-service/domain"
	mock_ports "video-provider/user-service/ports/mock"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
			authSvc := auth.NewAuth([]byte("test"))
			userService := app.NewUserService(mockRepo, mockHasher, authSvc)

			mockHasher.EXPECT().Hash(tc.password).Return([]byte(tc.password), nil).Times(1)
			mockRepo.EXPECT().
				Create(gomock.Any(), tc.user, []byte(tc.password)).
				Return(tc.id, tc.err).
				Times(1)

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
			authSvc := auth.NewAuth([]byte("test"))
			userService := app.NewUserService(mockRepo, mockHasher, authSvc)

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
	expectedUserID := uuid.MustParse("d9fa522f-0006-464f-8d68-356ba1d6ad7d")
	expectedHash := []byte("P@ssword123")
	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzU1NjQ5NTgsInVzZXJfaWQiOiI4MWNmMzNkNS03ZDIwLTQ1YzMtYjg3Mi1iMGNiYmMwNzI5NjEifQ.4HSvvKyphap68diVSgCUX862SLYgcaa_7Q6PQjQMJ28"
	jwtSecret := []byte("my-nigga-secret-key")

	testCases := []struct {
		testName       string
		inputEmail     string
		inputPas       []byte
		repoCalls      int
		repoErr        error
		hashCalls      int
		hashErr        error
		resErr         error
		resToken       string
		expectedUserID uuid.UUID
	}{
		{
			testName:       "Valid login",
			inputEmail:     "test@example.com",
			inputPas:       expectedHash,
			repoCalls:      1,
			repoErr:        nil,
			hashCalls:      1,
			hashErr:        nil,
			resErr:         nil,
			resToken:       expectedToken,
			expectedUserID: expectedUserID,
		},
		{
			testName:       "Empty email",
			inputEmail:     "",
			inputPas:       expectedHash,
			repoCalls:      0,
			repoErr:        nil,
			hashCalls:      0,
			hashErr:        nil,
			resErr:         shared.NewError(shared.ErrInvalidInput, "email is required", nil),
			resToken:       "",
			expectedUserID: uuid.Nil,
		},
		{
			testName:       "Empty password",
			inputEmail:     "test@example.com",
			inputPas:       nil,
			repoCalls:      0,
			repoErr:        nil,
			hashCalls:      0,
			hashErr:        nil,
			resErr:         shared.NewError(shared.ErrInvalidInput, "password is required", nil),
			resToken:       "",
			expectedUserID: uuid.Nil,
		},
		{
			testName:       "Db error (email not found)",
			inputEmail:     "unexpected@example.com",
			inputPas:       expectedHash,
			repoCalls:      1,
			repoErr:        errors.New("no such user"),
			hashCalls:      0,
			hashErr:        nil,
			resErr:         errors.New("no such user"),
			resToken:       "",
			expectedUserID: uuid.Nil,
		},
		{
			testName:   "Wrong password",
			inputEmail: "test@example.com",
			inputPas:   []byte("wrongpassword"),
			repoCalls:  1,
			repoErr:    nil,
			hashCalls:  1,
			hashErr:    errors.New("password mismatch"),
			resErr: shared.NewError(
				shared.ErrUnauthorized,
				"failed to compare password",
				nil,
			),
			resToken:       "",
			expectedUserID: uuid.Nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mock_ports.NewMockUserRepository(ctrl)
			mockHasher := mock_ports.NewMockPasswordHasher(ctrl)
			authSvc := auth.NewAuth(jwtSecret)

			mockRepo.EXPECT().GetPasswordHash(gomock.Any(), tc.inputEmail).
				Return(tc.expectedUserID, expectedHash, tc.repoErr).Times(tc.repoCalls)

			mockHasher.EXPECT().CompareHashAndPassword(expectedHash, tc.inputPas).
				Return(tc.hashErr).Times(tc.hashCalls)

			userService := &app.UserService{
				Repo:   mockRepo,
				Hasher: mockHasher,
				Auth:   authSvc,
			}

			token, err := userService.Login(context.Background(), tc.inputEmail, tc.inputPas)

			if tc.resErr != nil {
				assert.EqualError(t, err, tc.resErr.Error())
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Parse the JWT token
				parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					// Make sure the method is HMAC
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}

					// Return the secret used for signing
					return jwtSecret, nil
				})

				if err != nil {
					t.Fatalf("Failed to parse token: %v", err)
				}

				if !parsedToken.Valid {
					t.Fatal("Token is invalid")
				}

				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if !ok {
					t.Fatal("Failed to extract claims")
				}

				// Now check the claims
				assert.Equal(t, float64(time.Now().Add(1*time.Minute).Unix()), claims["exp"])
				assert.Equal(t, tc.expectedUserID.String(), claims["user_id"])
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
			authSvc := auth.NewAuth([]byte("test"))
			userService := app.NewUserService(mockRepo, mockHasher, authSvc)

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
