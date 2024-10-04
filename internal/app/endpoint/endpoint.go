package endpoint

import (
	"context"
	"github.com/Vic07Region/musicLibrary/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type Service interface {
	FetchSongs(ctx context.Context, params service.FetchSongsParam) ([]service.Song, error)
	FetchSongText(ctx context.Context, params service.FetchSongTextParam) ([]string, error)
	DeleteSong(ctx context.Context, song_id uuid.UUID) error
	EditSong(ctx context.Context, params service.EditSongParam) (*service.Song, error)
	NewSong(ctx context.Context, params service.NewSongParam) (*service.Song, error)
}

type Endpoint struct {
	s Service
}

func New(s Service) *Endpoint {
	return &Endpoint{
		s: s,
	}
}

// @BasePath /api/v1

// ListSongs lists all existing songs
// @Summary list songs
// @Schemes
// @Description fetching song list
// @Param   group      query     string     false  "Muse"
// @Param   song      query     string     false  "Supermassive Black Hole"
// @Param   text      query     string     false  "song text"
// @Param   limit      query     int     false  "10"
// @Param   offset      query     int     false "2"
// @Tags songs
// @Accept json
// @Produce json
// @Success 200 {array} endpoint.Song
// @Failure      400  {object}  endpoint.MessageError
// @Failure      500
// @Router /songs [get]
func (e *Endpoint) FetchSongsHandler(c *gin.Context) {
	var fetchParams service.FetchSongsParam

	if val, ok := c.GetQuery("group"); ok {
		fetchParams.Group_name = val
	}

	if val, ok := c.GetQuery("song"); ok {
		fetchParams.Song_name = val
	}

	if val, ok := c.GetQuery("text"); ok {
		fetchParams.Text = val
	}

	if val, ok := c.GetQuery("offset"); ok {
		if intval, err := strconv.Atoi(val); err == nil {
			fetchParams.Offset = int64(intval)
		}

	}

	if val, ok := c.GetQuery("limit"); ok {
		if intval, err := strconv.Atoi(val); err == nil {
			fetchParams.Limit = int64(intval)
		}

	}
	songs, err := e.s.FetchSongs(c.Request.Context(), fetchParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageError{Message: err.Error()})
	}
	c.JSON(http.StatusOK, songs)
}
