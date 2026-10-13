package redisclient

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func New(addr string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return &Client{rdb: rdb}, nil
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) Set(
	ctx context.Context,
	key, value string,
	ttl time.Duration,
) error {
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Get(
	ctx context.Context,
	key string,
) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *Client) Delete(
	ctx context.Context,
	key string,
) error {
	return c.rdb.Del(ctx, key).Err()
}
