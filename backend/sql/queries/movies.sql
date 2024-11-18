-- name: UploadMovie :exec
INSERT INTO movies (title, description, uploaded_at, user_id)
VALUES ($1, $2, NOW(), $3);

-- name: GetAllMovies :many
SELECT id, title, description
FROM movies;

-- name: GetMovieByTitle :one
SELECT id, description
FROM movies
WHERE title = $1;