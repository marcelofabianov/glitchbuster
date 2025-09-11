package port

import (
	"context"

	"github.com/marcelofabianov/order-api/internal/domain"
)

type CreateOrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
}

type OrderRepository interface {
	CreateOrderRepository
}
