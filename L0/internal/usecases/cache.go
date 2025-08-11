package usecases

import (
	"context"
	"sync"

	"wb-techschool/L0/internal/domain"
)

type OrderCache struct {
	mu    sync.RWMutex
	store map[string]domain.Order
}

func NewOrderCache() *OrderCache {
	return &OrderCache{store: make(map[string]domain.Order)}
}

func (c *OrderCache) Get(orderUID string) (domain.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	o, ok := c.store[orderUID]
	return o, ok
}

func (c *OrderCache) Set(o domain.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[o.Order_uid] = o
}

func (c *OrderCache) Warmup(ctx context.Context, repo OrderRepository) error {
	ids, err := repo.ListAllOrderUIDs(ctx)
	if err != nil {
		return err
	}
	for _, id := range ids {
		o, err := repo.GetOrderById(ctx, id)
		if err != nil {
			continue
		}
		c.Set(*o)
	}
	return nil
}
