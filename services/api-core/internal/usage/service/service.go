package service

import (
	"context"
	"log/slog"
	"time"

	usageRepo "org/api-core/internal/usage/repository"

	"github.com/google/uuid"
	"org/api-core/internal/db"
)

type Service struct {
	repo   usageRepo.Repository
	logger *slog.Logger
}

func NewService(repo usageRepo.Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func currentMonth() string {
	return time.Now().Format("2006-01")
}

func (s *Service) CanRunWorkflow(
	ctx context.Context,
	userID uuid.UUID,
	plan string,
) (bool, string, error) {
	s.logger.Info("checking workflow run eligibility",
		"user_id", userID.String(),
		"plan", plan,
	)

	month := currentMonth()

	usage, err := s.repo.GetMonthlyUsage(ctx, userID, month)
	if err != nil {
		// no usage record yet, allow
		s.logger.Info("no existing usage record, allowing workflow",
			"user_id", userID.String(),
		)
		return true, "", nil
	}

	limit := 100
	if plan != "free" {
		limit = 10000
	}

	runs := int32(0)
	if usage.WorkflowRuns.Valid {
		runs = usage.WorkflowRuns.Int32
	}

	if runs >= int32(limit) {
		s.logger.Warn("workflow limit reached",
			"user_id", userID.String(),
			"runs", runs,
			"limit", limit,
		)
		return false, "monthly workflow run limit reached", nil
	}

	s.logger.Info("workflow allowed",
		"user_id", userID.String(),
		"runs", runs,
		"limit", limit,
	)
	return true, "", nil
}

func (s *Service) IncrementWorkflowRun(
	ctx context.Context,
	userID uuid.UUID,
) error {
	s.logger.Info("incrementing workflow run",
		"user_id", userID.String(),
	)

	_, err := s.repo.IncrementWorkflowUsage(
		ctx,
		userID,
		currentMonth(),
	)

	if err != nil {
		s.logger.Error("failed to increment workflow run",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return err
	}

	s.logger.Info("workflow run incremented successfully",
		"user_id", userID.String(),
	)
	return nil
}

func (s *Service) CurrentMonth() string {
	return currentMonth()
}

func (s *Service) GetCurrentMonthUsage(
	ctx context.Context,
	userID uuid.UUID,
) (db.WorkflowUsage, error) {
	s.logger.Info("getting current month usage",
		"user_id", userID.String(),
	)

	usage, err := s.repo.GetMonthlyUsage(
		ctx,
		userID,
		currentMonth(),
	)

	if err != nil {
		s.logger.Warn("failed to get current month usage (may not exist)",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return usage, err
	}

	s.logger.Info("current month usage retrieved",
		"user_id", userID.String(),
		"workflow_runs", usage.WorkflowRuns,
	)
	return usage, nil
}