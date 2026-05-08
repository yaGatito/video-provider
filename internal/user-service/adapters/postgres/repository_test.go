package postgres_test

import (
	"context"
	"testing"

	"video-provider/user-service/adapters/postgres"
	"video-provider/user-service/adapters/postgres/sqlcgen"
	"video-provider/user-service/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) CreateUser(ctx context.Context, params sqlcgen.CreateUserParams) (uuid.UUID, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockQuerier) FindUserById(ctx context.Context, id uuid.UUID) (sqlcgen.FindUserByIdRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(sqlcgen.FindUserByIdRow), args.Error(1)
}

func (m *MockQuerier) FindUserByEmail(ctx context.Context, email string) (sqlcgen.FindUserByEmailRow, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(sqlcgen.FindUserByEmailRow), args.Error(1)
}

func (m *MockQuerier) UpdateUser(ctx context.Context, params sqlcgen.UpdateUserParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockQuerier) GetPassword(ctx context.Context, email string) (sqlcgen.GetPasswordRow, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(sqlcgen.GetPasswordRow), args.Error(1)
}

func TestPostgresUserRepository_Create(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := postgres.NewPostgresUserRepository(mockQuerier)

	user := domain.User{
		Name:     "John",
		LastName: "Doe",
		Email:    "john.doe@example.com",
	}
	password := []byte("password")

	mockId := uuid.New()
	mockQuerier.On("CreateUser", mock.Anything, sqlcgen.CreateUserParams{
		Name:      user.Name,
		Lastname:  user.LastName,
		Email:     user.Email,
		Password:  string(password),
		CreatedAt: pgtype.Timestamp{Time: user.CreatedAt, Valid: true},
		Status:    user.Status,
		IsAdmin:   user.IsAdmin,
	}).Return(mockId, nil)

	id, err := repo.Create(context.Background(), user, password)
	assert.NoError(t, err)
	assert.Equal(t, mockId, id)
}

func TestPostgresUserRepository_FindByID(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := postgres.NewPostgresUserRepository(mockQuerier)

	userID := uuid.New()

	mockRow := sqlcgen.FindUserByIdRow{
		ID:        userID,
		Name:      "John",
		Lastname:  "Doe",
		Email:     "john.doe@example.com",
		CreatedAt: pgtype.Timestamp{Time: domain.Nil.CreatedAt, Valid: true},
		Status:    domain.Nil.Status,
		IsAdmin:   domain.Nil.IsAdmin,
	}
	mockQuerier.On("FindUserById", mock.Anything, userID).Return(mockRow, nil)

	user, err := repo.FindByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
}

func TestPostgresUserRepository_FindByEmail(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := postgres.NewPostgresUserRepository(mockQuerier)

	email := "john.doe@example.com"

	mockRow := sqlcgen.FindUserByEmailRow{
		ID:        uuid.New(),
		Name:      "John",
		Lastname:  "Doe",
		Email:     email,
		CreatedAt: pgtype.Timestamp{Time: domain.Nil.CreatedAt, Valid: true},
		Status:    domain.Nil.Status,
		IsAdmin:   domain.Nil.IsAdmin,
	}
	mockQuerier.On("FindUserByEmail", mock.Anything, email).Return(mockRow, nil)

	user, err := repo.FindByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
}

func TestPostgresUserRepository_Update(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := postgres.NewPostgresUserRepository(mockQuerier)

	userID := uuid.New()
	user := domain.User{
		Name:     "John",
		LastName: "Doe",
		Email:    "john.doe@example.com",
	}

	mockQuerier.On("UpdateUser", mock.Anything, sqlcgen.UpdateUserParams{
		ID:       userID,
		Name:     user.Name,
		Lastname: user.LastName,
		Email:    user.Email,
	}).Return(nil)

	err := repo.Update(context.Background(), userID, user)
	assert.NoError(t, err)
}

func TestPostgresUserRepository_GetPasswordHash(t *testing.T) {
	mockQuerier := new(MockQuerier)
	repo := postgres.NewPostgresUserRepository(mockQuerier)

	email := "john.doe@example.com"
	hash := []byte("password")

	mockRow := sqlcgen.GetPasswordRow{
		ID:       uuid.New(),
		Password: string(hash),
	}
	mockQuerier.On("GetPassword", mock.Anything, email).Return(mockRow, nil)

	_, passwordHash, err := repo.GetPasswordHash(context.Background(), email)
	assert.NoError(t, err)
	assert.Equal(t, hash, passwordHash)
}
