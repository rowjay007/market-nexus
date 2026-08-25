package redis

import (
	"context"
	"errors"
	goredis "github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

func Connect(ctx context.Context, cfg Config) (*goredis.Client, error) {
	if cfg.Addr == "" {
		return nil, errors.New("redis address required")
	}
	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return client, nil
}
