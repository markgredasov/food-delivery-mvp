// Package simulate contains the restaurant-simulator's business logic:
// reacting to new orders (via webhook push or polling fallback) by
// auto-accepting them and progressing them through the post-acceptance
// status pipeline, plus the periodic menu-sync and polling loops.
package simulate

import (
	"context"
	"math/rand/v2"
	"time"

	"go.uber.org/zap"

	"avito-kitchen-restaurant-simulator/internal/client"
	"avito-kitchen-restaurant-simulator/internal/store"
)

// pipeline is the sequence of statuses the simulator advances an accepted
// order through, in order.
var pipeline = []string{"preparing", "ready", "in_delivery", "delivered"}

// Processor reacts to new orders discovered via webhook or polling.
type Processor struct {
	client *client.Client
	store  *store.Store
	log    *zap.Logger
}

// NewProcessor builds a Processor.
func NewProcessor(c *client.Client, s *store.Store, log *zap.Logger) *Processor {
	return &Processor{client: c, store: s, log: log}
}

// HandleNewOrder claims orderID (a no-op if already claimed by an earlier
// webhook delivery or poll) and, if this is the first sighting, spawns a
// goroutine that accepts it after a short random delay and then progresses
// it through the status pipeline.
func (p *Processor) HandleNewOrder(ctx context.Context, orderID string) {
	if !p.store.TryClaim(orderID) {
		return
	}
	p.log.Info("claimed new order", zap.String("order_id", orderID))
	go p.acceptAndProgress(ctx, orderID)
}

func (p *Processor) acceptAndProgress(ctx context.Context, orderID string) {
	sleep(ctx, randBetween(5*time.Second, 15*time.Second))

	if err := p.client.Accept(ctx, orderID); err != nil {
		p.log.Warn("failed to accept order", zap.String("order_id", orderID), zap.Error(err))
		return
	}
	p.log.Info("accepted order", zap.String("order_id", orderID))

	for _, status := range pipeline {
		sleep(ctx, randBetween(3*time.Second, 8*time.Second))
		if ctx.Err() != nil {
			return
		}
		if err := p.client.UpdateStatus(ctx, orderID, status); err != nil {
			p.log.Warn("failed to advance order status",
				zap.String("order_id", orderID), zap.String("status", status), zap.Error(err))
			return
		}
		p.log.Info("advanced order status", zap.String("order_id", orderID), zap.String("status", status))
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
