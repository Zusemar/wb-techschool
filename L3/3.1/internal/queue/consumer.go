package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"notifier/internal/domain"
	"notifier/internal/repo"
	"notifier/internal/sender"
)

type Consumer struct {
	consumer *rabbitmq.Consumer
}

func NewConsumer(
	client *rabbitmq.RabbitClient,
	r repo.Repository,
	s sender.Sender,
) *Consumer {
	handler := func(ctx context.Context, msg amqp091.Delivery) error {
		var n domain.Notification
		if err := json.Unmarshal(msg.Body, &n); err != nil {
			log.Printf("ERROR: malformed message body: %v", err)
			return nil // ACK to discard unprocessable message
		}

		// Skip if notification was cancelled after being queued
		current, err := r.GetByID(n.ID)
		if err != nil {
			log.Printf("ERROR: notification %s not found in store: %v", n.ID, err)
			return nil
		}
		if current.Status == domain.StatusCancelled {
			log.Printf("INFO: notification %s was cancelled, skipping send", n.ID)
			return nil
		}

		// Retry send with exponential backoff (3 attempts: 1s, 2s, 4s)
		sendErr := retry.DoContext(ctx, retry.Strategy{
			Attempts: 3,
			Delay:    1 * time.Second,
			Backoff:  2.0,
		}, func() error {
			return s.Send(ctx, &n)
		})

		if sendErr != nil {
			log.Printf("ERROR: notification %s failed after retries: %v", n.ID, sendErr)
			_ = r.UpdateStatus(n.ID, domain.StatusFailed)
			return nil // ACK — we've exhausted retries, don't requeue
		}

		_ = r.UpdateStatus(n.ID, domain.StatusSent)
		log.Printf("INFO: notification %s sent via %s", n.ID, n.Channel)
		return nil
	}

	c := rabbitmq.NewConsumer(client, rabbitmq.ConsumerConfig{
		Queue:         QueueName,
		ConsumerTag:   "notifier-consumer",
		AutoAck:       false,
		Workers:       3,
		PrefetchCount: 10,
		Nack:          rabbitmq.NackConfig{Requeue: false},
	}, handler)

	return &Consumer{consumer: c}
}

func (c *Consumer) Start(ctx context.Context) error {
	return c.consumer.Start(ctx)
}
