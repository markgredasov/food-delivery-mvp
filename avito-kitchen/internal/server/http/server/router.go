package server_http

import (
	"fmt"
	"net/http"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/middleware"
)

type APIVersion string

var (
	APIVersion1 = APIVersion("v1")
)

type APIVersionRouter struct {
	*http.ServeMux

	apiVersion APIVersion
}

func NewAPIVersionRouter(apiVersion APIVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, middleware.ChainMiddleware(route.Handler, route.Middleware...))
	}
}
