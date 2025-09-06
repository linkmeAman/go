package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func NewConnection() (*sql.DB, error) {
	// For local development without Docker, we'll use SQLite
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=saas_billing sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return db, nil
}
