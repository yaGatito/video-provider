package auth_test

import (
	"testing"
	"time"
	"video-provider/common/auth"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	authSvc := auth.NewAuth([]byte("test"))

	userID := uuid.New()
	exp := time.Now().Add(2 * time.Second).Unix()
	token, err := authSvc.CreateToken(userID, exp)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	authSvc := auth.NewAuth([]byte("test"))

	userID := uuid.Must(uuid.NewRandom())
	exp := time.Now().Add(1 * time.Second).Unix()
	token, err := authSvc.CreateToken(userID, exp)
	assert.NoError(t, err)

	// Valid token
	userIDFromToken, err := authSvc.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID.String(), userIDFromToken)

	// Expired token
	time.Sleep(2 * time.Second) // wait for token to expire
	userIDFromToken, err = authSvc.ValidateToken(token)
	assert.Equal(t, "", userIDFromToken)
	assert.Error(t, err)
}

func TestValidateToken_InvalidJWT(t *testing.T) {
	authSvc := auth.NewAuth([]byte("test"))

	// Invalid token
	_, err := authSvc.ValidateToken("invalid_token")
	assert.Error(t, err)
}

func TestValidateToken_InvalidClaims(t *testing.T) {
	authSvc := auth.NewAuth([]byte("test"))

	// Create a token with missing user_id
	claims := jwt.MapClaims{
		"exp": time.Now().Add(1 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(authSvc.Secret)
	assert.NoError(t, err)

	_, err = authSvc.ValidateToken(signedToken)
	assert.Error(t, err)
}
