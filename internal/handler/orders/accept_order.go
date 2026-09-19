package ordershandler

import (
	"net/http"

	"github.com/markgredasov/food-delivery-mvp/internal/logger"
	"github.com/markgredasov/food-delivery-mvp/internal/server/http/request"
	"github.com/markgredasov/food-delivery-mvp/internal/server/http/response"
)

// AcceptOrder godoc
// @Summary Принять заказ
// @Description Принимает заказ
// @Tags restaurant
// @Accept json
// @Produce json
// @Param X-Restaurant-ID header string true "UUID ресторана. Замена авторизации в MVP." Format(uuid)
// @Param 	id	 path string  true "Идентификатор заказа" Format(uuid)
// @Success 200 {object} OrderDTO "Принятый заказ"
// @Success 400 {object} response.ErrorResponse "Невалидный запрос"
// @Success 403 {object} response.ErrorResponse "Заказ не принадлежит данному ресторану"
// @Success 404 {object} response.ErrorResponse "Заказ не найден"
// @Success 409 {object} response.ErrorResponse "Конфликт при обновлении статуса заказа"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /restaurants/orders/{id}/accept [post]
func (h *handler) AcceptOrder(w http.ResponseWriter, r *http.Request) {
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

	acceptedOrder, err := h.orders.AcceptOrder(ctx, *restaurantID, orderID)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	rh.JSONResponse(toOrder(acceptedOrder), http.StatusOK)
}
