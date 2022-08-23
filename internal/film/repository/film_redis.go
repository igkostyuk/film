package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/igkostyuk/film/internal/models"
)

// Film redis repository
type FilmRedisRepo struct {
	redisClient *redis.Client
}

// News redis repository constructor
func NewFilmRedisRepository(redisClient *redis.Client) *FilmRedisRepo {
	return &FilmRedisRepo{redisClient: redisClient}
}

// Get film by key
func (r *FilmRedisRepo) Get(ctx context.Context, key string) (*models.Film, error) {
	filmBytes, err := r.redisClient.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "filmRedisRepo.GetFilmByName.redisClient.Get")
	}
	film := &models.Film{}
	if err = json.Unmarshal(filmBytes, film); err != nil {
		return nil, errors.Wrap(err, "filmRedisRepo.GetFilmByName.json.Unmarshal")
	}

	return film, nil
}

// Cache film item
func (r *FilmRedisRepo) Set(ctx context.Context, key string, film *models.Film, d time.Duration) error {
	filmBytes, err := json.Marshal(film)
	if err != nil {
		return errors.Wrap(err, "filmRedisRepo.SetFilm.json.Marshal")
	}
	if err = r.redisClient.Set(ctx, key, filmBytes, d).Err(); err != nil {
		return errors.Wrap(err, "filmRedisRepo.SetFilm.redisClient.Set")
	}

	return nil
}

// Get all cached keys
func (r *FilmRedisRepo) GetKeysTTL(ctx context.Context) (map[string]time.Duration, error) {
	iter := r.redisClient.Scan(ctx, 0, "", 0).Iterator()
	keys := make(map[string]time.Duration)
	for iter.Next(ctx) {
		key := iter.Val()

		d, err := r.redisClient.PTTL(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		keys[key] = d
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}
