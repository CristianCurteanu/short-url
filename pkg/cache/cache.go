package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type UrlCache interface {
	Set(ctx context.Context, key, val string) error
	Get(ctx context.Context, key string) (string, error)
	Close() error
}

type cache struct {
	db *redis.Client
}

func NewRedisCache(url, password string) UrlCache {
	return &cache{redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
	})}
}

func (c *cache) Set(ctx context.Context, key string, val string) error {
	return c.db.Set(ctx, key, val, time.Second*5).Err()
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	url, err := c.db.Get(ctx, key).Bytes()
	if err != nil {
		return "", err
	}

	return string(url), nil
}

func (c *cache) Close() error {
	return c.db.Close()
}
