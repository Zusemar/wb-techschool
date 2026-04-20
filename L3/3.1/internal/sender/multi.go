package sender

import (
	"context"
	"fmt"
	"log"

	"notifier/internal/domain"
)

type logSender struct{}

func (s *logSender) Send(_ context.Context, n *domain.Notification) error {
	log.Printf("[LOG] id=%s recipient=%q title=%q message=%q", n.ID, n.Recipient, n.Title, n.Message)
	return nil
}

type emailSender struct{}

func (s *emailSender) Send(_ context.Context, n *domain.Notification) error {
	log.Printf("[EMAIL] To: %s | Subject: %s | Body: %s", n.Recipient, n.Title, n.Message)
	return nil
}

type telegramSender struct{}

func (s *telegramSender) Send(_ context.Context, n *domain.Notification) error {
	log.Printf("[TELEGRAM] To: %s | %s: %s", n.Recipient, n.Title, n.Message)
	return nil
}

// MultiSender dispatches to the appropriate channel sender.
type MultiSender struct {
	senders map[domain.Channel]Sender
}

func NewMultiSender() *MultiSender {
	return &MultiSender{
		senders: map[domain.Channel]Sender{
			domain.ChannelLog:      &logSender{},
			domain.ChannelEmail:    &emailSender{},
			domain.ChannelTelegram: &telegramSender{},
		},
	}
}

func (m *MultiSender) Send(ctx context.Context, n *domain.Notification) error {
	s, ok := m.senders[n.Channel]
	if !ok {
		return fmt.Errorf("unknown channel: %s", n.Channel)
	}
	return s.Send(ctx, n)
}
