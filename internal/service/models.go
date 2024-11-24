package service

import "time"

type Song struct {
	ID          int       `json:"id" example:"1"`
	GroupName   string    `json:"group_name" example:"Muse"`
	SongName    string    `json:"song_name" example:"Supermassive Black Hole"`
	ReleaseDate time.Time `json:"release_date" example:"1987-07-03T00:00:00Z"`
	Link        string    `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}

type VerseSmall struct {
	VerseNumber int    `json:"verse_number"`
	VerseText   string `json:"verse_text"`
}
