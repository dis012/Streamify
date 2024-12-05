-- +goose Up
ALTER TABLE series_episode
ADD COLUMN series_path TEXT;

-- +goose Down
ALTER TABLE series_episode
DROP COLUMN series_path;