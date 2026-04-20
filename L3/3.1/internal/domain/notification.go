package domain

import (
	"errors"
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusQueued    Status = "queued"
	StatusSent      Status = "sent"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

type Channel string

const (
	ChannelLog      Channel = "log"
	ChannelEmail    Channel = "email"
	ChannelTelegram Channel = "telegram"
)

type Notification struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	Channel     Channel   `json:"channel"`
	Recipient   string    `json:"recipient"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var (
	ErrNotFound         = errors.New("notification not found")
	ErrAlreadySent      = errors.New("notification already sent")
	ErrAlreadyCancelled = errors.New("notification already cancelled")
)
