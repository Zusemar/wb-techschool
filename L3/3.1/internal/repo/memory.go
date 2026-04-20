package repo

import (
	"sync"
	"time"

	"notifier/internal/domain"
)

type MemoryRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.Notification
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		data: make(map[string]*domain.Notification),
	}
}

func (r *MemoryRepo) Save(n *domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[n.ID] = n
	return nil
}

func (r *MemoryRepo) GetByID(id string) (*domain.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.data[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *n
	return &cp, nil
}

func (r *MemoryRepo) UpdateStatus(id string, status domain.Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, ok := r.data[id]
	if !ok {
		return domain.ErrNotFound
	}
	n.Status = status
	n.UpdatedAt = time.Now()
	return nil
}

func (r *MemoryRepo) ListPending() ([]*domain.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.Notification
	for _, n := range r.data {
		if n.Status == domain.StatusPending {
			cp := *n
			result = append(result, &cp)
		}
	}
	return result, nil
}

func (r *MemoryRepo) ListAll() ([]*domain.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Notification, 0, len(r.data))
	for _, n := range r.data {
		cp := *n
		result = append(result, &cp)
	}
	return result, nil
}
