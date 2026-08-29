package order

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/address"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/money"
)

// Status is a step in the order lifecycle.
type Status string

// Order lifecycle statuses.
const (
	OrderStatusPending    Status = "pending"
	OrderStatusConfirmed  Status = "sent_to_restaurant"
	OrderStatusAccepted   Status = "accepted"
	OrderStatusPreparing  Status = "preparing"
	OrderStatusReady      Status = "ready"
	OrderStatusInDelivery Status = "in_delivery"
	OrderStatusDelivered  Status = "delivered"
	OrderStatusRejected   Status = "rejected_by_restaurant"
)

// Valid reports if s is a valid status.
func (s Status) Valid() bool {
	switch s {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusAccepted, OrderStatusPreparing,
		OrderStatusReady, OrderStatusInDelivery, OrderStatusDelivered, OrderStatusRejected:
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
	return s == OrderStatusInDelivery || s == OrderStatusRejected
}

// awaitingRestaurantResponse are the statuses eligible for the 5-minute
// auto-reject timeout and for the restaurant polling fallback.
func (s Status) awaitingRestaurantResponse() bool {
	return s == OrderStatusPending || s == OrderStatusConfirmed
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
		Status:          OrderStatusPending,
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
	o.Status = OrderStatusAccepted
	return nil
}

// Reject marks the order as rejected by the restaurant (or by the
// auto-reject timeout). Allowed from the same states as Accept.
func (o *Order) Reject() error {
	if !o.Status.awaitingRestaurantResponse() {
		return errs.Conflict("order in status " + string(o.Status) + " cannot be rejected")
	}
	o.Status = OrderStatusRejected
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
