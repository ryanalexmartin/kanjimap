package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Character struct {
	ID        string `json:"id"`
	Character string `json:"character"`
}

type CharacterCard struct {
	Username            string  `json:"username"`
	Character           string  `json:"character"`
	Learned             bool    `json:"learned"`
	CharacterID         string  `json:"characterId"`
	Frequency           int     `json:"frequency"`
	CumulativeFrequency float64 `json:"cumulativeFrequency"`
	Pinyin              string  `json:"pinyin"`
	English             string  `json:"english"`
}
