package repository

import (
	"context"

	"org/api-core/internal/db"

	"github.com/google/uuid"
)

type Repository interface {
	GetMonthlyUsage(ctx context.Context, userID uuid.UUID, month string) (db.WorkflowUsage, error)
	IncrementWorkflowUsage(ctx context.Context, userID uuid.UUID, month string) (db.WorkflowUsage, error)
}