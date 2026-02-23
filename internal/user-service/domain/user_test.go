package domain

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		ttname   string
		email    string
		name     string
		lastname string
		wantErr  bool
	}{
		{
			ttname:   "Valid user",
			email:    "test@example.com",
			name:     "John",
			lastname: "Doe",
			wantErr:  false,
		},
		{
			ttname:   "Invalid email format",
			email:    "invalid-email",
			name:     "John",
			lastname: "Doe",
			wantErr:  true,
		},
		{
			ttname:   "Empty name",
			email:    "test@example.com",
			name:     "",
			lastname: "Doe",
			wantErr:  true,
		},
		{
			ttname:   "Name with non-alphabetical characters",
			email:    "test@example.com",
			name:     "John123",
			lastname: "Doe",
			wantErr:  true,
		},
		{
			ttname:   "Empty lastname",
			email:    "test@example.com",
			name:     "John",
			lastname: "",
			wantErr:  true,
		},
		{
			ttname:   "Lastname with non-alphabetical characters",
			email:    "test@example.com",
			name:     "John",
			lastname: "Doe123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.ttname, func(t *testing.T) {
			got, err := NewUser(tt.email, tt.ttname, tt.lastname)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Email != tt.email {
				t.Errorf("NewUser().Email = %v, want %v", got.Email, tt.email)
			}
			if !tt.wantErr && got.Name != tt.ttname {
				t.Errorf("NewUser().Name = %v, want %v", got.Name, tt.ttname)
			}
			if !tt.wantErr && got.LastName != tt.lastname {
				t.Errorf("NewUser().LastName = %v, want %v", got.LastName, tt.lastname)
			}
			if !tt.wantErr && got.CreatedAt.IsZero() {
				t.Errorf("NewUser().CreatedAt is zero value")
			}
			if !tt.wantErr && got.IsAdmin != false {
				t.Errorf("NewUser().IsAdmin = %v, want false", got.IsAdmin)
			}
			if !tt.wantErr && got.Status != "active" {
				t.Errorf("NewUser().Status = %v, want 'active'", got.Status)
			}
		})
	}
}

func TestPassword_ValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "StrongPass123",
			wantErr:  false,
		},
		{
			name:     "Too short password",
			password: "pass",
			wantErr:  true,
		},
		{
			name:     "No uppercase letter",
			password: "lowercasepassword",
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Password(tt.password)
			err := got.ValidatePassword()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
