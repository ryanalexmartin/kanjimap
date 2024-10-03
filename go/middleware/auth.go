package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ryanalexmartin/kanjimap/utils"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		if reqToken == "" || !strings.HasPrefix(reqToken, "Bearer ") {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		reqToken = strings.TrimPrefix(reqToken, "Bearer ")

		token, err := utils.ValidateToken(reqToken)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add the token to the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "token", token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
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
