package entity

import "time"

type Song struct {
	ID          int64     `json:"id"`
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SongFilter struct {
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	Page        int
	PageSize    int
}

type SongText struct {
	ID          int64    `json:"id"`
	GroupName   string   `json:"group_name"`
	SongName    string   `json:"song_name"`
	Verses      []string `json:"verses"`
	TotalVerses int      `json:"total_verses"`
	Page        int      `json:"page"`
	PageSize    int      `json:"page_size"`
}
