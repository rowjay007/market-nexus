package kafka

import (
	"context"
	"encoding/json"
	"errors"
	kgo "github.com/segmentio/kafka-go"
	"time"
)

type Publisher struct {
	writer *kgo.Writer
}

type PublisherConfig struct {
	Brokers []string
	Topic   string
}

func NewPublisher(cfg PublisherConfig) (*Publisher, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka brokers required")
	}
	if cfg.Topic == "" {
		return nil, errors.New("kafka topic required")
	}
	w := &kgo.Writer{
		Addr:         kgo.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kgo.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kgo.RequireAll,
	}
	return &Publisher{writer: w}, nil
}

func (p *Publisher) PublishJSON(ctx context.Context, key string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kgo.Message{Key: []byte(key), Value: body})
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
