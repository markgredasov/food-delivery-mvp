package simulate

import (
	"context"
	"time"

	"go.uber.org/zap"

	"avito-kitchen-restaurant-simulator/internal/client"
	"avito-kitchen-restaurant-simulator/internal/logger"
)

// Menu is the simulator's fixed catalog of 3 dishes, synced to the main
// service on startup and periodically thereafter.
var Menu = []client.MenuItemInput{
	{ID: "e96808cf-7f59-46bc-a958-0f9392f08e0c", Name: "Пицца Маргарита", Price: "550.00", Available: true},
	{ID: "e96808cf-7f59-46bc-a958-0f9392f08e0b", Name: "Кока-кола 0.5л", Price: "120.00", Available: true},
	{ID: "e96808cf-7f59-46bc-a958-0f9392f08e0a", Name: "Тирамису", Price: "320.00", Available: true},
}

func strPtr(s string) *string { return &s }

// SyncMenuLoop pushes Menu to the main service immediately, then again every
// interval until ctx is cancelled.
func SyncMenuLoop(ctx context.Context, c *client.Client, interval time.Duration) {
	log := logger.FromContext(ctx)
	syncOnce := func() {
		if err := c.SyncMenu(ctx, Menu); err != nil {
			log.Warn("menu sync failed", zap.Error(err))
			return
		}
		log.Info("menu synced")
	}

	syncOnce()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			syncOnce()
		}
	}
}

// PollLoop periodically fetches the restaurant's pending/active orders and
// hands each to proc — a fallback in case a webhook push was missed, since
// HandleNewOrder is idempotent for orders already claimed.
func PollLoop(ctx context.Context, c *client.Client, proc *Processor, interval time.Duration) {
	log := logger.FromContext(ctx)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			orders, err := c.ListPendingOrders(ctx)
			if err != nil {
				log.Warn("poll failed", zap.Error(err))
				continue
			}
			for _, o := range orders.Orders {
				proc.HandleNewOrder(ctx, o.ID)
			}
		}
	}
}
