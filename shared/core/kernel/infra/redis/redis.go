package redis

import "github.com/redis/go-redis/v9"

type Client struct {
	*redis.Client
}

func New(opt redis.Options) *Client {
	return &Client{
		Client: redis.NewClient(&opt),
	}
}
