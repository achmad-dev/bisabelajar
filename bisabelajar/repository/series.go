package repository

import (
	"bisabelajar/domain"
	"bisabelajar/dto"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type SeriesRepository struct {
	pgdb  *sqlx.DB
	redis *redis.Client
}

func NewSeriesRepository(pgdb *sqlx.DB, redis *redis.Client) SeriesRepository {
	return SeriesRepository{pgdb: pgdb, redis: redis}
}

func (b *SeriesRepository) InsertSeries(seriesDto dto.SeriesDto) (int, error) {
	var id int
	query := `
        INSERT INTO bisabelajar.series (title, description)
        VALUES ($1, $2)
        RETURNING id
    `
	err := b.pgdb.QueryRow(query, seriesDto.Title, seriesDto.Description).Scan(&id)
	if err != nil {
		return 0, err
	}

	// After inserting into the database, update the cache
	series := domain.Series{
		ID:          id,
		Title:       seriesDto.Title,
		Description: seriesDto.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	seriesJson, err := json.Marshal(series)
	if err != nil {
		return 0, err
	}

	// Write-back to Redis cache
	ctx := context.Background()
	err = b.redis.Set(ctx, getRedisKey(id), seriesJson, time.Hour*24).Err()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (b *SeriesRepository) GetSeriesByID(id int) (*domain.Series, error) {
	ctx := context.Background()
	redisKey := getRedisKey(id)

	// Try to get the series from the cache
	seriesJson, err := b.redis.Get(ctx, redisKey).Result()
	if err == nil {
		var series domain.Series
		if err := json.Unmarshal([]byte(seriesJson), &series); err == nil {
			return &series, nil
		}
	}

	// If not found in cache, fetch from the database
	query := `
        SELECT id, title, description, created_at, updated_at
        FROM bisabelajar.series
        WHERE id = $1
    `
	var series domain.Series
	err = b.pgdb.Get(&series, query, id)
	if err != nil {
		return nil, err
	}

	// Update the cache with the retrieved series
	seriesJsonUpdated, err := json.Marshal(series)
	if err != nil {
		return nil, err
	}
	err = b.redis.Set(ctx, redisKey, seriesJsonUpdated, time.Hour*24).Err()
	if err != nil {
		return nil, err
	}

	return &series, nil
}

func getRedisKey(id int) string {
	return "series:" + strconv.Itoa(id)
}
