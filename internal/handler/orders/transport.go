package ordershandler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/order"
	server_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/server"
)

type handler struct {
	orders  orderService
	catalog catalogService
}

type orderService interface {
	CreateOrder(ctx context.Context, o order.Order) (order.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (order.Order, error)
	ListRestaurantOrders(ctx context.Context, restaurantID uuid.UUID) ([]order.Order, error)
	AcceptOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID) (order.Order, error)
	RejectOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID, reason string) (order.Order, error)
	UpdateStatus(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID, statusString string) (order.Order, error)
}

type catalogService interface {
	ReplaceMenu(ctx context.Context, restaurantID uuid.UUID, items []menu.MenuItem) ([]menu.MenuItem, error)
}

func New(orders orderService, catalog catalogService) *handler {
	return &handler{
		orders:  orders,
		catalog: catalog,
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
		{
			Method:  http.MethodPut,
			Path:    "/restaurants/menu",
			Handler: h.UpdateRestaurantMenu,
		},
	}
}
