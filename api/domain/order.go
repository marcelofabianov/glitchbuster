package domain

import "github.com/marcelofabianov/wisp"

type OrderStatus string

const (
	StatusPending    OrderStatus = "PENDING"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusPaid       OrderStatus = "PAID"
	StatusShipped    OrderStatus = "SHIPPED"
	StatusCancelled  OrderStatus = "CANCELLED"
	StatusFailed     OrderStatus = "FAILED"
)

type Order struct {
	ID          wisp.UUID
	CustomerID  wisp.UUID
	Status      OrderStatus
	TotalAmount wisp.Money
	wisp.Audit
}

func NewOrder(customerID wisp.UUID, totalAmount wisp.Money, createdBy wisp.AuditUser) (*Order, error) {
	id, err := wisp.NewUUID()
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:          id,
		CustomerID:  customerID,
		Status:      StatusPending,
		TotalAmount: totalAmount,
		Audit:       wisp.NewAudit(createdBy),
	}, nil
}
