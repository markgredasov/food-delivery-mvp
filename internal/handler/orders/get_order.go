package ordershandler

import (
	"net/http"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/logger"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/request"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/response"
)

// GetOrder godoc
// @Summary Получить информацию о заказе
// @Description Возвращает подробную информацию о заказе
// @Tags orders
// @Accept json
// @Produce json
// @Param 	id	 path string  true "Идентификатор заказа" Format(uuid)
// @Success 200 {object} OrderDTO "Информация о заказе"
// @Success 400 {object} response.ErrorResponse "Невалидный запрос"
// @Success 404 {object} response.ErrorResponse "Заказ не найден"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders/{id} [get]
func (h *handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	rh := response.NewHTTPResponseHandler(log, w)

	orderID, err := request.GetUUIDPathValue(r, "id")
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	order, err := h.orders.GetOrder(ctx, orderID)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	rh.JSONResponse(toOrder(order), http.StatusOK)
}
