package song

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"music-library/internal/models/song"
	"music-library/internal/storage/database"
	"music-library/pkg/standard/messages"
	"music-library/pkg/utils/logger"
	"net/http"
	"strconv"
	"strings"
)

// CreateSong godoc
// @Summary Create a new song
// @Description Create a new song with group and song name
// @Tags songs
// @Accept json
// @Produce json
// @Param group body string true "Group"
// @Param song body string true "Song"
// @Success 201 {object} song.Song
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/songs [post]
func (*songHandler) Create(context *gin.Context) {
	const op = "api.song.Create"
	var err error

	var body struct {
		Group string `json:"group" form:"group" binding:"required"`
		Song  string `json:"song" form:"song" binding:"required"`
	}

	logger.DebugLogOp(op, "request body received, parsing...")

	if err = context.BindJSON(&body); err != nil {
		logger.ErrorLogOp(op, fmt.Sprintf("cant proceed with body: %v", err.Error()))
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newSong := song.Song{
		Group:  body.Group,
		Song:   body.Song,
		Date:   nil,
		Lyrics: nil,
		Link:   nil,
	}

	if err = database.DB.Create(&newSong).Error; err != nil {
		logger.ErrorLogOp(op, fmt.Sprintf("cant create new song: %v", err.Error()))
		context.JSON(http.StatusInternalServerError, messages.InternalServerError)
		return
	}

	context.JSON(http.StatusCreated, gin.H{"data": &newSong})
}

// ListSongs godoc
// @Summary Get paginated list of songs
// @Description Get a list of songs with optional filtering by group and song name
// @Tags songs
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param group query string false "Filter by group"
// @Param song query string false "Filter by song"
// @Param showLyrics query string false "Show song lyrics if set to true"
// @Success 200 {object} map[string]interface{} "data, total, limit, page, pages"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/songs [get]
func (*songHandler) Index(context *gin.Context) {
	const op = "api.song.Index"

	var (
		songs  []song.Song
		total  int64
		limit  int
		page   int
		pages  int
		offset int
		err    error
	)

	pageStr := context.DefaultQuery("page", "1")
	limitStr := context.DefaultQuery("limit", "10")
	groupQuery := context.Query("group")
	songQuery := context.Query("song")
	showLyrics := context.Query("showLyrics")

	page, err = strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		logger.ErrorLogOp(op, "Invalid page number")
		return
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit number"})
		logger.ErrorLogOp(op, "Invalid limit number")
		return
	}

	offset = (page - 1) * limit

	db := database.DB.Model(&song.Song{})

	if groupQuery != "" {
		db = db.Where(`"group" ILIKE ?`, "%"+groupQuery+"%")
	}

	if songQuery != "" {
		db = db.Where(`"song" ILIKE ?`, "%"+songQuery+"%")
	}

	if showLyrics == "true" {
		db = db.Select("*")
	} else {
		db = db.Select(`id, "group", song, link, date, created_at, updated_at`)
	}

	if err = db.Count(&total).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": messages.InternalServerError})
		logger.Msg = fmt.Sprintf("Failed to count total: %v", err.Error())
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	pages = int(total) / limit
	if int(total)%limit != 0 {
		pages++
	}

	if err = db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": messages.InternalServerError})
		logger.Msg = fmt.Sprintf("failed to fetch data: %v", err.Error())
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data":  &songs,
		"total": total,
		"limit": limit,
		"page":  page,
		"pages": pages,
	})
}

// ShowSong godoc
// @Summary Get song by ID
// @Description Get details of a specific song by its UUID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Success 200 {object} song.Song
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/songs/{id} [get]
func (*songHandler) Show(context *gin.Context) {
	const op = "api.song.Show"

	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		logger.ErrorLogOp(op, messages.InternalServerError)
		logger.Msg = fmt.Sprintf("error validating uuid: %v", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": logger.Msg})
		return
	}

	var rawSong song.Song

	if err = database.DB.First(&rawSong, id).Error; err != nil {
		logger.Msg = "song not found"
		context.JSON(http.StatusNotFound, gin.H{"error": logger.Msg})
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": rawSong})
}

// UpdateSong godoc
// @Summary Update song by ID
// @Description Update the song details by its UUID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Param group body string true "Group"
// @Param song body string true "Song"
// @Param date body string true "Date"
// @Param lyrics body string true "Lyrics"
// @Param link body string true "Link"
// @Success 200 {object} song.Song
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/songs/{id} [put]
func (*songHandler) Update(context *gin.Context) {
	const op = "api.song.Update"

	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		logger.Msg = fmt.Sprintf("error validating uuid: %v", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": logger.Msg})
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	var body struct {
		Group  string `json:"group" form:"group" binding:"required"`
		Song   string `json:"song" form:"song" binding:"required"`
		Date   string `json:"date" form:"date" biding:"required"`
		Lyrics string `json:"lyrics" form:"lyrics" binding:"required"`
		Link   string `json:"link" form:"link" binding:"required"`
	}

	if err = context.BindJSON(&body); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Logger.Error(err.Error())
		return
	}

	var songData song.Song
	if err = database.DB.First(&songData, id).Error; err != nil {
		logger.Msg = "song not found"
		context.JSON(http.StatusNotFound, gin.H{"error": logger.Msg})
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	if err = database.DB.Model(&songData).Updates(&body).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": messages.InternalServerError})
		logger.Msg = fmt.Sprintf("failed to update song: %v", err.Error())
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	if err = database.DB.First(&songData, id).Error; err != nil {
		logger.Msg = "failed to fetch updated song"
		context.JSON(http.StatusInternalServerError, gin.H{"error": messages.InternalServerError})
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": &songData})
}

// GetSongLyrics godoc
// @Summary Get paginated song lyrics
// @Description Get lyrics of a song by its UUID, with pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Param page query int false "Page number"
// @Param limit query int false "Number of lines per page"
// @Success 200 {object} map[string]interface{} "data, page, limit, total"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/songs/{id}/lyrics [get]
func (*songHandler) Lyrics(context *gin.Context) {
	const op = "api.song.Lyrics"

	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		logger.Msg = fmt.Sprintf("error validating uuid: %v", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": messages.InternalServerError})
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	pageStr := context.DefaultQuery("page", "1")
	limitStr := context.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		logger.Msg = fmt.Sprintf("invalid page parameter: %v", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": logger.Msg})
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		logger.Msg = fmt.Sprintf("invalid limit parameter: %v", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": logger.Msg})
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	var songData song.Song
	if err = database.DB.First(&songData, id).Error; err != nil {
		logger.Msg = "song not found"
		context.JSON(http.StatusNotFound, gin.H{"error": logger.Msg})
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	lyricsLines := strings.Split(*songData.Lyrics, "\n")
	totalLines := len(lyricsLines)

	start := (page - 1) * limit
	if start >= totalLines {
		start = totalLines
	}
	end := start + limit
	if end > totalLines {
		end = totalLines
	}
	paginatedLyrics := lyricsLines[start:end]

	response := struct {
		Data  []string `json:"data"`
		Page  int      `json:"page"`
		Limit int      `json:"limit"`
		Total int      `json:"total"`
	}{
		Data:  paginatedLyrics,
		Page:  page,
		Limit: limit,
		Total: totalLines,
	}

	logger.DebugLogOp(op, fmt.Sprintf("%s: successfully fetched lyrics for song ID %s", op, id.String()))
	context.JSON(http.StatusOK, gin.H{"data": response})
}

// GetSongInfo godoc
// @Summary Get song info by group and song name
// @Description Get song lyrics, date, and link by specifying the group and song name
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string true "Group name"
// @Param song query string true "Song name"
// @Success 200 {object} song.Song
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/songs/info [get]
func (*songHandler) Info(context *gin.Context) {
	const op = "api.song.Info"

	groupQuery := context.Query("group")
	songQuery := context.Query("song")

	if groupQuery == "" || songQuery == "" {
		logger.Msg = "both group and song must be provided"
		logger.ErrorLogOp(op, logger.Msg)
		context.JSON(http.StatusBadRequest, gin.H{"error": logger.Msg})
		return
	}

	var songData song.Song
	if err := database.DB.Where("\"group\" = ? AND \"song\" = ?", groupQuery, songQuery).
		Select("lyrics, date, link, created_at, updated_at").
		First(&songData).Error; err != nil {
		logger.Msg = "song not found"
		context.JSON(http.StatusNotFound, gin.H{"error": logger.Msg})
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": &songData})
}

// DeleteSong godoc
// @Summary Delete a song by ID
// @Description Delete a specific song by its UUID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/songs/{id} [delete]
func (*songHandler) Delete(context *gin.Context) {
	const op = "api.song.Delete"

	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		logger.Msg = fmt.Sprintf("error validating uuid: %v", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": logger.Msg})
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	var songData song.Song
	if err = database.DB.First(&songData, id).Error; err != nil {
		logger.Msg = "song not found"
		context.JSON(http.StatusNotFound, gin.H{"error": logger.Msg})
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	if err = database.DB.Delete(&songData).Error; err != nil {
		logger.Msg = fmt.Sprintf("failed to delete song: %v", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": messages.InternalServerError})
		logger.ErrorLogOp(op, logger.Msg)
		return
	}

	context.JSON(http.StatusNoContent, gin.H{"message": "song deleted successfully"})
}
