package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// User represents a user in the database.
type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

// Database is a wrapper around the sql.DB connection.
type Database struct {
	*sql.DB
}

// NewDatabase initializes and returns a new database connection.
func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to the database successfully.")

	// Create the users table if it doesn't exist.
	createTableSQL := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL UNIQUE,
            password_hash TEXT NOT NULL
        );
    `
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("could not create users table: %w", err)
	}

	return &Database{DB: db}, nil
}

// CreateUser inserts a new user into the database.
func (db *Database) CreateUser(username, passwordHash string) (int64, error) {
	res, err := db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, passwordHash)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}
	return id, nil
}

// GetUserByUsername retrieves a user by their username.
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

// GetAllUsernames retrieves all usernames from the database.
func (db *Database) GetAllUsernames() ([]string, error) {
	rows, err := db.Query("SELECT username FROM users ORDER BY username")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("failed to scan username: %w", err)
		}
		usernames = append(usernames, username)
	}

	return usernames, nil
}
