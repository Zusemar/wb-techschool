package repo

import (
	"calendar/internal/domain"
	"sync"
)

type memoryEventRepository struct {
	events map[int]domain.Event // int - id события
	mutex  sync.RWMutex
	nextID int
}

func (m *memoryEventRepository) Create(event domain.Event) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.nextID++
	event.ID = m.nextID
	m.events[event.ID] = event

	return event.ID, nil
}

func (m *memoryEventRepository) GetByUser(user_id int) ([]domain.Event, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := []domain.Event{}
	for _, event := range m.events {
		if event.UserID == user_id {
			result = append(result, event)
		}
	}

	return result, nil
}

func (m *memoryEventRepository) Update(event domain.Event) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.events[event.ID]; !ok {
		return domain.ErrEventNotFound
	}

	m.events[event.ID] = event
	return nil
}

func (m *memoryEventRepository) Delete(id int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.events[id]; !ok {
		return domain.ErrEventNotFound
	}

	delete(m.events, id)
	return nil
}

func NewMemoryRepository() *memoryEventRepository {
	return &memoryEventRepository{
		events: make(map[int]domain.Event),
	}
}
