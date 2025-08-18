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
	reader Reader
	repo   usecases.OrderRepository
	cache  *usecases.OrderCache
}

// Reader is a minimal interface of kafka-go reader for testability
type Reader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
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

// NewConsumerWithReader allows injecting a mock reader in tests
func NewConsumerWithReader(r Reader, repo usecases.OrderRepository, cache *usecases.OrderCache) *Consumer {
	return &Consumer{reader: r, repo: repo, cache: cache}
}

func (c *Consumer) Close() error { return c.reader.Close() }

func (c *Consumer) Run(ctx context.Context) error {
	for {
		message, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("kafka read error: %v", err)
			time.Sleep(time.Second)
			continue
		}

		var order domain.Order
		if err := json.Unmarshal(message.Value, &order); err != nil {
			log.Printf("invalid message: %v", err)
			continue
		}

		if err := c.repo.CreateOrder(ctx, &order); err != nil {
			log.Printf("db error: %v", err)
			continue
		}
		c.cache.Set(order)
	}
}
