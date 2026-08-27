package cataloghandler

import (
	"context"
	"net/http"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
	server_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/server"
)

type handler struct {
	service catalogService
}

type catalogService interface {
	ListActiveRestaurants(ctx context.Context) ([]restaurant.Restaurant, error)
	GetRestaurant(ctx context.Context, restaurantID string) (restaurant.Restaurant, error)
	GetMenu(ctx context.Context, restaurantID string) ([]menu.MenuItem, error)
	ListCategories(ctx context.Context) ([]menu.Category, error)
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
		{
			Method:  http.MethodGet,
			Path:    "/restaurants/{id}",
			Handler: h.GetRestaurant,
		},
		{
			Method:  http.MethodGet,
			Path:    "/restaurants/{id}/menu",
			Handler: h.GetMenu,
		},
		{
			Method:  http.MethodGet,
			Path:    "/categories",
			Handler: h.ListCategories,
		},
	}
}
