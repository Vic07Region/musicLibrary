package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Vic07Region/musicLibrary/internal/connector/songinfo"
	"github.com/Vic07Region/musicLibrary/internal/database"
	"github.com/Vic07Region/musicLibrary/internal/lib/logger"
	"github.com/google/uuid"
	"strings"
	"time"
)

var (
	SongNotFoundError    = fmt.Errorf("Song is not found")
	NoSongsError         = fmt.Errorf("There are no songs that meet the request")
	BadRequestError      = fmt.Errorf("Bad reqiest")
	UnknowError          = fmt.Errorf("Unknow error")
	BadDataFormatError   = fmt.Errorf("Wrong release date format")
	SongInfoServiceError = fmt.Errorf("Song info service InternalServerError")
	RequestError         = fmt.Errorf("Request execution error")
	TimeOutError         = fmt.Errorf("Request timeout exceeded")
)

type MusicLibraryStorage interface {
	GetSong(ctx context.Context, song_id uuid.UUID) (*database.Song, error)
	GetSongs(ctx context.Context, params database.GetSongsParam) ([]database.Song, error)
	GetSongText(ctx context.Context, song_id uuid.UUID) (string, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
	UpdateSong(ctx context.Context, params database.UpdateSongParam) error
	CreateSong(ctx context.Context, params database.CreateSongParam) (*database.CreateSongResult, error)
}

type SongInfoSerice interface {
	FetchSongInfo(params songinfo.FetchSongInfoParam) (*songinfo.SongInfo, error)
}

type Service struct {
	storage MusicLibraryStorage
	songSrv SongInfoSerice
	log     *logger.Logger
}

func New(s MusicLibraryStorage, t SongInfoSerice, log *logger.Logger) *Service {
	return &Service{storage: s, songSrv: t, log: log}
}

type Song struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	GroupName   string    `json:"group_name" example:"Muse"`
	SongName    string    `json:"song_name" example:"Supermassive Black Hole"`
	ReleaseDate time.Time `json:"release_date" example:"1987-07-03T00:00:00Z"`
	Text        string    `json:"text" example:"string"`
	Link        string    `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
	CreatedAt   time.Time `json:"created_at" example:"2024-09-30T22:23:29.601031Z"`
}

type FetchSongsParam struct {
	Group_name string
	Song_name  string
	Text       string
	Limit      int64
	Offset     int64
}

func (s *Service) FetchSongs(ctx context.Context, params FetchSongsParam) ([]Song, error) {
	song_items, err := s.storage.GetSongs(ctx, database.GetSongsParam{
		Group_name: params.Group_name,
		Song_name:  params.Song_name,
		Text:       params.Text,
		Limit:      params.Limit,
		Offset:     params.Offset,
	})

	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			s.log.Debug("GetSongs: GetSongs", "error", err.Error())
			return nil, TimeOutError
		case sql.ErrNoRows:
			s.log.Debug("GetSongs: GetSongs", "error", err.Error())
			return nil, NoSongsError
		default:
			s.log.Debug("GetSongs: GetSongs", "error", err.Error())
			return nil, RequestError
		}

	}
	var songs []Song

	for _, item := range song_items {
		songs = append(songs, Song{
			ID:          item.ID,
			GroupName:   item.GroupName,
			SongName:    item.SongName,
			ReleaseDate: item.ReleaseDate,
			Text:        item.Text,
			Link:        item.Link,
			CreatedAt:   item.CreatedAt,
		})
	}

	return songs, err
}

type FetchSongTextParam struct {
	Song_id uuid.UUID
	Limit   int64
	Offset  int64
}

func (s *Service) FetchSongText(ctx context.Context, params FetchSongTextParam) ([]string, error) {

	song_text, err := s.storage.GetSongText(ctx, params.Song_id)
	if err != nil {
		s.log.Debug("FetchSongText - GetSongText", "Error", err.Error())
		if err == sql.ErrNoRows {
			return nil, SongNotFoundError
		}
		s.log.Debug("FetchSongText: GetSongText", "error", err.Error())
		return nil, RequestError
	}

	text_parts := strings.Split(song_text, "\n\n")

	var start, end int64

	if int64(len(text_parts))-1 >= params.Offset {
		start = params.Offset
	} else {
		start = int64(len(text_parts)) - 1
	}

	if int64(len(text_parts))-1 >= start+params.Limit {
		end = start + params.Limit
	} else {
		end = int64(len(text_parts))
	}

	return text_parts[start:end], nil
}

func (s *Service) DeleteSong(ctx context.Context, song_id uuid.UUID) error {
	err := s.storage.DeleteSong(ctx, song_id)
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			s.log.Debug("DeleteSong:", "error", err.Error())
			return TimeOutError
		case sql.ErrNoRows:
			s.log.Debug("DeleteSong:", "error", err.Error())
			return SongNotFoundError
		default:
			s.log.Debug("DeleteSong:", "error", err.Error())
			return RequestError
		}
	}
	return nil
}

type EditSongParam struct {
	Song_id     uuid.UUID
	GroupName   string
	SongName    string
	ReleaseDate time.Time
	Text        string
	Link        string
}

func (s *Service) EditSong(ctx context.Context, params EditSongParam) (*Song, error) {
	old_song, err := s.storage.GetSong(ctx, params.Song_id)
	if err != nil {
		s.log.Debug("EditSong - GetSong", "Error", err.Error())
		if err == sql.ErrNoRows {
			return nil, SongNotFoundError
		}
		return nil, UnknowError
	}
	var update_param database.UpdateSongParam
	updated_song := Song{
		ID:          old_song.ID,
		GroupName:   old_song.GroupName,
		SongName:    old_song.SongName,
		ReleaseDate: old_song.ReleaseDate,
		Text:        old_song.Text,
		Link:        old_song.Link,
		CreatedAt:   old_song.CreatedAt,
	}

	if params.GroupName != "" && params.GroupName != old_song.GroupName {
		update_param.GroupName = params.GroupName
		updated_song.GroupName = params.GroupName
	}

	if params.SongName != "" && params.SongName != old_song.SongName {
		update_param.SongName = params.SongName
		updated_song.SongName = params.SongName
	}

	if params.ReleaseDate.IsZero() && params.ReleaseDate != old_song.ReleaseDate {
		update_param.ReleaseDate = params.ReleaseDate
		updated_song.ReleaseDate = params.ReleaseDate
	}

	if params.Text != "" && params.Text != old_song.Text {
		update_param.Text = params.Text
		updated_song.Text = params.Text
	}

	if params.Link != "" && params.Link != old_song.Link {
		update_param.Link = params.Link
		updated_song.Link = params.Link
	}

	err = s.storage.UpdateSong(ctx, update_param)

	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			s.log.Debug("EditSong: UpdateSong", "error", err.Error())
			return nil, TimeOutError
		case sql.ErrNoRows:
			s.log.Debug("EditSong: UpdateSong", "error", err.Error())
			return nil, SongNotFoundError
		default:
			s.log.Debug("EditSong: UpdateSong", "error", err.Error())
			return nil, RequestError
		}
	}

	return &updated_song, nil
}

type NewSongParam struct {
	GroupName string `json:"group"`
	SongName  string `json:"song"`
}

func (s *Service) NewSong(ctx context.Context, params NewSongParam) (*Song, error) {
	song_info, err := s.songSrv.FetchSongInfo(songinfo.FetchSongInfoParam{
		GroupName: params.GroupName,
		SongName:  params.SongName,
	})
	if err != nil {
		s.log.Debug("NewSong - FetchSongInfo service", "Error", err.Error())
		return nil, err
	}

	release_date, err := time.Parse("02.01.2006", song_info.ReleaseDate)
	if err != nil {
		s.log.Debug("NewSong - parse release date", "Error", err.Error())
		return nil, BadDataFormatError
	}
	newSong, err := s.storage.CreateSong(ctx, database.CreateSongParam{
		GroupName:   params.GroupName,
		SongName:    params.SongName,
		ReleaseDate: release_date,
		Text:        song_info.Text,
		Link:        song_info.Link,
	})
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			s.log.Debug("NewSong: CreateSong", "error", err.Error())
			return nil, TimeOutError
		default:
			s.log.Debug("NewSong: CreateSong", "error", err.Error())
			return nil, RequestError
		}
	}
	//todo переосмыслить ошибки
	song := Song{
		ID:          newSong.ID,
		GroupName:   newSong.GroupName,
		SongName:    newSong.SongName,
		ReleaseDate: release_date,
		Text:        song_info.Text,
		Link:        song_info.Link,
		CreatedAt:   newSong.CreatedAt,
	}
	return &song, nil
}
