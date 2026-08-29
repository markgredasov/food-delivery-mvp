// Package simulate contains the restaurant-simulator's business logic:
// reacting to new orders (via webhook push or polling fallback) by
// auto-accepting them and progressing them through the post-acceptance
// status pipeline, plus the periodic menu-sync and polling loops.
package simulate

import (
	"context"

	"avito-kitchen-restaurant-simulator/internal/client"
)

// pipeline is the sequence of statuses the simulator advances an accepted
// order through, in order.
var pipeline = []string{"preparing", "ready", "in_delivery", "delivered"}

// Processor reacts to new orders discovered via webhook or polling.
type Processor struct {
	client Client
	store  Store
}

// Store is repository interface.
type Store interface {
	TryClaim(orderID string) bool
}

// Client is a small interface of HTTP client service uses to
// call back into the Avito.Kitchen
type Client interface {
	Accept(ctx context.Context, orderID string) error
	UpdateStatus(ctx context.Context, orderID, status string) error
}

// NewProcessor builds a Processor.
func NewProcessor(c *client.Client, s Store) *Processor {
	return &Processor{client: c, store: s}
}
