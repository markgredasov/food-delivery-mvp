package order

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/address"
	"github.com/markgredasov/food-delivery-mvp/internal/model/money"
)

// Status is a step in the order lifecycle.
type Status string

// Order lifecycle statuses.
const (
	StatusPending          Status = "pending"
	StatusSentToRestaurant Status = "sent_to_restaurant"
	StatusAccepted         Status = "accepted"
	StatusPreparing        Status = "preparing"
	StatusReady            Status = "ready"
	StatusInDelivery       Status = "in_delivery"
	StatusDelivered        Status = "delivered"
	StatusRejected         Status = "rejected_by_restaurant"
)

// pipelineTransitions defines the linear post-acceptance progression that
// UpdateStatus is allowed to advance through. Accept/Reject are handled by
// their own dedicated methods below, since they apply from two possible
// source statuses and carry their own business rules.
var pipelineTransitions = map[Status]Status{ //nolint:exhaustive // do not need to check all statuses
	StatusAccepted:   StatusPreparing,
	StatusPreparing:  StatusReady,
	StatusReady:      StatusInDelivery,
	StatusInDelivery: StatusDelivered,
}

// NewStatus builds and validates order Status.
func NewStatus(statusString string) (Status, error) {
	status := Status(statusString)
	if !status.Valid() {
		return "", errs.InvalidArgument(fmt.Sprintf("invalid status: %s", statusString))
	}
	return status, nil
}

// Valid reports if s is a valid status.
func (s Status) Valid() bool {
	switch s {
	case StatusPending, StatusSentToRestaurant, StatusAccepted, StatusPreparing,
		StatusReady, StatusInDelivery, StatusDelivered, StatusRejected:
		return true
	}
	return false
}

// String converts OrderStatus to string.
func (s Status) String() string {
	return string(s)
}

// IsTerminal checks if order is processed.
func (s Status) IsTerminal() bool {
	return s == StatusInDelivery || s == StatusRejected
}

// awaitingRestaurantResponse are the statuses eligible for the 5-minute
// auto-reject timeout and for the restaurant polling fallback.
func (s Status) awaitingRestaurantResponse() bool {
	return s == StatusPending || s == StatusSentToRestaurant
}

// OrderItem is a line item snapshot: it freezes the menu item's name and
// price at the moment the order was placed so later menu edits never affect
// existing orders.
type OrderItem struct {
	ID         uuid.UUID
	OrderID    uuid.UUID
	MenuItemID uuid.UUID
	Name       string
	Price      money.Money
	Quantity   int
}

// NewOrderItem builds and validates OrderItem.
func NewOrderItem(menuItemID uuid.UUID, quantity int) (OrderItem, error) {
	if quantity <= 0 {
		return OrderItem{}, errs.InvalidArgument("quantity must be greater than 0")
	}
	id := uuid.New()

	return OrderItem{
		ID:         id,
		MenuItemID: menuItemID,
		Quantity:   quantity,
	}, nil
}

// Total returns Price * Quantity.
func (i OrderItem) Total() money.Money {
	return i.Price.Mul(i.Quantity)
}

// Order is the aggregate root for a customer order placed against a single
// restaurant.
type Order struct {
	ID              uuid.UUID
	RestaurantID    uuid.UUID
	UserID          *uuid.UUID
	Status          Status
	DeliveryAddress address.Address
	TotalAmount     money.Money
	Items           []OrderItem
	Comment         *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// New constructs an order in OrderStatusPending with a server-computed total.
func New(id, restaurantID uuid.UUID, userID *uuid.UUID, deliveryAddr address.Address, items []OrderItem, comment *string) (Order, error) {
	if len(items) == 0 {
		return Order{}, errs.InvalidArgument("order must contain at least one item")
	}
	total := money.Zero
	for _, it := range items {
		if it.Quantity <= 0 {
			return Order{}, errs.InvalidArgument("item quantity must be positive")
		}
		total = total.Add(it.Total())
	}

	return Order{
		ID:              id,
		RestaurantID:    restaurantID,
		UserID:          userID,
		Status:          StatusPending,
		DeliveryAddress: deliveryAddr,
		TotalAmount:     total,
		Items:           items,
		Comment:         comment,
	}, nil
}

// Accept marks the order as accepted by the restaurant. Allowed from
// pending (webhook delivery failed, restaurant discovered the order via
// polling) or sent_to_restaurant (the normal push path).
func (o *Order) Accept() error {
	if !o.Status.awaitingRestaurantResponse() {
		return errs.Conflict("order in status " + string(o.Status) + " cannot be accepted")
	}
	o.Status = StatusAccepted
	return nil
}

// Reject marks the order as rejected by the restaurant (or by the
// auto-reject timeout). Allowed from the same states as Accept.
func (o *Order) Reject() error {
	if !o.Status.awaitingRestaurantResponse() {
		return errs.Conflict("order in status " + string(o.Status) + " cannot be rejected")
	}
	o.Status = StatusRejected
	return nil
}

// IsOwnedBy reports whether restaurantID is the restaurant this order was
// placed against — used to guard restaurant-side endpoints against acting
// on another restaurant's order.
func (o *Order) IsOwnedBy(restaurantID uuid.UUID) bool {
	return o.RestaurantID == restaurantID
}

// AwaitingRestaurantResponse reports whether the order is still waiting for
// the restaurant to accept/reject it.
func (o *Order) AwaitingRestaurantResponse() bool {
	return o.Status.awaitingRestaurantResponse()
}

// AdvanceTo moves the order forward through the post-acceptance pipeline
// (accepted -> preparing -> ready -> in_delivery -> delivered). It rejects
// any attempt to skip a step or to reach pending/sent_to_restaurant/accepted
// /rejected_by_restaurant through this generic path — those are reached via
// New/Accept/Reject respectively.
func (o *Order) AdvanceTo(next Status) error {
	if !next.Valid() {
		return errs.InvalidArgument("unknown order status " + string(next))
	}
	want, ok := pipelineTransitions[o.Status]
	if !ok || want != next {
		return errs.Conflict("cannot transition order from " + string(o.Status) + " to " + string(next))
	}
	o.Status = next
	return nil
}
