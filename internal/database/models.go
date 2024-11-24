package database

import (
	"time"
)

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Song struct {
	SongID      int       `json:"song_id"`
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date,omitempty"`
	Link        string    `json:"link,omitempty"`
}

type Verse struct {
	ID          int    `json:"id"`
	SongID      int    `json:"song_id"`
	VerseNumber int    `json:"verse_number"`
	VerseText   string `json:"verse_text"`
}

type VerseSmall struct {
	VerseNumber int    `json:"verse_number"`
	VerseText   string `json:"verse_text"`
}
