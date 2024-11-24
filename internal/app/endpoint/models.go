package endpoint

import (
	"time"
)

type Song struct {
	ID          int       `json:"id" example:"1"`
	GroupName   string    `json:"group" example:"Muse"`
	SongName    string    `json:"song" example:"Supermassive Black Hole"`
	ReleaseDate time.Time `json:"releaseDate" example:"1987-07-03T00:00:00Z"`
	Link        string    `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}

type MessageError struct {
	Message string `json:"message"`
}

type SongText struct {
	Text []string `json:"text"`
}

type UpdateSong struct {
	GroupName   string `json:"group"`
	SongName    string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Link        string `json:"link"`
}

type UpdateSongRequest struct {
	SongID      int        `json:"song_id"`
	SongName    *string    `json:"song_name,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	Link        *string    `json:"link,omitempty"`
}

type NewSong struct {
	GroupName string `json:"group" validate:"required"`
	SongName  string `json:"song" validate:"required"`
}

type UpdateVerseRequest struct {
	VerseNumber int    `json:"verse_number"`
	VerseText   string `json:"verse_text"`
}
