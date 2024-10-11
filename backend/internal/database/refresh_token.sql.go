// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: refresh_token.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_token (token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4
)
RETURNING token, created_at, updated_at, expires_at, revoked_at, user_id
`

type CreateRefreshTokenParams struct {
	Token     string
	ExpiresAt time.Time
	RevokedAt sql.NullTime
	UserID    int32
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken,
		arg.Token,
		arg.ExpiresAt,
		arg.RevokedAt,
		arg.UserID,
	)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.RevokedAt,
		&i.UserID,
	)
	return i, err
}
