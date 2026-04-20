package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wb-go/wbf/rabbitmq"
	"notifier/internal/domain"
)

const (
	ExchangeName = "notifications"
	RoutingKey   = "notify"
	QueueName    = "notifications"
)

type Publisher struct {
	pub *rabbitmq.Publisher
}

func NewPublisher(client *rabbitmq.RabbitClient) *Publisher {
	return &Publisher{
		pub: rabbitmq.NewPublisher(client, ExchangeName, "application/json"),
	}
}

func (p *Publisher) Publish(ctx context.Context, n *domain.Notification) error {
	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("marshal notification: %w", err)
	}
	return p.pub.Publish(ctx, body, RoutingKey)
}
