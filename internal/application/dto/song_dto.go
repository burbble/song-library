package dto

import (
	"song-library/internal/domain/entity"
	"time"
)

type CreateSongRequest struct {
	GroupName string `json:"group" binding:"required"`
	SongName  string `json:"song" binding:"required"`
}

type UpdateSongRequest struct {
	GroupName   string `json:"group_name" binding:"required"`
	SongName    string `json:"song_name" binding:"required"`
	ReleaseDate string `json:"release_date" binding:"required"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongResponse struct {
	ID          int64     `json:"id"`
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	ReleaseDate string    `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SongListRequest struct {
	GroupName   string `form:"group_name"`
	SongName    string `form:"song_name"`
	ReleaseDate string `form:"release_date"`
	Text        string `form:"text"`
	Link        string `form:"link"`
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
}

type SongListResponse struct {
	Songs      []SongResponse `json:"songs"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

type GetSongTextRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}

type SongTextResponse struct {
	ID          int64    `json:"id"`
	GroupName   string   `json:"group_name"`
	SongName    string   `json:"song_name"`
	Verses      []string `json:"verses"`
	TotalVerses int      `json:"total_verses"`
	Page        int      `json:"page"`
	PageSize    int      `json:"page_size"`
	TotalPages  int      `json:"total_pages"`
}

func ToSongResponse(song *entity.Song) SongResponse {
	return SongResponse{
		ID:          song.ID,
		GroupName:   song.GroupName,
		SongName:    song.SongName,
		ReleaseDate: song.ReleaseDate.Format("02-01-2006"),
		Text:        song.Text,
		Link:        song.Link,
		CreatedAt:   song.CreatedAt,
		UpdatedAt:   song.UpdatedAt,
	}
}
