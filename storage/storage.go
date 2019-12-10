package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type Manager interface {
	SetToken(token string, chat int64)
	GetToken(int64) (string, error)
	SetHost(host string, chat int64)
	GetHost(chat int64) (string, error)
	SetActivity(activity string, chat int64)
	GetActivity(chat int64) (string, error)
	ResetData(chat int64) error
}

type RedisStorage struct {
	redis      redisInstance
	passphrase string
}

func (r RedisStorage) SetActivity(activity string, chat int64) {
	r.redis.Set(fmt.Sprint(chat)+activitySuffix, activity, 0)
}

func (r RedisStorage) GetActivity(chat int64) (string, error) {
	return r.redis.Get(fmt.Sprint(chat) + activitySuffix).Result()
}

const (
	hostSuffix           = "_host"
	encryptedTokenSuffix = "_encrypted"
	activitySuffix       = "_activity"
)

func (r RedisStorage) ResetData(chat int64) error {
	chatString := fmt.Sprint(chat)
	return r.redis.Del(chatString+hostSuffix, chatString+encryptedTokenSuffix, chatString+activitySuffix).Err()
}

func (r RedisStorage) SetToken(token string, chat int64) {
	chatString := fmt.Sprint(chat)
	bytes, err := encrypt([]byte(token), r.passphrase)
	if err != nil {
		return
	}
	r.redis.Set(chatString+encryptedTokenSuffix, bytes, 0)
}

func (r RedisStorage) GetToken(chat int64) (string, error) {
	newToken, err := r.redis.Get(fmt.Sprint(chat) + encryptedTokenSuffix).Result()
	return r.decryptToken(newToken, err)
}

func (r RedisStorage) decryptToken(token string, err error) (string, error) {
	if err != nil {
		return "", err
	}
	decryptedToken, err := decrypt([]byte(token), r.passphrase)
	if err != nil {
		return "", err
	}
	return string(decryptedToken), nil
}

func createHash(key string) (string, error) {
	hasher := md5.New()
	_, err := hasher.Write([]byte(key))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	hash, err := createHash(passphrase)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(hash))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	hash, err := createHash(passphrase)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(hash))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, err
}

func (r RedisStorage) SetHost(host string, chat int64) {
	r.redis.Set(fmt.Sprint(chat)+hostSuffix, host, 0)
}

func (r RedisStorage) GetHost(chat int64) (string, error) {
	return r.redis.Get(fmt.Sprint(chat) + hostSuffix).Result()
}

func NewStorageInstance(urlEnvironment string, storagePassphareKey string) (Manager, error) {
	opt, err := redis.ParseURL(os.Getenv(urlEnvironment))
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(opt)
	_, err = redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	passphrase := os.Getenv(storagePassphareKey)
	if len(passphrase) == 0 {
		return nil, errors.New("passphase for encrypting tokens is not set")
	}
	return &RedisStorage{redisClient, passphrase}, err
}

type redisInstance interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Del(keys ...string) *redis.IntCmd
}
