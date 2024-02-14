package main

import (
	"database/sql"
	"encoding/json"
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

// This is used to get the character from its ID
type Character struct {
	ID        string `json:"id"`
	Character string `json:"character"`
}

// This is used to track the learned status of characters for a user
type CharacterCard struct {
	Username    string `json:"username"`
	Character   string `json:"character"`
	Learned     bool   `json:"learned"`
	CharacterID string `json:"characterId"` // "CharacterID" emphasizes that it's a string, not an int
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@/kanjimap")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Unable to connect to MySQL server or database does not exist: ", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/fetch-characters", fetchAllCharactersHandler)
	mux.HandleFunc("/learn-character", learnCharacter)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"}, // Replace with the origin of your Vue app
	}).Handler(mux)

	port := 8081
	fmt.Printf("Starting application on port %v \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
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
	// Create a new userID via an auto-incrementing column
	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
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
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		http.Error(w, "Unable to register user", http.StatusInternalServerError)
		fmt.Println("Unable to register user", err)
		return
	}
	w.Write([]byte("User successfully registered"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
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

func fetchAllCharactersHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Database error", err)
		return
	}
	fmt.Printf("Downloading characters for user %v with id %v\n", r.FormValue("username"), userID)

	rows, err := db.Query("SELECT characters.character_id, user_character_progress.learned FROM characters LEFT JOIN user_character_progress ON characters.character_id = user_character_progress.character_id AND user_character_progress.user_id = ?", userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Database error", err)
		return
	}
	defer rows.Close()

	var characterCards []CharacterCard
	for rows.Next() {
		var card CharacterCard
		var learned sql.NullBool
		err := rows.Scan(&card.CharacterID, &learned)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			fmt.Println("Database error", err)
			return
		}
		if !learned.Valid {
			continue
		}
		var character sql.NullString
		row, err := db.Query("SELECT chinese_character FROM characters WHERE character_id=?", card.CharacterID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			fmt.Println("Database error", err)
			return
		}
		defer row.Close()
		for row.Next() {
			err := row.Scan(&character)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				fmt.Println("Database error", err)
				return
			}
		}
		card.Username = username
		if character.Valid {
			card.Character = character.String
		} else {
			card.Character = "Character not found"
		}

		card.Learned = learned.Bool
		characterCards = append(characterCards, card)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characterCards)
}

func learnCharacter(w http.ResponseWriter, r *http.Request) {
	fmt.Println("learnCharacter")

	var character CharacterCard
	err := json.NewDecoder(r.Body).Decode(&character)
	if err != nil {
		http.Error(w, "Unable to decode JSON request", http.StatusBadRequest)
		fmt.Println("Unable to decode JSON request", err)
		return
	}

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE username=?", character.Username).Scan(&userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Database error", err)
		return
	}
	fmt.Println("UserID", character.CharacterID)

	var learned = character.Learned
	err = db.QueryRow("SELECT learned FROM user_character_progress WHERE character_id=? AND user_id=?", character.CharacterID, userID).Scan(&learned)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert a new row for the character
			_, err = db.Exec("INSERT INTO user_character_progress (user_id, character_id, learned) VALUES (?, ?, ?)", userID, character.CharacterID, character.Learned)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				fmt.Println("Database error", err)
				return
			}
			fmt.Println("Inserted new row for character", character.CharacterID)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
			fmt.Println("Database error", err)
			return
		}
	}

	fmt.Println("updating value of learned", character.Learned)

	_, err = db.Exec("UPDATE user_character_progress SET learned=? WHERE user_id=? AND character_id=?", character.Learned, userID, character.CharacterID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Database error", err)
		return
	}

	// Return the updated character card with the learned status from the database
	row := db.QueryRow("SELECT * FROM user_character_progress WHERE user_id=? AND character_id=?", userID, character.CharacterID)
	var card CharacterCard
	if err := row.Scan(&card.CharacterID, &card.Username, &card.Learned); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Database error", err)
		return
	}
	type JsonResponse struct {
		UserID      int    `json:"userId"`
		CharacterID string `json:"characterId"`
		Learned     bool   `json:"learned"`
	}
	jsonResponse := JsonResponse{
		UserID:      userID,
		CharacterID: character.CharacterID,
		Learned:     card.Learned,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)
}
