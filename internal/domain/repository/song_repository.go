package repository

import (
	"context"
	"song-library/internal/domain/entity"
)

type SongRepository interface {
	Create(ctx context.Context, song *entity.Song) error
	Update(ctx context.Context, song *entity.Song) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*entity.Song, error)
	List(ctx context.Context, filter *entity.SongFilter) ([]*entity.Song, int, error)
	GetSongTextByVerses(ctx context.Context, id int64, page, pageSize int) (*entity.SongText, error)
}
