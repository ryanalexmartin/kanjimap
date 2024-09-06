// main_test.go
package main

import (
  "net/http"
  "net/http/httptest"
  "net/url"
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

  mock.ExpectExec("^INSERT INTO users \\(username, password\\) VALUES \\(\\?, \\?\\)$").
  WithArgs("username", sqlmock.AnyArg()).
  WillReturnResult(sqlmock.NewResult(1, 1))

  // start the server
  ts := httptest.NewServer(http.HandlerFunc(registerHandler))
  defer ts.Close()

  // make a request to the server
  res, err := http.PostForm(ts.URL, map[string][]string{
    "username": {"username"},
    "password": {"password"},
  })
  if err != nil {
    t.Fatalf("an error '%s' was not expected when making a request to the server", err)
  }

  // assert the response
  if res.StatusCode != http.StatusOK {
    t.Errorf("expected status code to be %d, but got %d", http.StatusOK, res.StatusCode)
  }
}

func Test_loginHandler(t *testing.T) {
  // Create a new mock database
  mockdb, mock, err := sqlmock.New()
  if err != nil {
    t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
  }
  defer mockdb.Close()

  // Set up your expectations
  passwordHash, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
  if err != nil {
    t.Fatalf("an error '%s' occurred while generating password hash", err)
  }
  rows := sqlmock.NewRows([]string{"id", "username", "password"}).
  AddRow(1, "testuser", passwordHash)
  mock.ExpectQuery("^SELECT \\* FROM users WHERE username = \\?$").WithArgs("testuser").WillReturnRows(rows)

  db = mockdb

  // Create a request to pass to our handler
  req, err := http.NewRequest("POST", "", nil)
  if err != nil {
    t.Fatal(err)
  }

  // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
  rr := httptest.NewRecorder()
  handler := http.HandlerFunc(loginHandler)

  // Populate the request's form data
  req.Form = url.Values{}
  req.Form.Add("username", "testuser")
  req.Form.Add("password", "testpassword")

  // Call the handler function
  handler.ServeHTTP(rr, req)

  // Check the status code is what we expect
  if status := rr.Code; status != http.StatusOK {
    t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
  }

  // Check the response body is what we expect
  expected := `Logged in successfully`
  if rr.Body.String() != expected {
    t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
  }

  // Ensure all expectations were met
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
