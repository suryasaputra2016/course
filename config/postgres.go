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
	queryUserTable := `
		CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT,
		role VARCHAR(15)
		);`
	_, err := db.Exec(queryUserTable)
	if err != nil {
		return fmt.Errorf("creating users table: %w", err)
	}

	querySessionTable := `
		CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		user_id INT UNIQUE NOT NULL REFERENCES users(id),
		token_hash TEXT
		);`
	_, err = db.Exec(querySessionTable)
	if err != nil {
		return fmt.Errorf("creating sessions table: %w", err)
	}
	return nil
}
