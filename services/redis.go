package services

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string) (*RedisClient, error) {
	opts, err := redis.ParseURL(addr)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	if err != nil {
		return nil, err
	}

	return &RedisClient{
		client: client,
	}, nil
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(context.Background(), key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (value *redis.StringCmd, err error) {
	result := r.client.Get(context.Background(), key)

	if result.Err() != nil {
		return nil, err
	}

	return result, nil
}

func (r *RedisClient) Exist(key string) (exist bool, err error) {
	result := r.client.Exists(context.Background(), key)

	if result.Err() != nil {
		return false, err
	}

	return result.Val() >= 1, nil
}
