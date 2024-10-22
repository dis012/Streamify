-- +goose Up
CREATE TABLE series_episode (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,                      -- Title of the episode
    season INT NOT NULL,                      -- Season number
    episode INT NOT NULL,                     -- Episode number within the season
    uploaded_at TIMESTAMP NOT NULL,           -- When the episode was uploaded
    uploaded_by INT NOT NULL REFERENCES users(id), -- User who uploaded the episode
    series_id INT NOT NULL REFERENCES series(id) ON DELETE CASCADE -- Links episode to the series
);

-- +goose Down
DROP TABLE series_episode;