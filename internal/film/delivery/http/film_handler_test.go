package http

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	mock "github.com/igkostyuk/film/internal/film/mocks/usecase"
	"github.com/igkostyuk/film/internal/models"
)

func TestFilmHandler_GetByTitle(t *testing.T) {
	title := "title"
	ctx := context.TODO()

	testFilm := &models.Film{
		ID:          1,
		Title:       title,
		Description: "content",
	}

	type fields struct {
		FUsecase *mock.MockFilmUsecase
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		wantCode  int
		wantBody  string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "get film",
			prepare: func(f *fields) {
				f.FUsecase.EXPECT().GetByTitle(ctx, title).Return(testFilm, nil)
			},
			wantCode:  http.StatusOK,
			wantBody:  "{\"film_id\":1,\"title\":\"title\",\"description\":\"content\"}\n",
			assertion: assert.NoError,
		},
		{
			name: "film not found",
			prepare: func(f *fields) {
				f.FUsecase.EXPECT().GetByTitle(ctx, title).Return(testFilm, sql.ErrNoRows)
			},
			wantCode:  http.StatusNotFound,
			wantBody:  "{\"message\":\"not found\"}\n",
			assertion: assert.NoError,
		},
		{
			name: "unknown error",
			prepare: func(f *fields) {
				f.FUsecase.EXPECT().GetByTitle(ctx, title).Return(testFilm, errors.New("test error"))
			},
			wantCode:  http.StatusInternalServerError,
			wantBody:  "{\"message\":\"test error\"}\n",
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req, err := http.NewRequest(echo.GET, "/film", nil)
			assert.NoError(t, err)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				FUsecase: mock.NewMockFilmUsecase(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetParamNames("title")
			c.SetParamValues(title)

			handler := FilmHandler{
				FUsecase: f.FUsecase,
			}

			tt.assertion(t, handler.GetByTitle(c))
			assert.Equal(t, tt.wantCode, rec.Code)
			assert.Equal(t, tt.wantBody, rec.Body.String())
		})
	}
}
