package repository

import (
	"context"
	"org/api-core/internal/db"

	"github.com/google/uuid"
)

type SQLCUserRepository struct {
	q *db.Queries
}

func NewSQLCUserRepository(q *db.Queries) UserRepository {
	return &SQLCUserRepository{q: q}
}

func (r *SQLCUserRepository) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return r.q.CreateUser(ctx, params)
}

func (r *SQLCUserRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *SQLCUserRepository) UpdateUserPassword(ctx context.Context, params db.UpdateUserPasswordParams) error {
	return r.q.UpdateUserPassword(ctx, params)
}

func (r *SQLCUserRepository) VerifyUser(ctx context.Context, id uuid.UUID) error {
	return r.q.VerifyUser(ctx, id)
}

func (r *SQLCUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (db.GetUserByIDRow, error) {
	return r.q.GetUserByID(ctx, id)
}
