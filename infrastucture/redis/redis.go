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

func (r *RedisClient) Keys(context context.Context, prefix string) ([]string, error) {
	result := r.Client.Keys(context, prefix)

	if err := result.Err(); err != nil {
		return nil, err
	}

	return result.Val(), nil
}

func (r *RedisClient) HGet(context context.Context, key, field string) *redis.StringCmd {
	return r.Client.HGet(context, key, field)
}

func (r *RedisClient) HSet(context context.Context, key string, values ...interface{}) error {
	return r.Client.HSet(context, key, values...).Err()
}

func (r *RedisClient) HGetAll(context context.Context, key string) *redis.StringStringMapCmd {
	return r.Client.HGetAll(context, key)
}

func (r *RedisClient) Set(context context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(context, key, value, expiration).Err()
}

func (r *RedisClient) Get(context context.Context, key string) (value *redis.StringCmd, err error) {
	result := r.Client.Get(context, key)

	if result.Err() != nil {
		return nil, err
	}

	return result, nil
}

func (r *RedisClient) Delete(context context.Context, key string) error {
	return r.Client.Del(context, key).Err()
}

func (r *RedisClient) GetWithDefault(context context.Context, key string, target interface{}, expiration time.Duration, funcDefaultData func() interface{}) error {
	if exist, err := r.Exist(context, key); err != nil {
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

		if err := r.Set(context, key, bytes, expiration); err != nil {
			return err
		}
	}

	if value, err := r.Get(context, key); err != nil {
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

func (r *RedisClient) Exist(context context.Context, key string) (exist bool, err error) {
	result := r.Client.Exists(context, key)

	if result.Err() != nil {
		return false, err
	}

	return result.Val() >= 1, nil
}
