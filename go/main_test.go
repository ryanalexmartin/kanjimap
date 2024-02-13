// main_test.go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldCreateUser(t *testing.T) {
	fmt.Println("TestShouldCreateUser")
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
