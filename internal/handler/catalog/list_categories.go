package cataloghandler

import (
	"net/http"

	"github.com/markgredasov/food-delivery-mvp/internal/logger"
	"github.com/markgredasov/food-delivery-mvp/internal/server/http/response"
)

// ListCategories godoc
// @Summary Список категорий
// @Description Возвращает список зарегистрированных в сервисе категорий
// @Tags catalog
// @Accept json
// @Produce json
// @Success 200 {object} RestaurantDTO "Информация о заведении"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories [get]
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	rh := response.NewHTTPResponseHandler(log, w)

	categories, err := h.service.ListCategories(ctx)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}
	rh.JSONResponse(toCategories(categories), http.StatusOK)
}
