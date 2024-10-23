// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: series.sql

package database

import (
	"context"
	"database/sql"
)

const getAllSeasonEpisodes = `-- name: GetAllSeasonEpisodes :many
SELECT title, episode 
FROM series_episode
WHERE series_id = $1 AND season = $2
ORDER BY episode ASC
`

type GetAllSeasonEpisodesParams struct {
	SeriesID int32
	Season   int32
}

type GetAllSeasonEpisodesRow struct {
	Title   string
	Episode int32
}

func (q *Queries) GetAllSeasonEpisodes(ctx context.Context, arg GetAllSeasonEpisodesParams) ([]GetAllSeasonEpisodesRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllSeasonEpisodes, arg.SeriesID, arg.Season)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllSeasonEpisodesRow
	for rows.Next() {
		var i GetAllSeasonEpisodesRow
		if err := rows.Scan(&i.Title, &i.Episode); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllSeasons = `-- name: GetAllSeasons :many
SELECT DISTINCT season
FROM series_episode
WHERE series_id = $1
ORDER BY season ASC
`

func (q *Queries) GetAllSeasons(ctx context.Context, seriesID int32) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, getAllSeasons, seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var season int32
		if err := rows.Scan(&season); err != nil {
			return nil, err
		}
		items = append(items, season)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllSeries = `-- name: GetAllSeries :many
SELECT id, title, description
FROM series
`

type GetAllSeriesRow struct {
	ID          int32
	Title       string
	Description sql.NullString
}

func (q *Queries) GetAllSeries(ctx context.Context) ([]GetAllSeriesRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllSeries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllSeriesRow
	for rows.Next() {
		var i GetAllSeriesRow
		if err := rows.Scan(&i.ID, &i.Title, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSeriesByTitle = `-- name: GetSeriesByTitle :one
SELECT id, title, description
FROM series
WHERE title = $1
`

type GetSeriesByTitleRow struct {
	ID          int32
	Title       string
	Description sql.NullString
}

func (q *Queries) GetSeriesByTitle(ctx context.Context, title string) (GetSeriesByTitleRow, error) {
	row := q.db.QueryRowContext(ctx, getSeriesByTitle, title)
	var i GetSeriesByTitleRow
	err := row.Scan(&i.ID, &i.Title, &i.Description)
	return i, err
}

const uploadEpisode = `-- name: UploadEpisode :exec
INSERT INTO series_episode (title, season, episode, uploaded_at, uploaded_by, series_id)
VALUES ($1, $2, $3, NOW(), $4, $5)
`

type UploadEpisodeParams struct {
	Title      string
	Season     int32
	Episode    int32
	UploadedBy int32
	SeriesID   int32
}

func (q *Queries) UploadEpisode(ctx context.Context, arg UploadEpisodeParams) error {
	_, err := q.db.ExecContext(ctx, uploadEpisode,
		arg.Title,
		arg.Season,
		arg.Episode,
		arg.UploadedBy,
		arg.SeriesID,
	)
	return err
}

const uploadSeries = `-- name: UploadSeries :one
INSERT INTO series (title, description, uploaded_at, user_id)
VALUES ($1, $2, NOW(), $3)
RETURNING id
`

type UploadSeriesParams struct {
	Title       string
	Description sql.NullString
	UserID      int32
}

func (q *Queries) UploadSeries(ctx context.Context, arg UploadSeriesParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, uploadSeries, arg.Title, arg.Description, arg.UserID)
	var id int32
	err := row.Scan(&id)
	return id, err
}