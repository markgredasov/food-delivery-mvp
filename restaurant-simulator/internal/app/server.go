package app

import (
	"avito-kitchen-restaurant-simulator/internal/server/http/middleware"
	server_http "avito-kitchen-restaurant-simulator/internal/server/http/server"
)

func (a *App) initHTTPServer(routes []server_http.Route) {
	apiVersion1Router := server_http.NewAPIVersionRouter(server_http.APIVersion1)
	apiVersion1Router.RegisterRoutes(routes...)

	cfg := server_http.NewConfigMust()
	server := server_http.NewHTTPServer(cfg, a.logger,
		middleware.RequestID(),
		middleware.Logger(a.logger),
		middleware.Recovery(),
		middleware.Trace(),
	)

	server.RegisterAPIRouters(apiVersion1Router)

	a.server = server
}
