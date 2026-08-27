package cataloghandler

import (
	"context"
	"net/http"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
	server_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/server"
)

type handler struct {
	service catalogService
}

type catalogService interface {
	ListActiveRestaurants(ctx context.Context) ([]restaurant.Restaurant, error)
}

func New(service catalogService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Routes() []server_http.Route {
	return []server_http.Route{
		{
			Method:  http.MethodGet,
			Path:    "/restaurants",
			Handler: h.ListRestaurants,
		},
	}
}
