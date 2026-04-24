// internal/common/auth/auth.go
package auth

import (
	"fmt"
	"time"
	"video-provider/common/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const claimUserIDKey string = "user_id"
const claimTokenExpKey string = "exp"

type Tokenizer struct {
	config config.Config
}

func NewTokenizer(c config.Config) *Tokenizer {
	return &Tokenizer{config: c}
}

func (t *Tokenizer) CreateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		claimUserIDKey:   userID.String(),
		claimTokenExpKey: time.Now().Add(t.config.JsonConf.TokenExpTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(t.config.JwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (t *Tokenizer) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.config.JwtSecret, nil
	})

	if err != nil {
		return "", fmt.Errorf("error parsing token: %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	userID, ok := claims[claimUserIDKey].(string)
	if !ok {
		return "", fmt.Errorf("invalid user id")
	}

	expDateInterface, ok := claims[claimTokenExpKey]
	if !ok {
		return "", fmt.Errorf("invalid exp date")
	}

	expDate, ok := expDateInterface.(float64)
	if !ok {
		return "", fmt.Errorf("invalid exp date type, expected float64, got %T", expDateInterface)
	}

	isExpired := time.Unix(int64(expDate), 0).Before(time.Now())
	if isExpired {
		return "", fmt.Errorf("token expired")
	}

	return userID, nil
}
