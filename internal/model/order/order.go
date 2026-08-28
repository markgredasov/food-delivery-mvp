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
	OrderStatusConfirmed  Status = "confirmed"
	OrderStatusPreparing  Status = "preparing"
	OrderStatusReady      Status = "ready"
	OrderStatusInDelivery Status = "in_delivery"
	OrderStatusDelivered  Status = "delivered"
	OrderStatusCancelled  Status = "cancelled"
	OrderStatusRejected   Status = "rejected"
)

// Valid reports if s is a valid status.
func (s Status) Valid() bool {
	switch s {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusPreparing, OrderStatusReady,
		OrderStatusInDelivery, OrderStatusDelivered, OrderStatusCancelled, OrderStatusRejected:
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
