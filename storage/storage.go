package storage

import (
	"github.com/go-redis/redis"
	"os"
	"time"
)

type Manager interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
}

func NewStorageInstance(urlEnvironment string) (Manager, error) {
	opt, err := redis.ParseURL(os.Getenv(urlEnvironment))
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(opt)
	_, err = redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	return redisClient, err
}

