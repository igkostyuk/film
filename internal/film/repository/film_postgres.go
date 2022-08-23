package repository

import (
	"context"
	"test/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Film Repository
type FilmRepo struct {
	db *sqlx.DB
}

// Film repository constructor
func NewFilmRepository(db *sqlx.DB) *FilmRepo {
	return &FilmRepo{db: db}
}

// Get single film by title
func (r *FilmRepo) GetByTitle(ctx context.Context, title string) (*models.Film, error) {
	query := `SELECT film_id,title,description FROM film WHERE title = $1`

	film := &models.Film{}
	if err := r.db.GetContext(ctx, film, query, title); err != nil {
		return nil, errors.Wrap(err, "filmRepo.GetFilmByTitle.GetContext")
	}

	return film, nil
}
