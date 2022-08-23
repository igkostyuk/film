package redis

import (
	"context"
	"time"

	rd "github.com/go-redis/redis/v8"
)

type Config struct {
	Addr         string
	MinIdleConns int
	PoolSize     int
	PoolTimeout  int
	Password     string
	DB           int
}

// Returns new redis client
func NewRedisClient(cfg *Config) (*rd.Client, error) {
	redisHost := cfg.Addr

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := rd.NewClient(&rd.Options{
		Addr:         redisHost,
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
		Password:     cfg.Password,
		DB:           cfg.DB,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
