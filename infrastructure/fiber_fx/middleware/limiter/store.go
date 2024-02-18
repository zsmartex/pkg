package limiter

import (
	"context"
	"time"

	"github.com/zsmartex/pkg/v2/infrastructure/redis_fx"
)

type RedisStore struct {
	*redis_fx.Client
}

func (r *RedisStore) Get(key string) ([]byte, error) {
	value, err := r.Client.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}

	return value.Bytes()
}

func (r *RedisStore) Set(key string, val []byte, exp time.Duration) error {
	return r.Client.Set(context.Background(), key, val, exp)
}

func (r *RedisStore) Delete(key string) error {
	return r.Client.Delete(context.Background(), key)
}

func (r *RedisStore) Reset() error {
	return nil
}

func (r *RedisStore) Close() error {
	return r.Client.Close()
}
