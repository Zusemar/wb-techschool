package usecases

import (
	"context"
	"errors"
	"testing"

	"wb-techschool/L0/internal/domain"
)

type mockRepo struct {
	orders map[string]domain.Order
}

func (m *mockRepo) CreateOrder(ctx context.Context, o *domain.Order) error {
	m.orders[o.Order_uid] = *o
	return nil
}
func (m *mockRepo) GetOrderById(ctx context.Context, id string) (*domain.Order, error) {
	o, ok := m.orders[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &o, nil
}
func (m *mockRepo) UpdateOrder(ctx context.Context, o *domain.Order) error {
	m.orders[o.Order_uid] = *o
	return nil
}
func (m *mockRepo) DeleteOrder(ctx context.Context, id string) error {
	delete(m.orders, id)
	return nil
}
func (m *mockRepo) ListAllOrderUIDs(ctx context.Context) ([]string, error) {
	out := make([]string, 0, len(m.orders))
	for k := range m.orders {
		out = append(out, k)
	}
	return out, nil
}

func TestOrderCache_GetSet(t *testing.T) {
	c := NewOrderCache()
	o := domain.Order{Order_uid: "u1"}
	if _, ok := c.Get("u1"); ok {
		t.Fatalf("expected empty")
	}
	c.Set(o)
	got, ok := c.Get("u1")
	if !ok || got.Order_uid != "u1" {
		t.Fatalf("cache miss: %+v %v", got, ok)
	}
}

func TestOrderCache_Warmup(t *testing.T) {
	repo := &mockRepo{orders: map[string]domain.Order{
		"a": {Order_uid: "a"},
		"b": {Order_uid: "b"},
	}}
	c := NewOrderCache()
	if err := c.Warmup(context.Background(), repo); err != nil {
		t.Fatalf("warmup: %v", err)
	}
	if _, ok := c.Get("a"); !ok {
		t.Fatalf("no a in cache")
	}
	if _, ok := c.Get("b"); !ok {
		t.Fatalf("no b in cache")
	}
}

