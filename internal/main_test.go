package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// TestCreateUser tests the user creation handler.
func TestCreateUser(t *testing.T) {
	// Set up a temporary database for testing
	db, err := sql.Open("sqlite3", "test_data.db")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func() {
		db.Close()
		// Remove the temporary database after the test is done
		if err := removeTempDatabase("test_data.db"); err != nil {
			t.Fatalf("Failed to remove temporary database: %v", err)
		}
	}()

	// Create a new Gin router
	router := gin.Default()

	// Set up the routes for testing
	router.POST("/users", createUser)

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

	// Set the request content type
	req.Header.Set("Content-Type", "application/json")

	// Create a test HTTP recorder
	recorder := httptest.NewRecorder()

	// Serve the request to the router
	router.ServeHTTP(recorder, req)

	// Check the response status code
	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, recorder.Code)
	}

	fmt.Println(recorder.Body.String())

	// Parse the response body and check the user data
	var responseUser User
	err = json.Unmarshal(recorder.Body.Bytes(), &responseUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check the user data from the response
	assert.Equal(t, user.Name, responseUser.Name)
	assert.Equal(t, user.Email, responseUser.Email)
}

// removeTempDatabase removes the temporary database file after the test is done.
func removeTempDatabase(dbPath string) error {
	return os.Remove(dbPath)
}
