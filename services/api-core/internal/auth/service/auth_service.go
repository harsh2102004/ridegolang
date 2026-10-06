package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	authEmail "org/api-core/internal/auth/email"
	"org/api-core/internal/auth/repository"
	"org/api-core/internal/auth/security"
	"org/api-core/internal/auth/tokenstore"
	"org/api-core/internal/db"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewAuthService(userRepo repository.UserRepository, logger *slog.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		logger:   logger,
	}
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func (s *AuthService) Signup(ctx context.Context, email, password string) (string, error) {
	s.logger.Info("signup started",
		"email", email,
	)

	existing, _ := s.userRepo.GetUserByEmail(ctx, email)
	if existing.Email != "" {
		s.logger.Warn("signup blocked: email already registered",
			"email", email,
		)

		return "", errors.New("email already registered")
	}

	if len(password) < 6 {
		s.logger.Warn("signup blocked: weak password",
			"email", email,
		)

		return "", errors.New("password must be at least 6 characters")
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		s.logger.Error("failed to hash password",
			"email", email,
			"error", err.Error(),
		)

		return "", err
	}

	userID := uuid.New()

	_, err = s.userRepo.CreateUser(ctx, db.CreateUserParams{
		ID:         userID,
		Email:      email,
		Password:   hashedPassword,
		IsVerified: false,
	})
	if err != nil {
		s.logger.Error("failed to create user",
			"email", email,
			"error", err.Error(),
		)

		return "", err
	}

	token, err := security.GenerateAccessToken(userID.String())
	if err != nil {
		s.logger.Error("failed to generate verification token",
			"user_id", userID.String(),
			"error", err.Error(),
		)

		return "", err
	}

	err = tokenstore.StoreVerificationToken(token, userID.String(), 24*time.Hour)
	if err != nil {
		s.logger.Error("failed to store verification token",
			"user_id", userID.String(),
			"error", err.Error(),
		)

		return "", err
	}

	err = authEmail.SendVerificationEmail(email, token)
	if err != nil {
		s.logger.Error("failed to send verification email",
			"user_id", userID.String(),
			"email", email,
			"error", err.Error(),
		)

		return "", err
	}

	s.logger.Info("signup completed successfully",
		"user_id", userID.String(),
		"email", email,
	)

	return token, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email,
	password string,
) (*TokenPair, error) {

	s.logger.Info("login started",
		"email", email,
	)

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Warn("login failed: user not found",
			"email", email,
		)

		return nil, errors.New("invalid credentials")
	}

	err = security.CheckPassword(user.Password, password)
	if err != nil {
		s.logger.Warn("login failed: invalid password",
			"email", email,
		)

		return nil, errors.New("invalid credentials")
	}

	if !user.IsVerified {
		s.logger.Warn("login blocked: email not verified",
			"email", email,
		)

		return nil, errors.New("please verify your email first")
	}

	accessToken, err := security.GenerateAccessToken(user.ID.String())
	if err != nil {
		s.logger.Error("failed to generate access token",
			"user_id", user.ID.String(),
			"error", err.Error(),
		)

		return nil, err
	}

	refreshToken, err := security.GenerateRefreshToken(user.ID.String())
	if err != nil {
		s.logger.Error("failed to generate refresh token",
			"user_id", user.ID.String(),
			"error", err.Error(),
		)

		return nil, err
	}

	err = tokenstore.StoreRefreshToken(
		ctx,
		refreshToken,
		user.ID.String(),
	)
	if err != nil {
		s.logger.Error("failed to store refresh token",
			"user_id", user.ID.String(),
			"error", err.Error(),
		)

		return nil, err
	}

	s.logger.Info("login completed successfully",
		"user_id", user.ID.String(),
		"email", email,
	)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	s.logger.Info("verify email started")

	userID, err := tokenstore.GetUserIDByToken(token)
	if err != nil {
		s.logger.Warn("verify email failed: invalid or expired token",
			"error", err.Error(),
		)
		return errors.New("invalid or expired token")
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Error("verify email failed: invalid user ID format",
			"user_id", userID,
			"error", err.Error(),
		)
		return err
	}

	err = s.userRepo.VerifyUser(ctx, parsedID)
	if err != nil {
		s.logger.Error("verify email failed: database update error",
			"user_id", parsedID.String(),
			"error", err.Error(),
		)
		return err
	}

	err = tokenstore.DeleteToken(token)
	if err != nil {
		s.logger.Warn("verify email: failed to delete used token",
			"user_id", parsedID.String(),
			"error", err.Error(),
		)
		// Not a critical error, continue
	}

	s.logger.Info("verify email completed successfully",
		"user_id", parsedID.String(),
	)
	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	s.logger.Info("forgot password started",
		"email", email,
	)

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// Do not reveal whether the user exists
		s.logger.Info("forgot password: user not found or error, no action taken",
			"email", email,
			"error", err.Error(),
		)
		return nil
	}

	token, err := security.GenerateAccessToken(user.ID.String())
	if err != nil {
		s.logger.Error("forgot password: failed to generate token",
			"user_id", user.ID.String(),
			"email", email,
			"error", err.Error(),
		)
		return err
	}

	err = tokenstore.StorePasswordResetToken(
		ctx,
		token,
		user.ID.String(),
		15*time.Minute,
	)
	if err != nil {
		s.logger.Error("forgot password: failed to store reset token",
			"user_id", user.ID.String(),
			"email", email,
			"error", err.Error(),
		)
		return err
	}

	err = authEmail.SendPasswordResetEmail(user.Email, token)
	if err != nil {
		s.logger.Error("forgot password: failed to send reset email",
			"user_id", user.ID.String(),
			"email", email,
			"error", err.Error(),
		)
		return err
	}

	s.logger.Info("forgot password completed successfully",
		"user_id", user.ID.String(),
		"email", email,
	)
	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	s.logger.Info("reset password started")

	token = strings.TrimSpace(token)

	userID, err := tokenstore.GetPasswordResetToken(ctx, token)
	if err != nil {
		s.logger.Warn("reset password failed: invalid or expired token",
			"error", err.Error(),
		)
		return errors.New("invalid or expired token")
	}

	s.logger.Info("reset password: token validated",
		"user_id", userID,
	)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("reset password failed: password hashing error",
			"user_id", userID,
			"error", err.Error(),
		)
		return err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Error("reset password failed: invalid user ID format",
			"user_id", userID,
			"error", err.Error(),
		)
		return err
	}

	err = s.userRepo.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		Password: string(hashedPassword),
		ID:       userUUID,
	})
	if err != nil {
		s.logger.Error("reset password failed: database update error",
			"user_id", userUUID.String(),
			"error", err.Error(),
		)
		return err
	}

	err = tokenstore.DeletePasswordResetToken(ctx, token)
	if err != nil {
		s.logger.Warn("reset password: failed to delete used reset token",
			"user_id", userUUID.String(),
			"error", err.Error(),
		)
		// Non-critical, continue
	}

	s.logger.Info("reset password completed successfully",
		"user_id", userUUID.String(),
	)
	return nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	s.logger.Info("refresh token started")

	refreshToken = strings.TrimSpace(refreshToken)

	userID, err := tokenstore.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.Warn("refresh token failed: invalid token",
			"error", err.Error(),
		)
		return nil, errors.New("invalid refresh token")
	}

	// Delete old refresh token (rotation)
	err = tokenstore.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.Warn("refresh token: failed to delete old token",
			"user_id", userID,
			"error", err.Error(),
		)
		// Continue rotation, but log warning
	}

	newAccessToken, err := security.GenerateAccessToken(userID)
	if err != nil {
		s.logger.Error("refresh token: failed to generate new access token",
			"user_id", userID,
			"error", err.Error(),
		)
		return nil, err
	}

	newRefreshToken, err := security.GenerateRefreshToken(userID)
	if err != nil {
		s.logger.Error("refresh token: failed to generate new refresh token",
			"user_id", userID,
			"error", err.Error(),
		)
		return nil, err
	}

	err = tokenstore.StoreRefreshToken(ctx, newRefreshToken, userID)
	if err != nil {
		s.logger.Error("refresh token: failed to store new refresh token",
			"user_id", userID,
			"error", err.Error(),
		)
		return nil, err
	}

	s.logger.Info("refresh token completed successfully",
		"user_id", userID,
	)

	return &TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	s.logger.Info("logout started")

	refreshToken = strings.TrimSpace(refreshToken)

	err := tokenstore.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.Warn("logout: failed to delete refresh token",
			"error", err.Error(),
		)
		return err
	}

	s.logger.Info("logout completed successfully")
	return nil
}

func (s *AuthService) GetOrCreateGoogleUser(ctx context.Context, email string) (*TokenPair, error) {
	s.logger.Info("get or create google user started",
		"email", email,
	)

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Info("google user not found, creating new user",
				"email", email,
			)

			userID := uuid.New()
			user, err = s.userRepo.CreateUser(ctx, db.CreateUserParams{
				ID:         userID,
				Email:      email,
				Password:   "",
				IsVerified: true,
			})
			if err != nil {
				s.logger.Error("failed to create google user",
					"email", email,
					"error", err.Error(),
				)
				return nil, err
			}
			s.logger.Info("created new google user",
				"user_id", user.ID.String(),
				"email", email,
			)
		} else {
			s.logger.Error("failed to lookup google user",
				"email", email,
				"error", err.Error(),
			)
			return nil, err
		}
	} else {
		s.logger.Info("existing google user found",
			"user_id", user.ID.String(),
			"email", email,
		)
	}

	accessToken, err := security.GenerateAccessToken(user.ID.String())
	if err != nil {
		s.logger.Error("failed to generate access token for google user",
			"user_id", user.ID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	refreshToken, err := security.GenerateRefreshToken(user.ID.String())
	if err != nil {
		s.logger.Error("failed to generate refresh token for google user",
			"user_id", user.ID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	s.logger.Info("get or create google user completed successfully",
		"user_id", user.ID.String(),
		"email", email,
	)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (db.GetUserByIDRow, error) {
	s.logger.Info("get user by id started",
		"user_id", userID.String(),
	)

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user by id",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return user, err
	}

	s.logger.Info("get user by id completed successfully",
		"user_id", userID.String(),
	)
	return user, nil
}