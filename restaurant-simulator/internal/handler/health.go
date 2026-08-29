package handler

import (
	"avito-kitchen-restaurant-simulator/internal/logger"
	"avito-kitchen-restaurant-simulator/internal/server/http/response"
	"net/http"
)

func (h *handler) Health(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	rh := response.NewHTTPResponseHandler(log, w)
	rh.JSONResponse("", http.StatusOK)
}
