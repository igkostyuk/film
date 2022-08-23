package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"test/internal/models"
)

func Test_filmRepo_GetByTitle(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	filmRepo := NewFilmRepository(sqlxDB)

	t.Run("Get", func(t *testing.T) {
		id := 1
		title := "title"
		description := "content"

		rows := sqlmock.NewRows([]string{"film_id", "title", "description"}).AddRow(id, title, description)

		film := &models.Film{
			ID:          id,
			Title:       title,
			Description: description,
		}

		query := `SELECT film_id,title,description FROM film WHERE title = $1`
		mock.ExpectQuery(query).WithArgs(title).WillReturnRows(rows)

		returnedFilm, err := filmRepo.GetByTitle(context.Background(), title)

		require.NoError(t, err)
		require.NotNil(t, returnedFilm)
		require.Equal(t, film.Title, returnedFilm.Title)
	})
}
