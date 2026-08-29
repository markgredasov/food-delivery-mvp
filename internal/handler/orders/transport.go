package ordershandler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/order"
	server_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/server"
)

type handler struct {
	service orderService
}

type orderService interface {
	CreateOrder(ctx context.Context, o order.Order) (order.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (order.Order, error)
	ListRestaurantOrders(ctx context.Context, restaurantID uuid.UUID) ([]order.Order, error)
	AcceptOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID) (order.Order, error)
	RejectOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID, reason string) (order.Order, error)
	UpdateStatus(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID, statusString string) (order.Order, error)
}

func New(service orderService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Routes() []server_http.Route {
	return []server_http.Route{
		{
			Method:  http.MethodPost,
			Path:    "/orders",
			Handler: h.CreateOrder,
		},
		{
			Method:  http.MethodGet,
			Path:    "/orders/{id}",
			Handler: h.GetOrder,
		},
		{
			Method:  http.MethodGet,
			Path:    "/restaurants/orders",
			Handler: h.ListRestaurantOrders,
		},
		{
			Method:  http.MethodPost,
			Path:    "/restaurants/orders/{id}/accept",
			Handler: h.AcceptOrder,
		},
		{
			Method:  http.MethodPost,
			Path:    "/restaurants/orders/{id}/reject",
			Handler: h.RejectOrder,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/restaurants/orders/{id}/status",
			Handler: h.UpdateStatus,
		},
	}
}
