package service

import (
	"context"

	"org/api-core/internal/db"
	"org/api-core/internal/gemini/repository"

	"github.com/google/uuid"
)

type geminiService struct {
	repo repository.GeminiRepository
}

func NewGeminiService(
	repo repository.GeminiRepository,
) GeminiService {
	return &geminiService{
		repo: repo,
	}
}

func (s *geminiService) SaveKey(
	ctx context.Context,
	userID uuid.UUID,
	apiKey string,
) (db.UserGeminiKey, error) {

	return s.repo.CreateGeminiKey(
		ctx,
		db.CreateUserGeminiKeyParams{
			UserID: userID,
			ApiKey: apiKey,
		},
	)
}

func (s *geminiService) GetKey(
	ctx context.Context,
	userID uuid.UUID,
) (db.UserGeminiKey, error) {

	return s.repo.GetGeminiKeyByUserID(ctx, userID)
}

func (s *geminiService) UpdateKey(
	ctx context.Context,
	userID uuid.UUID,
	apiKey string,
) (db.UserGeminiKey, error) {

	return s.repo.UpdateGeminiKey(
		ctx,
		db.UpdateUserGeminiKeyParams{
			UserID: userID,
			ApiKey: apiKey,
		},
	)
}

func (s *geminiService) DeleteKey(
	ctx context.Context,
	userID uuid.UUID,
) error {

	return s.repo.DeleteGeminiKey(ctx, userID)
}

func (s *geminiService) HasKey(
	ctx context.Context,
	userID uuid.UUID,
) (bool, error) {

	return s.repo.UserHasGeminiKey(ctx, userID)
}

func (s *geminiService) GetUserAPIKey(
    ctx context.Context,
    userID uuid.UUID,
) (string, error) {
    return s.repo.GetUserAPIKey(ctx, userID)
}