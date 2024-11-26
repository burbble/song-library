package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"song-library/internal/application/dto"
	"song-library/internal/application/usecase"
	"song-library/internal/domain/entity"
	"song-library/internal/domain/repository"
	"song-library/pkg/logger"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SongHandler struct {
	useCase usecase.SongUseCase
	logger  *logger.Logger
}

func NewSongHandler(useCase usecase.SongUseCase, logger *logger.Logger) *SongHandler {
	return &SongHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// Create godoc
// @Summary Create a new song
// @Description Creates a new song based on group and title
// @Tags songs
// @Accept json
// @Produce json
// @Param request body dto.CreateSongRequest true "Song data"
// @Success 201 {object} dto.SongResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/songs [post]
func (h *SongHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	
	h.logger.Debug(ctx, "Starting song creation request processing")
	
	var req dto.CreateSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	h.logger.Debug(ctx, "Request data received", zap.Any("request", req))

	song, err := h.useCase.Create(ctx, &req)
	if err != nil {
		h.logger.Error(ctx, "Failed to create song", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	h.logger.Info(ctx, "Song successfully created", 
		zap.Int64("id", song.ID),
		zap.String("group", song.GroupName),
		zap.String("song", song.SongName))
		
	c.JSON(http.StatusCreated, song)
}

// Update godoc
// @Summary Update a song
// @Description Updates an existing song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param request body dto.UpdateSongRequest true "Update data"
// @Success 200 {object} dto.SongResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/songs/{id} [put]
func (h *SongHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error(ctx, "Failed to parse ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID"})
		return
	}
	
	h.logger.Debug(ctx, "Starting song update request processing", zap.Int64("id", id))

	var req dto.UpdateSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid data format"})
		return
	}

	releaseDate, err := time.Parse("02-01-2006", req.ReleaseDate)
	if err != nil {
		h.logger.Error(ctx, "Failed to parse release date", 
			zap.Error(err),
			zap.String("date_string", req.ReleaseDate))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid date format. Use format: DD-MM-YYYY (e.g., 01-01-2024)"})
		return
	}

	song := &entity.Song{
		ID:          id,
		GroupName:   req.GroupName,
		SongName:    req.SongName,
		ReleaseDate: releaseDate,
		Text:        req.Text,
		Link:        req.Link,
	}

	h.logger.Debug(ctx, "Request data received", zap.Any("song", song))

	updatedSong, err := h.useCase.Update(ctx, id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrSongNotFound) {
			h.logger.Warn(ctx, "Song not found", zap.Int64("id", id))
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "song not found"})
			return
		}
		h.logger.Error(ctx, "Failed to update song", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	h.logger.Info(ctx, "Song successfully updated", zap.Int64("id", id))
	c.JSON(http.StatusOK, updatedSong)
}

// Delete godoc
// @Summary Delete a song
// @Description Deletes a song by ID
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/songs/{id} [delete]
func (h *SongHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error(ctx, "Failed to parse ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID"})
		return
	}

	h.logger.Debug(ctx, "Starting song deletion request processing", zap.Int64("id", id))

	if err := h.useCase.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrSongNotFound) {
			h.logger.Warn(ctx, "Song not found during deletion attempt", zap.Int64("id", id))
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "song not found"})
			return
		}
		h.logger.Error(ctx, "Failed to delete song", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	h.logger.Info(ctx, "Song successfully deleted", zap.Int64("id", id))
	c.Status(http.StatusNoContent)
}

// Get godoc
// @Summary Get a song
// @Description Gets a song by ID
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} dto.SongResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/songs/{id} [get]
func (h *SongHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error(ctx, "Failed to parse ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID"})
		return
	}

	h.logger.Debug(ctx, "Starting song retrieval request processing", zap.Int64("id", id))

	song, err := h.useCase.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrSongNotFound) {
			h.logger.Warn(ctx, "Song not found", zap.Int64("id", id))
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "song not found"})
			return
		}
		h.logger.Error(ctx, "Failed to retrieve song", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	h.logger.Info(ctx, "Song successfully retrieved", zap.Int64("id", song.ID))
	c.JSON(http.StatusOK, song)
}

// List godoc
// @Summary List of songs
// @Description Gets a list of songs with filtering and pagination
// @Tags songs
// @Produce json
// @Param group_name query string false "Group name"
// @Param song_name query string false "Song name"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.SongListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/songs [get]
func (h *SongHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	
	h.logger.Debug(ctx, "Starting song list retrieval request processing")

	var req dto.SongListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind query parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	h.logger.Debug(ctx, "Filtering parameters received", zap.Any("filter", req))

	songs, err := h.useCase.List(ctx, &req)
	if err != nil {
		h.logger.Error(ctx, "Failed to retrieve song list", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	h.logger.Info(ctx, "Song list successfully retrieved", 
		zap.Int("total", songs.Total),
		zap.Int("page", songs.Page),
		zap.Int("page_size", songs.PageSize))

	c.JSON(http.StatusOK, songs)
}

// GetSongText godoc
// @Summary Get song text with pagination by verses
// @Description Returns the song text, split into verses with pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.SongTextResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/songs/{id}/text [get]
func (h *SongHandler) GetSongText(c *gin.Context) {
	ctx := c.Request.Context()
	
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error(ctx, "Failed to parse ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID format"})
		return
	}

	h.logger.Debug(ctx, "Starting song text retrieval request processing", zap.Int64("id", id))

	var req dto.GetSongTextRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind query parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	h.logger.Debug(ctx, "Pagination parameters received", 
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize))

	response, err := h.useCase.GetSongText(ctx, id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrSongNotFound) {
			h.logger.Warn(ctx, "Song text not found", zap.Int64("id", id))
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "song not found"})
			return
		}
		h.logger.Error(ctx, "Failed to retrieve song text", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	h.logger.Info(ctx, "Song text successfully retrieved", 
		zap.Int64("id", response.ID),
		zap.Int("total_verses", response.TotalVerses),
		zap.Int("page", response.Page))

	c.JSON(http.StatusOK, response)
}
