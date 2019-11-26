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
	ResetData(chat int64) error
}

type RedisStorage struct {
	redis        redisInstance
	passphaseKey string
}

const (
	hostSuffix           = "_host"
	tokenSuffix          = "_token"
	encryptedTokenSuffix = "_encrypted"
)

func (r RedisStorage) ResetData(chat int64) error {
	chatString := fmt.Sprint(chat)
	return r.redis.Del(chatString+tokenSuffix, chatString+hostSuffix, chatString+encryptedTokenSuffix).Err()
}

func (r RedisStorage) SetToken(token string, chat int64) {
	chatString := fmt.Sprint(chat)
	r.redis.Del(chatString + tokenSuffix)
	passphrase := os.Getenv(r.passphaseKey)
	if len(passphrase) == 0 {
		return
	}
	r.redis.Set(chatString+encryptedTokenSuffix, encrypt([]byte(token), passphrase), 0)
}

func (r RedisStorage) GetToken(chat int64) (string, error) {
	oldToken, _ := r.redis.Get(fmt.Sprint(chat) + tokenSuffix).Result()
	newToken, err := r.redis.Get(fmt.Sprint(chat) + encryptedTokenSuffix).Result()
	if len(oldToken) != 0 && len(newToken) == 0 {
		r.SetToken(oldToken, chat)
		token, err := r.redis.Get(fmt.Sprint(chat) + encryptedTokenSuffix).Result()
		return r.decryptToken(token, err)
	}
	return r.decryptToken(newToken, err)
}

func (r RedisStorage) decryptToken(token string, err error) (string, error) {
	if err != nil {
		return "", err
	}
	passphrase := os.Getenv(r.passphaseKey)
	if len(passphrase) == 0 {
		return "", errors.New("passphase for encrypting tokens is not set")
	}
	decryptedToken := string(decrypt([]byte(token), passphrase))
	return decryptedToken, nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	return gcm.Seal(nonce, nonce, data, nil)
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
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
	return &RedisStorage{redisClient, storagePassphareKey}, err
}

type redisInstance interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Del(keys ...string) *redis.IntCmd
}
