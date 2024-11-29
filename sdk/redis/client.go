package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Address        string
	Password       string
	DB             int
	DefaultChannel string
}

type Client struct {
	client         *redis.Client
	defaultChannel string
}

func NewClient(opts Config) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     opts.Address,
		Password: opts.Password,
		DB:       opts.DB,
	})

	return &Client{
		client:         client,
		defaultChannel: opts.DefaultChannel,
	}
}

func (r *Client) Publish(ctx context.Context, message interface{}) *redis.IntCmd {
	return r.client.Publish(ctx, r.defaultChannel, message)
}
func (r *Client) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}
