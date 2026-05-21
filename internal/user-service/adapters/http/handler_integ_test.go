package httpadp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"video-provider/pkg/auth"
	"video-provider/pkg/common"
	"video-provider/pkg/config"
	cryptoadp "video-provider/user-service/adapters/crypto"
	"video-provider/user-service/app"
	"video-provider/user-service/domain"
	"video-provider/user-service/ports"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestLoginIntegration(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewUserRepository()
	pwHasher := cryptoadp.NewBCryptPasswordHasher()
	tokenizer := auth.NewTokenizer(config.Config{
		JwtSecret: []byte("test"),
	})
	service := app.NewUserService(repo, pwHasher, tokenizer)

	handler := NewUserHandler(service, common.NewLogger(io.Discard, ""))
	router := mux.NewRouter()
	mockMiddleware := func(next http.Handler) http.Handler {
		return next
	}
	SetupRouter(router, handler, mockMiddleware, mockMiddleware, mockMiddleware)

	server := httptest.NewServer(router)
	defer server.Close()

	// Register step
	createUserResp, err := http.Post(
		server.URL+"/v1/users",
		"application/json",
		bytes.NewBuffer([]byte(`{
			"email":"test@example.com",
			"name":"User",
			"lastname":"Test",
			"password":"PeeSWORD!!11"
		}`)),
	)
	assert.NoError(t, err)
	defer createUserResp.Body.Close()
	assert.NotEmpty(t, createUserResp)
	var createUserRespBody createUserResponse
	err = json.NewDecoder(createUserResp.Body).Decode(&createUserRespBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, createUserRespBody)
	assert.Equal(t, http.StatusCreated, createUserResp.StatusCode)

	// Login step
	loginUserResp, err := http.Post(server.URL+"/v1/login", "application/json",
		bytes.NewBuffer([]byte(`
		{
			"email":"test@example.com",
			"password":"PeeSWORD!!11"
		}`)))
	assert.NoError(t, err)
	defer loginUserResp.Body.Close()
	assert.NotEmpty(t, loginUserResp)
	assert.Equal(t, http.StatusOK, loginUserResp.StatusCode)
	var authResponse authResponse
	err = json.NewDecoder(loginUserResp.Body).Decode(&authResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, authResponse.Token)

	// Profile step
	req, err := http.NewRequest("GET", server.URL+"/v1/users/"+createUserRespBody.UserID, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+authResponse.Token)
	getUserResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer getUserResp.Body.Close()
	assert.Equal(t, http.StatusOK, getUserResp.StatusCode)
	var userResp userResponse
	json.NewDecoder(getUserResp.Body).Decode(&userResp)
	assert.Equal(t, "test@example.com", userResp.Email)
	assert.Equal(t, "User", userResp.Name)
	assert.Equal(t, "Test", userResp.Lastname)
}

type MemoryUserRepository struct {
	Users     map[uuid.UUID]domain.User
	Emails    map[string]uuid.UUID
	Passwords map[string][]byte
	userLock  sync.RWMutex
}

var _ ports.UserRepository = (*MemoryUserRepository)(nil)

func NewUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		Users:     make(map[uuid.UUID]domain.User),
		Emails:    make(map[string]uuid.UUID),
		Passwords: make(map[string][]byte),
		userLock:  sync.RWMutex{},
	}
}

func (repo *MemoryUserRepository) Create(
	ctx context.Context,
	user domain.User,
	password []byte,
) (uuid.UUID, error) {
	repo.userLock.Lock()
	defer repo.userLock.Unlock()

	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}
	user.ID = id
	repo.Users[id] = user
	repo.Emails[user.Email] = id
	repo.Passwords[user.Email] = password

	return id, nil
}

func (repo *MemoryUserRepository) Update(
	ctx context.Context,
	id uuid.UUID,
	user domain.User,
) error {
	return nil
}

func (repo *MemoryUserRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	repo.userLock.RLock()
	defer repo.userLock.RUnlock()

	user, exists := repo.Users[id]
	if !exists {
		return domain.User{}, &common.Error{
			Code: http.StatusNotFound, Message: "user not found",
		}
	}
	return user, nil
}

func (repo *MemoryUserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (domain.User, error) {
	repo.userLock.RLock()
	defer repo.userLock.RUnlock()

	user, exists := repo.Users[repo.Emails[email]]
	if !exists {
		return domain.User{}, &common.Error{
			Code: http.StatusNotFound, Message: "user not found",
		}
	}
	return user, nil
}

func (repo *MemoryUserRepository) GetPasswordHash(
	ctx context.Context,
	email string,
) (uuid.UUID, []byte, error) {
	repo.userLock.RLock()
	defer repo.userLock.RUnlock()

	pass, exists := repo.Passwords[email]
	if !exists {
		return uuid.Nil, nil, &common.Error{
			Code: http.StatusNotFound, Message: "user not found",
		}
	}

	return repo.Users[repo.Emails[email]].ID, pass, nil
}
