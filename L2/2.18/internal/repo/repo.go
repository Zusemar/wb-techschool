package repo

import "calendar/internal/domain"

type EventRepository interface {
	Create(event domain.Event) (int, error)
	Get(id int) (domain.Event, error)
	GetByUser(userID int) ([]domain.Event, error)
	Update(event domain.Event) error
	Delete(id int) error
}
