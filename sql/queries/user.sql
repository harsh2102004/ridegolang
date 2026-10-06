-- name: CreateUser :one
INSERT INTO users (id, email, password, is_verified)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: VerifyUser :exec
UPDATE users
SET is_verified = true
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $1
WHERE id = $2;

-- name: CreateOAuthUser :one
INSERT INTO users (id, email, password, is_verified)
VALUES ($1, $2, '', true)
RETURNING id, email, password, is_verified, plan, created_at;


-- name: GetUserByID :one
SELECT id, email, plan, is_verified, created_at
FROM users
WHERE id = $1;