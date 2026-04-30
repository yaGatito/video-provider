package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"video-provider/common/auth"
	"video-provider/common/config"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthorizer(t *testing.T) {
	tokenizer := auth.NewTokenizer(config.Config{
		JwtSecret: []byte("test"),
	})
	authorizer := auth.NewAuthorizer(tokenizer)

	assert.NotNil(t, authorizer)
}

func TestAuth(t *testing.T) {
	t.Run("Valid Token", func(t *testing.T) {
		tokenizer := auth.NewTokenizer(config.Config{
			JwtSecret: []byte("test"),
		})
		authorizer := auth.NewAuthorizer(tokenizer)

		userID := uuid.New()
		token, err := tokenizer.CreateToken(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		req, err := http.NewRequest("GET", "/v1/users", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		rr := httptest.NewRecorder()

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDFromCtx := r.Context().Value(auth.ContextUserID)
			assert.NotNil(t, userIDFromCtx)
			assert.Equal(t, userID, userIDFromCtx.(uuid.UUID))
			w.WriteHeader(http.StatusOK)
		})

		authorizer.Auth(next).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		tokenizer := auth.NewTokenizer(config.Config{
			JwtSecret: []byte("test"),
		})
		authorizer := auth.NewAuthorizer(tokenizer)

		req, err := http.NewRequest("GET", "/protected", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer invalid_token")

		rr := httptest.NewRecorder()

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("This should not be called with an invalid token")
		})

		authorizer.Auth(next).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("No Token", func(t *testing.T) {
		tokenizer := auth.NewTokenizer(config.Config{
			JwtSecret: []byte("test"),
		})
		authorizer := auth.NewAuthorizer(tokenizer)

		req, err := http.NewRequest("GET", "/protected", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("This should not be called without a token")
		})

		authorizer.Auth(next).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Options Method", func(t *testing.T) {
		tokenizer := auth.NewTokenizer(config.Config{
			JwtSecret: []byte("test"),
		})
		authorizer := auth.NewAuthorizer(tokenizer)

		req, err := http.NewRequest("OPTIONS", "/protected", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("This should not be called for OPTIONS method")
		})

		authorizer.Auth(next).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})
}
