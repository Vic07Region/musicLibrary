package songinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

var (
	GroupNameRequiredError = fmt.Errorf("group name is required")
	SongNameRequiredError  = fmt.Errorf("song name is required")
	ServiceInternalError   = fmt.Errorf("Song Storage internal service error")
	ServiceBadRequestError = fmt.Errorf("Song Storage bad request")
	SerializeError         = fmt.Errorf("Song Storage bad request")
	ServiceUnknowError     = fmt.Errorf("Song Storage unknow error")
)

type SongStorage struct {
	base_url string
}

func New(base_url string) *SongStorage {
	return &SongStorage{base_url: base_url}
}

type SongInfo struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type FetchSongInfoParam struct {
	GroupName string
	SongName  string
}

func (s *SongStorage) FetchSongInfo(params FetchSongInfoParam) (*SongInfo, error) {
	if params.GroupName == "" {
		return nil, GroupNameRequiredError
	}

	if params.SongName == "" {
		return nil, SongNameRequiredError
	}

	resp, err := http.Get(path.Join(s.base_url, "/info"))
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer resp.Body.Close()
	var song SongInfo

	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(&song); err != nil {
			fmt.Println("Ошибка при декодировании ответа:", err)
			return nil, SerializeError
		}
	case http.StatusBadRequest:
		return nil, ServiceBadRequestError
	case http.StatusInternalServerError:
		return nil, ServiceInternalError
	default:
		return nil, ServiceUnknowError
	}

	return &song, err
}
