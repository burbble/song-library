package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"song-library/internal/application/dto"
	"song-library/internal/config"
	"song-library/internal/domain/entity"
	"song-library/internal/domain/repository"
	"song-library/pkg/logger"
)

type SongUseCase struct {
	repo   repository.SongRepository
	config *config.Config
}

func NewSongUseCase(repo repository.SongRepository, cfg *config.Config) *SongUseCase {
	return &SongUseCase{
		repo:   repo,
		config: cfg,
	}
}

type MusicInfoResponse struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (uc *SongUseCase) Create(ctx context.Context, req *dto.CreateSongRequest) (*dto.SongResponse, error) {
	log := logger.New("debug")
	
	log.Debug(ctx, "Starting song creation", 
		zap.String("group", req.GroupName),
		zap.String("song", req.SongName))

	musicInfo, err := uc.fetchMusicInfo(req.GroupName, req.SongName)
	if err != nil {
		log.Error(ctx, "Error getting song information", zap.Error(err))
		return nil, fmt.Errorf("error getting song information: %w", err)
	}

	releaseDate, err := time.Parse("02-01-2006", musicInfo.ReleaseDate)
	if err != nil {
		log.Error(ctx, "Error parsing release date", 
			zap.Error(err),
			zap.String("date", musicInfo.ReleaseDate))
		return nil, fmt.Errorf("error parsing release date: %w", err)
	}

	song := &entity.Song{
		GroupName:   req.GroupName,
		SongName:    req.SongName,
		ReleaseDate: releaseDate,
		Text:        musicInfo.Text,
		Link:        musicInfo.Link,
	}

	if err := uc.repo.Create(ctx, song); err != nil {
		log.Error(ctx, "Error creating song in DB", zap.Error(err))
		return nil, fmt.Errorf("error creating song: %w", err)
	}

	response := &dto.SongResponse{
		ID:          song.ID,
		GroupName:   song.GroupName,
		SongName:    song.SongName,
		ReleaseDate: song.ReleaseDate.Format("02-01-2006"),
		Text:        song.Text,
		Link:        song.Link,
		CreatedAt:   song.CreatedAt,
		UpdatedAt:   song.UpdatedAt,
	}

	log.Info(ctx, "Song successfully created", zap.Int64("id", song.ID))
	return response, nil
}

func (uc *SongUseCase) Update(ctx context.Context, id int64, req *dto.UpdateSongRequest) (*dto.SongResponse, error) {
	log := logger.New("debug")
	
	log.Debug(ctx, "Starting song update", zap.Int64("id", id))

	releaseDate, err := time.Parse("02-01-2006", req.ReleaseDate)
	if err != nil {
		log.Error(ctx, "Error parsing release date", 
			zap.Error(err),
			zap.String("date", req.ReleaseDate))
		return nil, fmt.Errorf("error parsing release date: %w", err)
	}

	song := &entity.Song{
		ID:          id,
		GroupName:   req.GroupName,
		SongName:    req.SongName,
		ReleaseDate: releaseDate,
		Text:        req.Text,
		Link:        req.Link,
	}

	if err := uc.repo.Update(ctx, song); err != nil {
		log.Error(ctx, "Error updating song", zap.Error(err))
		return nil, fmt.Errorf("error updating song: %w", err)
	}

	response := &dto.SongResponse{
		ID:          song.ID,
		GroupName:   song.GroupName,
		SongName:    song.SongName,
		ReleaseDate: song.ReleaseDate.Format("02-01-2006"),
		Text:        song.Text,
		Link:        song.Link,
		UpdatedAt:   song.UpdatedAt,
	}

	log.Info(ctx, "Song successfully updated", zap.Int64("id", id))
	return response, nil
}

func (uc *SongUseCase) Delete(ctx context.Context, id int64) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting song: %w", err)
	}
	return nil
}

func (uc *SongUseCase) Get(ctx context.Context, id int64) (*dto.SongResponse, error) {
	log := logger.New("debug")
	
	log.Debug(ctx, "Getting song by ID", zap.Int64("id", id))

	song, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		log.Error(ctx, "Error getting song", zap.Error(err))
		return nil, fmt.Errorf("error getting song: %w", err)
	}

	response := &dto.SongResponse{
		ID:          song.ID,
		GroupName:   song.GroupName,
		SongName:    song.SongName,
		ReleaseDate: song.ReleaseDate.Format("02-01-2006"),
		Text:        song.Text,
		Link:        song.Link,
		CreatedAt:   song.CreatedAt,
		UpdatedAt:   song.UpdatedAt,
	}

	log.Debug(ctx, "Song successfully retrieved", zap.Int64("id", id))
	return response, nil
}

func (uc *SongUseCase) fetchMusicInfo(group, song string) (*MusicInfoResponse, error) {
	log := logger.New("debug")
	ctx := context.Background()
	
	url := fmt.Sprintf("%s/info?group=%s&song=%s", uc.config.API.MusicInfoURL, group, song)
	log.Debug(ctx, "Sending request to external API", 
		zap.String("url", url),
		zap.String("group", group),
		zap.String("song", song))
	
	resp, err := http.Get(url)
	if err != nil {
		log.Error(ctx, "Error sending API request", zap.Error(err))
		return nil, fmt.Errorf("error sending API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error(ctx, "API returned unexpected status", 
			zap.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var musicInfo MusicInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&musicInfo); err != nil {
		log.Error(ctx, "Error decoding response", zap.Error(err))
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	log.Debug(ctx, "Successfully retrieved song information", 
		zap.Any("music_info", musicInfo))
	return &musicInfo, nil
}

func (uc *SongUseCase) GetSongText(ctx context.Context, id int64, req *dto.GetSongTextRequest) (*dto.SongTextResponse, error) {
	songText, err := uc.repo.GetSongTextByVerses(ctx, id, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("error getting song text: %w", err)
	}

	totalPages := (songText.TotalVerses + req.PageSize - 1) / req.PageSize

	return &dto.SongTextResponse{
		ID:          songText.ID,
		GroupName:   songText.GroupName,
		SongName:    songText.SongName,
		Verses:      songText.Verses,
		TotalVerses: songText.TotalVerses,
		Page:        songText.Page,
		PageSize:    songText.PageSize,
		TotalPages:  totalPages,
	}, nil
}

func (uc *SongUseCase) List(ctx context.Context, req *dto.SongListRequest) (*dto.SongListResponse, error) {
	filter := &entity.SongFilter{
		GroupName: req.GroupName,
		SongName:  req.SongName,
		Text:      req.Text,
		Link:      req.Link,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}

	if req.ReleaseDate != "" {
		releaseDate, err := time.Parse("2006-01-02", req.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("invalid release date format: %w", err)
		}
		filter.ReleaseDate = releaseDate
	}

	songs, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error getting song list: %w", err)
	}

	totalPages := (total + req.PageSize - 1) / req.PageSize

	var songResponses []dto.SongResponse
	for _, song := range songs {
		songResponses = append(songResponses, dto.ToSongResponse(song))
	}

	return &dto.SongListResponse{
		Songs:      songResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}
