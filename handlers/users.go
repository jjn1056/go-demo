package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jjn1056/go-demo/validators"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetUsers(c *gin.Context, db *sql.DB) {
	var users []User
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
	} else {
		c.JSON(http.StatusOK, users)
	}
}

func CreateUser(c *gin.Context, db *sql.DB) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var errors []string

	// Validate user data
	if user.Name == "" {
		errors = append(errors, "Name field cannot be empty")
	}

	if len(user.Name) < 2 {
		errors = append(errors, "Name must be at least 2 characters long")
	} else if len(user.Name) > 25 {
		errors = append(errors, "Name cannot exceed 25 characters")
	}

	if user.Email == "" {
		errors = append(errors, "Email field cannot be empty")
	} else if !validators.IsValidEmail(user.Email) {
		errors = append(errors, "Invalid email format")
	}

	// If there are validation errors, return them in the response
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Insert user data into the database
	result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	lastInsertedID, _ := result.LastInsertId()
	user.ID = int(lastInsertedID)

	c.JSON(http.StatusCreated, user)
}
