-- name: UpsertConnectedAccount :one
INSERT INTO connected_accounts (
    user_id,
    provider,
    provider_user_id,
    email,
    access_token,
    refresh_token,
    expires_at,
    scopes
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (user_id, provider)
DO UPDATE SET
    provider_user_id = EXCLUDED.provider_user_id,
    email = EXCLUDED.email,
    access_token = EXCLUDED.access_token,
    refresh_token = COALESCE(EXCLUDED.refresh_token, connected_accounts.refresh_token),
    expires_at = EXCLUDED.expires_at,
    scopes = EXCLUDED.scopes,
    updated_at = now()
RETURNING *;