package postgres

import (
	"context"
	"log"

	postgres "video-provider/internal/user-service/adapters/postgres/db"
	"video-provider/internal/user-service/domain"
	"video-provider/internal/user-service/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostgresUserRepository struct {
	q *postgres.Queries
}

// Ensure PostgresUserRepository implements ports.UserRepository
var _ ports.UserRepository = (*PostgresUserRepository)(nil)

func NewPostgresUserRepository(dbConn postgres.DBTX) *PostgresUserRepository {
	return &PostgresUserRepository{
		q: postgres.New(dbConn),
	}
}

func (r *PostgresUserRepository) Create(user domain.User, passwordHash string, passwordSalt string) (uuid.UUID, error) {
	params := postgres.CreateUserParams{
		Name:         user.Name,
		Lastname:     user.LastName,
		Email:        user.Email,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
		CreatedAt:    pgtype.Timestamp{Time: user.CreatedAt, Valid: true},
		Status:       user.Status,
		IsAdmin:      user.IsAdmin,
	}

	id, err := r.q.CreateUser(context.Background(), params)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return uuid.UUID{}, err
	}
	return id, nil
}

func (r *PostgresUserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	row, err := r.q.GetUserById(context.Background(), id)
	if err != nil {
		log.Printf("Error finding user by ID: %v", err)
		return nil, err
	}

	return &domain.User{
		ID:        row.ID,
		Name:      row.Name,
		LastName:  row.Lastname,
		Email:     row.Email,
		CreatedAt: row.CreatedAt.Time,
		Status:    row.Status,
		IsAdmin:   row.IsAdmin,
	}, nil
}

func (r *PostgresUserRepository) FindByEmail(email string) (*domain.User, error) {
	row, err := r.q.GetUserByEmail(context.Background(), email)
	if err != nil {
		log.Printf("Error finding user by email: %v", err)
		return nil, err
	}

	return &domain.User{
		ID:        row.ID,
		Name:      row.Name,
		LastName:  row.Lastname,
		Email:     row.Email,
		CreatedAt: row.CreatedAt.Time,
		Status:    row.Status,
		IsAdmin:   row.IsAdmin,
	}, nil
}

func (r *PostgresUserRepository) Update(user domain.User) error {
	params := postgres.UpdateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Lastname: user.LastName,
		Email:    user.Email,
		Status:   user.Status,
		IsAdmin:  user.IsAdmin,
	}

	err := r.q.UpdateUser(context.Background(), params)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}
	return nil
}
