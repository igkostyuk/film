package repository

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock "github.com/igkostyuk/film/internal/film/mocks/database"
	"github.com/igkostyuk/film/internal/models"
)

func Test_filmMemoryRepo_Get(t *testing.T) {
	title := "title"
	ctx := context.TODO()

	testFilm := &models.Film{
		ID:          1,
		Title:       title,
		Description: "content",
	}
	type fields struct {
		db *mock.MockMemoryDB
	}

	tests := []struct {
		name      string
		fields    fields
		prepare   func(f *fields)
		want      *models.Film
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "get film",
			prepare: func(f *fields) {
				f.db.EXPECT().Get(title).Return(testFilm, true)
			},
			want:      testFilm,
			assertion: assert.NoError,
		},
		{
			name: "film not found",
			prepare: func(f *fields) {
				f.db.EXPECT().Get(title).Return(nil, false)
			},
			want:      nil,
			assertion: assert.NoError,
		},
		{
			name: "invalid film type",
			prepare: func(f *fields) {
				f.db.EXPECT().Get(title).Return(struct{}{}, true)
			},
			want:      nil,
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				db: mock.NewMockMemoryDB(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}
			r := NewFilmMemoryRepository(f.db)
			got, err := r.Get(ctx, title)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_filmMemoryRepo_Set(t *testing.T) {
	title := "title"
	ctx := context.TODO()

	testFilm := &models.Film{
		ID:          1,
		Title:       title,
		Description: "content",
	}

	duration := time.Second

	type fields struct {
		db *mock.MockMemoryDB
	}

	tests := []struct {
		name      string
		fields    fields
		prepare   func(f *fields)
		want      *models.Film
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "set film",
			prepare: func(f *fields) {
				f.db.EXPECT().Set(title, testFilm, duration)
			},
			want:      testFilm,
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				db: mock.NewMockMemoryDB(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}
			r := NewFilmMemoryRepository(f.db)
			tt.assertion(t, r.Set(ctx, title, testFilm, duration))
		})
	}
}
