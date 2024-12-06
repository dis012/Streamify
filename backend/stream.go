package main

import (
	"log"
	"net/http"
	"os"
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
		if err != nil || !seriesPathResult.Valid {
			http.Error(w, "Series not found", http.StatusNotFound)
			return
		}
		seriesPath := seriesPathResult.String

		fileInfo, err := os.Stat(seriesPath)
		if err != nil {
			http.Error(w, "Error accessing series path", http.StatusInternalServerError)
			log.Println("Error accessing series path:", err)
			return
		}

		if fileInfo.IsDir() {
			// Serve the directory using http.FileServer
			http.StripPrefix("/api/stream/series/{episode_id}", http.FileServer(http.Dir(seriesPath))).ServeHTTP(w, r)
		} else {
			// Serve the file directly
			http.ServeFile(w, r, seriesPath)
		}
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
		if err != nil || !moviePathResult.Valid {
			http.Error(w, "Series not found", http.StatusNotFound)
			return
		}
		moviePath := moviePathResult.String

		fileInfo, err := os.Stat(moviePath)
		if err != nil {
			http.Error(w, "Error accessing series path", http.StatusInternalServerError)
			log.Println("Error accessing series path:", err)
			return
		}

		if fileInfo.IsDir() {
			// Serve the directory using http.FileServer
			http.StripPrefix("/api/stream/movie/{movie_id}", http.FileServer(http.Dir(moviePath))).ServeHTTP(w, r)
		} else {
			// Serve the file directly
			http.ServeFile(w, r, moviePath)
		}
	}
}
