package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ryanalexmartin/kanjimap/db"
	"github.com/ryanalexmartin/kanjimap/handlers"
	"github.com/ryanalexmartin/kanjimap/middleware"
)

func main() {
	// Initialize database connection
	err := db.Initialize()
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing database.  Did you start the db?  error: %w", err))
	}
	defer db.Close()

	// ... (keep the frontend directory check code)

	mux := http.NewServeMux()
	mux.Handle("/", middleware.LoggedFs(http.Dir("./frontend")))
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/fetch-characters", middleware.AuthMiddleware(handlers.FetchAllCharactersHandler))
	mux.HandleFunc("/learn-character", middleware.AuthMiddleware(handlers.LearnCharacter))
	mux.HandleFunc("/learned-characters", middleware.AuthMiddleware(handlers.LearnedCharactersHandler))

	// ... (keep the CORS setup code)

	var port int
	port, err = strconv.Atoi(os.Getenv("VUE_APP_API_PORT"))
	if err != nil {
		log.Fatal("PORT environment variable not set")
	}

	fmt.Printf("Starting application on port %v \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

// ... (keep the loggedFs handler)
