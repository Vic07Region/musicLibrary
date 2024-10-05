package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Vic07Region/musicLibrary/internal/lib/logger"
	"github.com/Vic07Region/musicLibrary/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

type Service interface {
	FetchSongs(ctx context.Context, params service.FetchSongsParam) ([]service.Song, error)
	FetchSongText(ctx context.Context, params service.FetchSongTextParam) ([]string, error)
	DeleteSong(ctx context.Context, song_id uuid.UUID) error
	EditSong(ctx context.Context, params service.EditSongParam) (*service.Song, error)
	NewSong(ctx context.Context, params service.NewSongParam) (*service.Song, error)
}

type Endpoint struct {
	s   Service
	log *logger.Logger
}

func New(s Service, log *logger.Logger) *Endpoint {
	return &Endpoint{
		s:   s,
		log: log,
	}
}

// @BasePath /api/v1

// ListSongs lists all existing songs
// @Summary List songs
// @Schemes
// @Description fetching song list
// @Param   group      query     string     false  "group name"	example(Muse)
// @Param   song      query     string     false  "song name"	example(Supermassive Black Hole)
// @Param   text      query     string     false  "song text"	example(song text)
// @Param   limit      query     int     false  "items limit"	example(10)
// @Param   offset      query     int     false "offset items"	example(2)
// @Tags Songs
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
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}
	var resp []Song
	for _, i := range songs {
		resp = append(resp, Song{
			ID:          i.ID,
			GroupName:   i.GroupName,
			SongName:    i.SongName,
			ReleaseDate: i.ReleaseDate,
			Text:        i.Text,
			Link:        i.Link,
			CreatedAt:   i.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// SongText text song
// @Summary Song text
// @Schemes
// @Description fetching song text
// @Tags Songs
// @Accept json
// @Produce json
// @Param   id      path     string     true  "Идентификатор сущности (UUID)" example(462b63b5-c101-424f-93a2-aa69997036e2)
// @Param   limit      query     int     false  "items limit"	example(10)
// @Param   offset      query     int     false "offset items"	example(2)
// @Success 200 {object} endpoint.SongText
// @Failure      400  {object}  endpoint.MessageError
// @Failure      404  {object}  endpoint.MessageError
// @Failure      500
// @Router /songs/{id} [get]
func (e *Endpoint) FetchSongTextHandler(c *gin.Context) {
	var fetchParams service.FetchSongTextParam

	song_id := c.Param("id")

	uuid_song, err := uuid.Parse(song_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageError{"wrong format id(uuid)"})
		return
	} else if uuid_song.Version() != 4 {
		c.JSON(http.StatusBadRequest, MessageError{"invalid version id(uuid_v4)"})
		return
	}

	fetchParams.Song_id = uuid_song

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

	texts, err := e.s.FetchSongText(c.Request.Context(), fetchParams)
	if err != nil {
		if err == service.SongNotFoundError {
			c.JSON(http.StatusNotFound, MessageError{err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}

	c.JSON(http.StatusBadRequest, SongText{Text: texts})
}

// @Summary Delete Song
// @Schemes
// @Description deleting song
// @Tags Songs
// @Accept json
// @Produce json
// @Param   id      path     string     true  "Идентификатор сущности (UUID)" example(462b63b5-c101-424f-93a2-aa69997036e2)
// @Success 204
// @Failure      400  {object}  endpoint.MessageError
// @Failure      404  {object}  endpoint.MessageError
// @Failure      500
// @Router /songs/{id} [delete]
func (e *Endpoint) DeleteSongHandler(c *gin.Context) {
	song_id := c.Param("id")

	uuid_song, err := uuid.Parse(song_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageError{"wrong format id(uuid)"})
		return
	} else if uuid_song.Version() != 4 {
		c.JSON(http.StatusBadRequest, MessageError{"invalid version id(uuid_v4)"})
		return
	}

	if err := e.s.DeleteSong(c.Request.Context(), uuid_song); err != nil {
		if err == service.SongNotFoundError {
			c.JSON(http.StatusNotFound, MessageError{err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}
	c.Status(http.StatusNoContent)

}

// @Summary Edit Song
// @Schemes
// @Description edit song
// @Tags Songs
// @Accept json
// @Produce json
// @Param   id      path     string     true  "Идентификатор сущности (UUID)" example(462b63b5-c101-424f-93a2-aa69997036e2)
// @Param request body endpoint.UpdateSong true "query params"
// @Success 200 {object} endpoint.Song
// @Failure      400  {object}  endpoint.MessageError
// @Failure      404  {object}  endpoint.MessageError
// @Failure      500
// @Router /songs/{id} [patch]
func (e *Endpoint) UpdateSongHandler(c *gin.Context) {
	song_id := c.Param("id")

	uuid_song, err := uuid.Parse(song_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageError{"wrong format id(uuid)"})
		return
	} else if uuid_song.Version() != 4 {
		c.JSON(http.StatusBadRequest, MessageError{"invalid version id(uuid_v4)"})
		return
	}

	var inputData map[string]interface{}
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}

	allowedKeys := map[string]bool{
		"group":       true,
		"song":        true,
		"releaseDate": true,
		"text":        true,
		"link":        true,
	}
	var wrongKey []string
	for key := range inputData {
		if _, ok := allowedKeys[key]; !ok {
			wrongKey = append(wrongKey, key)
		}
	}

	if len(wrongKey) > 0 {
		c.JSON(http.StatusBadRequest, MessageError{fmt.Sprintf("invalid fields %s", wrongKey)})
		return
	}

	var songData UpdateSong

	if err := mapToStruct(inputData, &songData); err != nil {
		c.JSON(http.StatusInternalServerError, MessageError{err.Error()})
		return
	}
	var releaseD time.Time
	if releaseDateStr, ok := inputData["releaseDate"].(string); ok {
		releaseDate, err := time.Parse("02.01.2006", releaseDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, MessageError{"Invalid date format. Use (DD-MM-YYYY)"})
			return
		}
		releaseD = releaseDate
	}

	song, err := e.s.EditSong(c.Request.Context(), service.EditSongParam{
		Song_id:     uuid_song,
		GroupName:   songData.GroupName,
		SongName:    songData.SongName,
		ReleaseDate: releaseD,
		Text:        songData.Text,
		Link:        songData.Link,
	})
	if err != nil {
		if err == service.SongNotFoundError {
			c.JSON(http.StatusNotFound, MessageError{err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}
	c.JSON(http.StatusOK, Song{
		ID:          song.ID,
		GroupName:   song.GroupName,
		SongName:    song.SongName,
		ReleaseDate: song.ReleaseDate,
		Text:        song.Text,
		Link:        song.Link,
		CreatedAt:   song.CreatedAt,
	})
}

// @Summary New song
// @Schemes
// @Description create new song
// @Tags Songs
// @Accept json
// @Produce json
// @Param request body endpoint.NewSong true "query params"
// @Success 201 {object} endpoint.Song
// @Failure      400  {object}  endpoint.MessageError
// @Failure      500
// @Router /songs/new [post]
func (e *Endpoint) NewSongHandler(c *gin.Context) {
	var songData NewSong

	validate := validator.New()

	if err := c.ShouldBindJSON(&songData); err != nil {
		c.JSON(http.StatusBadRequest, MessageError{"invalid data"})
		return
	}

	if err := validate.Struct(songData); err != nil {
		c.JSON(http.StatusBadRequest, MessageError{"invalid fields"})
		return
	}

	song, err := e.s.NewSong(c.Request.Context(), service.NewSongParam{
		GroupName: songData.GroupName,
		SongName:  songData.SongName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}
	c.JSON(http.StatusCreated, Song{
		ID:          song.ID,
		GroupName:   song.GroupName,
		SongName:    song.SongName,
		ReleaseDate: song.ReleaseDate,
		Text:        song.Text,
		Link:        song.Link,
		CreatedAt:   song.CreatedAt,
	})
}

func mapToStruct(m map[string]interface{}, val interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, val)
}

func (e *Endpoint) TestHandler(c *gin.Context) {
	group := c.Query("group")
	song := c.Query("song")
	if group == "" && song == "" {
		c.Status(http.StatusBadRequest)
	}
	c.JSON(http.StatusOK, gin.H{
		"releaseDate": "16.07.2006",
		"text":        "Ooh baby, don't you know I suffer?\\nOoh baby, can you hear me moan?\\nYou caught me under false pretenses\\nHow long before you let me go?\\n\\nOoh\\nYou set my soul alight\\nOoh\\nYou set my soul alight",
		"link":        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	})
}
