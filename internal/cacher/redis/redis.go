package redis

import (
	"context"
	"github.com/patyukin/bs-auth/internal/cacher"
	"github.com/redis/go-redis/v9"
	"time"
)

var _ cacher.Cacher = (*RedisClient)(nil)

type RedisClient struct {
	redis *redis.Client
}

func NewRedis(address string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
		Protocol: 3,  // specify 2 for RESP 2 or 3 for RESP 3
	})

	return &RedisClient{
		redis: client,
	}
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.redis.Get(ctx, key).Result()
}

func (c *RedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.redis.Set(ctx, key, value, expiration).Err()
}

func (c *RedisClient) Close() error {
	return c.redis.Close()
}

func (c *RedisClient) Delete(ctx context.Context, key string) error {
	return c.redis.Del(ctx, key).Err()
}

func (c *RedisClient) Exist(ctx context.Context, key string) (bool, error) {
	result, err := c.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}
