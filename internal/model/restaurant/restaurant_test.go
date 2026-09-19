package restaurant_test

import (
	"testing"

	"github.com/google/uuid"
	errs "github.com/markgredasov/food-delivery-mvp/internal/errors"
	"github.com/markgredasov/food-delivery-mvp/internal/model/address"
	"github.com/markgredasov/food-delivery-mvp/internal/model/restaurant"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    restaurant.Status
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid active status",
			input:   "active",
			want:    restaurant.RestaurantStatusActive,
			wantErr: false,
		},
		{
			name:    "valid inactive status",
			input:   "inactive",
			want:    restaurant.RestaurantStatusInactive,
			wantErr: false,
		},
		{
			name:    "valid blocked status",
			input:   "blocked",
			want:    restaurant.RestaurantStatusBlocked,
			wantErr: false,
		},
		{
			name:    "invalid status",
			input:   "unknown",
			want:    "",
			wantErr: true,
			errMsg:  "status 'unknown'",
		},
		{
			name:    "empty status",
			input:   "",
			want:    "",
			wantErr: true,
			errMsg:  "status ''",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := restaurant.NewStatus(tt.input)

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

func TestNewRestaurant(t *testing.T) {
	restaurantID := uuid.New()
	addr, _ := address.New("Moscow", "Tverskaya", "1", "10", nil)
	description := "Best pizza in town"

	tests := []struct {
		name         string
		id           uuid.UUID
		restName     string
		description  *string
		address      address.Address
		deliveryTime int
		status       restaurant.Status
		serviceURL   string
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid restaurant with all fields",
			id:           restaurantID,
			restName:     "Pizza House",
			description:  &description,
			address:      addr,
			deliveryTime: 30,
			status:       restaurant.RestaurantStatusActive,
			serviceURL:   "https://pizza-house.com/webhook",
			wantErr:      false,
		},
		{
			name:         "valid restaurant without description",
			id:           restaurantID,
			restName:     "Pizza House",
			description:  nil,
			address:      addr,
			deliveryTime: 30,
			status:       restaurant.RestaurantStatusActive,
			serviceURL:   "https://pizza-house.com/webhook",
			wantErr:      false,
		},
		{
			name:         "empty name",
			id:           restaurantID,
			restName:     "",
			description:  nil,
			address:      addr,
			deliveryTime: 30,
			status:       restaurant.RestaurantStatusActive,
			serviceURL:   "https://pizza-house.com/webhook",
			wantErr:      true,
			errMsg:       "name must not be empty",
		},
		{
			name:         "name too long",
			id:           restaurantID,
			restName:     string(make([]rune, 101)),
			description:  nil,
			address:      addr,
			deliveryTime: 30,
			status:       restaurant.RestaurantStatusActive,
			serviceURL:   "https://pizza-house.com/webhook",
			wantErr:      true,
			errMsg:       "name is too long",
		},
		{
			name:         "zero delivery time",
			id:           restaurantID,
			restName:     "Pizza House",
			description:  nil,
			address:      addr,
			deliveryTime: 0,
			status:       restaurant.RestaurantStatusActive,
			serviceURL:   "https://pizza-house.com/webhook",
			wantErr:      true,
			errMsg:       "delivery time must not be zero or less",
		},
		{
			name:         "negative delivery time",
			id:           restaurantID,
			restName:     "Pizza House",
			description:  nil,
			address:      addr,
			deliveryTime: -5,
			status:       restaurant.RestaurantStatusActive,
			serviceURL:   "https://pizza-house.com/webhook",
			wantErr:      true,
			errMsg:       "delivery time must not be zero or less",
		},
		{
			name:         "empty service URL",
			id:           restaurantID,
			restName:     "Pizza House",
			description:  nil,
			address:      addr,
			deliveryTime: 30,
			status:       restaurant.RestaurantStatusActive,
			serviceURL:   "",
			wantErr:      true,
			errMsg:       "service url must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rest, err := restaurant.New(
				tt.id,
				tt.restName,
				tt.description,
				tt.address,
				tt.deliveryTime,
				tt.status,
				tt.serviceURL,
			)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.id, rest.ID)
				assert.Equal(t, tt.restName, rest.Name)
				assert.Equal(t, tt.address, rest.Address)
				assert.Equal(t, tt.status, rest.Status)
				assert.Equal(t, tt.serviceURL, rest.ServiceURL)

				if tt.description != nil {
					assert.Equal(t, *tt.description, *rest.Description)
				} else {
					assert.Nil(t, rest.Description)
				}
			}
		})
	}
}

func TestRestaurant_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		status restaurant.Status
		want   bool
	}{
		{
			name:   "active restaurant",
			status: restaurant.RestaurantStatusActive,
			want:   true,
		},
		{
			name:   "inactive restaurant",
			status: restaurant.RestaurantStatusInactive,
			want:   false,
		},
		{
			name:   "blocked restaurant",
			status: restaurant.RestaurantStatusBlocked,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rest := createTestRestaurant(tt.status)
			assert.Equal(t, tt.want, rest.IsActive())
		})
	}
}

func TestRestaurant_EnsureAcceptsOrders(t *testing.T) {
	tests := []struct {
		name    string
		status  restaurant.Status
		wantErr bool
		errMsg  string
	}{
		{
			name:    "active restaurant accepts orders",
			status:  restaurant.RestaurantStatusActive,
			wantErr: false,
		},
		{
			name:    "inactive restaurant rejects orders",
			status:  restaurant.RestaurantStatusInactive,
			wantErr: true,
			errMsg:  "restaurant is not accepting orders",
		},
		{
			name:    "blocked restaurant rejects orders",
			status:  restaurant.RestaurantStatusBlocked,
			wantErr: true,
			errMsg:  "restaurant is not accepting orders",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rest := createTestRestaurant(tt.status)
			err := rest.EnsureAcceptsOrders()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.ErrorIs(t, err, errs.ErrConflict)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRestaurant_HasServiceURL(t *testing.T) {
	tests := []struct {
		name       string
		serviceURL string
		want       bool
	}{
		{
			name:       "has webhook URL",
			serviceURL: "https://pizza-house.com/webhook",
			want:       true,
		},
		{
			name:       "empty webhook URL",
			serviceURL: "",
			want:       false,
		},
		{
			name:       "whitespace only",
			serviceURL: "   ",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, _ := address.New("Moscow", "Tverskaya", "1", "10", nil)
			rest, _ := restaurant.New(
				uuid.New(),
				"Pizza House",
				nil,
				addr,
				30,
				restaurant.RestaurantStatusActive,
				tt.serviceURL,
			)
			assert.Equal(t, tt.want, rest.HasServiceURL())
		})
	}
}

func createTestRestaurant(status restaurant.Status) restaurant.Restaurant {
	addr, _ := address.New("Moscow", "Tverskaya", "1", "10", nil)
	rest, _ := restaurant.New(
		uuid.New(),
		"Pizza House",
		nil,
		addr,
		30,
		status,
		"https://pizza-house.com/webhook",
	)
	return rest
}
