package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT EXISTS").WithArgs("testuser").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(0))
	mock.ExpectExec("INSERT INTO users").WithArgs("testuser", sqlmock.AnyArg(), "test@example.com", "").WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("POST", "/register", strings.NewReader(url.Values{
		"username": {"testuser"},
		"password": {"testpassword"},
		"email":    {"test@example.com"},
	}.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "User successfully registered with ID: 1"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLoginHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	rows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(1, "testuser", hashedPassword)
	mock.ExpectQuery("SELECT id, username, password FROM users WHERE username").WithArgs("testuser").WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO user_tokens").WithArgs(1, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("POST", "/login", strings.NewReader(url.Values{
		"username": {"testuser"},
		"password": {"testpassword"},
	}.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(), "token") {
		t.Errorf("Handler returned unexpected body: token not found in response")
	}
}
