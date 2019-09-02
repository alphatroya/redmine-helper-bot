package storage

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type Manager interface {
	SetToken(token string, chat int64)
	GetToken(int64) (string, error)
	SetHost(host string, chat int64)
	GetHost(chat int64) (string, error)
	ResetData(chat int64) error
}

type RedisStorage struct {
	redis redisInstance
}

func (r RedisStorage) ResetData(chat int64) error {
	chatString := fmt.Sprint(chat)
	return r.redis.Del(chatString+"_token", chatString+"_host").Err()
}

func (r RedisStorage) SetToken(token string, chat int64) {
	r.redis.Set(fmt.Sprint(chat)+"_token", token, 0)
}

func (r RedisStorage) GetToken(chat int64) (string, error) {
	return r.redis.Get(fmt.Sprint(chat) + "_token").Result()
}

func (r RedisStorage) SetHost(host string, chat int64) {
	r.redis.Set(fmt.Sprint(chat)+"_host", host, 0)
}

func (r RedisStorage) GetHost(chat int64) (string, error) {
	return r.redis.Get(fmt.Sprint(chat) + "_host").Result()
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
	return &RedisStorage{redisClient}, err
}

type redisInstance interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Del(keys ...string) *redis.IntCmd
}
