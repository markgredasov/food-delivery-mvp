package order_test

import (
	"testing"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/address"
	"github.com/markgredasov/food-delivery-mvp/internal/model/money"
	"github.com/markgredasov/food-delivery-mvp/internal/model/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    order.Status
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid pending status",
			input:   "pending",
			want:    order.StatusPending,
			wantErr: false,
		},
		{
			name:    "valid accepted status",
			input:   "accepted",
			want:    order.StatusAccepted,
			wantErr: false,
		},
		{
			name:    "invalid status",
			input:   "invalid",
			want:    "",
			wantErr: true,
			errMsg:  "invalid status",
		},
		{
			name:    "empty status",
			input:   "",
			want:    "",
			wantErr: true,
			errMsg:  "invalid status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := order.NewStatus(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.ErrorIs(t, err, errs.ErrInvalidArgument)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, status)
			}
		})
	}
}

func TestStatusMethods(t *testing.T) {
	tests := []struct {
		name         string
		status       order.Status
		valid        bool
		isTerminal   bool
		awaitingResp bool
	}{
		{
			name:         "pending",
			status:       order.StatusPending,
			valid:        true,
			isTerminal:   false,
			awaitingResp: true,
		},
		{
			name:         "sent_to_restaurant",
			status:       order.StatusSentToRestaurant,
			valid:        true,
			isTerminal:   false,
			awaitingResp: true,
		},
		{
			name:         "accepted",
			status:       order.StatusAccepted,
			valid:        true,
			isTerminal:   false,
			awaitingResp: false,
		},
		{
			name:         "preparing",
			status:       order.StatusPreparing,
			valid:        true,
			isTerminal:   false,
			awaitingResp: false,
		},
		{
			name:         "ready",
			status:       order.StatusReady,
			valid:        true,
			isTerminal:   false,
			awaitingResp: false,
		},
		{
			name:         "in_delivery",
			status:       order.StatusInDelivery,
			valid:        true,
			isTerminal:   true,
			awaitingResp: false,
		},
		{
			name:         "delivered",
			status:       order.StatusDelivered,
			valid:        true,
			isTerminal:   false,
			awaitingResp: false,
		},
		{
			name:         "rejected",
			status:       order.StatusRejected,
			valid:        true,
			isTerminal:   true,
			awaitingResp: false,
		},
		{
			name:         "invalid status",
			status:       order.Status("unknown"),
			valid:        false,
			isTerminal:   false,
			awaitingResp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.status.Valid())
			assert.Equal(t, tt.isTerminal, tt.status.IsTerminal())

			ord := createTestOrder(tt.status)
			assert.Equal(t, tt.awaitingResp, ord.AwaitingRestaurantResponse())
		})
	}
}

func TestNewOrderItem(t *testing.T) {
	menuItemID := uuid.New()

	tests := []struct {
		name       string
		menuItemID uuid.UUID
		quantity   int
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "valid item with quantity 1",
			menuItemID: menuItemID,
			quantity:   1,
			wantErr:    false,
		},
		{
			name:       "valid item with quantity 5",
			menuItemID: menuItemID,
			quantity:   5,
			wantErr:    false,
		},
		{
			name:       "zero quantity",
			menuItemID: menuItemID,
			quantity:   0,
			wantErr:    true,
			errMsg:     "quantity must be greater than 0",
		},
		{
			name:       "negative quantity",
			menuItemID: menuItemID,
			quantity:   -1,
			wantErr:    true,
			errMsg:     "quantity must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := order.NewOrderItem(tt.menuItemID, tt.quantity)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, item.ID)
				assert.Equal(t, tt.menuItemID, item.MenuItemID)
				assert.Equal(t, tt.quantity, item.Quantity)
			}
		})
	}
}

func TestNewOrder(t *testing.T) {
	restaurantID := uuid.New()
	userID := uuid.New()
	addr, _ := address.New("Moscow", "Tverskaya", "1", "10", nil)
	menuItemID := uuid.New()

	validItems := []order.OrderItem{
		{MenuItemID: menuItemID, Quantity: 2, Price: money.MustFromString("10.50")},
		{MenuItemID: menuItemID, Quantity: 1, Price: money.MustFromString("5.00")},
	}

	tests := []struct {
		name      string
		restID    uuid.UUID
		userID    *uuid.UUID
		addr      address.Address
		items     []order.OrderItem
		comment   *string
		wantErr   bool
		errMsg    string
		wantTotal string
	}{
		{
			name:      "valid order with items and comment",
			restID:    restaurantID,
			userID:    &userID,
			addr:      addr,
			items:     validItems,
			comment:   new(string),
			wantErr:   false,
			wantTotal: "26.00",
		},
		{
			name:      "valid order without comment",
			restID:    restaurantID,
			userID:    &userID,
			addr:      addr,
			items:     validItems,
			comment:   nil,
			wantErr:   false,
			wantTotal: "26.00",
		},
		{
			name:      "empty items list",
			restID:    restaurantID,
			userID:    &userID,
			addr:      addr,
			items:     []order.OrderItem{},
			comment:   nil,
			wantErr:   true,
			errMsg:    "order must contain at least one item",
			wantTotal: "",
		},
		{
			name:      "item with zero quantity",
			restID:    restaurantID,
			userID:    &userID,
			addr:      addr,
			items:     []order.OrderItem{{MenuItemID: menuItemID, Quantity: 0}},
			comment:   nil,
			wantErr:   true,
			errMsg:    "item quantity must be positive",
			wantTotal: "",
		},
		{
			name:      "single item order",
			restID:    restaurantID,
			userID:    nil,
			addr:      addr,
			items:     []order.OrderItem{{MenuItemID: menuItemID, Quantity: 3, Price: money.MustFromString("10.50")}},
			comment:   nil,
			wantErr:   false,
			wantTotal: "31.50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ord, err := order.New(uuid.New(), tt.restID, tt.userID, tt.addr, tt.items, tt.comment)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, order.StatusPending, ord.Status)
			assert.Equal(t, tt.wantTotal, ord.TotalAmount.String())
			assert.Len(t, tt.items, len(ord.Items))

			if tt.comment != nil {
				assert.Equal(t, *tt.comment, *ord.Comment)
			} else {
				assert.Nil(t, ord.Comment)
			}

			if tt.userID != nil {
				assert.Equal(t, *tt.userID, *ord.UserID)
			} else {
				assert.Nil(t, ord.UserID)
			}
		})
	}
}

func TestOrderStatusTransitions(t *testing.T) {
	tests := []struct {
		name           string
		initialStatus  order.Status
		action         func(*order.Order) error
		expectedStatus order.Status
		wantErr        bool
		errMsg         string
	}{
		{
			name:           "accept from pending",
			initialStatus:  order.StatusPending,
			action:         func(o *order.Order) error { return o.Accept() },
			expectedStatus: order.StatusAccepted,
			wantErr:        false,
		},
		{
			name:           "accept from sent_to_restaurant",
			initialStatus:  order.StatusSentToRestaurant,
			action:         func(o *order.Order) error { return o.Accept() },
			expectedStatus: order.StatusAccepted,
			wantErr:        false,
		},
		{
			name:           "accept from accepted (should fail)",
			initialStatus:  order.StatusAccepted,
			action:         func(o *order.Order) error { return o.Accept() },
			expectedStatus: order.StatusAccepted,
			wantErr:        true,
			errMsg:         "cannot be accepted",
		},
		{
			name:           "reject from pending",
			initialStatus:  order.StatusPending,
			action:         func(o *order.Order) error { return o.Reject() },
			expectedStatus: order.StatusRejected,
			wantErr:        false,
		},
		{
			name:           "reject from sent_to_restaurant",
			initialStatus:  order.StatusSentToRestaurant,
			action:         func(o *order.Order) error { return o.Reject() },
			expectedStatus: order.StatusRejected,
			wantErr:        false,
		},
		{
			name:           "reject from accepted (should fail)",
			initialStatus:  order.StatusAccepted,
			action:         func(o *order.Order) error { return o.Reject() },
			expectedStatus: order.StatusAccepted,
			wantErr:        true,
			errMsg:         "cannot be rejected",
		},
		{
			name:           "advance accepted -> preparing",
			initialStatus:  order.StatusAccepted,
			action:         func(o *order.Order) error { return o.AdvanceTo(order.StatusPreparing) },
			expectedStatus: order.StatusPreparing,
			wantErr:        false,
		},
		{
			name:           "advance preparing -> ready",
			initialStatus:  order.StatusPreparing,
			action:         func(o *order.Order) error { return o.AdvanceTo(order.StatusReady) },
			expectedStatus: order.StatusReady,
			wantErr:        false,
		},
		{
			name:           "advance ready -> in_delivery",
			initialStatus:  order.StatusReady,
			action:         func(o *order.Order) error { return o.AdvanceTo(order.StatusInDelivery) },
			expectedStatus: order.StatusInDelivery,
			wantErr:        false,
		},
		{
			name:           "advance in_delivery -> delivered",
			initialStatus:  order.StatusInDelivery,
			action:         func(o *order.Order) error { return o.AdvanceTo(order.StatusDelivered) },
			expectedStatus: order.StatusDelivered,
			wantErr:        false,
		},
		{
			name:           "advance accepted -> ready (skip)",
			initialStatus:  order.StatusAccepted,
			action:         func(o *order.Order) error { return o.AdvanceTo(order.StatusReady) },
			expectedStatus: order.StatusAccepted,
			wantErr:        true,
			errMsg:         "cannot transition",
		},
		{
			name:           "advance pending -> preparing (skip)",
			initialStatus:  order.StatusPending,
			action:         func(o *order.Order) error { return o.AdvanceTo(order.StatusPreparing) },
			expectedStatus: order.StatusPending,
			wantErr:        true,
			errMsg:         "cannot transition",
		},
		{
			name:           "advance to invalid status",
			initialStatus:  order.StatusAccepted,
			action:         func(o *order.Order) error { return o.AdvanceTo(order.Status("invalid")) },
			expectedStatus: order.StatusAccepted,
			wantErr:        true,
			errMsg:         "unknown order status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ord := createTestOrder(tt.initialStatus)
			err := tt.action(&ord)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, ord.Status)
			}
		})
	}
}

func createTestOrder(status order.Status) order.Order {
	addr, _ := address.New("Moscow", "Tverskaya", "1", "10", nil)
	items := []order.OrderItem{
		{MenuItemID: uuid.New(), Quantity: 1, Price: money.MustFromString("10.00")},
	}
	ord, _ := order.New(uuid.New(), uuid.New(), nil, addr, items, nil)
	ord.Status = status
	return ord
}
