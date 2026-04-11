package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type contextKeyUserID string

const contextUserID contextKeyUserID = "USER_ID"

const jwtSecretEnvVar = "JWT_SECRET"

type Authorizer struct {
}

func (a Authorizer) GetJWTSecret() []byte {
	return []byte(os.Getenv(jwtSecretEnvVar))
}

func (a Authorizer) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		tokenString := bearer[len("Bearer "):]
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Use the secret from your UserService
			return a.GetJWTSecret(), nil
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

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextUserID, userID)))
	})
}
