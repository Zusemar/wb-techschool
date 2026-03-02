package usecases

import (
	"strings"
	"time"

	"calendar/internal/domain"
	"calendar/internal/repo"
)

type Calendar struct {
	repo repo.EventRepository
}

func NewCalendar(r repo.EventRepository) *Calendar {
	return &Calendar{repo: r}
}

func (c *Calendar) CreateEvent(userID int, date time.Time, text string) (int, error) {
	// валидация
	if userID <= 0 {
		return 0, domain.ErrInvalidUserID
	}
	if strings.TrimSpace(text) == "" {
		return 0, domain.ErrInvalidText
	}
	if date.IsZero() {
		return 0, domain.ErrInvalidDate
	}

	// создаем новую запись
	event := domain.Event{
		UserID: userID,
		Date:   date,
		Text:   text,
	}
	id, err := c.repo.Create(event)

	// возвращем id нового ивента
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *Calendar) UpdateEvent(id int, userID int, date time.Time, text string) error {
	if id <= 0 {
		return domain.ErrInvalidEventID
	}
	if userID <= 0 {
		return domain.ErrInvalidUserID
	}
	if strings.TrimSpace(text) == "" {
		return domain.ErrInvalidText
	}
	if date.IsZero() {
		return domain.ErrInvalidDate
	}
	event := domain.Event{
		ID:     id,
		UserID: userID,
		Date:   date,
		Text:   text,
	}
	err := c.repo.Update(event)
	if err != nil {
		return err
	}
	return nil
}

func (c *Calendar) DeleteEvent(id int) error {
	return c.repo.Delete(id)
}

func (c *Calendar) EventsForDay(userID int, date time.Time) ([]domain.Event, error) {
	if userID <= 0 {
		return nil, domain.ErrInvalidUserID
	}
	if date.IsZero() {
		return nil, domain.ErrInvalidDate
	}

	events, err := c.repo.GetByUser(userID)
	if err != nil {
		return nil, err
	}

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 1)

	var result []domain.Event
	for _, e := range events {
		if !e.Date.Before(start) && e.Date.Before(end) {
			result = append(result, e)
		}
	}

	return result, nil
}

func (c *Calendar) EventsForWeek(userID int, date time.Time) ([]domain.Event, error) {
	if userID <= 0 {
		return nil, domain.ErrInvalidUserID
	}
	if date.IsZero() {
		return nil, domain.ErrInvalidDate
	}

	events, err := c.repo.GetByUser(userID)
	if err != nil {
		return nil, err
	}

	// считаем, что неделя начинается с понедельника
	weekday := int(date.Weekday())
	offset := (weekday + 6) % 7
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).
		AddDate(0, 0, -offset)
	end := start.AddDate(0, 0, 7)

	var result []domain.Event
	for _, e := range events {
		if !e.Date.Before(start) && e.Date.Before(end) {
			result = append(result, e)
		}
	}

	return result, nil
}

func (c *Calendar) EventsForMonth(userID int, date time.Time) ([]domain.Event, error) {
	if userID <= 0 {
		return nil, domain.ErrInvalidUserID
	}
	if date.IsZero() {
		return nil, domain.ErrInvalidDate
	}

	events, err := c.repo.GetByUser(userID)
	if err != nil {
		return nil, err
	}

	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 1, 0)

	var result []domain.Event
	for _, e := range events {
		if !e.Date.Before(start) && e.Date.Before(end) {
			result = append(result, e)
		}
	}

	return result, nil
}
