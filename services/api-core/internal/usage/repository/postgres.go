package repository

import (
	"context"

	"org/api-core/internal/db"

	"github.com/google/uuid"
)

type postgresRepository struct {
	q *db.Queries
}

func NewPostgresRepository(q *db.Queries) Repository {
	return &postgresRepository{q: q}
}

func (r *postgresRepository) GetMonthlyUsage(
	ctx context.Context,
	userID uuid.UUID,
	month string,
) (db.WorkflowUsage, error) {
	return r.q.GetMonthlyUsage(ctx, db.GetMonthlyUsageParams{
		UserID: userID,
		Month: month,
	})
}

func (r *postgresRepository) IncrementWorkflowUsage(
	ctx context.Context,
	userID uuid.UUID,
	month string,
) (db.WorkflowUsage, error) {
	return r.q.IncrementWorkflowUsage(ctx, db.IncrementWorkflowUsageParams{
		ID: uuid.New(),
		UserID: userID,
		Month: month,
	})
}