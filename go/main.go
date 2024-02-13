package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@/kanjimap")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Try to ping the database
	err = db.Ping()
	if err != nil {
		log.Fatal("Unable to connect to MySQL server or database does not exist: ", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"}, // Replace with the origin of your Vue app
	}).Handler(mux)

	port := 8081
	fmt.Printf("Starting application on port %v \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received register request")
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Check if username already exists
	var userExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=?)", username).Scan(&userExists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Database error", err)
		return
	}
	if userExists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		fmt.Println("Username already exists", err)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		http.Error(w, "Unable to register user", http.StatusInternalServerError)
		fmt.Println("Unable to register user", err)
		return
	}

	w.Write([]byte("User successfully registered"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received login request")
	username := r.FormValue("username")
	password := r.FormValue("password")

	var user User
	row := db.QueryRow("SELECT * FROM users WHERE username = ?", username)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Error logging in", http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("Logged in successfully"))
}
