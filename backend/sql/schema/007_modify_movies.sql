-- +goose Up
ALTER TABLE movies
ADD COLUMN movie_path TEXT;

-- +goose Down
ALTER TABLE movies
DROP COLUMN movie_path;