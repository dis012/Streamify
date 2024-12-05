package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dis012/StreamingServer/internal/database"
)

type Movie struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (a *ApiConfig) UploadMovie(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for uploading a movie, it accepts movie Title and description and saves it to the database
		It also accepts movie file and saves it to the disk
	*/

	type param struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Email       string `json:"email"`
	}

	var movie param

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.dbQueries.GetUserByEmail(r.Context(), movie.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if !user.IsAdmin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	allMovies, err := a.dbQueries.GetAllMovies(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, m := range allMovies {
		if m.Title == movie.Title {
			http.Error(w, "Movie already exists", http.StatusBadRequest)
			return
		}
	}

	err = a.dbQueries.UploadMovie(r.Context(), database.UploadMovieParams{
		Title:       movie.Title,
		Description: sql.NullString{String: movie.Description, Valid: movie.Description != ""},
		UserID:      user.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(a.maxUploadSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["video"]
	file, err := files[0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	IsAllowedEx := IsAllowedExtension(files[0].Filename)
	if !IsAllowedEx {
		http.Error(w, "File extension not allowed", http.StatusBadRequest)
		return
	}

	if files[0].Size > a.maxUploadSize {
		http.Error(w, "File size too large", http.StatusBadRequest)
		return
	}

	// Define the file path
	filePath := a.data_path + "movies/" + movie.Title + "/" + files[0].Filename

	// Ensure the directory exists
	dir := a.data_path + "movies/" + movie.Title + "/"
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		http.Error(w, "Unable to create directories for file path", http.StatusInternalServerError)
		log.Println("Directory creation error:", err)
		return
	}

	movieId, err := a.dbQueries.GetMovieById(r.Context(), movie.Title)

	err = a.dbQueries.AddMoviePath(r.Context(), database.AddMoviePathParams{
		ID: movieId,
		MoviePath: sql.NullString{
			String: filePath,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error adding movie path:", err)
		return
	}

	// Save the file
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to create the file for writing", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}

	RespondWithJson(w, http.StatusCreated, Movie{
		Title:       movie.Title,
		Description: movie.Description,
	})
}
