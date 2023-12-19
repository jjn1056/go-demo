package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// TestCreateUser tests the user creation handler.
func TestCreateUser(t *testing.T) {
	// Set up a temporary database for testing
	db, dbErr := setupTestDatabase()
	if dbErr != nil {
		t.Fatalf("Failed to set up test database: %v", dbErr)
	}
	defer func() {
		dbErr = cleanupTestDatabase(db)
		if dbErr != nil {
			t.Fatalf("Failed to clean up test database: %v", dbErr)
		}
	}()

	// Create a new Gin router
	router := gin.Default()

	// Set up the routes for testing
	router.POST("/users", func(c *gin.Context) {
		CreateUser(c, db)
	})

	// Create a test user
	user := User{Name: "Test User", Email: "test@example.com"}
	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal test user data: %v", err)
	}

	// Create a request for the user creation
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Check the response status code
	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, recorder.Code)
	}

	// Parse the response body and check the user data
	var responseUser User
	err = json.Unmarshal(recorder.Body.Bytes(), &responseUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check the user data from the response
	if responseUser.Name != user.Name || responseUser.Email != user.Email {
		t.Errorf("Expected user data: %+v, but got: %+v", user, responseUser)
	}
}

// setupTestDatabase sets up a temporary database for testing.
func setupTestDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test_data.db")
	if err != nil {
		return nil, err
	}

	// Create the "users" table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// cleanupTestDatabase removes the temporary database after the test is done.
func cleanupTestDatabase(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		return err
	}

	// Remove the temporary database file after the test is done.
	return os.Remove("test_data.db")
}
