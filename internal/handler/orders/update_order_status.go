package ordershandler

import (
	"net/http"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/logger"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/request"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/response"
)

// UpdateStatus godoc
// @Summary Обновить статус заказа
// @Description Обновляет статус заказа на указанный
// @Tags restaurant
// @Accept json
// @Produce json
// @Param X-Restaurant-ID header string true "UUID ресторана. Замена авторизации в MVP." Format(uuid)
// @Param 	id	 path string  true "Идентификатор заказа" Format(uuid)
// @Param 	RequestBody	body UpdateOrderStatusDTO true "Тело запроса"
// @Success 200 {object} OrderDTO "Отклоненный заказ"
// @Success 400 {object} response.ErrorResponse "Невалидный запрос"
// @Success 403 {object} response.ErrorResponse "Заказ не принадлежит данному ресторану"
// @Success 404 {object} response.ErrorResponse "Заказ не найден"
// @Success 409 {object} response.ErrorResponse "Конфликт при обновлении статуса заказа (неверно задан следующий статус обновления)"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /restaurants/orders/{id}/status [patch]
func (h *handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	rh := response.NewHTTPResponseHandler(log, w)

	restaurantID, err := request.GetUUIDFromHeader(r, RestaurantIDHeaderName)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	orderID, err := request.GetUUIDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	var req UpdateOrderStatusDTO
	if err = request.Decode(r, &req); err != nil {
		rh.ErrorResponse(err)
		return
	}

	acceptedOrder, err := h.service.UpdateStatus(ctx, *restaurantID, orderID, req.Status)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	rh.JSONResponse(toOrder(acceptedOrder), http.StatusOK)
}
