package service

import (
	"context"

	"org/api-core/internal/db"

	"github.com/google/uuid"
)

type GeminiService interface {
    SaveKey(
        ctx context.Context,
        userID uuid.UUID,
        apiKey string,
    ) (db.UserGeminiKey, error)

    GetKey(
        ctx context.Context,
        userID uuid.UUID,
    ) (db.UserGeminiKey, error)

    UpdateKey(
        ctx context.Context,
        userID uuid.UUID,
        apiKey string,
    ) (db.UserGeminiKey, error)

    DeleteKey(
        ctx context.Context,
        userID uuid.UUID,
    ) error

    HasKey(
        ctx context.Context,
        userID uuid.UUID,
    ) (bool, error)

    GetUserAPIKey(
        ctx context.Context,
        userID uuid.UUID,
    ) (string, error)
}