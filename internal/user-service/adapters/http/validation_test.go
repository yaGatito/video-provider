package httpadp

// func TestCreateUserRequest_Validate(t *testing.T) {
// 	testCases := []struct {
// 		name        string
// 		request     createUserRequest
// 		expectedErr bool
// 	}{
// 		{
// 			name: "Valid request",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     "John",
// 				LastName: "Doe",
// 				Password: "SecurePass123!!!",
// 			},
// 			expectedErr: false,
// 		},
// 		{
// 			name: "Empty email",
// 			request: createUserRequest{
// 				Email:    "",
// 				Name:     "John",
// 				LastName: "Doe",
// 				Password: "SecurePass123!!!",
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Long email",
// 			request: createUserRequest{
// 				Email:    strings.Repeat("a", 101) + "@example.com",
// 				Name:     "John",
// 				LastName: "Doe",
// 				Password: "SecurePass123!!!",
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Empty name",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     "",
// 				LastName: "Doe",
// 				Password: "SecurePass123!!!",
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Short name",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     "J",
// 				LastName: "Doe",
// 				Password: "SecurePass123!!!",
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Long name",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     strings.Repeat("a", policy.MaxInputTextLen+1),
// 				LastName: "Doe",
// 				Password: "SecurePass123!!!",
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Empty lastname",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     "John",
// 				LastName: "",
// 				Password: "SecurePass123!!!",
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Short lastname",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     "John",
// 				LastName: "D",
// 				Password: "SecurePass123!!!",
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Long lastname",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     "John",
// 				LastName: strings.Repeat("a", policy.MaxInputTextLen+1),
// 				Password: "SecurePass123!!!",
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Short password",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     "John",
// 				LastName: "Doe",
// 				Password: "aA1!",
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Long password",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     "John",
// 				LastName: "Doe",
// 				Password: strings.Repeat("aA1!", 26),
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "No upper-case symbols in password",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     "John",
// 				LastName: "Doe",
// 				Password: strings.Repeat("a1!", 10),
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "No digits in password",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     "John",
// 				LastName: "Doe",
// 				Password: strings.Repeat("aA!", 10),
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "No spec chars in password",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     "John",
// 				LastName: "Doe",
// 				Password: strings.Repeat("aA1", 10),
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Unsupported spec chars in password",
// 			request: createUserRequest{
// 				Email:    "test@example.com",
// 				Name:     "John",
// 				LastName: "Doe",
// 				Password: strings.Repeat("aA1!/", 4),
// 			},
// 			expectedErr: true,
// 		},
// 	}

// 	t.Parallel()

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {

// 			err := tc.request.validate(NewUserValidator())
// 			if (err != nil) != tc.expectedErr {
// 				if tc.expectedErr {
// 					t.Errorf("Expected an error but got none")
// 				} else {
// 					t.Errorf("Unexpected error: %v", err)
// 				}
// 			}
// 		})
// 	}
// }
