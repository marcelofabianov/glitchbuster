package http

import "github.com/marcelofabianov/wisp"

type LineItem struct {
	ProductID wisp.UUID     `json:"product_id" validate:"required"`
	Quantity  wisp.Quantity `json:"quantity" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	CustomerID wisp.UUID  `json:"customer_id" validate:"required"`
	Items      []LineItem `json:"items" validate:"required,min=1,dive"`
}
