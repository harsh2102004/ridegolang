package repository

import (
	"context"
	"org/api-core/internal/db"
)

type SQLCConnectedAccountRepository struct {
	q *db.Queries
}

func NewSQLCConnectedAccountRepository(q *db.Queries) ConnectedAccountRepository {
	return &SQLCConnectedAccountRepository{q: q}
}

func (r *SQLCConnectedAccountRepository) UpsertConnectedAccount(
	ctx context.Context,
	params db.UpsertConnectedAccountParams,
) (db.ConnectedAccount, error) {
	return r.q.UpsertConnectedAccount(ctx, params)
}
