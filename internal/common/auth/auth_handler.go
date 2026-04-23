package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKeyUserID string

const contextUserID contextKeyUserID = "USER_ID"

const bearerHeaderPrefix string = "Bearer "

type Authorizer struct {
	auth *Auth
}

func NewAuthorizer(auth *Auth) Authorizer {
	return Authorizer{
		auth: auth,
	}
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

		userID, err := a.auth.ValidateToken(tokenString)
		if err != nil {
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
