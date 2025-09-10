package port

import (
	"context"

	"github.com/marcelofabianov/wisp"

	"github.com/marcelofabianov/glitchbuster-order-api/internal/domain"
)

// --- Product ---

type ProductService interface {
	GetProductPriceByID(ctx context.Context, productID wisp.UUID) (wisp.Money, error)
}

// --- Order ---

type OrderServiceInput struct {
	CustomerID wisp.UUID
	Items      []domain.LineItem
	CreatedBy  wisp.AuditUser
}

type OrderServiceOutput struct {
	OrderID wisp.NullableUUID
	Message string
	Status  bool
}

type OrderService interface {
	CreateOrder(ctx context.Context, input *OrderServiceInput) (*OrderServiceOutput, error)
}
