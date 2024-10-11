-- +goose Up
CREATE TABLE refresh_token (
    token TEXT NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    user_id INT NOT NULL REFERENCES users(id)
);

-- +goose Down
 DROP TABLE refresh_token;