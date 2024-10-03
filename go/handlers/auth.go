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

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if username == "" || email == "" || password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

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
	result, err := db.DB.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", username, hashedPassword, email)
	if err != nil {
		log.Printf("Unable to register user: %v", err)
		http.Error(w, "Unable to register user", http.StatusInternalServerError)
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

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	var user models.User

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

	// Fetch the user's ID from the database
	var userID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Error fetching user data", http.StatusInternalServerError)
		return
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"user_id":  userID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Store the token in the user_tokens table
	_, err = db.DB.Exec("INSERT INTO user_tokens (user_id, token, created_at) VALUES (?, ?, ?)", user.ID, tokenString, time.Now())
	if err != nil {
		log.Printf("Error storing token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Token stored for user: %s", username)

	// Return the token
	type JsonResponse struct {
		Token string `json:"token"`
	}
	jsonResponse := JsonResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)

	log.Printf("User %s logged in successfully", username)
}
