package endpoint

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time" //nolint:gci

	"github.com/Vic07Region/musicLibrary/internal/lib/logger" //nolint:gci
	"github.com/Vic07Region/musicLibrary/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Endpoint struct {
	s   service.MusicService
	log *logger.Logger
}

func New(s service.MusicService, log *logger.Logger) *Endpoint {
	return &Endpoint{
		s:   s,
		log: log,
	}
}

// @BasePath /api/v1

// @Summary List songs
// @Schemes
// @Description fetching song list
// @Param   group      query     string     false  "group name"	example(Muse)
// @Param   song      query     string     false  "song name"	example(Supermassive Black Hole)
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
	var fetchParams service.FetchSongsRequest

	if val, ok := c.GetQuery("group"); ok {
		fetchParams.GroupName = &val
	}

	if val, ok := c.GetQuery("song"); ok {
		fetchParams.SongName = &val
	}

	if val, ok := c.GetQuery("text"); ok {
		fetchParams.SongText = &val
	}

	if val, ok := c.GetQuery("releaseDate"); ok {
		rd, err := time.Parse(val, "02.01.2006")
		if err != nil {
			c.JSON(http.StatusBadRequest, MessageError{"wrong format releaseDate: example 02.01.2006 "})
		}
		fetchParams.ReleaseDate = &rd
	}

	if val, ok := c.GetQuery("offset"); ok {
		if intval, err := strconv.Atoi(val); err == nil {
			fetchParams.Offset = uint64(intval)
		}

	}

	if val, ok := c.GetQuery("limit"); ok {
		if intval, err := strconv.Atoi(val); err == nil {
			fetchParams.Limit = uint64(intval)
		}

	}
	songsResp, err := e.s.FetchSongs(c.Request.Context(), fetchParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}

	c.JSON(http.StatusOK, songsResp)
}

// @Summary Song text
// @Schemes
// @Description fetching song text
// @Param        id   path      int  true  "Song ID"
// @Param   limit      query     int     false  "items limit"	example(10)
// @Param   offset      query     int     false "offset items"	example(2)
// @Tags Songs
// @Accept json
// @Produce json
// @Success 200 {object} service.FetchVersesResponse
// @Failure      400  {object}  endpoint.MessageError
// @Failure      404  {object}  endpoint.MessageError
// @Failure      500
// @Router /songs/{id} [get]
func (e *Endpoint) FetchSongTextHandler(c *gin.Context) {
	var fetchParams service.FetchVersesRequest

	paramID := c.Param("id")
	songId, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageError{"wrong format id"})
	}

	fetchParams.SongID = songId

	if val, ok := c.GetQuery("offset"); ok {
		if intval, err := strconv.Atoi(val); err == nil {
			fetchParams.Offset = intval
		}

	}

	if val, ok := c.GetQuery("limit"); ok {
		if intval, err := strconv.Atoi(val); err == nil {
			fetchParams.Limit = intval
		}

	}

	verseResp, err := e.s.FetchVerses(c.Request.Context(), fetchParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}

	c.JSON(http.StatusOK, verseResp)
}

// @Summary Delete Song
// @Schemes
// @Description deleting song
// @Tags Songs
// @Accept json
// @Produce json
// @Param        id   path      int  true  "Song ID"
// @Success 	 200  {object}  service.DeleteSongResponse
// @Failure      400  {object}  endpoint.MessageError
// @Failure      404  {object}  endpoint.MessageError
// @Failure      500
// @Router /songs/{id} [delete]
func (e *Endpoint) DeleteSongHandler(c *gin.Context) {
	paramID := c.Param("id")

	songID, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageError{"wrong format id"})
	}
	resp, err := e.s.DeleteSong(c.Request.Context(), service.DeleteSongRequest{SongID: songID})
	if err != nil {
		if errors.Is(err, service.ErrSongNotFound) {
			c.JSON(http.StatusNotFound, MessageError{err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)

}

// @Summary Edit Song
// @Schemes
// @Description edit song
// @Tags Songs
// @Accept json
// @Produce json
// @Param        id   path      int  true  "Song ID"
// @Param request body endpoint.UpdateSong true "query params"
// @Success 200 {object} service.UpdateSongResponse
// @Failure      400  {object}  endpoint.MessageError
// @Failure      404  {object}  endpoint.MessageError
// @Failure      500
// @Router /songs/{id} [patch]
func (e *Endpoint) UpdateSongHandler(c *gin.Context) {

	paramID := c.Param("id")
	songID, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageError{"wrong format id"})
	}

	var request service.UpdateSongRequest
	request.SongID = songID

	var inputData map[string]interface{}
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}

	allowedKeys := map[string]bool{
		"group":       true,
		"song":        true,
		"releaseDate": true,
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

	if groupName, ok := inputData["group"].(string); ok {
		request.GroupName = &groupName
	}

	if songName, ok := inputData["song"].(string); ok {
		request.SongName = &songName
	}

	if releaseDateStr, ok := inputData["releaseDate"].(string); ok {
		releaseDate, err := time.Parse("02.01.2006", releaseDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, MessageError{"Invalid date format. Use (DD.MM.YYYY)"})
			return
		}
		request.ReleaseDate = &releaseDate
	}

	if link, ok := inputData["link"].(string); ok {
		request.Link = &link
	}

	resp, err := e.s.UpdateSong(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, service.ErrSongNotFound) {
			c.JSON(http.StatusNotFound, MessageError{err.Error()})
			return
		}
		if errors.Is(err, service.ErrGroupNotFound) {
			c.JSON(http.StatusNotFound, MessageError{err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Edit Song Verse
// @Schemes
// @Description edit song verse
// @Tags Songs
// @Accept json
// @Produce json
// @Param        id   path      int  true  "Song ID"
// @Param request body endpoint.UpdateVerseRequest true "query params"
// @Success 200 {object} service.UpdateVerseResponse
// @Failure      400  {object}  endpoint.MessageError
// @Failure      404  {object}  endpoint.MessageError
// @Failure      500
// @Router /songs/{id}/verse [patch]
func (e *Endpoint) UpdateSongVerseHandler(c *gin.Context) {

	paramID := c.Param("id")
	songID, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageError{"wrong format id"})
	}

	var request service.UpdateVerseRequest
	request.SongID = songID

	var inputData map[string]interface{}
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}

	allowedKeys := map[string]bool{
		"verseNumber": true,
		"verseText":   true,
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

	if verseNumber, ok := inputData["verseNumber"].(int); ok {
		request.VerseNumber = verseNumber
	}

	if verseText, ok := inputData["verseText"].(string); ok {
		request.VerseText = verseText
	}

	resp, err := e.s.UpdateVerse(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, service.ErrSongNotFound) {
			c.JSON(http.StatusNotFound, MessageError{err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, MessageError{err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// NewSongHandler @Summary New song
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

	song, err := e.s.NewSong(c.Request.Context(), service.NewSongRequest{
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
		Link:        song.Link,
	})
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
