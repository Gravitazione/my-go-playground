package redis

import (
	"context"
	"errors"
	"time"

	goredis "github.com/redis/go-redis/v9"

	apperror "github.com/Gravitazione/go-fiber-playground/internal/error"
)

const connectTimeout = 5 * time.Second

type Client struct {
	*goredis.Client
}

func New(ctx context.Context, url string) (*Client, error) {
	opts, err := goredis.ParseURL(url)
	if err != nil {
		return nil, apperror.Wrap(err, apperror.CodeInternal, "parse redis url")
	}

	client := goredis.NewClient(opts)

	pingCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		closeErr := client.Close()
		return nil, apperror.Wrap(errors.Join(err, closeErr), apperror.CodeServiceUnavailable, "ping redis")
	}

	return &Client{Client: client}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}
