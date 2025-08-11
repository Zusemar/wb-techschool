package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"wb-techschool/L0/internal/domain"
	"wb-techschool/L0/internal/usecases"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers  []string
	Topic    string
	GroupID  string
	MinBytes int
	MaxBytes int
}

type Consumer struct {
	reader *kafka.Reader
	repo   usecases.OrderRepository
	cache  *usecases.OrderCache
}

func NewConsumer(cfg Config, repo usecases.OrderRepository, cache *usecases.OrderCache) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
		Topic:    cfg.Topic,
		MinBytes: cfg.MinBytes,
		MaxBytes: cfg.MaxBytes,
	})
	return &Consumer{reader: r, repo: repo, cache: cache}
}

func (c *Consumer) Close() error { return c.reader.Close() }

func (c *Consumer) Run(ctx context.Context) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("kafka read error: %v", err)
			time.Sleep(time.Second)
			continue
		}

		var o domain.Order
		if err := json.Unmarshal(m.Value, &o); err != nil {
			log.Printf("invalid message: %v", err)
			continue
		}

		// Persist then cache
		if err := c.repo.CreateOrder(ctx, &o); err != nil {
			log.Printf("db error: %v", err)
			continue
		}
		c.cache.Set(o)
	}
}
