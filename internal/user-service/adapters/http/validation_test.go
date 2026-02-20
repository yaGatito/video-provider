package httpadp

import (
	"strings"
	"testing"
)

func TestCreateUserRequest_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		request     createUserRequest
		expectedErr bool
	}{
		{
			name: "Valid request",
			request: createUserRequest{
				Email:    "test@example.com",
				Name:     "John",
				LastName: "Doe",
				Password: "SecurePass123",
			},
			expectedErr: false,
		},
		{
			name: "Empty email",
			request: createUserRequest{
				Email:    "",
				Name:     "John",
				LastName: "Doe",
				Password: "SecurePass123",
			},
			expectedErr: true,
		},
		{
			name: "Long email",
			request: createUserRequest{
				Email:    strings.Repeat("a", 101) + "@example.com",
				Name:     "John",
				LastName: "Doe",
				Password: "SecurePass123",
			},
			expectedErr: true,
		},
		{
			name: "Empty name",
			request: createUserRequest{
				Email:    "test@example.com",
				Name:     "",
				LastName: "Doe",
				Password: "SecurePass123",
			},
			expectedErr: true,
		},
		{
			name: "Long name",
			request: createUserRequest{
				Email:    "test@example.com",
				Name:     strings.Repeat("a", 51),
				LastName: "Doe",
				Password: "SecurePass123",
			},
			expectedErr: true,
		},
		{
			name: "Empty lastname",
			request: createUserRequest{
				Email:    "test@example.com",
				Name:     "John",
				LastName: "",
				Password: "SecurePass123",
			},
			expectedErr: true,
		},
		{
			name: "Long lastname",
			request: createUserRequest{
				Email:    "test@example.com",
				Name:     "John",
				LastName: strings.Repeat("a", 101),
				Password: "SecurePass123",
			},
			expectedErr: true,
		},
		{
			name: "Short password",
			request: createUserRequest{
				Email:    "test@example.com",
				Name:     "John",
				LastName: "Doe",
				Password: "Pass",
			},
			expectedErr: true,
		},
		{
			name: "Long password",
			request: createUserRequest{
				Email:    "test@example.com",
				Name:     "John",
				LastName: "Doe",
				Password: strings.Repeat("a", 101),
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.request.validate()
			if (err != nil) != tc.expectedErr {
				if tc.expectedErr {
					t.Errorf("Expected an error but got none")
				} else {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
