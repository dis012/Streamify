package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dis012/StreamingServer/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const port = ":8080"

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is required")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("Secret string required for creating JWT tokens")
	}

	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		log.Fatal("Data path is required")
	}

	adminSecret := os.Getenv("ADMIN_SECRET")
	if adminSecret == "" {
		log.Fatal("Admin secret is required")
	}

	maxUploadSize := os.Getenv("MAX_UPLOAD_SIZE")
	if maxUploadSize == "" {
		log.Fatal("Max upload size is required")
	}

	maxUploadSizeInt, err := strconv.ParseInt(maxUploadSize, 10, 64)
	if err != nil {
		log.Fatalf("Error parsing max upload size: %v", err)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	dbQueries := database.New(db)

	apiCnf := &ApiConfig{
		dbQueries:     dbQueries,
		secret:        secret,
		admin_secret:  adminSecret,
		maxUploadSize: maxUploadSizeInt,
		data_path:     dataPath,
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	fileServer := http.StripPrefix("/data", http.FileServer(http.Dir(dataPath)))

	mux.Handle("/data/", fileServer)
	mux.HandleFunc("POST /api/register", apiCnf.RegisterUser)
	mux.HandleFunc("POST /api/login", apiCnf.LoginUser)
	mux.HandleFunc("POST /api/refresh", apiCnf.RefreshAccessToken)
	mux.HandleFunc("POST /api/revoke", apiCnf.RevokeRefreshToken)
	mux.HandleFunc("POST /api/uploadseries", apiCnf.UploadSeries)
	mux.HandleFunc("POST /api/uploadepisode", apiCnf.UploadEpisode)

	log.Printf("Server listening on port %s", port)
	log.Fatal(server.ListenAndServe())
}
