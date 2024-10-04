package database

import (
	"github.com/google/uuid"
	"time"
)

type Song struct {
	ID          uuid.UUID `json:"id"`
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
