package repo

import "notifier/internal/domain"

type Repository interface {
	Save(n *domain.Notification) error
	GetByID(id string) (*domain.Notification, error)
	UpdateStatus(id string, status domain.Status) error
	ListPending() ([]*domain.Notification, error)
	ListAll() ([]*domain.Notification, error)
}
