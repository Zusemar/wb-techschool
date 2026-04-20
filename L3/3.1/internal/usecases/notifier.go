package usecases

import (
	"time"

	"github.com/google/uuid"
	"notifier/internal/domain"
	"notifier/internal/repo"
)

type Notifier struct {
	repo repo.Repository
}

func NewNotifier(r repo.Repository) *Notifier {
	return &Notifier{repo: r}
}

type CreateRequest struct {
	Title       string
	Message     string
	Channel     domain.Channel
	Recipient   string
	ScheduledAt time.Time
}

func (n *Notifier) Create(req CreateRequest) (*domain.Notification, error) {
	channel := req.Channel
	if channel == "" {
		channel = domain.ChannelLog
	}

	notif := &domain.Notification{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Message:     req.Message,
		Channel:     channel,
		Recipient:   req.Recipient,
		ScheduledAt: req.ScheduledAt,
		Status:      domain.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := n.repo.Save(notif); err != nil {
		return nil, err
	}
	return notif, nil
}

func (n *Notifier) GetByID(id string) (*domain.Notification, error) {
	return n.repo.GetByID(id)
}

func (n *Notifier) Cancel(id string) error {
	notif, err := n.repo.GetByID(id)
	if err != nil {
		return err
	}
	switch notif.Status {
	case domain.StatusSent:
		return domain.ErrAlreadySent
	case domain.StatusCancelled:
		return domain.ErrAlreadyCancelled
	}
	return n.repo.UpdateStatus(id, domain.StatusCancelled)
}

func (n *Notifier) ListAll() ([]*domain.Notification, error) {
	return n.repo.ListAll()
}
