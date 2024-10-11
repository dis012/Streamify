package main

import (
	"github.com/dis012/StreamingServer/internal/database"
)

type ApiConfig struct {
	dbQueries *database.Queries
	secret    string
}
