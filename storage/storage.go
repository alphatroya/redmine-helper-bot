package storage

import (
	"github.com/go-redis/redis"
	"time"
)

type Manager interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
}

