package services

import (
	"context"
	"encoding/json"
	"reflect"
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
	result := r.client.Exists(context.Background(), key)

	if result.Err() != nil {
		return false, err
	}

	return result.Val() >= 1, nil
}
