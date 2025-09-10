package application

import (
	"context"

	"github.com/marcelofabianov/fault"
	"github.com/marcelofabianov/wisp"

	"github.com/marcelofabianov/glitchbuster-order-api/internal/domain"
	"github.com/marcelofabianov/glitchbuster-order-api/internal/port"
)

type OrderService struct {
	repo       port.CreateOrderRepository
	productSvc port.ProductService
}

func NewOrderService(repo port.OrderRepository, productSvc port.ProductService) *OrderService {
	return &OrderService{
		repo:       repo,
		productSvc: productSvc,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, input *port.OrderServiceInput) (*port.OrderServiceOutput, error) {
	var totalAmount wisp.Money

	for _, item := range input.Items {
		price, err := s.productSvc.GetProductPriceByID(ctx, item.ProductID)
		if err != nil {
			return nil, fault.Wrap(err,
				"failed to get product price for ID",
				fault.WithCode(fault.NotFound),
				fault.WithContext("product_id", item.ProductID),
			)
		}

		itemTotal, err := item.Quantity.MultiplyByMoney(price)
		if err != nil {
			return nil, fault.Wrap(err,
				"failed to calculate total for item",
				fault.WithCode(fault.Internal),
				fault.WithContext("item", item),
				fault.WithContext("price", price.String()),
				fault.WithContext("product_id", item.ProductID.String()),
			)
		}

		totalAmount, err = totalAmount.Add(itemTotal)
		if err != nil {
			return nil, fault.Wrap(err,
				"failed to calculate total amount",
				fault.WithCode(fault.Internal),
				fault.WithContext("total_amount", totalAmount.String()),
				fault.WithContext("product_id", item.ProductID.String()),
			)
		}
	}

	data := domain.NewOrderInput{
		CustomerID:  input.CustomerID,
		Items:       input.Items,
		TotalAmount: totalAmount,
		CreatedBy:   input.CreatedBy,
	}

	order, err := domain.NewOrder(data)
	if err != nil {
		return nil, fault.Wrap(err,
			"failed to create order domain entity",
			fault.WithCode(fault.Internal),
		)
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fault.Wrap(err,
			"failed to create order",
			fault.WithCode(fault.Internal),
		)
	}

	return &port.OrderServiceOutput{
		OrderID: wisp.NewNullableUUID(order.ID),
		Message: "Order created successfully",
		Status:  true,
	}, nil
}
