-- name: CreateRefreshToken :one
INSERT INTO refresh_token (token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4
)
RETURNING *;