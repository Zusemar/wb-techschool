package repo

import (
	"calendar/internal/domain"
	"sync"
)

type memoryEventRepository struct {
	events map[int]domain.Event
	mutex  sync.RWMutex
}

func NewMemoryRepository() EventRepository {
	return &memoryEventRepository{
		events: make(map[int]domain.Event),
		mutex:  sync.RWMutex{},
	}
}

func (m *memoryEventRepository) Create(event domain.Event) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	event.ID = len(m.events) + 1
	m.events[event.ID] = event
	return event.ID, nil
}

func (m *memoryEventRepository) Get(id int) (domain.Event, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	event, ok := m.events[id]
	if !ok {
		return domain.Event{}, domain.ErrEventNotFound
	}

	return event, nil
}

func (m *memoryEventRepository) GetByUser(userID int) ([]domain.Event, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var result []domain.Event
	for _, e := range m.events {
		if e.UserID == userID {
			result = append(result, e)
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
