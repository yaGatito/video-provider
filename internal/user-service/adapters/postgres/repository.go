package postgres

import (
	"context"
	"errors"

	"video-provider/pkg/common"
	"video-provider/user-service/adapters/postgres/sqlcgen"
	"video-provider/user-service/domain"
	"video-provider/user-service/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostgresUserRepository struct {
	q sqlcgen.Querier
}

// Ensure PostgresUserRepository implements ports.UserRepository
var _ ports.UserRepository = (*PostgresUserRepository)(nil)

func NewPostgresUserRepository(querier sqlcgen.Querier) *PostgresUserRepository {
	return &PostgresUserRepository{
		q: querier,
	}
}

func (r *PostgresUserRepository) Create(
	ctx context.Context,
	user domain.User,
	password []byte,
) (uuid.UUID, error) {
	params := sqlcgen.CreateUserParams{
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
		return uuid.UUID{}, &common.Error{
			Err:     err,
			Code:    common.ErrInternal,
			Message: "failed to create user"}
	}
	return id, nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row, err := r.q.FindUserById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Nil, &common.Error{
				Err:     err,
				Code:    common.ErrNotFound,
				Message: "user not found with ID " + id.String()}
		} else {
			return domain.Nil, &common.Error{
				Err:     err,
				Code:    common.ErrInternal,
				Message: "failed to retrieve user with ID " + id.String()}
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
			return domain.Nil, &common.Error{
				Err:     err,
				Code:    common.ErrNotFound,
				Message: "user not found with email " + email}
		} else {
			return domain.Nil, &common.Error{
				Err:     err,
				Code:    common.ErrInternal,
				Message: "failed to retrieve user with email " + email}
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
	params := sqlcgen.UpdateUserParams{
		ID:       id,
		Name:     user.Name,
		Lastname: user.LastName,
		Email:    user.Email,
	}

	err := r.q.UpdateUser(ctx, params)
	if err != nil {
		return &common.Error{
			Err:     err,
			Code:    common.ErrInternal,
			Message: "failed to update user"}
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
			return uuid.UUID{}, nil, &common.Error{
				Err:     err,
				Code:    common.ErrNotFound,
				Message: "not found combination of password and email",
				Details: email,
			}
		} else {
			return uuid.UUID{}, nil, &common.Error{
				Err:     err,
				Code:    common.ErrInternal,
				Message: "failed to retrieve password"}
		}
	}

	return row.ID, []byte(row.Password), nil
}
