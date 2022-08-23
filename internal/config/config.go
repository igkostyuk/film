package config

import (
	"test/pkg/postgres"
	"test/pkg/redis"
	"time"
)

type Config struct {
	RedisTTL  time.Duration `default:"4s"`
	MemoryTTL time.Duration `default:"1s"`

	Postgres *postgres.Config
	Redis    *redis.Config
}
