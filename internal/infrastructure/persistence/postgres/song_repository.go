package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"song-library/internal/domain/entity"
	"song-library/internal/domain/repository"
	"song-library/pkg/logger"
)

type SongRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewSongRepository(db *sql.DB, logger *logger.Logger) *SongRepository {
	return &SongRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SongRepository) Create(ctx context.Context, song *entity.Song) error {
	r.logger.Debug(ctx, "Starting song creation in DB", 
		zap.String("group", song.GroupName),
		zap.String("song", song.SongName))

	query := `
		INSERT INTO songs (group_name, song_name, release_date, text, link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx, query,
		song.GroupName,
		song.SongName,
		song.ReleaseDate,
		song.Text,
		song.Link,
	).Scan(&song.ID, &song.CreatedAt, &song.UpdatedAt)

	if err != nil {
		r.logger.Error(ctx, "Failed to create song in DB", zap.Error(err))
		return fmt.Errorf("failed to create record: %w", err)
	}

	r.logger.Info(ctx, "Song successfully created in DB", 
		zap.Int64("id", song.ID),
		zap.Time("created_at", song.CreatedAt))
	return nil
}

func (r *SongRepository) Update(ctx context.Context, song *entity.Song) error {

	r.logger.Debug(ctx, "Starting song update in DB",
		zap.Int64("id", song.ID))

	query := `
		UPDATE songs 
		SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at`

	result, err := r.db.ExecContext(
		ctx, query,
		song.GroupName,
		song.SongName,
		song.ReleaseDate,
		song.Text,
		song.Link,
		song.ID,
	)
	if err != nil {
		r.logger.Error(ctx, "Failed to update song in DB", zap.Error(err))
		return fmt.Errorf("error updating record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error(ctx, "Failed to get affected rows", zap.Error(err))
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Warn(ctx, "Song not found during update", zap.Int64("id", song.ID))
		return repository.ErrSongNotFound
	}

	r.logger.Info(ctx, "Song successfully updated in DB", zap.Int64("id", song.ID))
	return nil
}

func (r *SongRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM songs WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting song: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rows == 0 {
		return repository.ErrSongNotFound
	}

	return nil
}

func (r *SongRepository) GetByID(ctx context.Context, id int64) (*entity.Song, error) {
	query := `
		SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at
		FROM songs
		WHERE id = $1`

	song := &entity.Song{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&song.ID,
		&song.GroupName,
		&song.SongName,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
		&song.CreatedAt,
		&song.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, repository.ErrSongNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting song: %w", err)
	}

	return song, nil
}

func (r *SongRepository) GetSongTextByVerses(ctx context.Context, id int64, page, pageSize int) (*entity.SongText, error) {
	r.logger.Debug(ctx, "Starting to get song text by verses", 
		zap.Int64("id", id),
		zap.Int("page", page),
		zap.Int("page_size", pageSize))

	query := `
		SELECT id, group_name, song_name, text
		FROM songs
		WHERE id = $1`

	var song entity.SongText
	var fullText string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&song.ID,
		&song.GroupName,
		&song.SongName,
		&fullText,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn(ctx, "Song not found", zap.Int64("id", id))
			return nil, repository.ErrSongNotFound
		}
		r.logger.Error(ctx, "Error getting song", zap.Error(err))
		return nil, fmt.Errorf("error getting song: %w", err)
	}

	verses := strings.Split(fullText, "\n\n")
	song.TotalVerses = len(verses)
	
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(verses) {
		end = len(verses)
	}
	if start < len(verses) {
		song.Verses = verses[start:end]
	}

	song.Page = page
	song.PageSize = pageSize

	return &song, nil
}

func (r *SongRepository) List(ctx context.Context, filter *entity.SongFilter) ([]*entity.Song, int, error) {
	r.logger.Debug(ctx, "Starting song list retrieval", zap.Any("filter", filter))

	var conditions []string
	var args []interface{}
	argNum := 1

	if filter.GroupName != "" {
		conditions = append(conditions, fmt.Sprintf("group_name ILIKE $%d", argNum))
		args = append(args, "%"+filter.GroupName+"%")
		argNum++
	}
	if filter.SongName != "" {
		conditions = append(conditions, fmt.Sprintf("song_name ILIKE $%d", argNum))
		args = append(args, "%"+filter.SongName+"%")
		argNum++
	}
	if !filter.ReleaseDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("DATE(release_date) = DATE($%d)", argNum))
		args = append(args, filter.ReleaseDate)
		argNum++
	}
	if filter.Text != "" {
		conditions = append(conditions, fmt.Sprintf("text ILIKE $%d", argNum))
		args = append(args, "%"+filter.Text+"%")
		argNum++
	}
	if filter.Link != "" {
		conditions = append(conditions, fmt.Sprintf("link ILIKE $%d", argNum))
		args = append(args, "%"+filter.Link+"%")
		argNum++
	}

	query := `SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at 
			  FROM songs WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM songs WHERE 1=1`

	if len(conditions) > 0 {
		condStr := strings.Join(conditions, " AND ")
		query += " AND " + condStr
		countQuery += " AND " + condStr
	}

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		r.logger.Error(ctx, "Failed to count total records", zap.Error(err))
		return nil, 0, fmt.Errorf("error counting total records: %w", err)
	}

	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)

	r.logger.Debug(ctx, "Executing query to DB", 
		zap.String("query", query),
		zap.Any("args", args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error(ctx, "Failed to execute query", zap.Error(err))
		return nil, 0, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var songs []*entity.Song
	for rows.Next() {
		song := &entity.Song{}
		err := rows.Scan(
			&song.ID,
			&song.GroupName,
			&song.SongName,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
			&song.CreatedAt,
			&song.UpdatedAt,
		)
		if err != nil {
			r.logger.Error(ctx, "Failed to scan result", zap.Error(err))
			return nil, 0, fmt.Errorf("error scanning result: %w", err)
		}
		songs = append(songs, song)
	}

	r.logger.Info(ctx, "Song list successfully retrieved", 
		zap.Int("total", total),
		zap.Int("retrieved", len(songs)))
	return songs, total, nil
}
