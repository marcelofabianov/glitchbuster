package http

type LineItem struct {
	ProductID string  `json:"product_id" validate:"required,uuid7"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	CustomerID string     `json:"customer_id" validate:"required,uuid7"`
	Items      []LineItem `json:"items" validate:"required,min=1,dive"`
}
