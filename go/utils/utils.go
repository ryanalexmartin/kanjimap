package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/ryanalexmartin/kanjimap/db"

	"github.com/dgrijalva/jwt-go"
)

func ValidateToken(tokenString string) (*jwt.Token, error) {
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
		var count int
		err := db.DB.QueryRow("SELECT COUNT(*) FROM user_tokens WHERE token = ?", tokenString).Scan(&count)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, fmt.Errorf("token not found in database")
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token expired")
			}
		}

		return token, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func CleanupExpiredTokens() {
	_, err := db.DB.Exec("DELETE FROM user_tokens WHERE created_at < ?", time.Now().Add(-24*time.Hour))
	if err != nil {
		fmt.Printf("Error cleaning up expired tokens: %v", err)
	}
}
