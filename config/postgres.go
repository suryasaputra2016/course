package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// connect to postgres database
func ConnectPostgres() (*sql.DB, error) {
	dsn := "host=localhost port=5432 user=baloo password=junglebook dbname=coursedb sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening postgres: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	return db, nil
}

// close postgres database connection
func ClosePostgres(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		return fmt.Errorf("closing postgres: %w", err)
	}
	return nil
}

func PrepareTables(db *sql.DB) error {
	dbString := `
		CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT,
		passwordHash TEXT,
		role VARCHAR(15)
		);`
	_, err := db.Exec(dbString)
	if err != nil {
		return fmt.Errorf("creating users table: %w", err)
	}
	return nil
}
