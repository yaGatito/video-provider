package postgres

import (
	"context"
	"log"

	postgres "video-provider/internal/user-service/adapters/postgres/db"
	"video-provider/internal/user-service/domain"
	"video-provider/internal/user-service/ports"

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

func (r *PostgresUserRepository) Create(user *domain.User, passwordHash string, passwordSalt string) (int64, error) {
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
		return 0, err
	}
	return id, nil
}

func (r *PostgresUserRepository) FindByID(id int64) (*domain.User, error) {
	row, err := r.q.GetUser(context.Background(), id)
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
