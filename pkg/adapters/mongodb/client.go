package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Config struct {
	URI      string
	Database string
}

type Client struct {
	Raw      *mongo.Client
	Database *mongo.Database
}

func Connect(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.URI == "" {
		return nil, errors.New("mongodb uri required")
	}
	if cfg.Database == "" {
		return nil, errors.New("mongodb database required")
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}
	if err := c.Ping(ctx, nil); err != nil {
		_ = c.Disconnect(ctx)
		return nil, err
	}
	return &Client{Raw: c, Database: c.Database(cfg.Database)}, nil
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.Raw.Disconnect(ctx)
}
