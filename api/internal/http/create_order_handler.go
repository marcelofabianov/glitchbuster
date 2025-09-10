package http

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/marcelofabianov/fault"
	"github.com/marcelofabianov/wisp"

	"github.com/marcelofabianov/glitchbuster-order-api/internal/domain"
	"github.com/marcelofabianov/glitchbuster-order-api/internal/port"
)

const (
	UnitItem wisp.Unit = "item"
)

type CreateOrderHandler struct {
	validator *validator.Validate
	orderSvc  port.OrderService
}

func NewCreateOrderHandler(orderSvc port.OrderService) *CreateOrderHandler {
	return &CreateOrderHandler{
		validator: validator.New(),
		orderSvc:  orderSvc,
	}
}

func (h *CreateOrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	customerID, err := wisp.ParseUUID(req.CustomerID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid customer ID"})
	}

	// TODO: get user from auth context
	auditUser, err := wisp.NewAuditUser("marcelofabianov@gmail.com")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	input := &port.OrderServiceInput{
		CustomerID: customerID,
		Items:      make([]domain.LineItem, len(req.Items)),
		CreatedBy:  auditUser,
	}

	for i, item := range req.Items {
		productID, err := wisp.ParseUUID(item.ProductID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid product ID for item"})
		}

		quantity, err := wisp.NewQuantity(item.Quantity, UnitItem)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid quantity for item"})
		}

		input.Items[i] = domain.LineItem{
			ProductID: productID,
			Quantity:  quantity,
		}
	}

	output, err := h.orderSvc.CreateOrder(c.Context(), input)
	if err != nil {
		var fErr *fault.Error
		if errors.As(err, &fErr) {
			switch fErr.Code {
			case fault.NotFound:
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": fErr.Message})
			case fault.Invalid, fault.DomainViolation:
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": fErr.Message})
			case fault.Conflict:
				return c.Status(http.StatusConflict).JSON(fiber.Map{"error": fErr.Message})
			default:
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "An unexpected error occurred."})
			}
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "An unexpected error occurred."})
	}

	return c.Status(http.StatusCreated).JSON(output)
}
