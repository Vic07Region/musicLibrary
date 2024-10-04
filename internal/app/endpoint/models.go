package endpoint

import (
	"github.com/google/uuid"
	"time"
)

type Song struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	GroupName   string    `json:"group_name" example:"Muse"`
	SongName    string    `json:"song_name" example:"Supermassive Black Hole"`
	ReleaseDate time.Time `json:"release_date" example:"1987-07-03T00:00:00Z"`
	Text        string    `json:"text" example:"string"`
	Link        string    `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
	CreatedAt   time.Time `json:"created_at" example:"2024-09-30T22:23:29.601031Z"`
}

type MessageError struct {
	Message string `json:"message"`
}

type SongText struct {
	Text []string `json:"text"`
}
