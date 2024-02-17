package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
    "os"

    _ "github.com/go-sql-driver/mysql"

    "github.com/dgrijalva/jwt-go"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

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


func main() {
    // Get the JWT key from the environment
    jwtKey := os.Getenv("JWT_KEY")
    if jwtKey == "" {
        log.Fatal("JWT_KEY environment variable not set")
    }

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

    fs := http.FileServer(http.Dir("../vue/dist"))

	mux := http.NewServeMux()
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/fetch-characters", fetchAllCharactersHandler)
	mux.HandleFunc("/learn-character", learnCharacter)

    // serve the frontend
    mux.Handle("/", fs)

    // allow all origins
    handler := cors.AllowAll().Handler(mux)

	port := 8081
	fmt.Printf("Starting application on port %v \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
    email := r.FormValue("email")
    tokenString := ""

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
    result, err := db.Exec("INSERT INTO users (username, password, email, token) VALUES (?, ?, ?, ?)", username, hashedPassword, email, tokenString)
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
    email := ""
	var user User
    var tokenString string
	row := db.QueryRow("SELECT * FROM users WHERE username = ?", username)
    err := row.Scan(&user.ID, &user.Username, &user.Password, &email, &tokenString)
    if err != nil {
        http.Error(w, "Invalid username", http.StatusUnauthorized)
        fmt.Println("Invalid username", err)
        return
    }

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

    // Create a token
    token := jwt.New(jwt.SigningMethodHS256)
    tokenString, err = token.SignedString([]byte(os.Getenv("JWT_KEY")))
    if err != nil {
        http.Error(w, "Unable to sign token", http.StatusInternalServerError)
        fmt.Println("Unable to sign token", err)
        return
    }
    // save the token in the database
    _, err = db.Exec("UPDATE users SET token=? WHERE username=?", tokenString, username)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        fmt.Println("Database error", err)
        return
    }

    // Return the token
    type JsonResponse struct {
        Token string `json:"token"`
    }
    jsonResponse := JsonResponse{Token: tokenString}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(jsonResponse)

    fmt.Printf("User %v logged in successfully\n", username)
}

func fetchAllCharactersHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    reqToken := r.Header.Get("Authorization")

    token, err := jwt.Parse(reqToken[7:], func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("JWT_KEY")), nil
    })
    if err != nil {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        fmt.Println("Invalid token", err)
        return
    }
    if !token.Valid {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        fmt.Println("Invalid token", err)
        return
    }

    // get user id
	var userID int
    err = db.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Database error", err)
		return
	}
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
