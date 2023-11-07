package redis_fx

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
	"github.com/zsmartex/pkg/v2/config"
	"go.uber.org/fx"
)

type Client struct {
	client *redis.Client
}

type redisParams struct {
	fx.In

	Config config.Redis
}

func New(params redisParams) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", params.Config.Host, params.Config.Port),
		Password: params.Config.Password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "redis.Ping")
	}

	return &Client{
		client,
	}, nil
}

func (r *Client) Keys(context context.Context, prefix string) ([]string, error) {
	result := r.client.Keys(context, prefix)

	if err := result.Err(); err != nil {
		return nil, errors.Wrap(err, "redis.Keys")
	}

	return result.Val(), nil
}

func (r *Client) HGet(context context.Context, key, field string) *redis.StringCmd {
	return r.client.HGet(context, key, field)
}

func (r *Client) HSet(context context.Context, key string, values ...interface{}) error {
	err := r.client.HSet(context, key, values...).Err()
	if err != nil {
		return errors.Wrap(err, "redis.HSet")
	}

	return nil
}

func (r *Client) HExists(context context.Context, key, field string) (exist bool, err error) {
	result := r.client.HExists(context, key, field)

	if result.Err() != nil {
		return false, errors.Wrap(err, "redis.HExists")
	}

	return result.Val(), nil
}

func (r *Client) HGetAll(context context.Context, key string) *redis.MapStringStringCmd {
	return r.client.HGetAll(context, key)
}

func (r *Client) Set(context context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(context, key, value, expiration).Err()
	if err != nil {
		return errors.Wrap(err, "redis.Set")
	}

	return nil
}

func (r *Client) Get(context context.Context, key string) (value *redis.StringCmd, err error) {
	result := r.client.Get(context, key)
	if result.Err() != nil {
		return nil, errors.Wrap(err, "redis.Get")
	}

	return result, nil
}

func (r *Client) Delete(context context.Context, key string) error {
	err := r.client.Del(context, key).Err()
	if err != nil {
		return errors.Wrap(err, "redis.Delete")
	}

	return err
}

func (r *Client) GetWithDefault(context context.Context, key string, target interface{}, expiration time.Duration, funcDefaultData func() (interface{}, error)) error {
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

func (r *Client) Exist(context context.Context, key string) (exist bool, err error) {
	result := r.client.Exists(context, key)

	if result.Err() != nil {
		return false, errors.Wrap(err, "redis.Exists")
	}

	return result.Val() >= 1, nil
}

func (r *Client) Health(context context.Context) error {
	err := r.client.Ping(context).Err()
	if err != nil {
		return errors.Wrap(err, "redis.Health")
	}

	return nil
}

func (r *Client) Close() error {
	err := r.client.Close()
	if err != nil {
		return errors.Wrap(err, "redis.Close")
	}

	return nil
}
