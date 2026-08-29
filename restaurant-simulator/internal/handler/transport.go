package handler

import (
	server_http "avito-kitchen-restaurant-simulator/internal/server/http/server"
	"context"
	"net/http"
)

type handler struct {
	proc Processor
}

type Processor interface {
	HandleNewOrder(ctx context.Context, orderID string)
}

func New(proc Processor) *handler {
	return &handler{
		proc: proc,
	}
}

func (h *handler) Routes() []server_http.Route {
	return []server_http.Route{
		{
			Method:  http.MethodGet,
			Path:    "/health",
			Handler: h.Health,
		},
		{
			Method:  http.MethodPost,
			Path:    "/webhook/orders",
			Handler: h.Webhook,
		},
	}
}
