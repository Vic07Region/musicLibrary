package songinfo

import (
	"encoding/json"
	"fmt"
	"github.com/Vic07Region/musicLibrary/internal/lib/logger"
	"net/http"
)

var (
	GroupNameRequiredError = fmt.Errorf("group name is required")
	SongNameRequiredError  = fmt.Errorf("song name is required")
	ServiceInternalError   = fmt.Errorf("Song Storage internal service error")
	ServiceBadRequestError = fmt.Errorf("Song Storage bad request")
	SerializeError         = fmt.Errorf("Song Storage bad request")
	ServiceUnknowError     = fmt.Errorf("Song Storage unknow error")
)

type SongInfoSerice interface {
	FetchSongInfo(params FetchSongInfoParam) (*SongInfo, error)
}

type SongStorage struct {
	baseUrl string
	log     *logger.Logger
}

func New(baseUrl string, log *logger.Logger) *SongStorage {
	return &SongStorage{baseUrl: baseUrl, log: log}
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
		return nil, GroupNameRequiredError
	}

	if params.SongName == "" {
		return nil, SongNameRequiredError
	}

	url := s.baseUrl + "/info"

	resp, err := http.Get(url)
	if err != nil {
		s.log.Error(fmt.Sprintf("songinfo.FetchSongInfo http.Get(%s)", url),
			"error", err.Error(),
		)
		return nil, ServiceInternalError
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
			return nil, SerializeError
		}
	case http.StatusBadRequest:
		return nil, ServiceBadRequestError
	case http.StatusInternalServerError:
		return nil, ServiceInternalError
	default:
		s.log.Error(
			"songinfo.FetchSongInfo response",
			"error", err.Error(),
		)
		return nil, ServiceUnknowError
	}

	return &song, err
}
