package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

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

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	dbQueries := database.New(db)

	apiCnf := &ApiConfig{
		dbQueries: dbQueries,
		secret:    secret,
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

	log.Printf("Server listening on port %s", port)
	log.Fatal(server.ListenAndServe())
}
