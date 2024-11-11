package main

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

func IsValidVideoFormat(file io.Reader) bool {
	/*
		Checks if the video format is valid
		Returns true if the video format is valid, false otherwise
	*/

	// Create a buffer to read the first 512 bytes
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		return false
	}

	// Detect the MIME type of the file
	mimeType := http.DetectContentType(buf)
	switch mimeType {
	case "video/mp4", "video/x-matroska":
		// Reset the read pointer to the beginning if possible
		if seeker, ok := file.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		}
		return true
	default:
		// Reset the read pointer to the beginning if possible
		if seeker, ok := file.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		}
		return false
	}
}

func IsAllowedExtension(filename string) bool {
	/*
		Checks if the file extension is allowed
	*/
	allowedExtensions := map[string]bool{
		".mp4":  true,
		".mkv":  true,
		".mov":  true,
		".webm": true,
	}

	ext := strings.ToLower(filepath.Ext(filename))
	return allowedExtensions[ext]
}
