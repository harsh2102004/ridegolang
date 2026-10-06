package repository

import (
	"context"
	"org/api-core/internal/db"

	"github.com/google/uuid"
)

type SQLCGeminiRepository struct {
	q *db.Queries
}

func NewSQLCGeminiRepository(q *db.Queries) GeminiRepository {
	return &SQLCGeminiRepository{
		q: q,
	}
}

func (r *SQLCGeminiRepository) CreateGeminiKey(
	ctx context.Context,
	params db.CreateUserGeminiKeyParams,
) (db.UserGeminiKey, error) {
	return r.q.CreateUserGeminiKey(ctx, params)
}

func (r *SQLCGeminiRepository) GetGeminiKeyByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (db.UserGeminiKey, error) {
	return r.q.GetUserGeminiKeyByUserID(ctx, userID)
}

func (r *SQLCGeminiRepository) UpdateGeminiKey(
	ctx context.Context,
	params db.UpdateUserGeminiKeyParams,
) (db.UserGeminiKey, error) {
	return r.q.UpdateUserGeminiKey(ctx, params)
}

func (r *SQLCGeminiRepository) DeleteGeminiKey(
	ctx context.Context,
	userID uuid.UUID,
) error {
	return r.q.DeleteUserGeminiKey(ctx, userID)
}

func (r *SQLCGeminiRepository) UserHasGeminiKey(
	ctx context.Context,
	userID uuid.UUID,
) (bool, error) {
	return r.q.UserHasGeminiKey(ctx, userID)
}

func (r *SQLCGeminiRepository) GetUserAPIKey(
	ctx context.Context,
	userID uuid.UUID,
) (string, error) {
	key, err := r.q.GetUserGeminiKeyByUserID(
		ctx,
		userID,
	)
	if err != nil {
		return "", err
	}

	return key.ApiKey, nil
}