package app

import (
	cataloghandler "github.com/markgredasov/food-delivery-mvp/internal/handler/catalog"
	ordershandler "github.com/markgredasov/food-delivery-mvp/internal/handler/orders"
	categoryrepo "github.com/markgredasov/food-delivery-mvp/internal/repository/postgres/category"
	menurepo "github.com/markgredasov/food-delivery-mvp/internal/repository/postgres/menu"
	orderrepo "github.com/markgredasov/food-delivery-mvp/internal/repository/postgres/order"
	restaurantrepo "github.com/markgredasov/food-delivery-mvp/internal/repository/postgres/restaurant"
	server_http "github.com/markgredasov/food-delivery-mvp/internal/server/http/server"
	catalogservice "github.com/markgredasov/food-delivery-mvp/internal/service/catalog"
	orderservice "github.com/markgredasov/food-delivery-mvp/internal/service/order"
)

func (a *App) initFeatures() []server_http.Route {
	var routes []server_http.Route

	restaurantRepo := restaurantrepo.New(a.pool)
	menuRepo := menurepo.New(a.pool)
	categoryRepo := categoryrepo.New(a.pool)
	catalogService := catalogservice.New(restaurantRepo, menuRepo, categoryRepo, a.pool)
	catalogHandler := cataloghandler.New(catalogService)

	ordersRepo := orderrepo.New(a.pool)
	ordersService := orderservice.New(restaurantRepo, menuRepo, ordersRepo, a.pool, a.webhookClient)
	ordersHandler := ordershandler.New(ordersService, catalogService)

	routes = append(routes, catalogHandler.Routes()...)
	routes = append(routes, ordersHandler.Routes()...)

	return routes
}
