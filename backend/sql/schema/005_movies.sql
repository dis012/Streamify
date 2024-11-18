-- +goose Up
CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,                -- Title of the movie
    description TEXT,                          -- Optional description of the movie
    uploaded_at TIMESTAMP NOT NULL,            -- When the movie was uploaded
    user_id INT NOT NULL REFERENCES users(id)  -- User who uploaded the movie
);

-- +goose Down
DROP TABLE movies;