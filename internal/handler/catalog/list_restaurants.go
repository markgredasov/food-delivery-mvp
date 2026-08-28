package cataloghandler

import (
	"net/http"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/logger"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/response"
)

// ListRestaurants godoc
// @Summary Список активных заведений
// @Description Возвращает список активных заведений
// @Tags catalog
// @Accept json
// @Produce json
// @Success 200 {object} RestaurantsDTO "Список зведений"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /restaurants [get]
func (h *handler) ListRestaurants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	rh := response.NewHTTPResponseHandler(log, w)

	restaurants, err := h.service.ListActiveRestaurants(ctx)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}
	rh.JSONResponse(toRestaurants(restaurants), http.StatusOK)
}
