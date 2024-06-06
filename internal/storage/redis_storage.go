package storage

import (
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type RedisStorage struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisStorage(addr, password string) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisStorage{
		Client: client,
		Ctx:    context.Background(),
	}
}

func (r *RedisStorage) Get(key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}

func (r *RedisStorage) Set(key string, value string, expiration int) error {
	return r.Client.Set(r.Ctx, key, value, time.Duration(expiration)*time.Second).Err()
}

func (r *RedisStorage) Incr(key string) error {
	return r.Client.Incr(r.Ctx, key).Err()
}

func (r *RedisStorage) Clear() error {
	return r.Client.FlushDB(r.Ctx).Err()
}
