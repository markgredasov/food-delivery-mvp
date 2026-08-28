package app

import (
	cataloghandler "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/handler/catalog"
	ordershandler "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/handler/orders"
	categoryrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/repository/postgres/category"
	menurepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/repository/postgres/menu"
	orderrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/repository/postgres/order"
	restaurantrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/repository/postgres/restaurant"
	server_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/server"
	catalogservice "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/service/catalog"
	orderservice "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/service/order"
)

func (a *App) initFeatures() []server_http.Route {
	var routes []server_http.Route

	restaurantRepo := restaurantrepo.New(a.pool)
	menuRepo := menurepo.New(a.pool)
	categoryRepo := categoryrepo.New(a.pool)
	catalogService := catalogservice.New(restaurantRepo, menuRepo, categoryRepo)
	catalogHandler := cataloghandler.New(catalogService)

	ordersRepo := orderrepo.New(a.pool)
	ordersService := orderservice.New(restaurantRepo, menuRepo, ordersRepo, a.pool, a.webhookClient)
	ordersHandler := ordershandler.New(ordersService)

	routes = append(routes, catalogHandler.Routes()...)
	routes = append(routes, ordersHandler.Routes()...)

	return routes
}
