package songinfo

import (
	"encoding/json"
	"fmt"      //nolint:gci
	"net/http" //nolint:gci

	"github.com/Vic07Region/musicLibrary/internal/lib/logger" //nolint:gci
)

var (
	ErrGroupNameRequired = fmt.Errorf("group name is required")
	ErrSongNameRequired  = fmt.Errorf("song name is required")
	ErrServiceInternal   = fmt.Errorf("song Storage internal service error")
	ErrServiceBadRequest = fmt.Errorf("song Storage bad request")
	ErrSerialize         = fmt.Errorf("song Storage bad request")
	ErrServiceUnknow     = fmt.Errorf("song Storage unknow error")
)

type InfoSerice interface {
	FetchSongInfo(params FetchSongInfoParam) (*SongInfo, error)
}

type SongStorage struct {
	baseURL string
	log     *logger.Logger
}

func New(baseURL string, log *logger.Logger) InfoSerice {
	return &SongStorage{baseURL: baseURL, log: log}
}

type SongInfo struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type FetchSongInfoParam struct {
	GroupName string `json:"group_name"`
	SongName  string `json:"song_name"`
}

func (s *SongStorage) FetchSongInfo(params FetchSongInfoParam) (*SongInfo, error) {
	if params.GroupName == "" {
		s.log.Error("songinfo.FetchSongInfo | empty GroupName")
		return nil, ErrGroupNameRequired
	}

	if params.SongName == "" {
		s.log.Error("songinfo.FetchSongInfo | empty SongName")
		return nil, ErrSongNameRequired
	}
	url := s.baseURL
	resp, err := http.Get(url)
	if err != nil {
		s.log.Error(fmt.Sprintf("songinfo.FetchSongInfo http.Get(%s)", url),
			"error", err.Error(),
		)
		return nil, ErrServiceInternal
	}

	defer resp.Body.Close()

	var song SongInfo

	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(&song); err != nil {
			s.log.Error(
				"songinfo.FetchSongInfo | Song info service Decode body",
				"error", err.Error(),
			)
			return nil, ErrSerialize
		}
	case http.StatusBadRequest:
		return nil, ErrServiceBadRequest
	case http.StatusInternalServerError:
		return nil, ErrServiceInternal
	default:
		s.log.Error(
			"songinfo.FetchSongInfo response", "error", resp.Status)
		return nil, ErrServiceUnknow
	}

	return &song, err
}
