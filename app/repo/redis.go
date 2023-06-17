package repo

import (
	"context"
	"github.com/redis/go-redis/v9"
	"smart-contract-service/app/usecase"
	"time"
)

type RedisConnection struct {
	client *redis.Client
}

func NewRedisConnection(client *redis.Client) usecase.RedisRepository {
	return &RedisConnection{client: client}
}

func (r *RedisConnection) Set(key string, data string) (err error) {
	err = r.client.Set(context.Background(), key, data, 5*time.Minute).Err()
	return err
}

func (r *RedisConnection) Get(key string) (val string, err error) {
	val, err = r.client.Get(context.Background(), key).Result()
	return val, err
}
