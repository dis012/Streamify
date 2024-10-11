-- name: RegisterUser :one
INSERT INTO users (email, hashed_password, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING *;