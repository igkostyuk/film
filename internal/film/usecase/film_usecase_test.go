package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/igkostyuk/film/internal/config"
	mock "github.com/igkostyuk/film/internal/film/mocks/repository"
	"github.com/igkostyuk/film/internal/models"
)

func Test_filmUC_GetByTitle(t *testing.T) {
	type fields struct {
		memory *mock.MockCache
		redis  *mock.MockCache
		db     *mock.MockRepository
	}

	ctx := context.Background()
	title := "title"
	testFilm := &models.Film{
		ID:          1,
		Title:       title,
		Description: "content",
	}

	memoryTTL := time.Second
	redisTTL := time.Minute

	tests := []struct {
		name      string
		cfg       *config.Config
		prepare   func(f *fields)
		wantFilm  *models.Film
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "get film from in memory cache",
			prepare: func(f *fields) {
				f.memory.EXPECT().Get(ctx, title).Return(testFilm, nil)
			},
			wantFilm:  testFilm,
			assertion: assert.NoError,
		},
		{
			name: "get film form redis cache",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.memory.EXPECT().Get(ctx, title).Return(nil, nil),
					f.redis.EXPECT().Get(ctx, title).Return(testFilm, nil),
				)
			},
			wantFilm:  testFilm,
			assertion: assert.NoError,
		},
		{
			name: "film not exist in database",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.memory.EXPECT().Get(ctx, title).Return(nil, nil),
					f.redis.EXPECT().Get(ctx, title).Return(nil, nil),
					f.db.EXPECT().GetByTitle(ctx, title).Return(nil, nil),
				)
			},
			wantFilm:  nil,
			assertion: assert.NoError,
		},
		{
			name: "film exist in database",
			cfg:  &config.Config{MemoryTTL: memoryTTL, RedisTTL: redisTTL},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.memory.EXPECT().Get(ctx, title).Return(nil, nil),
					f.redis.EXPECT().Get(ctx, title).Return(nil, nil),
					f.db.EXPECT().GetByTitle(ctx, title).Return(testFilm, nil),
					f.memory.EXPECT().Set(ctx, title, testFilm, memoryTTL).Return(nil),
					f.redis.EXPECT().Set(ctx, title, testFilm, redisTTL).Return(nil),
				)
			},
			wantFilm:  testFilm,
			assertion: assert.NoError,
		},
		{
			name: "film exist in database and caches get errors",
			cfg:  &config.Config{MemoryTTL: memoryTTL, RedisTTL: redisTTL},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.memory.EXPECT().Get(ctx, title).Return(nil, errors.New("test error")),
					f.redis.EXPECT().Get(ctx, title).Return(nil, errors.New("test error")),
					f.db.EXPECT().GetByTitle(ctx, title).Return(testFilm, nil),
					f.memory.EXPECT().Set(ctx, title, testFilm, memoryTTL).Return(nil),
					f.redis.EXPECT().Set(ctx, title, testFilm, redisTTL).Return(nil),
				)
			},
			wantFilm:  testFilm,
			assertion: assert.NoError,
		},
		{
			name: "film exist in database, memory set error",
			cfg:  &config.Config{MemoryTTL: memoryTTL, RedisTTL: redisTTL},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.memory.EXPECT().Get(ctx, title).Return(nil, nil),
					f.redis.EXPECT().Get(ctx, title).Return(nil, nil),
					f.db.EXPECT().GetByTitle(ctx, title).Return(testFilm, nil),
					f.memory.EXPECT().Set(ctx, title, testFilm, memoryTTL).Return(errors.New("test error")),
					f.redis.EXPECT().Set(ctx, title, testFilm, redisTTL).Return(nil),
				)
			},
			wantFilm:  testFilm,
			assertion: assert.NoError,
		},
		{
			name: "film exist in database, redis set error",
			cfg:  &config.Config{MemoryTTL: memoryTTL, RedisTTL: redisTTL},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.memory.EXPECT().Get(ctx, title).Return(nil, nil),
					f.redis.EXPECT().Get(ctx, title).Return(nil, nil),
					f.db.EXPECT().GetByTitle(ctx, title).Return(testFilm, nil),
					f.memory.EXPECT().Set(ctx, title, testFilm, memoryTTL).Return(nil),
					f.redis.EXPECT().Set(ctx, title, testFilm, redisTTL).Return(errors.New("test error")),
				)
			},
			wantFilm:  testFilm,
			assertion: assert.NoError,
		},
		{
			name: "database error",
			cfg:  &config.Config{MemoryTTL: memoryTTL, RedisTTL: redisTTL},
			prepare: func(f *fields) {
				gomock.InOrder(
					f.memory.EXPECT().Get(ctx, title).Return(nil, nil),
					f.redis.EXPECT().Get(ctx, title).Return(nil, nil),
					f.db.EXPECT().GetByTitle(ctx, title).Return(testFilm, errors.New("test error")),
				)
			},
			wantFilm:  nil,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				memory: mock.NewMockCache(ctrl),
				redis:  mock.NewMockCache(ctrl),
				db:     mock.NewMockRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			observedZapCore, _ := observer.New(zap.InfoLevel)
			observedLogger := zap.New(observedZapCore)

			u := NewFilmUseCase(tt.cfg, f.memory, f.redis, f.db, observedLogger)
			gotFilm, err := u.GetByTitle(ctx, title)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantFilm, gotFilm)
		})
	}
}
