package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"test/pkg/memory"
	"test/pkg/postgres"
	"test/pkg/redis"

	"test/internal/config"
	filmDelivery "test/internal/film/delivery/http"
	filmRepo "test/internal/film/repository"
	filmUCase "test/internal/film/usecase"
)

func main() {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err.Error())
	}
	log.Fatal(run(&cfg))
}

func run(cfg *config.Config) error {
	md := memory.New(time.Second, time.Second)

	rd, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		return err
	}

	pg, err := postgres.NewPsqlDB(cfg.Postgres)
	if err != nil {
		return err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	e := echo.New()

	e.Server.ReadTimeout = 5 * time.Second
	e.Server.WriteTimeout = 5 * time.Second

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
				zap.Duration("duration", v.Latency),
			)

			return nil
		},
	}))

	fmd := filmRepo.NewFilmMemoryRepository(md)
	frd := filmRepo.NewFilmRedisRepository(rd)
	fpg := filmRepo.NewFilmRepository(pg)

	fu := filmUCase.NewFilmUseCase(cfg, fmd, frd, fpg, logger)
	filmDelivery.NewFilmHandler(e, fu)

	return e.Start(":9090")
}
