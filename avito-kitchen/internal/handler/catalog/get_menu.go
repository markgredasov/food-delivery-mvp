package cataloghandler

import (
	"net/http"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/logger"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/request"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/response"
)

// GetMenu godoc
// @Summary Меню заведения
// @Description Возвращает меню заведения
// @Tags catalog
// @Accept json
// @Produce json
// @Param 	id	 path string  true "Идентификатор заведения" Format(uuid)
// @Success 200 {object} MenuItemsDTO "Меню заведения"
// @Success 400 {object} response.ErrorResponse "Невалидный запрос (неправильный формат UUID ресторана)"
// @Success 404 {object} response.ErrorResponse "Заведение не найдено"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /restaurants/{id}/menu [get]
func (h *handler) GetMenu(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	rh := response.NewHTTPResponseHandler(log, w)

	restaurantID, err := request.GetUUIDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	m, err := h.service.GetMenu(ctx, restaurantID)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	rh.JSONResponse(toMenuItems(m), http.StatusOK)
}
