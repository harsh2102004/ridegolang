-- name: CreateUserGeminiKey :one
INSERT INTO user_gemini_keys (
    id,
    user_id,
    api_key,
    created_at,
    updated_at
)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW()
)
RETURNING *;


-- name: GetUserGeminiKeyByUserID :one
SELECT *
FROM user_gemini_keys
WHERE user_id = $1
LIMIT 1;


-- name: UpdateUserGeminiKey :one
UPDATE user_gemini_keys
SET
    api_key = $2,
    updated_at = NOW()
WHERE user_id = $1
RETURNING *;


-- name: DeleteUserGeminiKey :exec
DELETE FROM user_gemini_keys
WHERE user_id = $1;


-- name: UserHasGeminiKey :one
SELECT EXISTS (
    SELECT 1
    FROM user_gemini_keys
    WHERE user_id = $1
);