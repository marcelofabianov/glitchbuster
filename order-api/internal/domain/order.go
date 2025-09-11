package domain

import "github.com/marcelofabianov/wisp"

type LineItem struct {
	ProductID wisp.UUID
	Quantity  wisp.Quantity
}

type OrderStatus string

const (
	StatusPending    OrderStatus = "PENDING"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusPaid       OrderStatus = "PAID"
	StatusShipped    OrderStatus = "SHIPPED"
	StatusCancelled  OrderStatus = "CANCELLED"
	StatusFailed     OrderStatus = "FAILED"
)

type NewOrderInput struct {
	CustomerID  wisp.UUID
	Items       []LineItem
	TotalAmount wisp.Money
	CreatedBy   wisp.AuditUser
}

type Order struct {
	ID          wisp.UUID
	CustomerID  wisp.UUID
	Items       []LineItem
	TotalAmount wisp.Money
	Status      OrderStatus
	wisp.Audit
}

func NewOrder(input NewOrderInput) (*Order, error) {
	id, err := wisp.NewUUID()
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:          id,
		CustomerID:  input.CustomerID,
		Status:      StatusPending,
		TotalAmount: input.TotalAmount,
		Items:       input.Items,
		Audit:       wisp.NewAudit(input.CreatedBy),
	}, nil
}
