package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"database/sql"
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"github.com/ryanalexmartin/kanjimap/db"
	"github.com/ryanalexmartin/kanjimap/handlers"
	"github.com/ryanalexmartin/kanjimap/middleware"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	// Setup
	err := db.Initialize()
	if err != nil {
		panic(err)
	}

	// Run tests
	code := m.Run()

	// Teardown (if needed)

	os.Exit(code)
}

func TestDatabaseInitialization(t *testing.T) {
	err := db.Initialize()
	if err != nil {
		t.Errorf("Database initialization failed: %v", err)
	}
}

func TestHTTPRoutes(t *testing.T) {
	// Ensure database is initialized for these tests
	err := db.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", middleware.LoggedFs(http.Dir("./frontend")))
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/fetch-characters", middleware.AuthMiddleware(handlers.FetchAllCharactersHandler))
	mux.HandleFunc("/learn-character", middleware.AuthMiddleware(handlers.LearnCharacter))
	mux.HandleFunc("/learned-characters", middleware.AuthMiddleware(handlers.LearnedCharactersHandler))

	testCases := []struct {
		name           string
		route          string
		method         string
		expectedStatus int
	}{
		{"Home", "/", "GET", http.StatusFound},                   // Changed from http.StatusOK to http.StatusFound
		{"Register", "/register", "POST", http.StatusBadRequest}, // Changed to expect 400 (Bad Request, no username or password)
		{"Login", "/login", "POST", http.StatusBadRequest},       // Changed to expect 400 (Bad Request, no username or password)
		{"Fetch Characters (Unauthorized)", "/fetch-characters", "GET", http.StatusUnauthorized},
		{"Learn Character (Unauthorized)", "/learn-character", "POST", http.StatusUnauthorized},
		{"Learned Characters (Unauthorized)", "/learned-characters", "GET", http.StatusUnauthorized},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, tc.route, nil)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
				t.Logf("Response body: %s", rr.Body.String()) // Add this line for debugging
			}
		})
	}
}

// Helper function to clean up the test user and associated tokens
func cleanupTestUser(t *testing.T, username string) {
	err := db.Initialize()
	if err != nil {
		t.Logf("Failed to initialize database: %v", err)
		return
	}

	// Start a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		t.Logf("Failed to start transaction: %v", err)
		return
	}
	defer tx.Rollback() // Rollback if not committed

	// Check if user exists
	var userID int
	err = tx.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err == sql.ErrNoRows {
		t.Logf("User %s not found, nothing to clean up", username)
		return
	} else if err != nil {
		t.Logf("Error checking for user existence: %v", err)
		return
	}

	// Delete associated tokens first
	_, err = tx.Exec("DELETE FROM user_tokens WHERE user_id = ?", userID)
	if err != nil {
		t.Logf("Failed to delete user tokens: %v", err)
		return
	}

	// Then delete the user
	_, err = tx.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		t.Logf("Failed to delete user: %v", err)
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		t.Logf("Failed to commit transaction: %v", err)
		return
	}

	t.Logf("Successfully cleaned up user %s and associated tokens", username)
}

func TestValidRegistration(t *testing.T) {

	err := db.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.RegisterHandler)

	cleanupTestUser(t, "testuser")

	// Change this to use form data
	formData := "username=testuser&email=testuser@example.com&password=testpassword"
	body := strings.NewReader(formData)
	req, _ := http.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Response body: %s", rr.Body.String())
	}

	// Optionally, verify that the user was actually created in the database
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", "testuser").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 user in database, got %d", count)
	}
	cleanupTestUser(t, "testuser")
}

func TestValidLogin(t *testing.T) {
	err := db.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Create a test user
	testUsername := "testuser"
	testPassword := "testpassword"
	testEmail := "testuser@example.com"

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	_, err = db.DB.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)",
		testUsername, string(hashedPassword), testEmail)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/login", handlers.LoginHandler)

	formData := fmt.Sprintf("username=%s&password=%s", testUsername, testPassword)
	body := strings.NewReader(formData)
	req, _ := http.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Response body: %s", rr.Body.String())
	}
	cleanupTestUser(t, "testuser")
}

func TestAuthorizedAccess(t *testing.T) {
	err := db.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Create a test user
	testUsername := "testuser"
	testPassword := "testpassword"
	testEmail := "testuser@example.com"

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	_, err = db.DB.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)",
		testUsername, string(hashedPassword), testEmail)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Perform login to get JWT
	loginURL := "/login"
	loginData := fmt.Sprintf("username=%s&password=%s", testUsername, testPassword)
	loginReq, _ := http.NewRequest("POST", loginURL, strings.NewReader(loginData))
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginRR := httptest.NewRecorder()

	http.HandlerFunc(handlers.LoginHandler).ServeHTTP(loginRR, loginReq)

	if loginRR.Code != http.StatusOK {
		t.Fatalf("Login failed: %v", loginRR.Body.String())
	}

	var loginResp struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(loginRR.Body.Bytes(), &loginResp)
	if err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}

	fmt.Println("Token: ", loginResp.Token)

	// Decode and verify the token
	token, err := jwt.Parse(loginResp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse JWT token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("Failed to get claims from token")
	}

	if _, ok := claims["user_id"]; !ok {
		t.Fatalf("user_id claim not found in token")
	}

	fmt.Printf("Token claims: %+v\n", claims)

	// Make the authorized request
	req, err := http.NewRequest("GET", "/fetch-characters?username="+testUsername, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.FetchAllCharactersHandler)

	// Create a new context with both the token and user_id
	ctx := context.WithValue(req.Context(), "token", &jwt.Token{
		Claims: jwt.MapClaims{
			"username": testUsername,
			"user_id":  float64(1), // Assuming user_id is 1 for the test user
		},
	})
	ctx = context.WithValue(ctx, "user_id", float64(1))

	middleware.AuthMiddleware(handler).ServeHTTP(rr, req.WithContext(ctx))

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Response body: %s", rr.Body.String())
	}

	cleanupTestUser(t, testUsername)
}

func getUserIDFromUsername(username string) (int, error) {
	fmt.Printf("Querying for user with username: %s\n", username)
	var userID int
	err := db.DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		fmt.Printf("Error querying for user: %v\n", err)
		return 0, err
	}
	fmt.Printf("Found user with ID: %d\n", userID)
	return userID, nil
}
