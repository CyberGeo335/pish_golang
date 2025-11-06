package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache struct {
	rdb *redis.Client
}

func New(addr string) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	return &Cache{rdb: rdb}
}

func (c *Cache) Set(key string, value string, ttl time.Duration) error {
	return c.rdb.Set(context.Background(), key, value, ttl).Err()
}

func (c *Cache) Get(key string) (string, error) {
	return c.rdb.Get(context.Background(), key).Result()
}

func (c *Cache) TTL(key string) (time.Duration, error) {
	return c.rdb.TTL(context.Background(), key).Result()
}
