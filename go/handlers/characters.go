package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ryanalexmartin/kanjimap/db"
	"github.com/ryanalexmartin/kanjimap/models"

	"github.com/dgrijalva/jwt-go"
)

func FetchAllCharactersHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	token := r.Context().Value("token").(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	if claims["username"] != username {
		http.Error(w, "Username mismatch", http.StatusUnauthorized)
		return
	}

	var userID int
	err := db.DB.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Database error", err)
		return
	}

	rows, err := db.DB.Query(`
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

	var characterCards []models.CharacterCard
	for rows.Next() {
		var card models.CharacterCard
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

func LearnCharacter(w http.ResponseWriter, r *http.Request) {
	var character models.CharacterCard
	err := json.NewDecoder(r.Body).Decode(&character)
	if err != nil {
		http.Error(w, "Unable to decode JSON request", http.StatusBadRequest)
		fmt.Println("Unable to decode JSON request", err)
		return
	}

	var userID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE username=?", character.Username).Scan(&userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Database error", err)
		return
	}

	var learned bool
	err = db.DB.QueryRow("SELECT learned FROM user_character_progress WHERE character_id=? AND user_id=?", character.CharacterID, userID).Scan(&learned)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = db.DB.Exec("INSERT INTO user_character_progress (user_id, character_id, learned) VALUES (?, ?, ?)", userID, character.CharacterID, character.Learned)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				fmt.Println("Database error", err)
				return
			}
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
			fmt.Println("Database error", err)
			return
		}
	}

	_, err = db.DB.Exec("UPDATE user_character_progress SET learned=? WHERE user_id=? AND character_id=?", character.Learned, userID, character.CharacterID)
	if err != nil {
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
		Learned:     character.Learned,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)
}

func LearnedCharactersHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	token := r.Context().Value("token").(*jwt.Token)
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

	var userID int
	err := db.DB.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	rows, err := db.DB.Query(`
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(learnedCharacters)
}
