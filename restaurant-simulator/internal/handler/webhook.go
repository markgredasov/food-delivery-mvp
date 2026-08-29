package handler

import (
	"avito-kitchen-restaurant-simulator/internal/client"
	errs "avito-kitchen-restaurant-simulator/internal/errors"
	"avito-kitchen-restaurant-simulator/internal/logger"
	"avito-kitchen-restaurant-simulator/internal/server/http/request"
	"avito-kitchen-restaurant-simulator/internal/server/http/response"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (h *handler) Webhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	rh := response.NewHTTPResponseHandler(log, w)

	var payload client.WebhookPayload
	if err := request.Decode(r, &payload); err != nil {
		err = fmt.Errorf("malformed webhook payload: %w", err)
		rh.ErrorResponse(err)
		return
	}

	if payload.OrderID == "" {
		rh.ErrorResponse(errs.InvalidRequest("no order id provided"))
		return
	}

	log.Info("webhook received", zap.String("order_id", payload.OrderID))
	h.proc.HandleNewOrder(ctx, payload.OrderID)
	rh.JSONResponse("", http.StatusOK)
}
