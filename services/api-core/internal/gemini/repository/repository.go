package repository

import (
	"context"
	"org/api-core/internal/db"

	"github.com/google/uuid"
)

type GeminiRepository interface {
	CreateGeminiKey(ctx context.Context, params db.CreateUserGeminiKeyParams) (db.UserGeminiKey, error)
	GetGeminiKeyByUserID(ctx context.Context, userID uuid.UUID) (db.UserGeminiKey, error)
	UpdateGeminiKey(ctx context.Context, params db.UpdateUserGeminiKeyParams) (db.UserGeminiKey, error)
	DeleteGeminiKey(ctx context.Context, userID uuid.UUID) error
	UserHasGeminiKey(ctx context.Context, userID uuid.UUID) (bool, error)
    
    GetUserAPIKey(ctx context.Context, userID uuid.UUID) (string, error)
}