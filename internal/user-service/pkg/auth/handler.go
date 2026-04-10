package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const jwtSecretEnvVar = "JWT_SECRET"

var secretBytes []byte

const contextUserID = "user_id"
const bearerHeaderPref = "Bearer "

var GetJWTSecret = func() []byte {
	res := []byte(os.Getenv(jwtSecretEnvVar))
	if len(res) == 0 {
		fmt.Println("failed to load secret")
	}
	return res
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if len(bearer) <= len(bearerHeaderPref) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tokenString := bearer[len(bearerHeaderPref):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Use the secret from your UserService
			return GetJWTSecret(), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check if the token is valid
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Extract the user ID from the claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if _, err := uuid.Parse(userID); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextUserID, userID)))
	})
}
