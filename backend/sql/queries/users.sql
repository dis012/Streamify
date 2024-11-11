-- name: RegisterUser :one
INSERT INTO users (email, hashed_password, is_admin, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;