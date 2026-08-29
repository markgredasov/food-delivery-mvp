package simulate

import (
	"context"
	"time"

	"go.uber.org/zap"

	"avito-kitchen-restaurant-simulator/internal/client"
)

// Menu is the simulator's fixed catalog of 3 dishes, synced to the main
// service on startup and periodically thereafter.
var Menu = []client.MenuItemInput{
	{Name: "Пицца Маргарита", Price: "550.00", Available: true, CategoryID: strPtr("22222222-2222-2222-2222-222222222221")},
	{Name: "Кока-кола 0.5л", Price: "120.00", Available: true, CategoryID: strPtr("22222222-2222-2222-2222-222222222222")},
	{Name: "Тирамису", Price: "320.00", Available: true, CategoryID: strPtr("22222222-2222-2222-2222-222222222223")},
}

func strPtr(s string) *string { return &s }

// SyncMenuLoop pushes Menu to the main service immediately, then again every
// interval until ctx is cancelled.
func SyncMenuLoop(ctx context.Context, c *client.Client, interval time.Duration, log *zap.Logger) {
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
func PollLoop(ctx context.Context, c *client.Client, proc *Processor, interval time.Duration, log *zap.Logger) {
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
