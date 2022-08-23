package config

import (
	"time"

	"github.com/igkostyuk/film/pkg/postgres"
	"github.com/igkostyuk/film/pkg/redis"
)

type Config struct {
	RedisTTL  time.Duration `default:"4s"`
	MemoryTTL time.Duration `default:"1s"`

	Postgres *postgres.Config
	Redis    *redis.Config
}
