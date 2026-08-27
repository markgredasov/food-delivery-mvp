package server_http

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/docs"
)

func (h *HTTPServer) RegisterSwagger() {
	h.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		),
	)

	h.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			//nolint:errcheck // not needed
			w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}
