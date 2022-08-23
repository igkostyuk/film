package http

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"test/internal/models"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

//go:generate mockgen -destination=../../mocks/usecase/film.go -package=mocks_usecase . FilmUsecase
type FilmUsecase interface {
	GetByTitle(context.Context, string) (*models.Film, error)
	GetKeysTTL(ctx context.Context) (map[string]models.Title, error)
}

// FilmHandler  represent the http handler for film
type FilmHandler struct {
	FUsecase FilmUsecase
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewFilmHandler(e *echo.Echo, us FilmUsecase) {
	handler := &FilmHandler{
		FUsecase: us,
	}
	e.GET("/film/:title", handler.GetByTitle)
	e.GET("/film/metric", handler.GetMetric)
}

// GetByName will get film by given name
func (a *FilmHandler) GetByTitle(c echo.Context) error {
	title := c.Param("title")
	ctx := c.Request().Context()
	film, err := a.FUsecase.GetByTitle(ctx, title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, ResponseError{Message: "not found"})
		}

		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, film)
}

func (a *FilmHandler) GetMetric(c echo.Context) error {
	ctx := c.Request().Context()
	m, err := a.FUsecase.GetKeysTTL(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, m)
}
