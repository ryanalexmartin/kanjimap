package main

import (
  "database/sql"
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"
  "github.com/joho/godotenv"

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
    var err error
    // set env vars from .env file
    err = godotenv.Load("../vue/.env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // check if the secret key is set
    if os.Getenv("SECRET_KEY") == "" {
        log.Fatal("SECRET_KEY environment variable not set")
    }

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

    // Determine if we're in a local development environment
    isLocalDev := os.Getenv("VUE_APP_URL") == "http://localhost"

    var handler http.Handler
    if isLocalDev {
        // Allow all origins in local development
        c := cors.New(cors.Options{
            AllowedOrigins: []string{"*"},
            AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
            AllowedHeaders: []string{"*"},
        })
        handler = c.Handler(mux)
    } else {
        // Use more restrictive CORS settings in production
        c := cors.New(cors.Options{
            AllowedOrigins: []string{os.Getenv("VUE_APP_URL")},
            AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
            AllowedHeaders: []string{"Authorization", "Content-Type"},
        })
        handler = c.Handler(mux)
    }

    var port int
    port, err = strconv.Atoi(os.Getenv("VUE_APP_PORT"))
    if err != nil {
        log.Fatal("PORT environment variable not set")
    }

    fmt.Printf("Starting application on port %v \n", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received registration request")
    
    username := r.FormValue("username")
    password := r.FormValue("password")
    email := r.FormValue("email")
    
    fmt.Printf("Registration attempt for username: %s, email: %s\n", username, email)

    var userExists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=?)", username).Scan(&userExists)
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
    result, err := db.Exec("INSERT INTO users (username, password, email, token) VALUES (?, ?, ?, ?)", username, hashedPassword, email, "")
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
  username := r.FormValue("username")
  password := r.FormValue("password")
  var user User
  var tokenString string

  row := db.QueryRow("SELECT id, username, password FROM users WHERE username=?", username)

  err := row.Scan(&user.ID, &user.Username, &user.Password)
  if err != nil {
    http.Error(w, "Invalid username", http.StatusUnauthorized)
    fmt.Println("Invalid username", err)
    return
  }

  if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
    http.Error(w, "Invalid password", http.StatusUnauthorized)
    return
  }

  // Create a new token via random generation
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "username": username,
    "exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
  })

  // Sign and get the complete encoded token as a string
  tokenString, err = token.SignedString([]byte(os.Getenv("SECRET_KEY")))

  if err != nil {
    http.Error(w, "Unable to sign token", http.StatusInternalServerError)
    fmt.Println("Unable to sign token", err)
    return
  }

  // Invalidate old token and save the new token in the database
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
  if reqToken == "" || !strings.HasPrefix(reqToken, "Bearer ") {
    http.Error(w, "Invalid token format", http.StatusUnauthorized)
    return
  }
  reqToken = reqToken[7:] // remove "Bearer " from the token

  // Check if jwt matches the one in the database 
  var tokenString string
  err := db.QueryRow("SELECT token FROM users WHERE username=?", username).Scan(&tokenString)
  if err != nil {
    http.Error(w, "Invalid token", http.StatusUnauthorized)
    fmt.Println("Invalid token", err)
    return
  }

  if reqToken != tokenString {
    http.Error(w, "Invalid token", http.StatusUnauthorized)
    fmt.Println("Invalid token: token mismatch")
    return
  }

  // Verify the token
  token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return []byte(os.Getenv("SECRET_KEY")), nil
  })

  if err != nil {
    http.Error(w, "Invalid token", http.StatusUnauthorized)
    fmt.Println("Invalid token:", err)
    return
  }

  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    if exp, ok := claims["exp"].(float64); ok {
      if time.Now().Unix() > int64(exp) {
        http.Error(w, "Token expired", http.StatusUnauthorized)
        fmt.Println("Token expired")
        return
      }
    }
  } else {
    http.Error(w, "Invalid token claims", http.StatusUnauthorized)
    fmt.Println("Invalid token claims")
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
