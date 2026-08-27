package app

import (
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/middleware"
	server_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/server"
)

func (a *App) initHTTPServer(routes []server_http.Route) {
	apiVersion1Router := server_http.NewAPIVersionRouter(server_http.APIVersion1)
	apiVersion1Router.RegisterRoutes(routes...)

	serverConfig := server_http.NewConfigMust()
	server := server_http.NewHTTPServer(
		serverConfig,
		a.logger,
		middleware.RequestID(),
		middleware.Logger(a.logger),
		middleware.Recovery(),
		middleware.Trace(),
		middleware.CORS(),
	)
	server.RegisterAPIRouters(apiVersion1Router)
	server.RegisterSwagger()
	a.server = server
}
