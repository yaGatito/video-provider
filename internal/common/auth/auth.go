// internal/common/auth/auth.go
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const claimUserIDKey string = "user_id"
const claimTokenExpKey string = "exp"

type Auth struct {
	Secret []byte
}

func NewAuth(secret []byte) *Auth {
	return &Auth{Secret: secret}
}

func (a *Auth) CreateToken(userID uuid.UUID, expDate int64) (string, error) {
	claims := jwt.MapClaims{
		claimUserIDKey:   userID.String(),
		claimTokenExpKey: expDate,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(a.Secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *Auth) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return a.Secret, nil
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

	fmt.Println("Exp Date (float64):", expDate)
	fmt.Println("Current Time:", time.Now().Unix())

	isExpired := time.Unix(int64(expDate), 0).Before(time.Now())
	if isExpired {
		return "", fmt.Errorf("token expired")
	}

	return userID, nil
}