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
  "path/filepath"

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
  Username            string  `json:"username"`
  Character           string  `json:"character"`
  Learned             bool    `json:"learned"`
  CharacterID         string  `json:"characterId"`
  Frequency           int     `json:"frequency"`
  CumulativeFrequency float64 `json:"cumulativeFrequency"`
  Pinyin              string  `json:"pinyin"`
  English             string  `json:"english"`
}

func main() {
  dbHost := os.Getenv("DB_HOST")
  dbUser := os.Getenv("DB_USER")
  dbPassword := os.Getenv("DB_PASSWORD")
  dbName := os.Getenv("DB_NAME")
  // secretKey := os.Getenv("SECRET_KEY") // todo - convert this to variable

  dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
  var err error
  db, err = sql.Open("mysql", dsn)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    log.Fatal("Unable to connect to MySQL server or database does not exist: ", err)
  }

  // Test the database connection
  err = db.Ping()
  if err != nil {
    log.Fatalf("Unable to connect to MySQL server: %v", err)
  }
  log.Println("Successfully connected to the database")

  log.Println("Now let's try to initialize our frontend.")

  // Check if the frontend directory exists
  frontendDir := "/vue/dist"
  if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
    log.Fatalf("Frontend directory does not exist: %v", err)
  }

  // Check if the index.html file exists
  indexPath := filepath.Join(frontendDir, "index.html")
  if _, err := os.Stat(indexPath); os.IsNotExist(err) {
    log.Fatalf("index.html not found in frontend directory: %v", err)
  }

  // Create the file server
  fs := http.FileServer(http.Dir(frontendDir))

  // Wrap the file server with a handler that logs requests
  loggedFs := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Printf("Serving file: %s", r.URL.Path)
    fs.ServeHTTP(w, r)
  })

  // Determine if we're in a local development environment
  isLocalDev := os.Getenv("VUE_APP_API_URL") == "http://localhost"

  mux := http.NewServeMux()
  mux.Handle("/", loggedFs)
  mux.HandleFunc("/register", registerHandler)
  mux.HandleFunc("/login", loginHandler)
  mux.HandleFunc("/fetch-characters", fetchAllCharactersHandler)
  mux.HandleFunc("/learn-character", learnCharacter)
  mux.HandleFunc("/learned-characters", learnedCharactersHandler)


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
  port, err = strconv.Atoi(os.Getenv("VUE_APP_API_PORT"))
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

    log.Printf("Login attempt for user: %s", username)

    row := db.QueryRow("SELECT id, username, password FROM users WHERE username=?", username)

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
    result, err := db.Exec("INSERT INTO user_tokens (user_id, token) VALUES (?, ?)", user.ID, tokenString)
    if err != nil {
        log.Printf("Failed to store token for user: %s, error: %v", username, err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("Token stored for user: %s, rows affected: %d", username, rowsAffected)

    // Rotate tokens (keep only the 5 most recent tokens)
    result, err = db.Exec(`
        DELETE FROM user_tokens 
        WHERE user_id = ? AND id NOT IN (
            SELECT id FROM (
                SELECT id FROM user_tokens 
                WHERE user_id = ? 
                ORDER BY created_at DESC 
                LIMIT 5
            ) AS t
        )
    `, user.ID, user.ID)
    if err != nil {
        log.Printf("Error rotating tokens for user: %s, error: %v", username, err)
    } else {
        rowsAffected, _ := result.RowsAffected()
        log.Printf("Token rotation completed for user: %s, rows deleted: %d", username, rowsAffected)
    }

    // Return the token
    type JsonResponse struct {
        Token string `json:"token"`
    }
    jsonResponse := JsonResponse{Token: tokenString}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(jsonResponse)

    log.Printf("User %s logged in successfully", username)
}




func validateToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("SECRET_KEY")), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        // Check if the token exists in the database
        var count int
        err := db.QueryRow("SELECT COUNT(*) FROM user_tokens WHERE token = ?", tokenString).Scan(&count)
        if err != nil {
            return nil, err
        }
        if count == 0 {
            return nil, fmt.Errorf("token not found in database")
        }

        // Check expiration
        if exp, ok := claims["exp"].(float64); ok {
            if time.Now().Unix() > int64(exp) {
                return nil, fmt.Errorf("token expired")
            }
        }

        return token, nil
    }

    return nil, fmt.Errorf("invalid token")
}

func learnedCharactersHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("learnedCharactersHandler called")

    // Enable CORS for this endpoint
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

    // Handle preflight request
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    // Ensure the request method is GET
    if r.Method != "GET" {
        log.Printf("Invalid method: %s", r.Method)
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    username := r.URL.Query().Get("username")
    log.Printf("Received request for username: %s", username)

    reqToken := r.Header.Get("Authorization")
    if reqToken == "" || !strings.HasPrefix(reqToken, "Bearer ") {
        http.Error(w, "Invalid token format", http.StatusUnauthorized)
        return
    }
    reqToken = strings.TrimPrefix(reqToken, "Bearer ")

    // Validate the token
    token, err := validateToken(reqToken)
    if err != nil {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        log.Printf("Token validation failed: %v", err)
        return
    }

    // Extract username from token
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        http.Error(w, "Invalid token claims", http.StatusUnauthorized)
        return
    }
    tokenUsername, ok := claims["username"].(string)
    if !ok || tokenUsername != username {
        http.Error(w, "Username mismatch", http.StatusUnauthorized)
        return
    }

    // Get user ID
    var userID int
    err = db.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&userID)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        log.Printf("Database error: %v", err)
        return
    }

    // Fetch learned characters
    rows, err := db.Query(`
    SELECT c.chinese_character
    FROM user_character_progress ucp
    JOIN characters c ON ucp.character_id = c.character_id
    WHERE ucp.user_id = ? AND ucp.learned = true
    `, userID)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        log.Printf("Database error: %v", err)
        return
    }
    defer rows.Close()

    var learnedCharacters []string
    for rows.Next() {
        var char string
        if err := rows.Scan(&char); err != nil {
            http.Error(w, "Database error", http.StatusInternalServerError)
            log.Printf("Database error: %v", err)
            return
        }
        learnedCharacters = append(learnedCharacters, char)
    }

    log.Printf("Fetched %d learned characters for user %s", len(learnedCharacters), username)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(learnedCharacters)
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

    // Validate the token
    token, err := validateToken(reqToken)
    if err != nil {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        log.Printf("Token validation failed: %v", err)
        return
    }

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

  // Update the query to use the new schema
  rows, err := db.Query(`
  SELECT c.character_id, ucp.learned, c.chinese_character, cm.frequency, cm.cumulative_frequency, cm.pinyin, cm.english
  FROM characters c
  LEFT JOIN user_character_progress ucp ON c.character_id = ucp.character_id AND ucp.user_id = ?
  LEFT JOIN character_metadata cm ON c.chinese_character = cm.chinese_character
  `, userID)
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
    var character sql.NullString
    var frequency sql.NullInt64
    var cumulativeFrequency sql.NullFloat64
    var pinyin, english sql.NullString

    err := rows.Scan(&card.CharacterID, &learned, &character, &frequency, &cumulativeFrequency, &pinyin, &english)
    if err != nil {
      http.Error(w, "Database error", http.StatusInternalServerError)
      fmt.Println("Database error", err)
      return
    }

    card.Username = username
    if character.Valid {
      card.Character = character.String
    } else {
      card.Character = "Character not found"
    }
    card.Learned = learned.Bool
    if frequency.Valid {
      card.Frequency = int(frequency.Int64)
    }
    if cumulativeFrequency.Valid {
      card.CumulativeFrequency = cumulativeFrequency.Float64
    }
    if pinyin.Valid {
      card.Pinyin = pinyin.String
    }
    if english.Valid {
      card.English = english.String
    }

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
