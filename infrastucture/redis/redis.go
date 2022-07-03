package redis

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	*redis.Client
}

func New(addr string) (*RedisClient, error) {
	opts, err := redis.ParseURL(addr)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	if err != nil {
		return nil, err
	}

	return &RedisClient{
		client,
	}, nil
}

func (r *RedisClient) Keys(prefix string) ([]string, error) {
	result := r.Client.Keys(context.Background(), prefix)

	if err := result.Err(); err != nil {
		return nil, err
	}

	return result.Val(), nil
}

func (r *RedisClient) HGet(key, field string) *redis.StringCmd {
	return r.Client.HGet(context.Background(), key, field)
}

func (r *RedisClient) HSet(key string, values ...interface{}) error {
	return r.Client.HSet(context.Background(), key, values...).Err()
}

func (r *RedisClient) HGetAll(key string) *redis.StringStringMapCmd {
	return r.Client.HGetAll(context.Background(), key)
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(context.Background(), key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (value *redis.StringCmd, err error) {
	result := r.Client.Get(context.Background(), key)

	if result.Err() != nil {
		return nil, err
	}

	return result, nil
}

func (r *RedisClient) Delete(key string) error {
	return r.Client.Del(context.Background(), key).Err()
}

func (r *RedisClient) GetWithDefault(key string, target interface{}, expiration time.Duration, funcDefaultData func() interface{}) error {
	if exist, err := r.Exist(key); err != nil {
		return err
	} else if !exist {
		data := funcDefaultData()

		val := reflect.ValueOf(target)
		if val.Kind() != reflect.Ptr {
			panic("some: check must be a pointer")
		}

		val.Elem().Set(reflect.ValueOf(data))

		bytes, err := json.Marshal(data)
		if err != nil {
			return err
		}

		if err := r.Set(key, bytes, expiration); err != nil {
			return err
		}
	}

	if value, err := r.Get(key); err != nil {
		return err
	} else if value == nil {
		return nil
	} else {
		bytes, err := value.Bytes()
		if err != nil {
			return err
		}

		return json.Unmarshal(bytes, target)
	}
}

func (r *RedisClient) Exist(key string) (exist bool, err error) {
	result := r.Client.Exists(context.Background(), key)

	if result.Err() != nil {
		return false, err
	}

	return result.Val() >= 1, nil
}
