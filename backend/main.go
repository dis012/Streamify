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

	mux.HandleFunc("POST /api/register", apiCnf.RegisterUser)

	log.Printf("Server listening on port %s", port)
	log.Fatal(server.ListenAndServe())
}
