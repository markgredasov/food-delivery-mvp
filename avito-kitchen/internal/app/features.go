package app

import (
	cataloghandler "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/handler/catalog"
	menurepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/repository/postgres/menu"
	restaurantrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/repository/postgres/restaurant"
	server_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/server"
	catalogservice "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/service/catalog"
)

func (a *App) initFeatures() []server_http.Route {
	var routes []server_http.Route

	catalogRoutes := a.initFeatureCatalog()
	routes = append(routes, catalogRoutes...)

	return routes
}

func (a *App) initFeatureCatalog() []server_http.Route {
	restaurantRepo := restaurantrepo.New(a.pool)
	menuRepo := menurepo.New(a.pool)
	service := catalogservice.New(restaurantRepo, menuRepo)
	handler := cataloghandler.New(service)

	return handler.Routes()
}
