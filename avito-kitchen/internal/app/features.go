package app

import server_http "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/server"

func (a *App) initFeatures() ([]server_http.Route, error) {
	var routes []server_http.Route

	initFuncs := []func() ([]server_http.Route, error){}

	for _, initFunc := range initFuncs {
		featureRoutes, err := initFunc()
		if err != nil {
			return routes, err
		}
		routes = append(routes, featureRoutes...)
	}

	return routes, nil
}
