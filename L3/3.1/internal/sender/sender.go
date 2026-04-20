package sender

import (
	"context"

	"notifier/internal/domain"
)

type Sender interface {
	Send(ctx context.Context, n *domain.Notification) error
}
