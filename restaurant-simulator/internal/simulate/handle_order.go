package simulate

import (
	"avito-kitchen-restaurant-simulator/internal/logger"
	"context"
	"math/rand/v2"
	"time"

	"go.uber.org/zap"
)

// HandleNewOrder claims orderID (a no-op if already claimed by an earlier
// webhook delivery or poll) and, if this is the first sighting, spawns a
// goroutine that accepts it after a short random delay and then progresses
// it through the status pipeline.
func (p *Processor) HandleNewOrder(ctx context.Context, orderID string) {
	log := logger.FromContext(ctx)
	if !p.store.TryClaim(orderID) {
		return
	}
	log.Info("claimed new order", zap.String("order_id", orderID))
	go func() {
		bgCtx := context.Background()
		bgCtx = context.WithValue(bgCtx, "logger", log)
		p.acceptAndProgress(bgCtx, orderID)
	}()

}

func (p *Processor) acceptAndProgress(ctx context.Context, orderID string) {
	log := logger.FromContext(ctx)
	sleep(ctx, randBetween(5*time.Second, 15*time.Second))

	if err := p.client.Accept(ctx, orderID); err != nil {
		log.Warn("failed to accept order", zap.String("order_id", orderID), zap.Error(err))
		return
	}
	log.Info("accepted order", zap.String("order_id", orderID))

	for _, status := range pipeline {
		sleep(ctx, randBetween(3*time.Second, 8*time.Second))
		if ctx.Err() != nil {
			return
		}
		if err := p.client.UpdateStatus(ctx, orderID, status); err != nil {
			log.Warn("failed to advance order status",
				zap.String("order_id", orderID), zap.String("status", status), zap.Error(err))
			return
		}
		log.Info("advanced order status", zap.String("order_id", orderID), zap.String("status", status))
	}
}

// randBetween returns a random duration in [lo, hi). Used only to jitter
// simulated accept/progression delays, not for anything security-sensitive.
func randBetween(lo, hi time.Duration) time.Duration {
	if hi <= lo {
		return lo
	}
	return lo + time.Duration(rand.Int64N(int64(hi-lo))) //nolint:gosec
}

func sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}
