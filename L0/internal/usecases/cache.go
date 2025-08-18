package usecases

import (
	"context"
	"sync"

	"wb-techschool/L0/internal/domain"
)

type OrderCache struct {
	mu         sync.RWMutex
	store      map[string]domain.Order
	orderQueue []string // FIFO очередь для отслеживания порядка
	maxSize    int      // максимальный размер кэша
}

func NewOrderCache() *OrderCache {
	return &OrderCache{store: make(map[string]domain.Order), orderQueue: make([]string, 0, 1000), maxSize: 1000}
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
	if _, exists := c.store[o.Order_uid]; !exists {
		c.orderQueue = append(c.orderQueue, o.Order_uid)
		if len(c.orderQueue) > c.maxSize {
			oldest := c.orderQueue[0]
			c.orderQueue = c.orderQueue[1:]
			delete(c.store, oldest)
		}
	}
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
