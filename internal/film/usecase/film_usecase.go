package usecase

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/igkostyuk/film/internal/config"
	"github.com/igkostyuk/film/internal/models"
)

//go:generate mockgen -destination=../mocks/repository/cache.go -package=mocks_repository . Cache
//go:generate mockgen -destination=../mocks/repository/repository.go -package=mocks_repository . Repository

type Cache interface {
	Get(ctx context.Context, key string) (*models.Film, error)
	Set(ctx context.Context, key string, film *models.Film, d time.Duration) error
	GetKeysTTL(ctx context.Context) (map[string]time.Duration, error)
}

type Repository interface {
	GetByTitle(ctx context.Context, title string) (*models.Film, error)
}

// Film UseCase
type FilmUC struct {
	cfg *config.Config

	memory Cache
	redis  Cache
	db     Repository

	log *zap.Logger
}

// Film UseCase constructor
func NewFilmUseCase(cfg *config.Config, memory Cache, redis Cache, db Repository, log *zap.Logger) *FilmUC {
	return &FilmUC{cfg: cfg, memory: memory, redis: redis, db: db, log: log}
}

func (u *FilmUC) GetByTitle(ctx context.Context, title string) (film *models.Film, err error) {
	film, err = u.memory.Get(ctx, title)
	if err != nil {
		u.log.Error("memory cache get", zap.String("title", title), zap.Error(err))
	}
	if film != nil {
		u.log.Info("memory cache get", zap.String("title", title))

		return film, nil
	}

	film, err = u.redis.Get(ctx, title)
	if err != nil {
		u.log.Error("redis cache get", zap.String("title", title), zap.Error(err))
	}
	if film != nil {
		u.log.Info("redis cache get", zap.String("title", title))

		return film, nil
	}

	film, err = u.db.GetByTitle(ctx, title)
	if err != nil {
		return nil, err
	}
	if film == nil {
		return nil, nil
	}
	u.log.Info("repository cache get", zap.String("title", title))

	err = u.memory.Set(ctx, title, film, u.cfg.MemoryTTL)
	if err != nil {
		u.log.Error("memory cache set film", zap.String("title", title), zap.Error(err))
	}

	err = u.redis.Set(ctx, title, film, u.cfg.RedisTTL)
	if err != nil {
		u.log.Error("redis cache set film", zap.String("title", title), zap.Error(err))
	}

	return film, nil
}

func (u *FilmUC) GetKeysTTL(ctx context.Context) (map[string]models.Title, error) {
	rk, err := u.redis.GetKeysTTL(ctx)
	if err != nil {
		return nil, err
	}
	mk, err := u.memory.GetKeysTTL(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]models.Title)
	for k, v := range rk {
		m[k] = models.Title{Key: k, RedisTTL: v / time.Millisecond, MemoryTTL: mk[k] / time.Millisecond}
	}

	return m, nil
}
