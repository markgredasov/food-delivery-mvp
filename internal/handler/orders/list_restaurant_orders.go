package ordershandler

import (
	"net/http"

	"github.com/markgredasov/food-delivery-mvp/internal/logger"
	"github.com/markgredasov/food-delivery-mvp/internal/server/http/request"
	"github.com/markgredasov/food-delivery-mvp/internal/server/http/response"
)

const RestaurantIDHeaderName = "X-Restaurant-ID"

// ListRestaurantOrders godoc
// @Summary Заказы заведения (pending + активные)
// @Description Возвращает активные заказы заведения
// @Tags restaurant
// @Accept json
// @Produce json
// @Param X-Restaurant-ID header string true "UUID ресторана. UUID заведения, от имени которого выполняется запрос (замена авторизации в MVP)." Format(uuid)
// @Success 200 {object} OrdersDTO "Информация об активных заказах"
// @Success 400 {object} response.ErrorResponse "Невалидный запрос (неверный формат UUID)"
// @Success 404 {object} response.ErrorResponse "Ресторан не найден"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /restaurants/orders [get]
func (h *handler) ListRestaurantOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	rh := response.NewHTTPResponseHandler(log, w)

	restaurantID, err := request.GetUUIDFromHeader(r, RestaurantIDHeaderName)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	orders, err := h.orders.ListRestaurantOrders(ctx, *restaurantID)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	rh.JSONResponse(toOrders(orders), http.StatusOK)
}
