package ordershandler

import (
	"net/http"

	"github.com/markgredasov/food-delivery-mvp/internal/logger"
	"github.com/markgredasov/food-delivery-mvp/internal/server/http/request"
	"github.com/markgredasov/food-delivery-mvp/internal/server/http/response"
)

// UpdateRestaurantMenu godoc
// @Summary Обновить меню целиком
// @Description Обновляет меню заведения целиком и возвращает обновленное меню. Блюда, которые отсутствуют в запросе, но присутствовали в меню помечаются как недоступные (available = false)
// @Tags restaurant
// @Accept json
// @Produce json
// @Param X-Restaurant-ID header string true "UUID ресторана. Замена авторизации в MVP." Format(uuid)
// @Param 	RequestBody	body UpdateMenuDTO true "Тело запроса"
// @Success 200 {object} UpdateMenuDTO "Обновленное меню"
// @Success 400 {object} response.ErrorResponse "Невалидный запрос"
// @Success 404 {object} response.ErrorResponse "Заведение не найдено"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /restaurants/menu [put]
func (h *handler) UpdateRestaurantMenu(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	rh := response.NewHTTPResponseHandler(log, w)

	restaurantID, err := request.GetUUIDFromHeader(r, RestaurantIDHeaderName)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	var req UpdateMenuDTO
	if err = request.Decode(r, &req); err != nil {
		rh.ErrorResponse(err)
		return
	}

	newItems, err := updateMenuToModel(req)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	items, err := h.catalog.ReplaceMenu(ctx, *restaurantID, newItems)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	rh.JSONResponse(updateMenuToDTO(items), http.StatusOK)
}
