package repository

import (
	"context"
	"test/internal/models"
	"time"
)

//go:generate mockgen -destination=../mocks/database/memory.go -package=mocks_database . MemoryDB

type MemoryDB interface {
	Set(k string, x interface{}, d time.Duration)
	Get(k string) (interface{}, bool)
	ItemsTTL() map[string]time.Duration
}

// Film Repository
type FilmMemoryRepo struct {
	db MemoryDB
}

// Film repository constructor
func NewFilmMemoryRepository(db MemoryDB) *FilmMemoryRepo {
	return &FilmMemoryRepo{db: db}
}

// Get single film by title
func (r *FilmMemoryRepo) Get(_ context.Context, key string) (*models.Film, error) {
	x, found := r.db.Get(key)
	if !found {
		return nil, nil
	}
	film, ok := x.(*models.Film)
	if !ok {
		return nil, nil
	}

	return film, nil
}

func (r *FilmMemoryRepo) Set(_ context.Context, key string, film *models.Film, d time.Duration) error {
	r.db.Set(key, film, d)

	return nil
}
func (r *FilmMemoryRepo) GetKeysTTL(_ context.Context) (map[string]time.Duration, error) {
	return r.db.ItemsTTL(), nil
}
