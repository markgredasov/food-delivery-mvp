package order

import (
	"time"

	"github.com/google/uuid"
	errs "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/errors"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/address"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/money"
)

// OrderStatus is a step in the order lifecycle.
type OrderStatus string

// Order lifecycle statuses.
const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusPreparing  OrderStatus = "preparing"
	OrderStatusReady      OrderStatus = "ready"
	OrderStatusInDelivery OrderStatus = "in_delivery"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRejected   OrderStatus = "rejected"
)

// Valid reports if s is a valid status.
func (s OrderStatus) Valid() bool {
	switch s {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusPreparing, OrderStatusReady,
		OrderStatusInDelivery, OrderStatusDelivered, OrderStatusCancelled, OrderStatusRejected:
		return true
	}
	return false
}

// IsTerminal checks if order is processed.
func (s OrderStatus) IsTerminal() bool {
	return s == OrderStatusInDelivery || s == OrderStatusCancelled || s == OrderStatusRejected
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

// Total returns Price * Quantity.
func (i OrderItem) Total() money.Money {
	return i.Price.Mul(i.Quantity)
}

// Order is the aggregate root for a customer order placed against a single
// restaurant.
type Order struct {
	ID              uuid.UUID
	RestaurantID    uuid.UUID
	UserID          *string
	Status          OrderStatus
	DeliveryAddress address.Address
	TotalAmount     money.Money
	Items           []OrderItem
	Comment         *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// New constructs an order in OrderStatusPending with a server-computed total.
func New(id, restaurantID uuid.UUID, userID *string, deliveryAddr address.Address, items []OrderItem, comment *string) (Order, error) {
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
