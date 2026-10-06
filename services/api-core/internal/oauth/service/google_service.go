package service

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"

	"org/api-core/internal/db"
	"org/api-core/internal/oauth/repository"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type OAuthService struct {
	accountRepo repository.ConnectedAccountRepository
	logger      *slog.Logger
}

func NewOAuthService(accountRepo repository.ConnectedAccountRepository, logger *slog.Logger) *OAuthService {
	return &OAuthService{
		accountRepo: accountRepo,
		logger:      logger,
	}
}

func (s *OAuthService) SaveGoogleAccount(
	ctx context.Context,
	userID uuid.UUID,
	token *oauth2.Token,
	scopes []string,
) error {
	s.logger.Info("saving google account",
		"user_id", userID.String(),
	)

	_, err := s.accountRepo.UpsertConnectedAccount(ctx, db.UpsertConnectedAccountParams{
		UserID:      userID,
		Provider:    "google",
		AccessToken: token.AccessToken,
		RefreshToken: sql.NullString{
			String: token.RefreshToken,
			Valid:  token.RefreshToken != "",
		},
		ExpiresAt: sql.NullTime{
			Time:  token.Expiry,
			Valid: !token.Expiry.IsZero(),
		},
		Scopes: sql.NullString{
			String: strings.Join(scopes, " "),
			Valid:  len(scopes) > 0,
		},
	})
	if err != nil {
		s.logger.Error("failed to save google account",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return err
	}

	s.logger.Info("google account saved successfully",
		"user_id", userID.String(),
	)
	return nil
}