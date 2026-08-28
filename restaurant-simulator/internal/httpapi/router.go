// Package httpapi is the restaurant-simulator's inbound HTTP surface: just
// the webhook endpoint the main service pushes new orders to.
package httpapi

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"avito-kitchen-restaurant-simulator/internal/client"
	"avito-kitchen-restaurant-simulator/internal/simulate"
)

// NewRouter builds the simulator's HTTP handler.
func NewRouter(proc *simulate.Processor, log *zap.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /webhook/orders", handleWebhook(proc, log))
	mux.HandleFunc("GET /health", handleHealth)
	return mux
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// handleWebhook accepts a pushed order and hands it to proc. It always
// acknowledges with 200 once the payload is valid JSON — acceptance is
// decided asynchronously — and relies on Processor.HandleNewOrder's claim
// check for idempotency against duplicate deliveries/retries.
func handleWebhook(proc *simulate.Processor, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload client.WebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Warn("malformed webhook payload", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if payload.OrderID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Info("webhook received", zap.String("order_id", payload.OrderID))
		proc.HandleNewOrder(r.Context(), payload.OrderID)
		w.WriteHeader(http.StatusOK)
	}
}
