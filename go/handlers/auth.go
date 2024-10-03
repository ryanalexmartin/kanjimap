package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ryanalexmartin/kanjimap/db"
	"github.com/ryanalexmartin/kanjimap/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received registration request")

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	fmt.Printf("Registration attempt for username: %s, email: %s\n", username, email)

	var userExists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=?)", username).Scan(&userExists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Database error", err)
		return
	}
	if userExists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		fmt.Println("Username already exists")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	result, err := db.DB.Exec("INSERT INTO users (username, password, email, token) VALUES (?, ?, ?, ?)", username, hashedPassword, email, "")
	if err != nil {
		http.Error(w, "Unable to register user", http.StatusInternalServerError)
		fmt.Println("Unable to register user", err)
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Unable to get user ID", http.StatusInternalServerError)
		fmt.Println("Unable to get user ID", err)
		return
	}
	w.Write([]byte(fmt.Sprintf("User successfully registered with ID: %d", userID)))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if db.DB == nil {
		log.Println("Database connection is nil")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	var user models.User
	var tokenString string

	log.Printf("Login attempt for user: %s", username)

	row := db.DB.QueryRow("SELECT id, username, password FROM users WHERE username=?", username)

	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		log.Printf("Invalid username for: %s, error: %v", username, err)
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Invalid password for user: %s", username)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err = token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Printf("Unable to sign token for user: %s, error: %v", username, err)
		http.Error(w, "Unable to sign token", http.StatusInternalServerError)
		return
	}

	// Store the new token
	result, err := db.DB.Exec("INSERT INTO user_tokens (user_id, token) VALUES (?, ?)", user.ID, tokenString)
	if err != nil {
		log.Printf("Failed to store token for user: %s, error: %v", username, err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Token stored for user: %s, rows affected: %d", username, rowsAffected)

	// Return the token
	type JsonResponse struct {
		Token string `json:"token"`
	}
	jsonResponse := JsonResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)

	log.Printf("User %s logged in successfully", username)
}
