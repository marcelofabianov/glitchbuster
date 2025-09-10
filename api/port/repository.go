package port

import "github.com/marcelofabianov/glitchbuster-order-api/domain"

type OrderRepository interface {
	Create(order *domain.Order) error
}
