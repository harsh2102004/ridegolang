package repository

import (
	"context"
	"org/api-core/internal/db"
)

type ConnectedAccountRepository interface {
	UpsertConnectedAccount(ctx context.Context, params db.UpsertConnectedAccountParams) (db.ConnectedAccount, error)
}
