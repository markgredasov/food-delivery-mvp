package ordershandler

import (
	"errors"
	"net/http"

	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/logger"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/request"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/server/http/response"
)

const UserIDHeaderName = "X-User-ID"

// CreateOrder godoc
// @Summary Создать заказ
// @Description Создает заказ
// @Tags orders
// @Accept json
// @Produce json
// @Param 		RequestBody		body		CreateOrderDTORequest		true	"Тело запроса"
// @Success 201 {object} CreateOrderDTOResponse "Созданный заказ"
// @Success 400 {object} response.ErrorResponse "Невалидный запрос"
// @Success 404 {object} response.ErrorResponse "Не найдено"
// @Success 409 {object} response.ErrorResponse "Конфликт"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders [post]
func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	rh := response.NewHTTPResponseHandler(log, w)

	var req CreateOrderDTORequest
	if err := request.Decode(r, &req); err != nil {
		rh.ErrorResponse(err)
		return
	}

	userID, err := request.GetUserIDFromHeader(r, UserIDHeaderName)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			userID = nil
		} else {
			rh.ErrorResponse(err)
			return
		}
	}

	reqOrder, err := createOrderToModel(req, userID)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	createdOrder, err := h.service.CreateOrder(ctx, reqOrder)
	if err != nil {
		rh.ErrorResponse(err)
		return
	}

	rh.JSONResponse(createOrderToDTO(createdOrder), http.StatusCreated)
}
