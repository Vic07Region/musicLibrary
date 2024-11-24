package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Vic07Region/musicLibrary/internal/connector/songinfo"
	"github.com/Vic07Region/musicLibrary/internal/database"
	"github.com/Vic07Region/musicLibrary/internal/lib/logger"
	"strings"
	"time"
)

var (
	SongNotFoundError    = fmt.Errorf("Song is not found")
	GroupNotFoundError   = fmt.Errorf("Group is not found")
	NoSongsError         = fmt.Errorf("There are no songs that meet the request")
	BadRequestError      = fmt.Errorf("Bad reqiest")
	UnknowError          = fmt.Errorf("Unknow error")
	BadDataFormatError   = fmt.Errorf("Wrong release date format")
	SongInfoserviceError = fmt.Errorf("Song info service InternalServerError")
	RequestError         = fmt.Errorf("Request execution error")
	TimeOutError         = fmt.Errorf("Request timeout exceeded")
)

type MusicService interface {
	FetchSongs(ctx context.Context, request FetchSongsRequest) (*FetchSongsResponse, error)
	FetchVerses(ctx context.Context, request FetchVersesRequest) (*FetchVersesResponse, error)
	DeleteSong(ctx context.Context, request DeleteSongRequest) (*DeleteSongResponse, error)
	UpdateSong(ctx context.Context, request UpdateSongRequest) (UpdateSongResponse, error)
	UpdateVerse(ctx context.Context, request UpdateVerseRequest) (UpdateVerseResponse, error)
	NewSong(ctx context.Context, request NewSongRequest) (*Song, error)
}

type Service struct {
	storage database.Storage
	songSrv songinfo.SongInfoSerice
	log     *logger.Logger
	debug   bool
}

func New(s database.Storage, t songinfo.SongInfoSerice, log *logger.Logger, debug bool) *Service {
	return &Service{storage: s, songSrv: t, log: log, debug: debug}
}

type FetchSongsRequest struct {
	GroupName   *string    `json:"group_name,omitempty" form:"group_name"`
	SongName    *string    `json:"song_name,omitempty" form:"song_name"`
	ReleaseDate *time.Time `json:"release_date,omitempty" form:"release_date"`
	SongText    *string    `json:"song_text,omitempty" form:"song_text"`
	Limit       uint64     `json:"limit" form:"limit"`
	Offset      uint64     `json:"offset" form:"offset"`
}

type FetchSongsResponse struct {
	Songs      []Song `json:"songs"`
	TotalCount int    `json:"total_count"`
}

func (s *Service) FetchSongs(ctx context.Context, request FetchSongsRequest) (*FetchSongsResponse, error) {
	if s.debug {
		s.log.Info("service.FetchSongs | request data", "request", request)
	}

	songList, err := s.storage.GetSongs(ctx, database.GetSongsRequest{
		GroupName:   request.GroupName,
		SongName:    request.SongName,
		SongText:    request.SongText,
		ReleaseDate: request.ReleaseDate,
		Limit:       request.Limit,
		Offset:      request.Offset,
	})

	if err != nil {
		s.log.Error("service.FetchSongs: GetSongs", "error", err.Error())
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, TimeOutError
		case errors.Is(err, sql.ErrNoRows):
			return nil, NoSongsError
		default:
			return nil, RequestError
		}
	}

	totalCount, err := s.storage.CountSongs(ctx)
	if err != nil {
		s.log.Error("service.FetchSongs: CountSongs", "error", err.Error())
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, TimeOutError
		default:
			return nil, RequestError
		}
	}

	var songs []Song

	for _, item := range songList {
		songs = append(songs, Song{
			ID:          item.SongID,
			GroupName:   item.GroupName,
			SongName:    item.SongName,
			ReleaseDate: item.ReleaseDate,
			Link:        item.Link,
		})
	}

	if s.debug {
		s.log.Info("service.FetchSongs | request data", "request", request)
		s.log.Info("service.FetchSongs | response data", "songs", songs, "totalCount", totalCount)
	}

	return &FetchSongsResponse{
		Songs:      songs,
		TotalCount: totalCount,
	}, err
}

type FetchVersesRequest struct {
	SongID int `json:"song_id" form:"song_id"`
	Limit  int `json:"limit" form:"limit"`
	Offset int `json:"offset" form:"offset"`
}

type FetchVersesResponse struct {
	Verses     []VerseSmall `json:"verses"`
	TotalCount int          `json:"total_count"`
}

func (s *Service) FetchVerses(ctx context.Context, request FetchVersesRequest) (*FetchVersesResponse, error) {
	verses, err := s.storage.GetVerses(ctx, database.GetVersesRequest{
		SongID: request.SongID,
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		s.log.Error("service.FetchVerses: GetVerses", "error", err.Error())
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, TimeOutError
		default:
			return nil, RequestError
		}
	}

	songVersesCount, err := s.storage.CountVerses(ctx, request.SongID)
	if err != nil {
		s.log.Error("service.FetchVerses: CountVerses", "error", err.Error())
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, TimeOutError
		default:
			return nil, RequestError
		}
	}

	var vs []VerseSmall
	for _, v := range verses {
		vs = append(vs, VerseSmall{
			VerseNumber: v.VerseNumber,
			VerseText:   v.VerseText,
		})
	}

	if s.debug {
		s.log.Info("service.FetchVerses | request data", "request", request)
		s.log.Info("service.FetchVerses | response data", "songs", vs, "totalCount", songVersesCount)
	}

	return &FetchVersesResponse{
		Verses:     vs,
		TotalCount: songVersesCount,
	}, nil
}

type DeleteSongRequest struct {
	SongID int `json:"song_id"`
}

type DeleteSongResponse struct {
	Success bool `json:"success"`
}

func (s *Service) DeleteSong(ctx context.Context, request DeleteSongRequest) (*DeleteSongResponse, error) {
	var reponse DeleteSongResponse

	err := s.storage.DeleteSong(ctx, request.SongID)
	if err != nil {
		s.log.Error("service.DeleteSong | DeleteSong", "error", err.Error())
		switch {
		case errors.Is(err, context.DeadlineExceeded):

			return &reponse, TimeOutError
		case errors.Is(err, sql.ErrNoRows):
			return &reponse, SongNotFoundError
		default:
			return &reponse, RequestError
		}
	}

	reponse.Success = true

	if s.debug {
		s.log.Info("service.DeleteSong | request data", "request", request)
		s.log.Info("service.DeleteSong | response data", "success", reponse.Success)
	}

	return &reponse, nil
}

type UpdateSongRequest struct {
	SongID      int        `json:"song_id"`
	GroupName   *string    `json:"group_id"`
	SongName    *string    `json:"song_name,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	Link        *string    `json:"link,omitempty"`
}

type UpdateSongResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateSong(ctx context.Context, request UpdateSongRequest) (UpdateSongResponse, error) {
	var result UpdateSongResponse
	songParam := database.UpdateSongRequest{
		SongID:      request.SongID,
		SongName:    request.SongName,
		ReleaseDate: request.ReleaseDate,
		Link:        request.Link,
	}

	if request.GroupName != nil {
		groupID, err := s.storage.GetGroupID(ctx, *request.GroupName)
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				return result, TimeOutError
			case errors.Is(err, sql.ErrNoRows):
				return result, GroupNotFoundError
			default:
				return result, RequestError
			}
		}
		songParam.GroupID = &groupID

	}

	err := s.storage.UpdateSong(ctx, songParam)
	if err != nil {
		s.log.Error("service.UpdateSong | UpdateSong", "error", err.Error())
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return result, TimeOutError
		case errors.Is(err, sql.ErrNoRows):
			return result, SongNotFoundError
		default:
			return result, RequestError
		}
	}
	result.Success = true

	if s.debug {
		s.log.Info("service.UpdateSong | request data", "request", request)
		s.log.Info("service.UpdateSong | response data", "success", result.Success)
	}

	return result, err
}

type UpdateVerseRequest struct {
	SongID      int    `json:"song_id"`
	VerseNumber int    `json:"verse_number"`
	VerseText   string `json:"verse_text"`
}

type UpdateVerseResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateVerse(ctx context.Context, request UpdateVerseRequest) (UpdateVerseResponse, error) {
	var result UpdateVerseResponse
	err := s.storage.UpdateVerse(ctx, database.UpdateVerseRequest{
		SongID:      request.SongID,
		VerseNumber: request.VerseNumber,
		VerseText:   request.VerseText,
	})
	if err != nil {
		s.log.Error("service.UpdateVerse | UpdateVerse", "error", err.Error())
		switch {
		case errors.Is(err, context.DeadlineExceeded):

			return result, TimeOutError
		case errors.Is(err, sql.ErrNoRows):
			return result, SongNotFoundError
		default:
			return result, RequestError
		}
	}
	result.Success = true

	if s.debug {
		s.log.Info("service.UpdateVerse | request data", "request", request)
		s.log.Info("service.UpdateVerse | response data", "success", result.Success)
	}

	return result, err
}

type NewSongRequest struct {
	GroupName string `json:"group"`
	SongName  string `json:"song"`
}

func (s *Service) NewSong(ctx context.Context, request NewSongRequest) (*Song, error) {
	songInfo, err := s.songSrv.FetchSongInfo(songinfo.FetchSongInfoParam{
		GroupName: request.GroupName,
		SongName:  request.SongName,
	})
	if err != nil {
		s.log.Error("service.NewSong | FetchSongInfo", "error", err.Error())
		return nil, err
	}

	release_date, err := time.Parse("02.01.2006", songInfo.ReleaseDate)
	if err != nil {
		s.log.Warn("service.NewSong | parse release date", "Error", err.Error())
		return nil, BadDataFormatError
	}

	verseSplit := strings.Split(songInfo.Text, "\n\n")

	var verses []database.VerseSmall

	for idx, verse := range verseSplit {
		verses = append(verses, database.VerseSmall{
			VerseNumber: idx,
			VerseText:   verse,
		})
	}

	newSong, err := s.storage.AddSong(ctx, database.AddSongRequest{
		GroupName:   request.GroupName,
		SongName:    request.SongName,
		ReleaseDate: release_date,
		Verses:      verses,
		Link:        songInfo.Link,
	})
	if err != nil {
		s.log.Error("service.NewSong | CreateSong", "error", err.Error())
		switch err {
		case context.DeadlineExceeded:
			return nil, TimeOutError
		default:
			return nil, RequestError
		}
	}
	//todo переосмыслить ошибки

	song := Song{
		ID:          int(newSong.SongID),
		GroupName:   request.GroupName,
		SongName:    request.SongName,
		ReleaseDate: release_date,
		Link:        songInfo.Link,
	}

	if s.debug {
		s.log.Info("service.NewSong | request data", "request", request)
		s.log.Info("service.NewSong | response data", "song", song)
	}

	return &song, nil
}
