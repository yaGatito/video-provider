package postgres

import (
	"context"
	"errors"

	"pkg/shared"
	postgres "user-service/adapters/postgres/db"
	"user-service/domain"
	"user-service/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *PostgresUserRepository) Create(
	ctx context.Context,
	user domain.User,
	password []byte,
) (uuid.UUID, error) {
	params := postgres.CreateUserParams{
		Name:      user.Name,
		Lastname:  user.LastName,
		Email:     user.Email,
		Password:  string(password),
		CreatedAt: pgtype.Timestamp{Time: user.CreatedAt, Valid: true},
		Status:    user.Status,
		IsAdmin:   user.IsAdmin,
	}

	id, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return uuid.UUID{}, shared.NewError(shared.ErrInternal, "failed to create user", err)
	}
	return id, nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row, err := r.q.FindUserById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, shared.NewError(
				shared.ErrNotFound,
				"user not found with ID "+id.String(),
				err,
			)
		} else {
			return domain.User{}, shared.NewError(shared.ErrInternal, "failed to retrieve user with ID "+id.String(), err)
		}
	}

	return domain.User{
		ID:        row.ID,
		Name:      row.Name,
		LastName:  row.Lastname,
		Email:     row.Email,
		CreatedAt: row.CreatedAt.Time,
		Status:    row.Status,
		IsAdmin:   row.IsAdmin,
	}, nil
}

func (r *PostgresUserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (domain.User, error) {
	row, err := r.q.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, shared.NewError(
				shared.ErrNotFound,
				"user not found with email "+email,
				err,
			)
		} else {
			return domain.User{}, shared.NewError(shared.ErrInternal, "failed to retrieve user with email "+email, err)
		}
	}

	return domain.User{
		ID:        row.ID,
		Name:      row.Name,
		LastName:  row.Lastname,
		Email:     row.Email,
		CreatedAt: row.CreatedAt.Time,
		Status:    row.Status,
		IsAdmin:   row.IsAdmin,
	}, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, id uuid.UUID, user domain.User) error {
	params := postgres.UpdateUserParams{
		ID:       id,
		Name:     user.Name,
		Lastname: user.LastName,
		Email:    user.Email,
	}

	err := r.q.UpdateUser(ctx, params)
	if err != nil {
		return shared.NewError(shared.ErrInternal, "failed to update user", err)
	}
	return nil
}

func (r *PostgresUserRepository) GetPasswordHash(
	ctx context.Context,
	email string,
) (uuid.UUID, []byte, error) {
	row, err := r.q.GetPassword(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, nil, shared.NewError(
				shared.ErrNotFound,
				"not found password and email combination for email: "+email,
				err,
			)
		} else {
			return uuid.UUID{}, nil, shared.NewError(shared.ErrInternal, "failed to retrieve password", err)
		}
	}

	return row.ID, []byte(row.Password), nil
}
