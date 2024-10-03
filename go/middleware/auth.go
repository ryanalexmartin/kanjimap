package middleware

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/ryanalexmartin/kanjimap/db"
	"github.com/ryanalexmartin/kanjimap/utils"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.ExtractToken(r)
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid user_id in token claims", http.StatusUnauthorized)
			return
		}

		// Check if the token exists in the database
		err = db.DB.QueryRow("SELECT user_id FROM user_tokens WHERE token = ?", tokenString).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// Add the username to the request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func LoggedFs(fs http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/index.html", http.StatusFound)
			return
		}
		http.FileServer(fs).ServeHTTP(w, r)
	})
}
