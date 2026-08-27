package cataloghandler

import (
	"net/http"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/logger"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/request"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/response"
)

// GetRestaurant godoc
// @Summary Детали заведения
// @Description Возвращает информацию о заведении
// @Tags catalog
// @Accept json
// @Produce json
// @Param 	id	 path string  true "Идентификатор заведения" Format(uuid)
// @Success 200 {object} RestaurantDTOResponse "Информация о заведении"
// @Success 400 {object} response.ErrorResponse "Невалидный запрос (неправильный формат UUID ресторана)"
// @Success 404 {object} response.ErrorResponse "Заведение не найдено"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /restaurants/{id} [get]
func (h *handler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	rh := response.NewHTTPResponseHandler(log, w)

	restaurantID, err := request.GetUUIDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err, err.Error())
		return
	}

	restaurant, err := h.service.GetRestaurant(ctx, restaurantID)
	if err != nil {
		rh.ErrorResponse(err, err.Error())
		return
	}
	rh.JSONResponse(toRestaurant(restaurant), http.StatusOK)
}
