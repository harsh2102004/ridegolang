package repository

import (
	"context"
	"org/api-core/internal/db"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
    GetUserByID(ctx context.Context, id uuid.UUID) (db.GetUserByIDRow, error)
	UpdateUserPassword(ctx context.Context, params db.UpdateUserPasswordParams) error
	VerifyUser(ctx context.Context, id uuid.UUID) error
}
