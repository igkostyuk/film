package repository

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"

	"github.com/igkostyuk/film/internal/models"
)

func Test_filmRedisRepo(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	filmRedisRepo := NewFilmRedisRepository(client)

	film := &models.Film{
		ID:          1,
		Title:       "title",
		Description: "content",
	}

	returnedFilm, err := filmRedisRepo.Get(context.Background(), film.Title)
	require.Nil(t, returnedFilm)
	require.NoError(t, err)

	err = filmRedisRepo.Set(context.Background(), film.Title, film, time.Second)
	require.NoError(t, err)

	returnedFilm, err = filmRedisRepo.Get(context.Background(), film.Title)
	require.NotNil(t, returnedFilm)
	require.NoError(t, err)
}
