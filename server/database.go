package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

// User represents a user record in the database.
type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

// Database is a struct that holds the database connection.
type Database struct {
	*sql.DB
}

// NewDatabase creates and initializes a new Database connection.
func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Create the users table if it doesn't already exist.
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("could not create users table: %w", err)
	}

	log.Println("Database connection successful and table initialized.")
	return &Database{DB: db}, nil
}

// CreateUser inserts a new user into the database.
func (db *Database) CreateUser(username, passwordHash string) (int64, error) {
	res, err := db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, passwordHash)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}
	return res.LastInsertId()
}

// GetUserByUsername retrieves a user from the database by their username.
func (db *Database) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	row := db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username)
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("could not scan user row: %w", err)
	}
	return user, nil
}
