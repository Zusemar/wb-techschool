package usecases

import (
	"context"
	"wb-techschool/L0/internal/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrderById(ctx context.Context, id string) (*domain.Order, error)
	UpdateOrder(ctx context.Context, order *domain.Order) error
	DeleteOrder(ctx context.Context, id string) error
	ListAllOrderUIDs(ctx context.Context) ([]string, error)
}
