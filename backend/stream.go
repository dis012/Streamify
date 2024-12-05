package main

import (
	"net/http"
	"strconv"
)

func (a *ApiConfig) StreamEpisodeRequestMiddleware() http.HandlerFunc {
	/*
		Method of ApiConfig struct that will handle the streaming of a video.
		It accepts the video's name and returns the video file.
	*/
	return func(w http.ResponseWriter, r *http.Request) {
		// Get series id
		seriesID, err := strconv.Atoi(r.URL.Query().Get("series_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get series path
		seriesPathResult, err := a.dbQueries.GetSeriesPath(r.Context(), int32(seriesID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !seriesPathResult.Valid {
			http.Error(w, "Invalid series path", http.StatusInternalServerError)
			return
		}
		seriesPath := seriesPathResult.String

		// Add path to the file server
		fileServer := http.StripPrefix("/api/stream/series/{episode_id}", http.FileServer(http.Dir(seriesPath)))

		// Serve the file
		fileServer.ServeHTTP(w, r)
	}
}

func (a *ApiConfig) StreamMovieRequestMiddleware() http.HandlerFunc {
	/*
		Method of ApiConfig struct that will handle the streaming of a video.
		It accepts the video's name and returns the video file.
	*/
	return func(w http.ResponseWriter, r *http.Request) {
		// Get series id
		movieID, err := strconv.Atoi(r.URL.Query().Get("movie_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get series path
		moviePathResult, err := a.dbQueries.GetMoviePath(r.Context(), int32(movieID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !moviePathResult.Valid {
			http.Error(w, "Invalid series path", http.StatusInternalServerError)
			return
		}
		moviePath := moviePathResult.String

		// Add path to the file server
		fileServer := http.StripPrefix("/api/stream/movie/{movie_id}", http.FileServer(http.Dir(moviePath)))

		// Serve the file
		fileServer.ServeHTTP(w, r)
	}
}
