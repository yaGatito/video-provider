package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKeyUserID string

const contextUserID contextKeyUserID = "USER_ID"

const jwtSecretEnvVar string = "JWT_SECRET"

const bearerHeaderPrefix string = "Bearer "

type Authorizer struct{}

func (a Authorizer) GetJWTSecret() []byte {
	return []byte(os.Getenv(jwtSecretEnvVar))
}

func (a Authorizer) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if len(bearer) <= 7 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tokenString, ok := strings.CutPrefix(bearer, bearerHeaderPrefix)
		if tokenString == "" || !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return a.GetJWTSecret(), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

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

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextUserID, userID)))
	})
}
