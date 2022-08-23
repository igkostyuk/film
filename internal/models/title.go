package models

import "time"

type Title struct {
	Key       string
	MemoryTTL time.Duration
	RedisTTL  time.Duration
}
