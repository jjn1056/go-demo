package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	// Connect to the database
	dbConn, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	db = dbConn

	// Create the "users" table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal("Failed to create the table:", err)
	}
}

// GetDB returns the database instance.
func GetDB() *sql.DB {
	return db
}
