// main_test.go
package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "net/url"
    "strings"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    _ "github.com/go-sql-driver/mysql"
    "golang.org/x/crypto/bcrypt"
)

func Test_registerHandler(t *testing.T) {
    mockdb, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer mockdb.Close()
    db = mockdb

    mock.ExpectQuery("^SELECT EXISTS\\(SELECT 1 FROM users WHERE username=\\?\\)$").
        WithArgs("username").
        WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(0))

    mock.ExpectExec("^INSERT INTO users \\(username, password, email, token\\) VALUES \\(\\?, \\?, \\?, \\?\\)$").
        WithArgs("username", sqlmock.AnyArg(), "test@example.com", "").
        WillReturnResult(sqlmock.NewResult(1, 1))

    ts := httptest.NewServer(http.HandlerFunc(registerHandler))
    defer ts.Close()

    res, err := http.PostForm(ts.URL, url.Values{
        "username": {"username"},
        "password": {"password"},
        "email":    {"test@example.com"},
    })
    if err != nil {
        t.Fatalf("an error '%s' was not expected when making a request to the server", err)
    }

    if res.StatusCode != http.StatusOK {
        t.Errorf("expected status code to be %d, but got %d", http.StatusOK, res.StatusCode)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}

func Test_loginHandler(t *testing.T) {
    mockdb, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer mockdb.Close()

    passwordHash, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
    if err != nil {
        t.Fatalf("an error '%s' occurred while generating password hash", err)
    }
    rows := sqlmock.NewRows([]string{"id", "username", "password"}).
        AddRow(1, "testuser", passwordHash)
    mock.ExpectQuery("^SELECT id, username, password FROM users WHERE username=\\?$").
        WithArgs("testuser").
        WillReturnRows(rows)

    mock.ExpectExec("^UPDATE users SET token=\\? WHERE username=\\?$").
        WithArgs(sqlmock.AnyArg(), "testuser").
        WillReturnResult(sqlmock.NewResult(1, 1))

    db = mockdb

    req, err := http.NewRequest("POST", "/login", strings.NewReader("username=testuser&password=testpassword"))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(loginHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    var response map[string]string
    json.Unmarshal(rr.Body.Bytes(), &response)

    if _, exists := response["token"]; !exists {
        t.Errorf("handler returned unexpected body: token not found in response")
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}

// func TestGetAllCharacterCards(t *testing.T) {
// 	mockdb, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer mockdb.Close()

// 	rows := sqlmock.NewRows([]string{"id", "character", "learned"}).
// 		AddRow(1, "A", true).
// 		AddRow(2, "B", false)
// 	mock.ExpectQuery("^SELECT \\* FROM character_cards$").WillReturnRows(rows)

// 	db = mockdb

// 	// Create a request to pass to our handler
// 	req, err := http.NewRequest("GET", "", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(fetchAllCharactersHandler)

// 	// Call the handler function
// 	handler.ServeHTTP(rr, req)

// 	// Check the status code is what we expect
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	// Check the response body is what we expect
// 	expected := `[{"id":1,"character":"A","learned":true},{"id":2,"character":"B","learned":false}]`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
// 	}

// 	// Ensure all expectations were met
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expectations: %s", err)
// 	}
// }
