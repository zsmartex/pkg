package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
}

func New(host string, port int, password string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis.Ping")
	}

	return &RedisClient{
		client,
	}, nil
}

func (r *RedisClient) Keys(context context.Context, prefix string) ([]string, error) {
	result := r.Client.Keys(context, prefix)

	if err := result.Err(); err != nil {
		return nil, errors.Wrap(err, "redis.Keys")
	}

	return result.Val(), nil
}

func (r *RedisClient) HGet(context context.Context, key, field string) *redis.StringCmd {
	return r.Client.HGet(context, key, field)
}

func (r *RedisClient) HSet(context context.Context, key string, values ...interface{}) error {
	err := r.Client.HSet(context, key, values...).Err()
	if err != nil {
		return errors.Wrap(err, "redis.HSet")
	}

	return nil
}

func (r *RedisClient) HExists(context context.Context, key, field string) (exist bool, err error) {
	result := r.Client.HExists(context, key, field)

	if result.Err() != nil {
		return false, errors.Wrap(err, "redis.HExists")
	}

	return result.Val(), nil
}

func (r *RedisClient) HGetAll(context context.Context, key string) *redis.MapStringStringCmd {
	return r.Client.HGetAll(context, key)
}

func (r *RedisClient) Set(context context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.Client.Set(context, key, value, expiration).Err()
	if err != nil {
		return errors.Wrap(err, "redis.Set")
	}

	return nil
}

func (r *RedisClient) Get(context context.Context, key string) (value *redis.StringCmd, err error) {
	result := r.Client.Get(context, key)
	if result.Err() != nil {
		return nil, errors.Wrap(err, "redis.Get")
	}

	return result, nil
}

func (r *RedisClient) Delete(context context.Context, key string) error {
	err := r.Client.Del(context, key).Err()
	if err != nil {
		return errors.Wrap(err, "redis.Delete")
	}

	return err
}

func (r *RedisClient) GetWithDefault(context context.Context, key string, target interface{}, expiration time.Duration, funcDefaultData func() (interface{}, error)) error {
	if exist, err := r.Exist(context, key); err != nil {
		return err
	} else if !exist {
		data, err := funcDefaultData()
		if err != nil {
			return errors.Wrap(err, "redis.GetWithDefault")
		}

		val := reflect.ValueOf(target)
		if val.Kind() != reflect.Ptr {
			return errors.New("data must be a pointer")
		}

		val.Elem().Set(reflect.ValueOf(data))

		bytes, err := json.Marshal(data)
		if err != nil {
			return errors.Wrap(err, "redis.GetWithDefault")
		}

		if err := r.Set(context, key, bytes, expiration); err != nil {
			return errors.Wrap(err, "redis.GetWithDefault")
		}
	}

	if value, err := r.Get(context, key); err != nil {
		return errors.Wrap(err, "redis.GetWithDefault")
	} else if value == nil {
		return nil
	} else {
		bytes, err := value.Bytes()
		if err != nil {
			return errors.Wrap(err, "redis.GetWithDefault")
		}

		return json.Unmarshal(bytes, target)
	}
}

func (r *RedisClient) Exist(context context.Context, key string) (exist bool, err error) {
	result := r.Client.Exists(context, key)

	if result.Err() != nil {
		return false, errors.Wrap(err, "redis.Exists")
	}

	return result.Val() >= 1, nil
}

func (r *RedisClient) Health(context context.Context) error {
	err := r.Client.Ping(context).Err()
	if err != nil {
		return errors.Wrap(err, "redis.Ping")
	}

	return nil
}

func (r *RedisClient) Close() error {
	err := r.Client.Close()
	if err != nil {
		return errors.Wrap(err, "redis.Ping")
	}

	return nil
}
