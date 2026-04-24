package auth_test

import (
	"testing"
	"time"
	"video-provider/common/auth"
	"video-provider/common/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	exp := 1 * time.Second
	c := config.Config{
		JwtSecret: []byte("test"),
		JsonConf: config.JsonServiceConfig{
			TokenExpTime: exp,
		},
	}
	tokenizer := auth.NewTokenizer(c)

	userID := uuid.New()
	token, err := tokenizer.CreateToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	exp := 1 * time.Second
	c := config.Config{
		JwtSecret: []byte("test"),
		JsonConf: config.JsonServiceConfig{
			TokenExpTime: exp,
		},
	}
	tokenizer := auth.NewTokenizer(c)

	userID := uuid.Must(uuid.NewRandom())
	token, err := tokenizer.CreateToken(userID)
	assert.NoError(t, err)

	// Valid token
	userIDFromToken, err := tokenizer.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID.String(), userIDFromToken)

	// Expired token
	time.Sleep(2 * time.Second) // wait for token to expire
	userIDFromToken, err = tokenizer.ValidateToken(token)
	assert.Equal(t, "", userIDFromToken)
	assert.Error(t, err)
}

func TestValidateToken_InvalidJWT(t *testing.T) {
	tokenizer := auth.NewTokenizer(config.Config{JwtSecret: []byte("test")})

	// Invalid token
	_, err := tokenizer.ValidateToken("invalid_token")
	assert.Error(t, err)
}

func TestValidateToken_InvalidClaims(t *testing.T) {
	c := config.Config{
		JwtSecret: []byte("test"),
	}
	tokenizer := auth.NewTokenizer(c)

	// Create a token with missing user_id
	claims := jwt.MapClaims{
		"exp": time.Now().Add(1 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(c.JwtSecret)
	assert.NoError(t, err)

	_, err = tokenizer.ValidateToken(signedToken)
	assert.Error(t, err)
}
