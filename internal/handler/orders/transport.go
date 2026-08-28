package ordershandler

import (
	"context"
	"net/http"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/order"
	server_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/server"
)

type handler struct {
	service orderService
}

type orderService interface {
	CreateOrder(ctx context.Context, o order.Order) (order.Order, error)
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
	}
}
