package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dis012/StreamingServer/internal/database"
)

type Series struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Episode struct {
	Title   []string `json:"title"`
	Season  int32    `json:"season"`
	Episode []int32  `json:"episode"`
}

func (a *ApiConfig) UploadSeries(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for uploading a series
		It acceps series Title and description and saves it to the database
	*/

	type param struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Email       string `json:"email"`
	}

	var series param

	err := json.NewDecoder(r.Body).Decode(&series)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.dbQueries.GetUserByEmail(r.Context(), series.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if !user.IsAdmin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	allSeries, err := a.dbQueries.GetAllSeries(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, serie := range allSeries {
		if series.Title == serie.Title {
			http.Error(w, "Series already exists", http.StatusBadRequest)
			return
		}
	}

	err = a.dbQueries.UploadSeries(r.Context(), database.UploadSeriesParams{
		Title:       series.Title,
		Description: sql.NullString{String: series.Description, Valid: series.Description != ""},
		UserID:      user.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondWithJson(w, http.StatusCreated, Series{
		Title:       series.Title,
		Description: series.Description,
	})
}

func (a *ApiConfig) UploadEpisode(w http.ResponseWriter, r *http.Request) {
	/*
		Handler for uploading an episode
		It accepts episode title, season, episode number, series id and saves it to the database
		It uploads all the files to the data path
	*/

	type param struct {
		Titles      []string `json:"titles"`
		Season      int32    `json:"season"`
		Episodes    []int32  `json:"episode"`
		SeriesTitle string   `json:"series_title"`
		Email       string   `json:"email"`
	}

	err := r.ParseMultipartForm(a.maxUploadSize)
	if err != nil {
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the metadata from the form
	metadataStr := r.FormValue("metadata")
	if metadataStr == "" {
		http.Error(w, "Missing metadata in form", http.StatusBadRequest)
		return
	}

	// Parse the JSON metadata
	var episode param
	err = json.Unmarshal([]byte(metadataStr), &episode)
	if err != nil {
		http.Error(w, "Error parsing metadata JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(episode.Titles) != len(episode.Episodes) {
		http.Error(w, "Titles and Episodes must be of same length", http.StatusBadRequest)
		return
	}

	user, err := a.dbQueries.GetUserByEmail(r.Context(), episode.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if !user.IsAdmin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	seriesId, err := a.dbQueries.GetSeriesByTitle(r.Context(), episode.SeriesTitle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for i, title := range episode.Titles {
		err = a.dbQueries.UploadEpisode(r.Context(), database.UploadEpisodeParams{
			Title:      title,
			Season:     episode.Season,
			Episode:    episode.Episodes[i],
			UploadedBy: user.ID,
			SeriesID:   seriesId.ID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Parse the multipart form to handle multiple file uploads
	err = r.ParseMultipartForm(a.maxUploadSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the files with the field name "video" (adjust as needed)
	files := r.MultipartForm.File["video"]
	if len(files) != len(episode.Titles) {
		http.Error(w, "The number of video files must match the number of titles", http.StatusBadRequest)
		return
	}

	// Loop through each uploaded file
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Unable to open file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		/*
			isValidMime := IsValidVideoFormat(file)
			if !isValidMime {
				http.Error(w, "Invalid video format", http.StatusBadRequest)
				return
			}
		*/

		isAllowedExt := IsAllowedExtension(fileHeader.Filename)
		if !isAllowedExt {
			http.Error(w, "Invalid file extension", http.StatusBadRequest)
			return
		}

		// Check file size
		if fileHeader.Size > a.maxUploadSize {
			http.Error(w, "File size too large", http.StatusBadRequest)
			return
		}

		// Define the file path
		filePath := a.data_path + "series/" + episode.SeriesTitle + "/" + strconv.Itoa(int(episode.Season)) + "/" + fileHeader.Filename

		// Ensure the directory exists
		dir := a.data_path + "series/" + episode.SeriesTitle + "/" + strconv.Itoa(int(episode.Season)) + "/"
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			http.Error(w, "Unable to create directories for file path", http.StatusInternalServerError)
			log.Println("Directory creation error:", err)
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
	}

	RespondWithJson(w, http.StatusCreated, Episode{
		Title:   episode.Titles,
		Season:  episode.Season,
		Episode: episode.Episodes,
	})
}
