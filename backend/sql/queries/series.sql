-- name: UploadSeries :exec
INSERT INTO series (title, description, uploaded_at, user_id)
VALUES ($1, $2, NOW(), $3);

-- name: UploadEpisode :exec
INSERT INTO series_episode (title, season, episode, uploaded_at, uploaded_by, series_id)
VALUES ($1, $2, $3, NOW(), $4, $5);

-- name: GetAllSeries :many
SELECT id, title, description
FROM series;

-- name: GetSeriesByTitle :one
SELECT id, description
FROM series
WHERE title = $1;

-- name: GetAllSeasons :many
SELECT DISTINCT season
FROM series_episode
WHERE series_id = $1
ORDER BY season ASC;

-- name: GetAllSeasonEpisodes :many
SELECT title, episode 
FROM series_episode
WHERE series_id = $1 AND season = $2
ORDER BY episode ASC;

-- name: GetEpisode :one
SELECT id
FROM series_episode
WHERE title = $1 AND season = $2 AND episode = $3;

-- name: AddSeriesPath :exec
UPDATE series_episode
SET series_path = $1
WHERE id = $2;

-- name: GetSeriesPath :one
SELECT series_path
FROM series_episode
WHERE id = $1;