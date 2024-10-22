-- +goose Up
CREATE TABLE series (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,                     -- Title of the series
    description TEXT,                        -- Optional description of the series
    uploaded_at TIMESTAMP NOT NULL,           -- When the series was uploaded
    user_id INT NOT NULL REFERENCES users(id) -- User who uploaded the series
);

-- +goose Down
DROP TABLE series;